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

package main

import (
	"embed"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

//go:embed testdata/*
var testdata embed.FS

func readTestdataFile(t *testing.T, name string) []byte {
	data, err := testdata.ReadFile(name)
	require.NoError(t, err)
	return data
}

func boolPointer(b bool) *bool {
	return &b
}

func Test_parseArgs(t *testing.T) {
	type args struct {
		args []string
	}
	tests := []struct {
		name    string
		args    args
		want    arguments
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "minimum required args with no error",
			args: args{
				args: []string{
					"address=localhost:8080",
					"api-key=xxxxxxxxxxxxxx",
					"token=xxxxxxxxxxxxxxxx",
				},
			},
			want: arguments{
				Address:            "localhost:8080",
				APIKey:             "xxxxxxxxxxxxxx",
				Token:              "xxxxxxxxxxxxxxxx",
				Timeout:            defaultTimeout,
				PipedHandleTimeout: defaultTimeout,
			},
			wantErr: assert.NoError,
		},
		{
			name: "minimum required args and specified timeout arg with no error",
			args: args{
				args: []string{
					"address=localhost:8080",
					"api-key=xxxxxxxxxxxxxx",
					"token=xxxxxxxxxxxxxxxx",
					"timeout=10m",
				},
			},
			want: arguments{
				Address:            "localhost:8080",
				APIKey:             "xxxxxxxxxxxxxx",
				Token:              "xxxxxxxxxxxxxxxx",
				Timeout:            10 * time.Minute,
				PipedHandleTimeout: defaultTimeout,
			},
			wantErr: assert.NoError,
		},
		{
			name: "minimum required args and specified piped-handle-timeout arg with no error",
			args: args{
				args: []string{
					"address=localhost:8080",
					"api-key=xxxxxxxxxxxxxx",
					"token=xxxxxxxxxxxxxxxx",
					"piped-handle-timeout=10m",
				},
			},
			want: arguments{
				Address:            "localhost:8080",
				APIKey:             "xxxxxxxxxxxxxx",
				Token:              "xxxxxxxxxxxxxxxx",
				Timeout:            defaultTimeout,
				PipedHandleTimeout: 10 * time.Minute,
			},
			wantErr: assert.NoError,
		},
		{
			name: "minimum required args and specified timeout and piped-handle-timeout arg with no error",
			args: args{
				args: []string{
					"address=localhost:8080",
					"api-key=xxxxxxxxxxxxxx",
					"token=xxxxxxxxxxxxxxxx",
					"timeout=12m",
					"piped-handle-timeout=15m",
				},
			},
			want: arguments{
				Address:            "localhost:8080",
				APIKey:             "xxxxxxxxxxxxxx",
				Token:              "xxxxxxxxxxxxxxxx",
				Timeout:            12 * time.Minute,
				PipedHandleTimeout: 15 * time.Minute,
			},
			wantErr: assert.NoError,
		},
		{
			name: "missing required args (address) returns error",
			args: args{
				args: []string{
					"api-key=xxxxxxxxxxxxxx",
					"token=xxxxxxxxxxxxxxxx",
				},
			},
			want:    arguments{},
			wantErr: assert.Error,
		},
		{
			name: "missing required args (api-key) returns error",
			args: args{
				args: []string{
					"address=localhost:8080",
					"token=xxxxxxxxxxxxxxxx",
				},
			},
			want:    arguments{},
			wantErr: assert.Error,
		},
		{
			name: "missing required args (token) returns error",
			args: args{
				args: []string{
					"address=localhost:8080",
					"api-key=xxxxxxxxxxxxxx",
				},
			},
			want:    arguments{},
			wantErr: assert.Error,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseArgs(tt.args.args)
			if tt.wantErr(t, err, fmt.Sprintf("parseArgs(%v)", tt.args.args)) {
				return
			}
			assert.Equalf(t, tt.want, got, "parseArgs(%v)", tt.args.args)
		})
	}
}
