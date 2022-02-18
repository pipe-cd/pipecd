package filedb

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/pipe-cd/pipecd/pkg/datastore"
	"github.com/pipe-cd/pipecd/pkg/filestore"
	"github.com/pipe-cd/pipecd/pkg/filestore/gcs"
	"github.com/pipe-cd/pipecd/pkg/model"
)

func newFileStore(ctx context.Context) filestore.Store {
	opts := []gcs.Option{
		gcs.WithCredentialsFile("/Users/s12228/.config/gcloud/application_default_credentials.json"),
	}
	s, _ := gcs.NewStore(ctx, "filestore-db-test", opts...)
	return s
}

type fakeEventCollection struct {
}

func (e *fakeEventCollection) Kind() string {
	return "Event"
}

func (e *fakeEventCollection) Factory() datastore.Factory {
	return func() interface{} {
		return &model.Event{}
	}
}

func (e *fakeEventCollection) ListInUsedShards() []datastore.Shard {
	return []datastore.Shard{
		datastore.AgentShard,
	}
}

func (e *fakeEventCollection) GetUpdatableShard() (datastore.Shard, error) {
	return datastore.AgentShard, nil
}

func TestFileDBGet(t *testing.T) {
	ctx := context.Background()
	db, err := NewFileDB(newFileStore(ctx), []Option{}...)
	assert.Nil(t, err)
	assert.NotNil(t, db)
	assert.NotNil(t, db.backend)

	// ctrl := gomock.NewController(t)
	// col := datastore.NewMockCollection(ctrl)
	// col.EXPECT().Factory().Return(func() interface{} {
	// 	return &model.Event{}
	// })
	// col.EXPECT().Kind().Return("Event")
	col := &fakeEventCollection{}

	env := &model.Event{}
	err = db.Get(ctx, col, "077babca-c175-4da4-81e8-d5196f8697e0", env)
	fmt.Printf("%v\n", env)
	assert.Nil(t, err)
	assert.Equal(t, env.Id, "077babca-c175-4da4-81e8-d5196f8697e0")
	assert.Equal(t, env.Data, "v1.1")
}
