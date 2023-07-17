// Copyright (c) 2023 Tiago Melo. All rights reserved.
// Use of this source code is governed by the MIT License that can be found in
// the LICENSE file.
//
// Package product provides the business logic and data operations for the product catalog.
// It includes functions for creating, getting, updating, deleting, and listing products.
package product

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/tiagomelo/golang-grpc-mongodb-arbitrary-data/api/proto/gen/productcatalog"
	"github.com/tiagomelo/golang-grpc-mongodb-arbitrary-data/store"
	"github.com/tiagomelo/golang-grpc-mongodb-arbitrary-data/store/product/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

const collectionName = "products"

type Cursor interface {
	Decode(interface{}) error
	Err() error
	Close(context.Context) error
	Next(context.Context) bool
}

type cursorWrapper struct {
	*mongo.Cursor
}

// For ease of unit testing.
var (
	uuidProvider         = uuid.NewString
	insertIntoCollection = func(ctx context.Context, collection *mongo.Collection, document interface{}) (*mongo.InsertOneResult, error) {
		return collection.InsertOne(ctx, document)
	}
	find = func(ctx context.Context, collection *mongo.Collection, filter interface{}) (Cursor, error) {
		cur, err := collection.Find(ctx, filter)
		return &cursorWrapper{cur}, err
	}
	findOne = func(ctx context.Context, collection *mongo.Collection, filter interface{}, p *models.Product) error {
		sr := collection.FindOne(ctx, filter)
		return sr.Decode(p)
	}
	updateOne = func(ctx context.Context, collection *mongo.Collection, filter interface{}, update interface{}) (*mongo.UpdateResult, error) {
		return collection.UpdateOne(ctx, filter, update)
	}
	deleteOne = func(ctx context.Context, collection *mongo.Collection, filter interface{}) (*mongo.DeleteResult, error) {
		return collection.DeleteOne(ctx, filter)
	}
)

// Get retrieves a product from the database by uuid.
func Get(ctx context.Context, db *store.MongoDb, req *productcatalog.GetProductRequest) (*models.Product, error) {
	coll := db.Client.Database(db.DatabaseName).Collection(collectionName)
	var product models.Product
	err := findOne(ctx, coll, bson.M{"uuid": req.GetUuid()}, &product)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf(`product with uuid "%s" does not exist`, req.GetUuid())
		}
		return nil, errors.Wrapf(err, `getting product with uuid "%s"`, req.GetUuid())
	}
	return &product, nil
}

// Create creates a new product in the database.
func Create(ctx context.Context, db *store.MongoDb, newProduct *models.Product) (*models.Product, error) {
	coll := db.Client.Database(db.DatabaseName).Collection(collectionName)
	newProduct.Uuid = uuidProvider()
	_, err := insertIntoCollection(ctx, coll, newProduct)
	if err != nil {
		return nil, errors.Wrap(err, "inserting product")
	}
	return newProduct, nil
}

// Update updates a product in the database.
func Update(ctx context.Context, db *store.MongoDb, productToUpdate *models.Product) (*models.Product, error) {
	coll := db.Client.Database(db.DatabaseName).Collection(collectionName)
	_, err := updateOne(ctx, coll, bson.M{"uuid": productToUpdate.Uuid}, bson.M{"$set": productToUpdate})
	if err != nil {
		return nil, errors.Wrapf(err, `updating product with uuid "%s"`, productToUpdate.Uuid)
	}
	return productToUpdate, nil
}

// Delete deletes a product from the database by uuid.
func Delete(ctx context.Context, db *store.MongoDb, req *productcatalog.DeleteProductRequest) (*productcatalog.DeleteProductResponse, error) {
	coll := db.Client.Database(db.DatabaseName).Collection(collectionName)
	_, err := deleteOne(ctx, coll, bson.M{"uuid": req.Uuid})
	if err != nil {
		return nil, errors.Wrapf(err, `deleting product with uuid "%s"`, req.Uuid)
	}
	return &productcatalog.DeleteProductResponse{Result: "success"}, nil
}

// List lists all products in the database.
func List(ctx context.Context, db *store.MongoDb, req *productcatalog.ListProductsRequest) ([]*models.Product, error) {
	coll := db.Client.Database(db.DatabaseName).Collection(collectionName)
	cur, err := find(ctx, coll, bson.M{})
	if err != nil {
		return nil, errors.Wrap(err, "finding products")
	}
	defer cur.Close(ctx)
	var products []*models.Product
	for cur.Next(ctx) {
		var product models.Product
		if err = cur.Decode(&product); err != nil {
			return nil, errors.Wrap(err, "decoding product")
		}
		products = append(products, &product)
	}
	if err := cur.Err(); err != nil {
		return nil, errors.Wrap(err, "cursor error")
	}
	return products, nil
}
