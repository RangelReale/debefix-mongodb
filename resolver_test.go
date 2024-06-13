package debefix_mongodb

import (
	"context"
	"testing"
	"testing/fstest"

	"github.com/rrgmc/debefix"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/integration/mtest"
	"gotest.tools/v3/assert"
)

func TestResolve(t *testing.T) {
	provider := debefix.NewFSFileProvider(fstest.MapFS{
		"users.dbf.yaml": &fstest.MapFile{
			Data: []byte(`tables:
  tags:
    config:
      table_name: "public.tags"
    rows:
      - _id: !expr generated
        _refid: !refid "all"
        tag_name: "All"
      - _id: !expr generated
        _refid: !refid "half"
        tag_name: "Half"
  posts:
    config:
      table_name: "public.posts"
      depends: ["tags"]
    rows:
      - _id: !expr generated
        _refid: !refid "post_1"
        title: "First post"
      - _id: !expr generated
        _refid: !refid "post_2"
        title: "Second post"
  post_tags:
    config:
      table_name: "public.post_tags"
    rows:
      - post_id: !expr "refid:posts:post_1:_id"
        tag_id: !expr "refid:tags:all:_id"
      - post_id: !expr "refid:posts:post_2:_id"
        tag_id: !expr "refid:tags:half:_id"
`),
		},
	})

	data, err := debefix.Load(provider)
	assert.NilError(t, err)

	mt := mtest.New(t, mtest.NewOptions().DatabaseName("debefix-test").ClientType(mtest.Mock))
	defer mt.Close()

	mt.Run("test", func(mt *mtest.T) {
		var expectedResponses []bson.D
		for i := 0; i < 6; i++ {
			expectedResponses = append(expectedResponses, mtest.CreateSuccessResponse())
		}

		mt.AddMockResponses(expectedResponses...)

		_, err = Resolve(context.Background(), mt.DB, data)
		assert.NilError(mt, err)
	})
}
