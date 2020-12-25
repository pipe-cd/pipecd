package datastore

import (
	"context"
	"fmt"
	"testing"

	"github.com/pipe-cd/pipe/pkg/insight"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestGetInsightMilestone(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	testcases := []struct {
		name    string
		id      string
		ds      DataStore
		wantErr bool
	}{
		{
			name: "successful fetch from datastore",
			id:   "id",
			ds: func() DataStore {
				ds := NewMockDataStore(ctrl)
				ds.EXPECT().
					Get(gomock.Any(), insightModelKind, insightMilestone, &insight.Milestone{}).
					Return(nil)
				return ds
			}(),
			wantErr: false,
		},
		{
			name: "failed fetch from datastore",
			id:   "id",
			ds: func() DataStore {
				ds := NewMockDataStore(ctrl)
				ds.EXPECT().
					Get(gomock.Any(), insightModelKind, insightMilestone, &insight.Milestone{}).
					Return(fmt.Errorf("err"))
				return ds
			}(),
			wantErr: true,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			s := NewInsightMilestoneStore(tc.ds)
			_, err := s.GetInsightMilestone(context.Background())
			assert.Equal(t, tc.wantErr, err != nil)
		})
	}
}
