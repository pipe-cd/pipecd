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

package apikeyverifier

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"

	"github.com/pipe-cd/pipecd/pkg/model"
)

type fakeAPIKeyGetter struct {
	calls   int
	apiKeys map[string]*model.APIKey
}

func (g *fakeAPIKeyGetter) Get(_ context.Context, id string) (*model.APIKey, error) {
	g.calls++
	p, ok := g.apiKeys[id]
	if ok {
		msg := proto.Clone(p)
		return msg.(*model.APIKey), nil
	}
	return nil, fmt.Errorf("not found")
}

type fakeRedisHashCache struct{}

func (f *fakeRedisHashCache) Put(k string, v interface{}) error {
	return nil
}

func TestVerify(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var id1 = "test-api-key"
	key1, hash1, err := model.GenerateAPIKey(id1)
	require.NoError(t, err)

	var id2 = "disabled-api-key"
	key2, hash2, err := model.GenerateAPIKey(id2)
	require.NoError(t, err)

	apiKeyGetter := &fakeAPIKeyGetter{
		apiKeys: map[string]*model.APIKey{
			id1: {
				Id:        id1,
				Name:      id1,
				KeyHash:   hash1,
				ProjectId: "test-project",
			},
			id2: {
				Id:        id2,
				Name:      id2,
				KeyHash:   hash2,
				ProjectId: "test-project",
				Disabled:  true,
			},
		},
	}
	fakeRedisHashCache := &fakeRedisHashCache{}
	v := NewVerifier(ctx, apiKeyGetter, fakeRedisHashCache, zap.NewNop())

	// Not found key.
	notFoundKey, _, err := model.GenerateAPIKey("not-found-api-key")
	require.NoError(t, err)

	apiKey, err := v.Verify(ctx, notFoundKey)
	require.Nil(t, apiKey)
	require.NotNil(t, err)
	assert.Equal(t, "unable to find API key not-found-api-key from datastore, not found", err.Error())
	require.Equal(t, 1, apiKeyGetter.calls)

	// Found key but it was disabled.
	apiKey, err = v.Verify(ctx, key2)
	require.Nil(t, apiKey)
	require.NotNil(t, err)
	assert.Equal(t, "the api key disabled-api-key was already disabled", err.Error())
	require.Equal(t, 2, apiKeyGetter.calls)

	// Found key but invalid secret.
	apiKey, err = v.Verify(ctx, fmt.Sprintf("%s.invalidhash", id1))
	require.Nil(t, apiKey)
	require.NotNil(t, err)
	assert.Equal(t, "invalid api key test-api-key: wrong api key test-api-key.invalidhash: crypto/bcrypt: hashedPassword is not the hash of the given password", err.Error())
	require.Equal(t, 3, apiKeyGetter.calls)

	// OK.
	apiKey, err = v.Verify(ctx, key1)
	assert.Equal(t, id1, apiKey.Name)
	assert.Nil(t, err)
	require.Equal(t, 3, apiKeyGetter.calls)
}
