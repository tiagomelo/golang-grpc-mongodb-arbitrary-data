// Copyright (c) 2023 Tiago Melo. All rights reserved.
// Use of this source code is governed by the MIT License that can be found in
// the LICENSE file.
package product

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tiagomelo/golang-grpc-mongodb-arbitrary-data/api/proto/gen/productcatalog"
	"github.com/tiagomelo/golang-grpc-mongodb-arbitrary-data/store"
	"github.com/tiagomelo/golang-grpc-mongodb-arbitrary-data/store/product/models"
	"go.mongodb.org/mongo-driver/mongo"
	"google.golang.org/protobuf/proto"
)

func TestCreate(t *testing.T) {
	uuidProvider = func() string {
		return "uuid"
	}
	testCases := []struct {
		name                     string
		input                    *models.Product
		mockInsertIntoCollection func(ctx context.Context, collection *mongo.Collection, document interface{}) (*mongo.InsertOneResult, error)
		expectedOutput           *models.Product
		expectedError            error
	}{
		{
			name: "happy path",
			input: &models.Product{
				Name:        "name",
				Description: "description",
				Price:       1,
				Attributes: map[string]interface{}{
					"attr": "value",
				},
			},
			expectedOutput: &models.Product{
				Uuid:        "uuid",
				Name:        "name",
				Description: "description",
				Price:       1,
				Attributes: map[string]interface{}{
					"attr": "value",
				},
			},
			mockInsertIntoCollection: func(ctx context.Context, collection *mongo.Collection, document interface{}) (*mongo.InsertOneResult, error) {
				return &mongo.InsertOneResult{}, nil
			},
		},
		{
			name: "error",
			input: &models.Product{
				Name:        "name",
				Description: "description",
				Price:       1,
				Attributes: map[string]interface{}{
					"attr": "value",
				},
			},
			mockInsertIntoCollection: func(ctx context.Context, collection *mongo.Collection, document interface{}) (*mongo.InsertOneResult, error) {
				return nil, errors.New("random error")
			},
			expectedError: errors.New("inserting product: random error"),
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			insertIntoCollection = tc.mockInsertIntoCollection
			output, err := Create(context.TODO(), &store.MongoDb{DatabaseName: "db", Client: &mongo.Client{}}, tc.input)
			if err != nil {
				if tc.expectedError == nil {
					t.Fatalf("expected no error, got %v", err)
				}
				require.Equal(t, tc.expectedError.Error(), err.Error())
			} else {
				if tc.expectedError != nil {
					t.Fatalf("expected error %v, got nil", tc.expectedError)
				}
				require.Equal(t, tc.expectedOutput, output)
			}
		})
	}
}

func TestGet(t *testing.T) {
	testCases := []struct {
		name           string
		mockFindOne    func(ctx context.Context, collection *mongo.Collection, filter interface{}, p *models.Product) error
		expectedOutput *models.Product
		expectedError  error
	}{
		{
			name: "happy path",
			mockFindOne: func(ctx context.Context, collection *mongo.Collection, filter interface{}, p *models.Product) error {
				p.Uuid = "uuid"
				p.Name = "name"
				p.Description = "description"
				p.Price = 1
				p.Attributes = map[string]interface{}{
					"color": "blue",
					"size":  12.0,
				}
				return nil
			},
			expectedOutput: &models.Product{
				Uuid:        "uuid",
				Name:        "name",
				Description: "description",
				Price:       1,
				Attributes: map[string]interface{}{
					"color": "blue",
					"size":  12.0,
				},
			},
		},
		{
			name: "error",
			mockFindOne: func(ctx context.Context, collection *mongo.Collection, filter interface{}, p *models.Product) error {
				return errors.New("random error")
			},
			expectedError: errors.New(`getting product with uuid "uuid": random error`),
		},
		{
			name: "document not found",
			mockFindOne: func(ctx context.Context, collection *mongo.Collection, filter interface{}, p *models.Product) error {
				return mongo.ErrNoDocuments
			},
			expectedError: errors.New(`product with uuid "uuid" does not exist`),
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			findOne = tc.mockFindOne
			output, err := Get(context.TODO(), &store.MongoDb{DatabaseName: "db", Client: &mongo.Client{}}, &productcatalog.GetProductRequest{Uuid: "uuid"})
			if err != nil {
				if tc.expectedError == nil {
					t.Fatalf("expected no error, got %v", err)
				}
				require.Equal(t, tc.expectedError.Error(), err.Error())
			} else {
				if tc.expectedError != nil {
					t.Fatalf("expected error %v, got nil", tc.expectedError)
				}
				require.Equal(t, tc.expectedOutput, output)
			}
		})
	}
}

