// Copyright 2023 The PipeCD Authors.
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

package rpcclient

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

type option struct {
	tls                          bool
	certFile                     string
	requestValidationInterceptor bool
	options                      []grpc.DialOption
}

type DialOption func(*option)

func WithBlock() DialOption {
	return func(o *option) {
		o.options = append(o.options, grpc.WithBlock())
	}
}

func WithTLS(certFile string) DialOption {
	return func(o *option) {
		o.tls = true
		o.certFile = certFile
	}
}

func WithTransportCredentials(creds credentials.TransportCredentials) DialOption {
	return func(o *option) {
		o.options = append(o.options, grpc.WithTransportCredentials(creds))
	}
}

func WithInsecure() DialOption {
	return func(o *option) {
		o.options = append(o.options, grpc.WithInsecure())
	}
}

func WithRequestValidationInterceptor() DialOption {
	return func(o *option) {
		o.requestValidationInterceptor = true
	}
}

func WithPerRPCCredentials(creds credentials.PerRPCCredentials) DialOption {
	return func(o *option) {
		o.options = append(o.options, grpc.WithPerRPCCredentials(creds))
	}
}

func WithMaxRecvMsgSize(m int) DialOption {
	return func(o *option) {
		o.options = append(o.options, grpc.WithMaxMsgSize(m))
	}
}

func DialOptions(opts ...DialOption) ([]grpc.DialOption, error) {
	o := &option{
		options: []grpc.DialOption{},
	}
	for _, opt := range opts {
		opt(o)
	}
	if o.tls {
		cred, err := credentials.NewClientTLSFromFile(o.certFile, "")
		if err != nil {
			return nil, err
		}
		o.options = append(o.options, grpc.WithTransportCredentials(cred))
	}
	if o.requestValidationInterceptor {
		o.options = append(o.options, grpc.WithUnaryInterceptor(RequestValidationUnaryClientInterceptor()))
	}
	return o.options, nil
}

func DialContext(ctx context.Context, addr string, opts ...DialOption) (*grpc.ClientConn, error) {
	options, err := DialOptions(opts...)
	if err != nil {
		return nil, err
	}
	return grpc.DialContext(ctx, addr, options...)
}
