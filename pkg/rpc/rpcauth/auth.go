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

package rpcauth

import (
	"context"
	"fmt"
	"strings"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

// CredentialsType represents the type of credentials
// was set inside gRPC metadata.
type CredentialsType string

const (
	// IDTokenCredentials represents JWT IDToken for a web user.
	IDTokenCredentials CredentialsType = "ID-TOKEN"
	// PipedTokenCredentials represents a generated token for
	// authenticating between Piped and control-plane.
	PipedTokenCredentials CredentialsType = "PIPED-TOKEN"
	// APIKeyCredentials represents a generated key for
	// authenticating between pipectl/external-service and control-plane.
	APIKeyCredentials CredentialsType = "API-KEY"
	// UnknownCredentials represents an unsupported credentials.
	UnknownCredentials CredentialsType = "UNKNOWN"
)

// Credentials contains the type of credentials and credentials data.
type Credentials struct {
	Type CredentialsType
	Data string
}

// MakePipedToken builds a piped token can be used as data of Credentials.
func MakePipedToken(projectID, pipedID, pipedKey string) string {
	return fmt.Sprintf("%s,%s,%s", projectID, pipedID, pipedKey)
}

// ParsePipedToken parses a piped token and returns projectID, pipedID, pipedKey.
func ParsePipedToken(token string) (projectID, pipedID, pipedKey string, err error) {
	return parsePipedToken(token)
}

func parsePipedToken(token string) (projectID, pipedID, pipedKey string, err error) {
	parts := strings.Split(token, ",")
	if len(parts) != 3 {
		err = fmt.Errorf("malformed piped token")
		return
	}
	projectID = parts[0]
	if projectID == "" {
		err = fmt.Errorf("malformed piped token: projectID was empty")
		return
	}
	pipedID = parts[1]
	if pipedID == "" {
		err = fmt.Errorf("malformed piped token: pipedID was empty")
		return
	}
	pipedKey = parts[2]
	if pipedKey == "" {
		err = fmt.Errorf("malformed piped token: pipedKey was empty")
		return
	}
	return
}

func extractCredentials(ctx context.Context) (creds Credentials, err error) {
	creds.Type = UnknownCredentials
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		err = status.Error(codes.Unauthenticated, "missing credentials")
		return
	}

	rawCredentials := md["authorization"]
	if len(rawCredentials) == 0 {
		err = status.Error(codes.Unauthenticated, "missing credentials in authorization")
		return
	}

	subs := strings.Split(rawCredentials[0], " ")
	if len(subs) != 2 {
		err = status.Error(codes.Unauthenticated, "credentials is malformed")
		return
	}

	switch CredentialsType(subs[0]) {
	case IDTokenCredentials:
		creds.Data = subs[1]
		creds.Type = IDTokenCredentials

	case PipedTokenCredentials:
		creds.Data = subs[1]
		creds.Type = PipedTokenCredentials

	case APIKeyCredentials:
		creds.Data = subs[1]
		creds.Type = APIKeyCredentials

	default:
		err = status.Error(codes.Unauthenticated, "unsupported credentials type")
	}

	if creds.Data == "" {
		err = status.Error(codes.Unauthenticated, "credentials is malformed")
	}
	return
}

func extractCookie(ctx context.Context) (map[string]string, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "missing credentials")
	}
	rawCookie := md["cookie"]
	if len(rawCookie) == 0 {
		return nil, status.Error(codes.Unauthenticated, "missing cookie")
	}
	cs := strings.Split(rawCookie[0], ";")
	cookie := make(map[string]string, len(cs))
	for _, c := range cs {
		subs := strings.Split(strings.TrimSpace(c), "=")
		if len(subs) != 2 {
			return nil, status.Error(codes.Unauthenticated, "cookie is malformed")
		}
		cookie[subs[0]] = subs[1]
	}
	return cookie, nil
}
