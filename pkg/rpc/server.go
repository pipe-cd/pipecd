// Copyright 2020 The Pipe Authors.
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
	"fmt"
	"net"
	"time"

	"go.opencensus.io/plugin/ocgrpc"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	"github.com/kapetaniosci/pipe/pkg/jwt"
	"github.com/kapetaniosci/pipe/pkg/rpc/rpcauth"
)

// Service represents a gRPC service will be registered to server.
type Service interface {
	Register(server *grpc.Server)
}

// Server used to register gRPC services then start and serve incoming requests.
type Server struct {
	port       int
	tls        bool
	certFile   string
	keyFile    string
	services   []Service
	grpcServer *grpc.Server
	logger     *zap.Logger

	authUnaryInterceptor              grpc.UnaryServerInterceptor
	serviceKeyAuthUnaryInterceptor    grpc.UnaryServerInterceptor
	jwtAuthUnaryInterceptor           grpc.UnaryServerInterceptor
	requestValidationUnaryInterceptor grpc.UnaryServerInterceptor
	authStreamInterceptor             grpc.StreamServerInterceptor
}

// Option defines a function to set configurable field of Server.
type Option func(*Server)

// WithPort sets grpc port number.
func WithPort(port int) Option {
	return func(s *Server) {
		s.port = port
	}
}

// WithAuthUnaryInterceptor sets an interceptor for extracting token in imcoming metadata.
func WithAuthUnaryInterceptor() Option {
	return func(s *Server) {
		s.authUnaryInterceptor = rpcauth.UnaryServerInterceptor()
	}
}

// WithServiceKeyAuthUnaryInterceptor sets an interceptor for checking service key.
func WithServiceKeyAuthUnaryInterceptor(key string, logger *zap.Logger) Option {
	return func(s *Server) {
		s.serviceKeyAuthUnaryInterceptor = rpcauth.ServiceKeyUnaryServerInterceptor(key, logger)
	}
}

// WithJwtAuthUnaryInterceptor sets an interceprot for checking JWT token.
func WithJwtAuthUnaryInterceptor(verifier jwt.Verifier, authorizer rpcauth.Authorizer, logger *zap.Logger) Option {
	return func(s *Server) {
		s.jwtAuthUnaryInterceptor = rpcauth.JwtUnaryServerInterceptor(verifier, authorizer, logger)
	}
}

// WithRequestValidationUnaryInterceptor sets an interceptor for validating request payload.
func WithRequestValidationUnaryInterceptor() Option {
	return func(s *Server) {
		s.requestValidationUnaryInterceptor = RequestValidationUnaryServerInterceptor()
	}
}

// WithAuthStreamInterceptor sets an interceptor for extracting token in stream.
func WithAuthStreamInterceptor() Option {
	return func(s *Server) {
		s.authStreamInterceptor = rpcauth.StreamServerInterceptor()
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

// WithLogger sets logger to server.
func WithLogger(logger *zap.Logger) Option {
	return func(s *Server) {
		s.logger = logger
	}
}

// NewServer creates a new server for handling gPRC services.
func NewServer(service Service, opts ...Option) *Server {
	s := &Server{
		logger: zap.NewNop(),
	}
	for _, opt := range opts {
		opt(s)
	}
	s.logger = s.logger.Named("rpc-server")
	if service == nil {
		s.logger.Fatal("service must not be nil")
	}
	s.services = append(s.services, service)
	return s
}

// Run starts running gRPC server for handling incoming requests.
func (s *Server) Run() error {
	opts := []grpc.ServerOption{
		grpc.StatsHandler(&ocgrpc.ServerHandler{}),
	}
	// If tls option is enabled we load and use certificate and
	// key files from specified paths.
	if s.tls {
		creds, err := credentials.NewServerTLSFromFile(s.certFile, s.keyFile)
		if err != nil {
			s.logger.Fatal("failed to load tls certificate file", zap.Error(err))
		}
		opts = append(opts, grpc.Creds(creds))
	} else {
		s.logger.Info("grpc server will be run without tls")
	}
	// Builds a chain of enabled interceptors.
	var unaryInterceptors []grpc.UnaryServerInterceptor
	if s.authUnaryInterceptor != nil {
		unaryInterceptors = append(unaryInterceptors, s.authUnaryInterceptor)
	}
	if s.serviceKeyAuthUnaryInterceptor != nil {
		unaryInterceptors = append(unaryInterceptors, s.serviceKeyAuthUnaryInterceptor)
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
	if s.authStreamInterceptor != nil {
		opts = append(opts, grpc.StreamInterceptor(s.authStreamInterceptor))
	}
	s.grpcServer = grpc.NewServer(opts...)

	// Register all registered services.
	for _, service := range s.services {
		service.Register(s.grpcServer)
	}
	// Start open a tcp connection.
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

// Stop stops running gRPC server gracefully.
func (s *Server) Stop(timeout time.Duration) {
	ch := make(chan struct{})
	go func() {
		s.logger.Info("gracefulStop is running")
		s.grpcServer.GracefulStop()
		close(ch)
	}()

	select {
	case <-ch:
		s.logger.Info("gracefulStop completed before timing out")
	case <-time.After(timeout):
		s.logger.Info("force server to stop")
		s.grpcServer.Stop()
	}
}