func TestUpdate(t *testing.T) {
	testCases := []struct {
		name           string
		input          *models.Product
		mockUpdateOne  func(ctx context.Context, collection *mongo.Collection, filter interface{}, update interface{}) (*mongo.UpdateResult, error)
		expectedOutput *models.Product
		expectedError  error
	}{
		{
			name: "happy path",
			input: &models.Product{
				Uuid:        "uuid",
				Name:        "name",
				Description: "description",
				Price:       1,
				Attributes: map[string]interface{}{
					"color": "blue",
					"size":  12.0,
				},
			},
			mockUpdateOne: func(ctx context.Context, collection *mongo.Collection, filter, update interface{}) (*mongo.UpdateResult, error) {
				return &mongo.UpdateResult{}, nil
			},
			expectedOutput: &models.Product{
				Uuid:        "uuid",
				Name:        "name",
				Description: "description",
				Price:       1,
				Attributes: map[string]interface{}{
					"color": "blue",
					"size":  12.0,
				},
			},
		},
		{
			name: "error",
			input: &models.Product{
				Uuid:        "uuid",
				Name:        "name",
				Description: "description",
				Price:       1,
				Attributes: map[string]interface{}{
					"color": "blue",
					"size":  12.0,
				},
			},
			mockUpdateOne: func(ctx context.Context, collection *mongo.Collection, filter, update interface{}) (*mongo.UpdateResult, error) {
				return nil, errors.New("random error")
			},
			expectedError: errors.New(`updating product with uuid "uuid": random error`),
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			updateOne = tc.mockUpdateOne
			output, err := Update(context.TODO(), &store.MongoDb{DatabaseName: "db", Client: &mongo.Client{}}, tc.input)
			if err != nil {
				if tc.expectedError == nil {
					t.Fatalf("expected no error, got %v", err)
				}
				require.Equal(t, tc.expectedError.Error(), err.Error())
			} else {
				if tc.expectedError != nil {
					t.Fatalf("expected error %v, got nil", tc.expectedError)
				}
				require.Equal(t, tc.expectedOutput, output)
			}
		})
	}
}

func TestDelete(t *testing.T) {
	testCases := []struct {
		name           string
		mockDeleteOne  func(ctx context.Context, collection *mongo.Collection, filter interface{}) (*mongo.DeleteResult, error)
		expectedOutput *productcatalog.DeleteProductResponse
		expectedError  error
	}{
		{
			name: "happy path",
			mockDeleteOne: func(ctx context.Context, collection *mongo.Collection, filter interface{}) (*mongo.DeleteResult, error) {
				return nil, nil
			},
			expectedOutput: &productcatalog.DeleteProductResponse{
				Result: "success",
			},
		},
		{
			name: "error",
			mockDeleteOne: func(ctx context.Context, collection *mongo.Collection, filter interface{}) (*mongo.DeleteResult, error) {
				return nil, errors.New("random error")
			},
			expectedError: errors.New(`deleting product with uuid "uuid": random error`),
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			deleteOne = tc.mockDeleteOne
			output, err := Delete(context.TODO(), &store.MongoDb{DatabaseName: "db", Client: &mongo.Client{}}, &productcatalog.DeleteProductRequest{Uuid: "uuid"})
			if err != nil {
				if tc.expectedError == nil {
					t.Fatalf("expected no error, got %v", err)
				}
				require.Equal(t, tc.expectedError.Error(), err.Error())
			} else {
				if tc.expectedError != nil {
					t.Fatalf("expected error %v, got nil", tc.expectedError)
				}
				require.True(t, proto.Equal(tc.expectedOutput, output))
			}
		})
	}
}

