package debefix_mongodb

import (
	"context"
	"testing"

	"github.com/rrgmc/debefix/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/integration/mtest"
	"gotest.tools/v3/assert"
)

var (
	tableTags     = debefix.TableName("public.tags")
	tablePosts    = debefix.TableName("public.posts")
	tablePostTags = debefix.TableName("public.post_tags")
)

func TestResolve(t *testing.T) {
	data := debefix.NewData()

	data.AddValues(tableTags,
		debefix.MapValues{
			"tag_id":   2,
			"_refid":   debefix.SetValueRefID("all"),
			"tag_name": "All",
		},
		debefix.MapValues{
			"tag_id":   5,
			"_refid":   debefix.SetValueRefID("half"),
			"tag_name": "Half",
		},
	)

	data.AddValues(tablePosts,
		debefix.MapValues{
			"post_id": 1,
			"_refid":  debefix.SetValueRefID("post_1"),
			"title":   "First post",
		},
		debefix.MapValues{
			"post_id": 2,
			"_refid":  debefix.SetValueRefID("post_2"),
			"title":   "Second post",
		},
	)

	data.AddDependencies(tablePosts, tableTags)

	data.AddValues(debefix.TableName(tablePostTags),
		debefix.MapValues{
			"post_id": debefix.ValueRefID(tablePosts, "post_1", "post_id"),
			"tag_id":  debefix.ValueRefID(tableTags, "all", "tag_id"),
		},
		debefix.MapValues{
			"post_id": debefix.ValueRefID(tablePosts, "post_2", "post_id"),
			"tag_id":  debefix.ValueRefID(tableTags, "half", "tag_id"),
		},
	)

	mt := mtest.New(t, mtest.NewOptions().DatabaseName("debefix-test").ClientType(mtest.Mock))
	defer mt.Close()

	mt.Run("test", func(mt *mtest.T) {
		var expectedResponses []bson.D
		for i := 0; i < 6; i++ {
			expectedResponses = append(expectedResponses, mtest.CreateSuccessResponse())
		}

		mt.AddMockResponses(expectedResponses...)

		_, err := Resolve(context.Background(), mt.DB, data)
		assert.NilError(mt, err)
	})
}
