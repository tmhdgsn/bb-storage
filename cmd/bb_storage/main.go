package main

import (
	"log"
	"net/http"
	"os"

	remoteexecution "github.com/bazelbuild/remote-apis/build/bazel/remote/execution/v2"
	"github.com/buildbarn/bb-storage/pkg/ac"
	"github.com/buildbarn/bb-storage/pkg/blobstore/completenesschecking"
	blobstore "github.com/buildbarn/bb-storage/pkg/blobstore/configuration"
	"github.com/buildbarn/bb-storage/pkg/builder"
	"github.com/buildbarn/bb-storage/pkg/cas"
	"github.com/buildbarn/bb-storage/pkg/configuration"
	bb_grpc "github.com/buildbarn/bb-storage/pkg/grpc"
	"github.com/buildbarn/bb-storage/pkg/opencensus"
	"github.com/buildbarn/bb-storage/pkg/util"
	"github.com/gorilla/mux"

	"google.golang.org/genproto/googleapis/bytestream"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func main() {
	if len(os.Args) != 2 {
		log.Fatal("Usage: bb-storage bb-storage.jsonnet")
	}

	storageConfiguration, err := configuration.GetStorageConfiguration(os.Args[1])
	if err != nil {
		log.Fatalf("Failed to read configuration from %s: %s", os.Args[1], err)
	}

	if storageConfiguration.Jaeger != nil {
		opencensus.Initialize(storageConfiguration.Jaeger)
	}

	// Storage access.
	contentAddressableStorageBlobAccess, actionCache, err := blobstore.CreateBlobAccessObjectsFromConfig(
		storageConfiguration.Blobstore,
		int(storageConfiguration.MaximumMessageSizeBytes))
	if err != nil {
		log.Fatal("Failed to create blob access: ", err)
	}

	// If this instance of bb-storage has access to all data (as in,
	// it's not a single shard within a distributed setup), it can
	// be configured to verify that all objects referenced by
	// ActionResults are present in the Content Addressable Storage.
	// Such validation is required by Bazel.
	if storageConfiguration.VerifyActionResultCompleteness {
		actionCache = completenesschecking.NewCompletenessCheckingBlobAccess(
			actionCache,
			cas.NewBlobAccessContentAddressableStorage(
				contentAddressableStorageBlobAccess,
				int(storageConfiguration.MaximumMessageSizeBytes)),
			contentAddressableStorageBlobAccess,
			100,
			int(storageConfiguration.MaximumMessageSizeBytes))
	}

	// Let GetCapabilities() work, even for instances that don't
	// have a scheduler attached to them, but do allow uploading
	// results into the Action Cache.
	schedulers := map[string]builder.BuildQueue{}
	allowActionCacheUpdatesForInstances := map[string]bool{}
	if len(storageConfiguration.AllowAcUpdatesForInstances) > 0 {
		fallback := builder.NewNonExecutableBuildQueue()
		for _, instance := range storageConfiguration.AllowAcUpdatesForInstances {
			schedulers[instance] = fallback
			allowActionCacheUpdatesForInstances[instance] = true
		}
	}

	// Backends capable of compiling.
	for name, endpoint := range storageConfiguration.Schedulers {
		scheduler, err := bb_grpc.NewGRPCClientFromConfiguration(endpoint)
		if err != nil {
			log.Fatal("Failed to create scheduler RPC client: ", err)
		}
		schedulers[name] = builder.NewForwardingBuildQueue(scheduler)
	}
	buildQueue := builder.NewDemultiplexingBuildQueue(func(instance string) (builder.BuildQueue, error) {
		scheduler, ok := schedulers[instance]
		if !ok {
			return nil, status.Errorf(codes.InvalidArgument, "Unknown instance name")
		}
		return scheduler, nil
	})

	go func() {
		log.Fatal(
			"gRPC server failure: ",
			bb_grpc.NewGRPCServersFromConfigurationAndServe(
				storageConfiguration.GrpcServers,
				func(s *grpc.Server) {
					remoteexecution.RegisterActionCacheServer(s, ac.NewActionCacheServer(actionCache, allowActionCacheUpdatesForInstances, int(storageConfiguration.MaximumMessageSizeBytes)))
					remoteexecution.RegisterContentAddressableStorageServer(s, cas.NewContentAddressableStorageServer(contentAddressableStorageBlobAccess))
					bytestream.RegisterByteStreamServer(s, cas.NewByteStreamServer(contentAddressableStorageBlobAccess, 1<<16))
					remoteexecution.RegisterCapabilitiesServer(s, buildQueue)
					remoteexecution.RegisterExecutionServer(s, buildQueue)
				}))
	}()

	// Web server for metrics and profiling.
	router := mux.NewRouter()
	util.RegisterAdministrativeHTTPEndpoints(router)
	log.Fatal(http.ListenAndServe(storageConfiguration.MetricsListenAddress, router))
}
