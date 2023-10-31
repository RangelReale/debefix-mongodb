package debefix_mongodb

import (
	"context"
	"io/fs"

	"github.com/rrgmc/debefix"
	"go.mongodb.org/mongo-driver/mongo"
)

// Generate loads files and inserts records in the database, returning the resolved data.
func Generate(ctx context.Context, fileProvider debefix.FileProvider, db *mongo.Database,
	options ...debefix.GenerateOption) (*debefix.Data, error) {
	return debefix.Generate(fileProvider, ResolverFunc(ctx, db), options...)
}

// GenerateFS is a version of Generate that loads from a fs.FS, returning the resolved data.
func GenerateFS(ctx context.Context, fs fs.FS, db *mongo.Database,
	options ...debefix.GenerateOption) (*debefix.Data, error) {
	return debefix.GenerateFS(fs, ResolverFunc(ctx, db), options...)
}

// GenerateDirectory is a version of Generate that loads from a directory name, returning the resolved data.
func GenerateDirectory(ctx context.Context, rootDir string, db *mongo.Database,
	options ...debefix.GenerateOption) (*debefix.Data, error) {
	return debefix.GenerateDirectory(rootDir, ResolverFunc(ctx, db), options...)
}
