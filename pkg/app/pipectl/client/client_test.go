// Copyright 2026 The PipeCD Authors.
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
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestValidateAddress(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		address string
		wantErr string // empty means the address is expected to be valid
	}{
		// host:port
		{
			name:    "host and port",
			address: "localhost:9000",
		},
		{
			name:    "IPv4 and port",
			address: "127.0.0.1:443",
		},
		{
			name:    "IPv6 and port",
			address: "[::1]:9000",
		},
		{
			name:    "fully qualified domain and port",
			address: "api.pipecd.dev:443",
		},
		{
			name:    "named port",
			address: "localhost:http",
		},

		// Resolver URIs. Their endpoint syntax is resolver-specific, so an
		// address without a port is still valid here.
		{
			name:    "dns resolver with port",
			address: "dns:///localhost:9000",
		},
		{
			name:    "dns resolver without port",
			address: "dns:///localhost",
		},
		{
			name:    "dns resolver with authority",
			address: "dns://8.8.8.8/localhost:9000",
		},
		{
			name:    "unix resolver",
			address: "unix:///var/run/pipecd.sock",
		},
		{
			name:    "unix resolver with relative path",
			address: "unix:pipecd.sock",
		},
		{
			name:    "passthrough resolver",
			address: "passthrough:///localhost:9000",
		},

		// Missing port: the default passthrough resolver hands the address
		// straight to net.Dial, which requires one.
		{
			name:    "host without port",
			address: "localhost",
			wantErr: "must be in host:port form",
		},

		// A scheme that no registered resolver backs is not a gRPC target,
		// even though it parses as a URL.
		{
			name:    "http URL",
			address: "http://localhost:9000",
			wantErr: `scheme "http" is not a supported gRPC resolver`,
		},
		{
			name:    "https URL",
			address: "https://api.pipecd.dev",
			wantErr: `scheme "https" is not a supported gRPC resolver`,
		},
		{
			name:    "tcp scheme",
			address: "tcp://localhost:9000",
			wantErr: `scheme "tcp" is not a supported gRPC resolver`,
		},

		// Malformed input.
		{
			name:    "no scheme before separator",
			address: "://bad",
			wantErr: "host must not be empty",
		},
		{
			name:    "empty host",
			address: ":9000",
			wantErr: "host must not be empty",
		},
		{
			name:    "trailing path after port",
			address: "localhost:9000/path",
			wantErr: "is not a valid port",
		},
		{
			name:    "non numeric port",
			address: "localhost:abc",
			wantErr: "is not a valid port",
		},
		{
			name:    "port out of range",
			address: "localhost:99999",
			wantErr: "is not a valid port",
		},
		{
			name:    "negative port",
			address: "localhost:-1",
			wantErr: "is not a valid port",
		},
		{
			name:    "too many colons",
			address: "localhost:9000:extra",
			wantErr: "must be in host:port form",
		},
		{
			name:    "IPv6 without brackets",
			address: "::1:9000",
			wantErr: "must be in host:port form",
		},
		{
			name:    "bracketed IPv6 without port",
			address: "[::1]",
			wantErr: "must be in host:port form",
		},
		// SplitHostPort accepts an empty port and LookupPort maps it to 0, so
		// these would otherwise slip through and be dialed against port 0.
		{
			name:    "trailing colon without port",
			address: "localhost:",
			wantErr: "port must be between 1 and 65535",
		},
		{
			name:    "port zero",
			address: "localhost:0",
			wantErr: "port must be between 1 and 65535",
		},
		// A copy-paste stray space makes the host unresolvable.
		{
			name:    "leading whitespace",
			address: " localhost:9000",
			wantErr: "host must not contain whitespace",
		},
		{
			name:    "whitespace inside host",
			address: "local host:9000",
			wantErr: "host must not contain whitespace",
		},
		{
			name:    "carriage return from a pasted value",
			address: "localhost\r:9000",
			wantErr: "host must not contain whitespace",
		},
		{
			name:    "empty address",
			address: "",
			wantErr: "address must be set",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			err := validateAddress(tt.address)
			if tt.wantErr == "" {
				assert.NoError(t, err)
				return
			}
			require.Error(t, err)
			assert.ErrorContains(t, err, tt.wantErr)
			// The offending value is echoed back so the user can see what was
			// parsed. It is quoted with %q, which escapes control characters.
			if tt.address != "" {
				assert.ErrorContains(t, err, fmt.Sprintf("%q", tt.address))
			}
		})
	}
}

func TestOptionsValidate(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		options Options
		wantErr string // empty means the options are expected to be valid
	}{
		{
			name:    "valid with api key",
			options: Options{Address: "localhost:9000", APIKey: "key"},
		},
		{
			name:    "valid with api key file",
			options: Options{Address: "localhost:9000", APIKeyFile: "/path/to/key"},
		},
		{
			name:    "valid with resolver URI",
			options: Options{Address: "dns:///localhost:9000", APIKey: "key"},
		},
		{
			name:    "missing address",
			options: Options{APIKey: "key"},
			wantErr: "address must be set",
		},
		{
			name:    "malformed address",
			options: Options{Address: "localhost", APIKey: "key"},
			wantErr: "must be in host:port form",
		},
		{
			name:    "address is validated before credentials",
			options: Options{Address: "localhost"},
			wantErr: "must be in host:port form",
		},
		{
			name:    "missing credentials",
			options: Options{Address: "localhost:9000"},
			wantErr: "either api-key or api-key-file must be set",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			err := tt.options.Validate()
			if tt.wantErr == "" {
				assert.NoError(t, err)
				return
			}
			require.Error(t, err)
			assert.ErrorContains(t, err, tt.wantErr)
		})
	}
}

// NewClient must reject a malformed address up front rather than spending the
// dial timeout on it and reporting a transport error, which is the whole point
// of validating the flag.
func TestNewClientRejectsMalformedAddressWithoutDialing(t *testing.T) {
	t.Parallel()

	opts := &Options{Address: "localhost", APIKey: "key"}

	start := time.Now()
	client, err := opts.NewClient(context.Background())
	elapsed := time.Since(start)

	require.Error(t, err)
	assert.Nil(t, client)
	assert.ErrorContains(t, err, "must be in host:port form")
	// NewClient dials with a 5s timeout, so returning promptly proves the
	// address was rejected before any connection was attempted.
	assert.Less(t, elapsed, time.Second)
}
