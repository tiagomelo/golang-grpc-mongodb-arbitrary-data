// Copyright (c) 2023 Tiago Melo. All rights reserved.
// Use of this source code is governed by the MIT License that can be found in
// the LICENSE file.
//
// Package mapper provides functions for converting between Protobuf messages
// and MongoDB models in the context of a product catalog.
// The functions in this package handle the conversion of product data between
// the Protobuf representation used in the API and the MongoDB model representation
// used in the data store.
package mapper

import (
	"github.com/pkg/errors"
	"github.com/tiagomelo/golang-grpc-mongodb-arbitrary-data/api/proto/gen/productcatalog"
	"github.com/tiagomelo/golang-grpc-mongodb-arbitrary-data/store/product/models"
	"google.golang.org/protobuf/types/known/structpb"
)

// For ease of unit testing.
var structpbNewValue = structpb.NewValue

// ProductProtobufToProductModel converts a Protobuf Product message to a MongoDB Product model.
func ProductProtobufToProductModel(product *productcatalog.Product) (*models.Product, error) {
	dbProduct := &models.Product{
		Uuid:        product.Uuid,
		Name:        product.Name,
		Description: product.Description,
		Price:       product.Price,
	}
	attributes := make(map[string]interface{})
	for k, p := range product.Attributes {
		attributes[k] = p.AsInterface()
	}
	dbProduct.Attributes = attributes
	return dbProduct, nil
}

// ProductModelToProductProtobuf converts a MongoDB Product model to a Protobuf Product message.
func ProductModelToProductProtobuf(dbProduct *models.Product) (*productcatalog.Product, error) {
	product := &productcatalog.Product{
		Uuid:        dbProduct.Uuid,
		Name:        dbProduct.Name,
		Description: dbProduct.Description,
		Price:       dbProduct.Price,
	}
	var err error
	attributes := make(map[string]*structpb.Value)
	for k, p := range dbProduct.Attributes {
		attributes[k], err = structpbNewValue(p)
		if err != nil {
			return nil, errors.Wrapf(err, `parsing attribute "%s"`, k)
		}
	}
	product.Attributes = attributes
	return product, nil
}

// ProductModelListToListProductsResponse converts a list of MongoDB Product models to a Protobuf ListProductsResponse message.
func ProductModelListToListProductsResponse(dbProducts []*models.Product) (*productcatalog.ListProductsResponse, error) {
	response := &productcatalog.ListProductsResponse{}
	products := []*productcatalog.Product{}
	for _, dbProduct := range dbProducts {
		product, err := ProductModelToProductProtobuf(dbProduct)
		if err != nil {
			return nil, err
		}
		products = append(products, product)
	}
	response.Products = products
	return response, nil
}
