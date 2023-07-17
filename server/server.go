// Copyright (c) 2023 Tiago Melo. All rights reserved.
// Use of this source code is governed by the MIT License that can be found in
// the LICENSE file.
//
// Package server implements the gRPC server for the product catalog service.
// It provides functions to handle CRUD operations for products.
//
// The server package is responsible for setting up the gRPC server,
// registering the product catalog service, and routing incoming gRPC
// requests to the corresponding functions in the product package.
package server

import (
	"context"

	"github.com/pkg/errors"
	"github.com/tiagomelo/golang-grpc-mongodb-arbitrary-data/api/proto/gen/productcatalog"
	"github.com/tiagomelo/golang-grpc-mongodb-arbitrary-data/mapper"
	"github.com/tiagomelo/golang-grpc-mongodb-arbitrary-data/store"
	"github.com/tiagomelo/golang-grpc-mongodb-arbitrary-data/store/product"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

// server implements the ProductCatalogServiceServer interface.
// It handles the gRPC requests and delegates the actual processing to
// the corresponding functions in the product package.
type server struct {
	productcatalog.UnimplementedProductCatalogServiceServer
	GrpcSrv *grpc.Server
	db      *store.MongoDb
}

// New creates a new instance of the server with the provided database client.
// It sets up the gRPC server, registers the product catalog service,
// and initializes reflection for gRPC server debugging.
func New(db *store.MongoDb) *server {
	grpcServer := grpc.NewServer()
	srv := &server{
		GrpcSrv: grpcServer,
		db:      db}
	productcatalog.RegisterProductCatalogServiceServer(grpcServer, srv)
	reflection.Register(grpcServer)
	return srv
}

// CreateProduct creates a new product in the catalog.
// It delegates the actual creation logic to the product package's Create function.
func (s *server) CreateProduct(ctx context.Context, in *productcatalog.Product) (*productcatalog.Product, error) {
	newProduct, err := mapper.ProductProtobufToProductModel(in)
	if err != nil {
		return nil, err
	}
	createdProduct, err := product.Create(ctx, s.db, newProduct)
	if err != nil {
		return nil, err
	}
	protoResponse, err := mapper.ProductModelToProductProtobuf(createdProduct)
	if err != nil {
		return nil, err
	}
	return protoResponse, nil
}

// GetProduct retrieves a product by its ID from the catalog.
// It delegates the actual retrieval logic to the product package's Get function.
func (s *server) GetProduct(ctx context.Context, in *productcatalog.GetProductRequest) (*productcatalog.Product, error) {
	product, err := product.Get(ctx, s.db, in)
	if err != nil {
		return nil, errors.Wrapf(err, "getting product with uuid %s", in.Uuid)
	}
	protoResponse, err := mapper.ProductModelToProductProtobuf(product)
	if err != nil {
		return nil, err
	}
	return protoResponse, nil
}

// UpdateProduct updates an existing product in the catalog.
// It delegates the actual update logic to the product package's Update function.
func (s *server) UpdateProduct(ctx context.Context, in *productcatalog.Product) (*productcatalog.Product, error) {
	productToUpdate, err := mapper.ProductProtobufToProductModel(in)
	if err != nil {
		return nil, err
	}
	updatedProduct, err := product.Update(ctx, s.db, productToUpdate)
	if err != nil {
		return nil, err
	}
	protoResponse, err := mapper.ProductModelToProductProtobuf(updatedProduct)
	if err != nil {
		return nil, err
	}
	return protoResponse, nil
}

// DeleteProduct deletes a product from the catalog.
// It delegates the actual deletion logic to the product package's Delete function.
func (s *server) DeleteProduct(ctx context.Context, in *productcatalog.DeleteProductRequest) (*productcatalog.DeleteProductResponse, error) {
	resp, err := product.Delete(ctx, s.db, in)
	if err != nil {
		return nil, errors.Wrapf(err, "deleting product with uuid %s", in.Uuid)
	}
	return resp, nil
}

// ListProducts lists all the products in the catalog.
// It delegates the actual listing logic to the product package's ListProducts function.
func (s *server) ListProducts(ctx context.Context, in *productcatalog.ListProductsRequest) (*productcatalog.ListProductsResponse, error) {
	products, err := product.List(ctx, s.db, in)
	if err != nil {
		return nil, errors.Wrap(err, "listing products")
	}
	protoResponse, err := mapper.ProductModelListToListProductsResponse(products)
	if err != nil {
		return nil, err
	}
	return protoResponse, nil
}
