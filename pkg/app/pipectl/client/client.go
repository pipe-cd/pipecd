// Copyright 2024 The PipeCD Authors.
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

package client

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"net"
	"net/url"
	"strings"
	"time"
	"unicode"

	"github.com/spf13/cobra"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/resolver"

	"github.com/pipe-cd/pipecd/pkg/app/server/service/apiservice"
	"github.com/pipe-cd/pipecd/pkg/model"
	"github.com/pipe-cd/pipecd/pkg/rpc/rpcauth"
	"github.com/pipe-cd/pipecd/pkg/rpc/rpcclient"
)

type Options struct {
	Address    string
	APIKey     string
	APIKeyFile string
	Insecure   bool
	CertFile   string
}

func (o *Options) RegisterPersistentFlags(cmd *cobra.Command) {
	cmd.PersistentFlags().StringVar(&o.Address, "address", o.Address, "The address to control-plane api.")
	cmd.PersistentFlags().StringVar(&o.APIKey, "api-key", o.APIKey, "The API key used while authenticating with control-plane.")
	cmd.PersistentFlags().StringVar(&o.APIKeyFile, "api-key-file", o.APIKeyFile, "Path to the file containing API key used while authenticating with control-plane.")
	cmd.PersistentFlags().BoolVar(&o.Insecure, "insecure", o.Insecure, "Whether disabling transport security while connecting to control-plane.")
	cmd.PersistentFlags().StringVar(&o.CertFile, "cert-file", o.CertFile, "The path to the TLS certificate file.")
}

// validateAddress checks that addr is a target that gRPC is able to dial.
//
// It mirrors the way grpc-go itself parses a target: an address whose scheme is
// backed by a registered resolver is a resolver URI, and its endpoint syntax is
// resolver-specific (dns:///host defaults to port 443, unix:///path has no port
// at all), so it is left to that resolver. Any other address is handed to the
// default passthrough resolver and therefore has to be a dialable host:port.
func validateAddress(addr string) error {
	// Guarded here as well as in Validate so that the function is correct on its
	// own, independently of the order its callers happen to run their checks in.
	if addr == "" {
		return errors.New("address must be set")
	}

	if u, err := url.Parse(addr); err == nil && u.Scheme != "" {
		// Note that a scheme always wins over a host of the same name, so
		// "dns:9000" is the dns resolver with endpoint "9000" rather than port
		// 9000 of a host named "dns". That is grpc-go's own reading of it.
		if resolver.Get(u.Scheme) != nil {
			return nil
		}
		// An address written as a URL whose scheme is backed by no resolver,
		// such as http:// or tcp://, is never a valid target. Report the scheme
		// rather than letting it fall through and be blamed on the port.
		if strings.Contains(addr, "://") {
			return fmt.Errorf("invalid address %q: scheme %q is not a supported gRPC resolver", addr, u.Scheme)
		}
	}

	host, port, err := net.SplitHostPort(addr)
	if err != nil {
		return fmt.Errorf("invalid address %q: must be in host:port form (e.g. localhost:9000) or use a gRPC scheme such as dns:///", addr)
	}
	if host == "" {
		return fmt.Errorf("invalid address %q: host must not be empty", addr)
	}
	// Only the host of a host:port address is checked, so that socket paths
	// behind a resolver scheme stay free to contain whitespace.
	if strings.IndexFunc(host, unicode.IsSpace) >= 0 {
		return fmt.Errorf("invalid address %q: host must not contain whitespace", addr)
	}
	// LookupPort rather than strconv.Atoi because net.Dial accepts named ports
	// such as "localhost:http"; it reads the local services database, and does
	// no network I/O.
	//
	// SplitHostPort accepts an empty port ("localhost:") and LookupPort maps it
	// to 0, so both need rejecting here: a client can never dial port 0.
	p, err := net.LookupPort("tcp", port)
	if err != nil {
		return fmt.Errorf("invalid address %q: %q is not a valid port", addr, port)
	}
	if p == 0 {
		return fmt.Errorf("invalid address %q: port must be between 1 and 65535", addr)
	}
	return nil
}

func (o *Options) Validate() error {
	if o.Address == "" {
		return errors.New("address must be set")
	}
	if err := validateAddress(o.Address); err != nil {
		return err
	}
	if o.APIKey == "" && o.APIKeyFile == "" {
		return errors.New("either api-key or api-key-file must be set")
	}
	return nil
}

func (o *Options) NewClient(ctx context.Context) (apiservice.Client, error) {
	if err := o.Validate(); err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	var creds credentials.PerRPCCredentials
	var err error

	if o.APIKey != "" {
		creds = rpcclient.NewPerRPCCredentials(o.APIKey, rpcauth.APIKeyCredentials, !o.Insecure)
	} else {
		creds, err = rpcclient.NewPerRPCCredentialsFromFile(o.APIKeyFile, rpcauth.APIKeyCredentials, !o.Insecure)
		if err != nil {
			return nil, err
		}
	}

	options := []rpcclient.DialOption{
		rpcclient.WithBlock(),
		rpcclient.WithPerRPCCredentials(creds),
	}

	if !o.Insecure {
		if o.CertFile != "" {
			options = append(options, rpcclient.WithTLS(o.CertFile))
		} else {
			config := &tls.Config{}
			options = append(options, rpcclient.WithTransportCredentials(credentials.NewTLS(config)))
		}
	} else {
		options = append(options, rpcclient.WithInsecure())
	}

	client, err := apiservice.NewClient(ctx, o.Address, options...)
	if err != nil {
		return nil, err
	}

	return client, nil
}

func getCommand(ctx context.Context, cli apiservice.Client, cmdID string) (*model.Command, error) {
	req := &apiservice.GetCommandRequest{
		CommandId: cmdID,
	}
	resp, err := cli.GetCommand(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp.Command, nil
}
