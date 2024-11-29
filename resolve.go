package debefix_mongodb

import (
	"context"

	"github.com/rrgmc/debefix-db/v2"
	"github.com/rrgmc/debefix/v2"
	"go.mongodb.org/mongo-driver/mongo"
)

// Resolve inserts data on the MongoDB database, and returns the resolved data.
func Resolve(ctx context.Context, db *mongo.Database, data *debefix.Data, options ...debefix.ResolveOption) (*debefix.ResolvedData, error) {
	return debefix.Resolve(ctx, data, ResolveFunc(ctx, db), options...)
}

// ResolveFunc is the debefix.ResolveCallback used by Resolve.
func ResolveFunc(ctx context.Context, mdb *mongo.Database) debefix.ResolveCallback {
	return db.ResolveFunc(ResolveDBCallback(ctx, mdb))
}

// ResolveDBCallback is a db.ResolveDBCallback to generate MongoDB collection records.
func ResolveDBCallback(ctx context.Context, mdb *mongo.Database) db.ResolveDBCallback {
	return func(ctx context.Context, resolveInfo db.ResolveDBInfo, fields map[string]any, returnFields map[string]any) (returnValues map[string]any, err error) {
		collection := mdb.Collection(resolveInfo.TableID.TableName())
		res, err := collection.InsertOne(ctx, fields)
		if err != nil {
			return nil, err
		}
		if _, ok := returnFields["_id"]; ok {
			return map[string]any{"_id": res.InsertedID}, nil
		}
		return nil, nil
	}
}
