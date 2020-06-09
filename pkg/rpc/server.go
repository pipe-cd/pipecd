// Copyright 2020 The PipeCD Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package rpc

import (
	"context"
	"fmt"
	"net"
	"time"

	"go.opencensus.io/plugin/ocgrpc"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/reflection"

	"github.com/kapetaniosci/pipe/pkg/jwt"
	"github.com/kapetaniosci/pipe/pkg/rpc/rpcauth"
)

// Service represents a gRPC service will be registered to server.
type Service interface {
	Register(server *grpc.Server)
}

// Server used to register gRPC services then start and serve incoming requests.
type Server struct {
	port                 int
	tls                  bool
	certFile             string
	keyFile              string
	services             []Service
	grpcServer           *grpc.Server
	gracePeriod          time.Duration
	enabelGRPCReflection bool
	logger               *zap.Logger

	pipedKeyAuthUnaryInterceptor      grpc.UnaryServerInterceptor
	pipedKeyAuthStreamInterceptor     grpc.StreamServerInterceptor
	jwtAuthUnaryInterceptor           grpc.UnaryServerInterceptor
	requestValidationUnaryInterceptor grpc.UnaryServerInterceptor
}

// Option defines a function to set configurable field of Server.
type Option func(*Server)

// WithPort sets grpc port number.
func WithPort(port int) Option {
	return func(s *Server) {
		s.port = port
	}
}

// WithPipedTokenAuthUnaryInterceptor sets an interceptor for validating piped key.
func WithPipedTokenAuthUnaryInterceptor(verifier rpcauth.PipedTokenVerifier, logger *zap.Logger) Option {
	return func(s *Server) {
		s.pipedKeyAuthUnaryInterceptor = rpcauth.PipedTokenUnaryServerInterceptor(verifier, logger)
	}
}

// WithPipedTokenAuthStreamInterceptor sets an interceptor for validating piped key.
func WithPipedTokenAuthStreamInterceptor(verifier rpcauth.PipedTokenVerifier, logger *zap.Logger) Option {
	return func(s *Server) {
		s.pipedKeyAuthStreamInterceptor = rpcauth.PipedTokenStreamServerInterceptor(verifier, logger)
	}
}

// WithJWTAuthUnaryInterceptor sets an interceprot for checking JWT token.
func WithJWTAuthUnaryInterceptor(verifier jwt.Verifier, authorizer rpcauth.RBACAuthorizer, logger *zap.Logger) Option {
	return func(s *Server) {
		s.jwtAuthUnaryInterceptor = rpcauth.JWTUnaryServerInterceptor(verifier, authorizer, logger)
	}
}

// WithRequestValidationUnaryInterceptor sets an interceptor for validating request payload.
func WithRequestValidationUnaryInterceptor() Option {
	return func(s *Server) {
		s.requestValidationUnaryInterceptor = RequestValidationUnaryServerInterceptor()
	}
}

// WithTLS configures TLS files.
func WithTLS(certFile, keyFile string) Option {
	return func(s *Server) {
		s.tls = true
		s.certFile = certFile
		s.keyFile = keyFile
	}
}

// WithService appends gPRC service to server.
func WithService(service Service) Option {
	return func(s *Server) {
		s.services = append(s.services, service)
	}
}

// WithGracePeriod sets maximum time to wait for gracefully shutdown.
func WithGracePeriod(d time.Duration) Option {
	return func(s *Server) {
		s.gracePeriod = d
	}
}

// WithLogger sets logger to server.
func WithLogger(logger *zap.Logger) Option {
	return func(s *Server) {
		s.logger = logger.Named("rpc-server")
	}
}

// WithGRPCReflection enables gRPC reflection service for debugging.
func WithGRPCReflection() Option {
	return func(s *Server) {
		s.enabelGRPCReflection = true
	}
}

// NewServer creates a new server for handling gPRC services.
func NewServer(service Service, opts ...Option) *Server {
	s := &Server{
		gracePeriod: 15 * time.Second,
		logger:      zap.NewNop(),
	}
	for _, opt := range opts {
		opt(s)
	}
	s.services = append(s.services, service)
	if len(s.services) == 0 {
		s.logger.Fatal("at least one service must be specified")
	}
	if err := s.init(); err != nil {
		s.logger.Fatal(err.Error())
	}
	return s
}

// Run starts running gRPC server for handling incoming requests.
func (s *Server) Run(ctx context.Context) error {
	doneCh := make(chan error, 1)
	ctx, cancel := context.WithCancel(ctx)

	go func() {
		defer cancel()
		doneCh <- s.run()
	}()

	<-ctx.Done()
	s.stop()
	return <-doneCh
}

func (s *Server) init() error {
	opts := []grpc.ServerOption{
		grpc.StatsHandler(&ocgrpc.ServerHandler{}),
	}
	// If tls option is enabled we load and use certificate and
	// key files from specified paths.
	if s.tls {
		creds, err := credentials.NewServerTLSFromFile(s.certFile, s.keyFile)
		if err != nil {
			return fmt.Errorf("failed to load tls certificate file: %v", err)
		}
		opts = append(opts, grpc.Creds(creds))
	} else {
		s.logger.Info("grpc server will be run without tls")
	}
	// Builds a chain of enabled interceptors.
	var unaryInterceptors []grpc.UnaryServerInterceptor
	if s.pipedKeyAuthUnaryInterceptor != nil {
		unaryInterceptors = append(unaryInterceptors, s.pipedKeyAuthUnaryInterceptor)
	}
	if s.jwtAuthUnaryInterceptor != nil {
		unaryInterceptors = append(unaryInterceptors, s.jwtAuthUnaryInterceptor)
	}
	if s.requestValidationUnaryInterceptor != nil {
		unaryInterceptors = append(unaryInterceptors, s.requestValidationUnaryInterceptor)
	}
	if len(unaryInterceptors) > 0 {
		c := ChainUnaryServerInterceptors(unaryInterceptors...)
		opts = append(opts, grpc.UnaryInterceptor(c))
	}
	if s.pipedKeyAuthStreamInterceptor != nil {
		opts = append(opts, grpc.StreamInterceptor(s.pipedKeyAuthStreamInterceptor))
	}
	s.grpcServer = grpc.NewServer(opts...)

	// Register all registered services.
	for _, service := range s.services {
		service.Register(s.grpcServer)
	}
	if s.enabelGRPCReflection {
		reflection.Register(s.grpcServer)
	}

	return nil
}

func (s *Server) run() error {
	// Start listening at the specified port.
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", s.port))
	if err != nil {
		s.logger.Error("failed to listen", zap.Error(err))
		return err
	}
	// Start running gRPC server for serving.
	s.logger.Info(fmt.Sprintf("grpc server is running on %s", lis.Addr().String()))
	err = s.grpcServer.Serve(lis)
	if err != nil && err != grpc.ErrServerStopped {
		s.logger.Error("failed to serve", zap.Error(err))
		return err
	}
	return nil
}

// stop stops running gRPC server gracefully.
func (s *Server) stop() {
	ch := make(chan struct{})
	go func() {
		s.logger.Info("gracefulStop is running")
		s.grpcServer.GracefulStop()
		close(ch)
	}()

	select {
	case <-ch:
		s.logger.Info("gracefulStop completed before timing out")
	case <-time.After(s.gracePeriod):
		s.logger.Info("force server to stop")
		s.grpcServer.Stop()
	}
}
