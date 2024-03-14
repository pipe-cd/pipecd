package main

import (
	"context"
	"time"

	"github.com/pipe-cd/pipecd/pkg/app/pipedv1/plugin/platform/kubernetes/planner"
	"github.com/pipe-cd/pipecd/pkg/cli"
	"github.com/pipe-cd/pipecd/pkg/rpc"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
)

type server struct {
	apiPort     int
	gracePeriod time.Duration
	tls         bool
	certFile    string
	keyFile     string

	enableGRPCReflection bool
}

// NewServerCommand creates a new cobra command for executing api server.
func NewServerCommand() *cobra.Command {
	s := &server{
		apiPort:     10000,
		gracePeriod: 30 * time.Second,
	}
	cmd := &cobra.Command{
		Use:   "server",
		Short: "Start running server.",
		RunE:  cli.WithContext(s.run),
	}

	cmd.Flags().IntVar(&s.apiPort, "api-port", s.apiPort, "The port number used to run a grpc server for external apis.")
	cmd.Flags().DurationVar(&s.gracePeriod, "grace-period", s.gracePeriod, "How long to wait for graceful shutdown.")

	cmd.Flags().BoolVar(&s.tls, "tls", s.tls, "Whether running the gRPC server with TLS or not.")
	cmd.Flags().StringVar(&s.certFile, "cert-file", s.certFile, "The path to the TLS certificate file.")
	cmd.Flags().StringVar(&s.keyFile, "key-file", s.keyFile, "The path to the TLS key file.")

	// For debugging early in development
	cmd.Flags().BoolVar(&s.enableGRPCReflection, "enable-grpc-reflection", s.enableGRPCReflection, "Whether to enable the reflection service or not.")

	return cmd
}

func (s *server) run(ctx context.Context, input cli.Input) (runErr error) {
	// Make a cancellable context.
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	group, ctx := errgroup.WithContext(ctx)

	// Start a gRPC server for handling external API requests.
	{
		var (
			service = planner.NewPlannerService(input.Logger)
			opts    = []rpc.Option{
				rpc.WithPort(s.apiPort),
				rpc.WithGracePeriod(s.gracePeriod),
				rpc.WithLogger(input.Logger),
				rpc.WithLogUnaryInterceptor(input.Logger),
				rpc.WithRequestValidationUnaryInterceptor(),
			}
		)
		if s.tls {
			opts = append(opts, rpc.WithTLS(s.certFile, s.keyFile))
		}
		if s.enableGRPCReflection {
			opts = append(opts, rpc.WithGRPCReflection())
		}
		if input.Flags.Metrics {
			opts = append(opts, rpc.WithPrometheusUnaryInterceptor())
		}

		server := rpc.NewServer(service, opts...)
		group.Go(func() error {
			return server.Run(ctx)
		})
	}

	if err := group.Wait(); err != nil {
		input.Logger.Error("failed while running", zap.Error(err))
		return err
	}
	return nil
}
