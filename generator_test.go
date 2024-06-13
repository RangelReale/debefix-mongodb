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

func TestGenerate(t *testing.T) {
	providerData := fstest.MapFS{
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
`),
		},
	}

	mt := mtest.New(t, mtest.NewOptions().DatabaseName("debefix-test").ClientType(mtest.Mock))
	defer mt.Close()

	mt.Run("test", func(mt *mtest.T) {
		var expectedResponses []bson.D
		for i := 0; i < 2; i++ {
			expectedResponses = append(expectedResponses, mtest.CreateSuccessResponse())
		}

		mt.AddMockResponses(expectedResponses...)

		_, err := Generate(context.Background(), debefix.NewFSFileProvider(providerData), mt.DB)
		assert.NilError(mt, err)
	})

	// same test using FS

	mt.Run("test", func(mt *mtest.T) {
		var expectedResponses []bson.D
		for i := 0; i < 2; i++ {
			expectedResponses = append(expectedResponses, mtest.CreateSuccessResponse())
		}

		mt.AddMockResponses(expectedResponses...)

		_, err := GenerateFS(context.Background(), providerData, mt.DB)
		assert.NilError(mt, err)
	})
}
