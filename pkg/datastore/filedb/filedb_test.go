package filedb

import (
	"context"
	"encoding/json"
	"testing"
	"time"

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

func TestFileDBGet(t *testing.T) {
	ctx := context.Background()
	db, err := NewFileDB(newFileStore(ctx), []Option{}...)
	assert.Nil(t, err)
	assert.NotNil(t, db)
	assert.NotNil(t, db.backend)

	env := &model.Event{}
	err = db.Get(ctx, datastore.EventModelKind, "077babca-c175-4da4-81e8-d5196f8697e0", env)
	assert.Nil(t, err)
	assert.Equal(t, env.Id, "077babca-c175-4da4-81e8-d5196f8697e0")
	assert.Equal(t, env.Data, "v1.1")
}

func _TestFileDBPut(t *testing.T) {
	ctx := context.Background()
	db, err := NewFileDB(newFileStore(ctx), []Option{}...)
	assert.Nil(t, err)
	assert.NotNil(t, db)
	assert.NotNil(t, db.backend)

	data := `{
"data": "v1.2",
"event_key": "552614db6e72c994dcd0338d43bdba9b10a3a94af037fc219bb0b16c989a7d5e",
"id": "077babca-c175-4da4-81e8-d5196f8697e0",
"labels": {
  "app": "foo",
  "env": "dev"
},
"name": "image-update",
"project_id": "pipecd",
"created_at": 1621330102,
"updated_at": 1621330102
}`
	env := &model.Event{}
	err = json.Unmarshal([]byte(data), env)
	assert.Nil(t, err)
	err = db.Put(ctx, datastore.EventModelKind, "077babca-c175-4da4-81e8-d5196f8697e0", env)
	assert.Nil(t, err)
}

func _TestFileDBCreate(t *testing.T) {
	ctx := context.Background()
	db, err := NewFileDB(newFileStore(ctx), []Option{}...)
	assert.Nil(t, err)
	assert.NotNil(t, db)
	assert.NotNil(t, db.backend)

	data1 := `{
"data": "v1.1",
"event_key": "552614db6e72c994dcd0338d43bdba9b10a3a94af037fc219bb0b16c989a7d5e",
"id": "557537c8-7775-11ec-90d6-0242ac120003",
"labels": {
  "app": "foo",
  "env": "dev"
},
"name": "image-update",
"project_id": "pipecd",
"created_at": 1621330102,
"updated_at": 1621330102
}`
	data2 := `{
"data": "v1.2",
"event_key": "552614db6e72c994dcd0338d43bdba9b10a3a94af037fc219bb0b16c989a7d5e",
"id": "557537c8-7775-11ec-90d6-0242ac120003",
"labels": {
  "app": "foo",
  "env": "dev"
},
"name": "image-update",
"project_id": "pipecd",
"created_at": 1621330102,
"updated_at": 1621330102
}`
	go func(fd *FileDB) {
		env := &model.Event{}
		err = json.Unmarshal([]byte(data1), env)
		assert.Nil(t, err)
		err = db.Create(ctx, datastore.EventModelKind, "557537c8-7775-11ec-90d6-0242ac120003", env)
		assert.Nil(t, err)
	}(db)
	go func(fd *FileDB) {
		env := &model.Event{}
		err = json.Unmarshal([]byte(data2), env)
		assert.Nil(t, err)
		err = db.Create(ctx, datastore.EventModelKind, "557537c8-7775-11ec-90d6-0242ac120003", env)
		assert.Nil(t, err)
	}(db)
	time.Sleep(time.Second)
}