func TestList(t *testing.T) {
	testCases := []struct {
		name           string
		mockFind       func(ctx context.Context, collection *mongo.Collection, filter interface{}) (Cursor, error)
		expectedOutput []*models.Product
		expectedError  error
	}{
		{
			name: "happy path",
			mockFind: func(ctx context.Context, collection *mongo.Collection, filter interface{}) (Cursor, error) {
				data := []models.Product{
					{
						Uuid:        "id",
						Name:        "name",
						Description: "description",
						Price:       1,
						Attributes: map[string]interface{}{
							"color": "blue",
							"size":  12.0,
						},
					},
					{
						Uuid:        "id2",
						Name:        "name2",
						Description: "description2",
						Price:       12,
						Attributes: map[string]interface{}{
							"color": "red",
							"size":  1.0,
						},
					},
				}
				return &MockCursor{data: data}, nil
			},
			expectedOutput: []*models.Product{
				{
					Uuid:        "id",
					Name:        "name",
					Description: "description",
					Price:       1,
					Attributes: map[string]interface{}{
						"color": "blue",
						"size":  12.0,
					},
				},
				{
					Uuid:        "id2",
					Name:        "name2",
					Description: "description2",
					Price:       12,
					Attributes: map[string]interface{}{
						"color": "red",
						"size":  1.0,
					},
				},
			},
		},
		{
			name: "error when finding products",
			mockFind: func(ctx context.Context, collection *mongo.Collection, filter interface{}) (Cursor, error) {
				return nil, errors.New("random error")
			},
			expectedError: errors.New("finding products: random error"),
		},
		{
			name: "error when decoding product",
			mockFind: func(ctx context.Context, collection *mongo.Collection, filter interface{}) (Cursor, error) {
				data := []models.Product{
					{
						Uuid:        "id",
						Name:        "name",
						Description: "description",
						Price:       1,
						Attributes: map[string]interface{}{
							"color": "blue",
							"size":  12.0,
						},
					},
				}
				return &MockCursor{data: data, decodeErr: errors.New("random error")}, nil
			},
			expectedError: errors.New("decoding product: random error"),
		},
		{
			name: "error in cursor",
			mockFind: func(ctx context.Context, collection *mongo.Collection, filter interface{}) (Cursor, error) {
				data := []models.Product{
					{
						Uuid:        "id",
						Name:        "name",
						Description: "description",
						Price:       1,
						Attributes: map[string]interface{}{
							"color": "blue",
							"size":  12.0,
						},
					},
				}
				return &MockCursor{data: data, err: errors.New("random error")}, nil
			},
			expectedError: errors.New("cursor error: random error"),
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			find = tc.mockFind
			output, err := List(context.TODO(), &store.MongoDb{DatabaseName: "db", Client: &mongo.Client{}}, &productcatalog.ListProductsRequest{})
			if err != nil {
				if tc.expectedError == nil {
					t.Fatalf("expected no error, got %v", err)
				}
				require.Equal(t, tc.expectedError.Error(), err.Error())
			} else {
				if tc.expectedError != nil {
					t.Fatalf("expected error %v, got nil", tc.expectedError)
				}
				require.Equal(t, tc.expectedOutput, output)
			}
		})
	}
}

type MockCursor struct {
	data      []models.Product
	index     int
	decodeErr error
	err       error
}

func (m *MockCursor) Next(ctx context.Context) bool {
	if m.index < len(m.data) {
		m.index++
		return true
	}
	return false
}

func (m *MockCursor) Decode(val interface{}) error {
	if m.decodeErr != nil {
		return m.decodeErr
	}
	product, ok := val.(*models.Product)
	if !ok {
		return errors.New("Decode type not *models.Product")
	}
	if m.index <= 0 || m.index > len(m.data) {
		return errors.New("No data to decode")
	}
	*product = m.data[m.index-1]
	return nil
}

func (m *MockCursor) Err() error {
	return m.err
}

func (m *MockCursor) Close(ctx context.Context) error {
	return nil
}
