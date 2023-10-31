package debefix_mongodb

import (
	"context"
	"slices"

	"github.com/rrgmc/debefix"
	"github.com/rrgmc/debefix/db"
	"go.mongodb.org/mongo-driver/mongo"
)

// Resolve inserts data on the MongoDB database, and returns the resolved data.
func Resolve(ctx context.Context, db *mongo.Database, data *debefix.Data, options ...debefix.ResolveOption) (*debefix.Data, error) {
	return debefix.Resolve(data, ResolverFunc(ctx, db), options...)
}

// ResolverFunc is the debefix.ResolveCallback used by Resolve.
func ResolverFunc(ctx context.Context, mdb *mongo.Database) debefix.ResolveCallback {
	return db.ResolverFunc(ResolverDBCallback(ctx, mdb))
}

// ResolverDBCallback is a db.ResolverDBCallback to generate MongoDB collection records.
func ResolverDBCallback(ctx context.Context, db *mongo.Database) db.ResolverDBCallback {
	return func(tableName string, fields map[string]any, returnFieldNames []string) (map[string]any, error) {
		collection := db.Collection(tableName)

		res, err := collection.InsertOne(ctx, fields)
		if err != nil {
			return nil, err
		}

		if slices.Contains(returnFieldNames, "_id") {
			return map[string]any{"_id": res.InsertedID}, nil
		}

		return nil, nil
	}
}
