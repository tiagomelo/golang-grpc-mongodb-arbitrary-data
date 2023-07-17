// Copyright (c) 2023 Tiago Melo. All rights reserved.
// Use of this source code is governed by the MIT License that can be found in
// the LICENSE file.
package server

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tiagomelo/golang-grpc-mongodb-arbitrary-data/api/proto/gen/productcatalog"
	"github.com/tiagomelo/golang-grpc-mongodb-arbitrary-data/config"
	"github.com/tiagomelo/golang-grpc-mongodb-arbitrary-data/store"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/structpb"
)

var (
	ctx context.Context
	db  *store.MongoDb
)

const host = "localhost:4444"

func TestMain(m *testing.M) {
	ctx = context.Background()
	const envFilePath = "../.env"
	cfg, err := config.Read(envFilePath)
	if err != nil {
		fmt.Println("error when reading config for integration tests:", err)
		os.Exit(1)
	}
	db, err = store.Connect(ctx, cfg.MongodbTestHostName, cfg.MongodbTestDatabase, cfg.MongodbTestPort)
	if err != nil {
		fmt.Println("error when connecting to MongoDB:", err)
		os.Exit(1)
	}
	lis, err := net.Listen("tcp", host)
	if err != nil {
		fmt.Printf("Failed to listen: %v\n", err)
		os.Exit(1)
	}
	defer lis.Close()
	srv := New(db)
	go func() {
		grpcServer := grpc.NewServer()
		productcatalog.RegisterProductCatalogServiceServer(grpcServer, srv)
		reflection.Register(grpcServer)
		log.Println("Server started")
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("Server error: %v", err)
		}
	}()
	exitVal := m.Run()
	if err := db.Database(cfg.MongodbTestDatabase).Drop(ctx); err != nil {
		fmt.Println("error when dropping test MongoDB:", err)
		os.Exit(1)
	}
	os.Exit(exitVal)
}

func TestProduct(t *testing.T) {
	conn, err := grpc.Dial(host, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatalf("Failed to dial server: %v", err)
	}
	defer conn.Close()
	client := productcatalog.NewProductCatalogServiceClient(conn)

	_newProduct := newProduct()
	_newProduct2 := newProduct()

	// Create two products.
	t.Run("Create", func(t *testing.T) {
		response, err := client.CreateProduct(ctx, _newProduct)
		require.Nil(t, err)
		require.NotNil(t, response)
		require.Equal(t, _newProduct.Name, response.Name)
		require.Equal(t, _newProduct.Price, response.Price)
		require.Equal(t, _newProduct.Description, response.Description)
		for k, v := range response.Attributes {
			require.Equal(t, _newProduct.Attributes[k].AsInterface(), v.AsInterface())
		}

		response2, err2 := client.CreateProduct(ctx, _newProduct2)
		require.Nil(t, err2)
		require.NotNil(t, response2)
		require.Equal(t, _newProduct2.Name, response2.Name)
		require.Equal(t, _newProduct2.Price, response2.Price)
		require.Equal(t, _newProduct2.Description, response2.Description)
		for k, v := range response2.Attributes {
			require.Equal(t, _newProduct2.Attributes[k].AsInterface(), v.AsInterface())
		}

		_newProduct.Uuid = response.Uuid
		_newProduct2.Uuid = response2.Uuid
	})

	// Get the first product.
	t.Run("Get", func(t *testing.T) {
		response, err := client.GetProduct(ctx, &productcatalog.GetProductRequest{Uuid: _newProduct.Uuid})
		require.Nil(t, err)
		require.NotNil(t, response)
		require.True(t, proto.Equal(_newProduct, response))
	})

	// List the products.
	t.Run("List", func(t *testing.T) {
		response, err := client.ListProducts(ctx, &productcatalog.ListProductsRequest{})
		require.Nil(t, err)
		require.NotNil(t, response)
		require.True(t, proto.Equal(products(_newProduct.Uuid, _newProduct2.Uuid), response))
	})

	// Update the second product.
	t.Run("Update", func(t *testing.T) {
		_updatedProduct := updatedProduct(_newProduct2.Uuid)
		response, err := client.UpdateProduct(ctx, _updatedProduct)
		require.Nil(t, err)
		require.NotNil(t, response)
		require.True(t, proto.Equal(_updatedProduct, response))
	})

	// Delete the first product.
	t.Run("Delete", func(t *testing.T) {
		response, err := client.DeleteProduct(ctx, &productcatalog.DeleteProductRequest{Uuid: _newProduct.Uuid})
		require.Nil(t, err)
		require.NotNil(t, response)
		require.True(t, proto.Equal(deletedProductResponse(), response))
	})

	// List the products again. There should be only the updated product.
	t.Run("List", func(t *testing.T) {
		_updatedProduct := updatedProduct(_newProduct2.Uuid)
		response, err := client.ListProducts(ctx, &productcatalog.ListProductsRequest{})
		require.Nil(t, err)
		require.NotNil(t, response)
		require.True(t, proto.Equal(_updatedProduct, response.Products[0]))
	})
}

func newProduct() *productcatalog.Product {
	return &productcatalog.Product{
		Name:        "Test Product Name",
		Description: "Test Product Description",
		Price:       9.99,
		Attributes: map[string]*structpb.Value{
			"color": structpb.NewStringValue("blue"),
			"size":  structpb.NewNumberValue(12),
		},
	}
}

func updatedProduct(id string) *productcatalog.Product {
	return &productcatalog.Product{
		Uuid:        id,
		Name:        "Test Product Name updated",
		Description: "Test Product Description",
		Price:       9.99,
		Attributes: map[string]*structpb.Value{
			"color": structpb.NewStringValue("red"),
			"size":  structpb.NewNumberValue(15),
		},
	}
}

func products(productId1, productId2 string) *productcatalog.ListProductsResponse {
	return &productcatalog.ListProductsResponse{
		Products: []*productcatalog.Product{
			{
				Uuid:        productId1,
				Name:        "Test Product Name",
				Description: "Test Product Description",
				Price:       9.99,
				Attributes: map[string]*structpb.Value{
					"color": structpb.NewStringValue("blue"),
					"size":  structpb.NewNumberValue(12),
				},
			},
			{
				Uuid:        productId2,
				Name:        "Test Product Name",
				Description: "Test Product Description",
				Price:       9.99,
				Attributes: map[string]*structpb.Value{
					"color": structpb.NewStringValue("blue"),
					"size":  structpb.NewNumberValue(12),
				},
			},
		},
	}
}

func deletedProductResponse() *productcatalog.DeleteProductResponse {
	return &productcatalog.DeleteProductResponse{
		Result: "success",
	}
}
