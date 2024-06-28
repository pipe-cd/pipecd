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

package pipedverifier

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/protobuf/proto"

	"github.com/pipe-cd/pipecd/pkg/config"
	"github.com/pipe-cd/pipecd/pkg/model"
)

type fakeProjectGetter struct {
	calls    int
	projects map[string]*model.Project
}

func (g *fakeProjectGetter) Get(_ context.Context, id string) (*model.Project, error) {
	g.calls++
	p, ok := g.projects[id]
	if ok {
		msg := proto.Clone(p)
		return msg.(*model.Project), nil
	}
	return nil, fmt.Errorf("not found")
}

type fakePipedGetter struct {
	calls  int
	pipeds map[string]*model.Piped
}

func (g *fakePipedGetter) Get(_ context.Context, id string) (*model.Piped, error) {
	g.calls++
	p, ok := g.pipeds[id]
	if ok {
		msg := proto.Clone(p)
		return msg.(*model.Piped), nil
	}
	return nil, fmt.Errorf("not found")
}

func TestVerify(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	hashGenerator := func(k string) string {
		h, err := bcrypt.GenerateFromPassword([]byte(k), bcrypt.DefaultCost)
		require.NoError(t, err)
		return string(h)
	}
	projectGetter := &fakeProjectGetter{
		projects: map[string]*model.Project{
			"project-1": {
				Id: "project-1",
			},
		},
	}
	pipedGetter := &fakePipedGetter{
		pipeds: map[string]*model.Piped{
			"piped-0-1": {
				Id:        "piped-0-1",
				ProjectId: "project-0",
				Keys: []*model.PipedKey{
					{
						Hash: hashGenerator("piped-key-0-1"),
					},
				},
			},
			"piped-1-1": {
				Id:        "piped-1-1",
				ProjectId: "project-1",
				Keys: []*model.PipedKey{
					{
						Hash: hashGenerator("piped-key-1-1"),
					},
				},
			},
			"piped-1-2": {
				Id:        "piped-1-2",
				ProjectId: "project-non-existence",
				Keys: []*model.PipedKey{
					{
						Hash: hashGenerator("piped-key-1-2"),
					},
				},
			},
			"piped-1-3": {
				Id:        "piped-1-3",
				ProjectId: "project-1",
				Keys: []*model.PipedKey{
					{
						Hash: hashGenerator("piped-key-1-3"),
					},
				},
				Disabled: true,
			},
		},
	}
	v := NewVerifier(
		ctx,
		&config.ControlPlaneSpec{
			Projects: []config.ControlPlaneProject{
				{
					ID: "project-0",
				},
			},
		},
		projectGetter,
		pipedGetter,
		zap.NewNop(),
	)

	// A piped from a project that was specified in the control-plane configuration.
	err := v.Verify(ctx, "project-0", "piped-0-1", "piped-key-0-1")
	assert.Equal(t, nil, err)
	require.Equal(t, 0, projectGetter.calls)
	require.Equal(t, 1, pipedGetter.calls)

	// Non-existence project.
	err = v.Verify(ctx, "project-not-found", "piped-1-1", "piped-key-1-1")
	assert.Equal(t, fmt.Errorf("project project-not-found for piped piped-1-1 was not found"), err)
	require.Equal(t, 1, projectGetter.calls)
	require.Equal(t, 1, pipedGetter.calls)

	// Found piped but project id was not correct.
	err = v.Verify(ctx, "project-1", "piped-1-2", "piped-key-1-2")
	assert.Equal(t, fmt.Errorf("the project of piped piped-1-2 is not matched, expected=project-1, got=project-non-existence"), err)
	require.Equal(t, 2, projectGetter.calls)
	require.Equal(t, 2, pipedGetter.calls)

	// Found piped but it was disabled.
	err = v.Verify(ctx, "project-1", "piped-1-3", "piped-key-1-3")
	assert.Equal(t, fmt.Errorf("piped piped-1-3 was already disabled"), err)
	require.Equal(t, 2, projectGetter.calls)
	require.Equal(t, 3, pipedGetter.calls)

	piped13 := pipedGetter.pipeds["piped-1-3"]
	piped13.Disabled = false
	err = v.Verify(ctx, "project-1", "piped-1-3", "piped-key-1-3")
	assert.Equal(t, nil, err)
	require.Equal(t, 2, projectGetter.calls)
	require.Equal(t, 4, pipedGetter.calls)

	// OK.
	err = v.Verify(ctx, "project-1", "piped-1-1", "piped-key-1-1")
	assert.Equal(t, nil, err)
	require.Equal(t, 2, projectGetter.calls)
	require.Equal(t, 5, pipedGetter.calls)

	// Wrong piped key.
	err = v.Verify(ctx, "project-1", "piped-1-1", "piped-key-1-1-wrong")
	assert.Equal(t, fmt.Errorf("the key of piped piped-1-1 is not matched, crypto/bcrypt: hashedPassword is not the hash of the given password"), err)
	require.Equal(t, 2, projectGetter.calls)
	require.Equal(t, 6, pipedGetter.calls)

	// The invalid key should be cached.
	err = v.Verify(ctx, "project-1", "piped-1-1", "piped-key-1-1-wrong")
	assert.Equal(t, fmt.Errorf("the key of piped piped-1-1 is not matched, crypto/bcrypt: hashedPassword is not the hash of the given password"), err)
	require.Equal(t, 2, projectGetter.calls)
	require.Equal(t, 6, pipedGetter.calls)

	// OK for newly added key.
	piped11 := pipedGetter.pipeds["piped-1-1"]
	piped11.Keys = append(piped11.Keys, &model.PipedKey{
		Hash: hashGenerator("piped-key-1-1-new"),
	})
	err = v.Verify(ctx, "project-1", "piped-1-1", "piped-key-1-1-new")
	assert.Equal(t, nil, err)
	require.Equal(t, 2, projectGetter.calls)
	require.Equal(t, 7, pipedGetter.calls)
}
