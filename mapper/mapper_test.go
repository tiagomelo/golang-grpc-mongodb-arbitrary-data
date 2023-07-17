// Copyright (c) 2023 Tiago Melo. All rights reserved.
// Use of this source code is governed by the MIT License that can be found in
// the LICENSE file.
package mapper

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tiagomelo/golang-grpc-mongodb-arbitrary-data/api/proto/gen/productcatalog"
	"github.com/tiagomelo/golang-grpc-mongodb-arbitrary-data/store/product/models"
	"google.golang.org/protobuf/types/known/structpb"
)

func TestProductProtobufToProductModel(t *testing.T) {
	p := &productcatalog.Product{
		Uuid:        "1",
		Name:        "Product1",
		Description: "Product Description",
		Price:       10.0,
		Attributes:  map[string]*structpb.Value{"Color": structpb.NewStringValue("Blue")},
	}

	dbProduct, err := ProductProtobufToProductModel(p)
	assert.Nil(t, err)
	assert.Equal(t, "1", dbProduct.Uuid)
	assert.Equal(t, "Product1", dbProduct.Name)
	assert.Equal(t, "Product Description", dbProduct.Description)
	assert.Equal(t, float32(10.0), dbProduct.Price)
	assert.Equal(t, map[string]interface{}{"Color": "Blue"}, dbProduct.Attributes)
}

func TestProdutcModelToProductProtobuf(t *testing.T) {
	testCases := []struct {
		name                 string
		input                *models.Product
		mockStructpbNewValue func(v interface{}) (*structpb.Value, error)
		expectedOutput       *productcatalog.Product
		expectedError        error
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
			expectedOutput: &productcatalog.Product{
				Uuid:        "uuid",
				Name:        "name",
				Description: "description",
				Price:       1,
				Attributes: map[string]*structpb.Value{
					"color": structpb.NewStringValue("blue"),
					"size":  structpb.NewNumberValue(12.0),
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
			mockStructpbNewValue: func(v interface{}) (*structpb.Value, error) {
				return nil, errors.New("random error")
			},
			expectedError: errors.New(`parsing attribute "color": random error`),
		},
	}
	originalStructpbNewValue := structpbNewValue
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.mockStructpbNewValue != nil {
				structpbNewValue = tc.mockStructpbNewValue
			} else {
				structpbNewValue = originalStructpbNewValue
			}
			defer func() { structpbNewValue = originalStructpbNewValue }()
			output, err := ProductModelToProductProtobuf(tc.input)
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

func TestProductModelListToListProductsResponse(t *testing.T) {
	testCases := []struct {
		name                 string
		input                []*models.Product
		mockStructpbNewValue func(v interface{}) (*structpb.Value, error)
		expectedOutput       *productcatalog.ListProductsResponse
		expectedError        error
	}{
		{
			name: "happy path",
			input: []*models.Product{
				{
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
			expectedOutput: &productcatalog.ListProductsResponse{
				Products: []*productcatalog.Product{
					{
						Uuid:        "uuid",
						Name:        "name",
						Description: "description",
						Price:       1,
						Attributes: map[string]*structpb.Value{
							"color": structpb.NewStringValue("blue"),
							"size":  structpb.NewNumberValue(12.0),
						},
					},
				},
			},
		},
		{
			name: "error",
			input: []*models.Product{
				{
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
			mockStructpbNewValue: func(v interface{}) (*structpb.Value, error) {
				return nil, errors.New("random error")
			},
			expectedError: errors.New(`parsing attribute "color": random error`),
		},
	}
	originalStructpbNewValue := structpbNewValue
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.mockStructpbNewValue != nil {
				structpbNewValue = tc.mockStructpbNewValue
			} else {
				structpbNewValue = originalStructpbNewValue
			}
			defer func() { structpbNewValue = originalStructpbNewValue }()
			output, err := ProductModelListToListProductsResponse(tc.input)
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
