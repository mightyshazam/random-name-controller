package main

import (
	"context"
	"flag"
	"random_name_controller/pkg/server"

	"go.uber.org/zap"
)

var (
	flags = struct {
		listenAddress          string
		containerImage         string
		serviceaccount         string
		usePodIdentity         bool
		podIdentityLabel       string
		podIdentity            string
		secretStoreCredentials string
		masterUrl              string
		useVmIdentity          bool
	}{}
)

func init() {
	flag.StringVar(&flags.listenAddress, "listen-address", ":8080", "--listen-address 8080")
}

func main() {
	logger, _ := zap.NewProduction()
	host := server.New(&server.HostArgs{
		Logger:        logger,
		ListenAddress: flags.listenAddress,
	})
	if err := host.Run(context.Background()); err != nil {
		logger.Error("exited unexpectedly", zap.Error(err))
	}
}
