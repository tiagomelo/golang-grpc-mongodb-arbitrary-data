// Copyright (c) 2023 Tiago Melo. All rights reserved.
// Use of this source code is governed by the MIT License that can be found in
// the LICENSE file.
//
// Package models provides the data models used in the application.
package models

// Product represents a product with its associated attributes.
type Product struct {
	Uuid        string                 `bson:"uuid"`
	Name        string                 `bson:"name"`
	Description string                 `bson:"description"`
	Price       float32                `bson:"price"`
	Attributes  map[string]interface{} `bson:"attributes"`
}
