package insightstore

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/pipe-cd/pipe/pkg/filestore/filestoretest"
	"github.com/stretchr/testify/assert"

	"github.com/pipe-cd/pipe/pkg/filestore"
	"github.com/pipe-cd/pipe/pkg/insight"
)

func TestStore_LoadApplicationsCount(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	store := filestoretest.NewMockStore(ctrl)

	fs := Store{
		filestore: store,
	}

	tests := []struct {
		name        string
		projectID   string
		content     string
		readerErr   error
		want        *insight.ApplicationsCount
		expectedErr error
	}{
		{
			name:        "file not found in filestore",
			projectID:   "pid1",
			content:     "",
			readerErr:   filestore.ErrNotFound,
			expectedErr: filestore.ErrNotFound,
		},
		{
			name:      "success",
			projectID: "pid1",
			content: `{
				"accumulated_to": 1609459200,
				"counts": {
					"CLOUDRUN": {
						"deploying": 1,
						"deleted": 2,
						"disabled": 3
					}
				}
			}`,
			want: &insight.ApplicationsCount{
				AccumulatedTo: 1609459200,
				Counts: map[string]insight.Count{
					"CLOUDRUN": {
						Deploying: 1,
						Deleted:   2,
						Disabled:  3,
					},
				},
			},
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			path := determineFilePath(tc.projectID)
			obj := filestore.Object{
				Content: []byte(tc.content),
			}
			store.EXPECT().GetObject(context.TODO(), path).Return(obj, tc.readerErr)
			ac, err := fs.LoadApplicationsCount(context.TODO(), tc.projectID)
			if err != nil {
				if tc.expectedErr == nil {
					assert.NoError(t, err)
					return
				}
				assert.Error(t, err, tc.expectedErr)
				return
			}
			assert.Equal(t, tc.want, ac)
		})
	}
}
