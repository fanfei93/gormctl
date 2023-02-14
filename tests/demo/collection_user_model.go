package model

import (
    "gorm.io/gorm"
)

var _ CollectionUserModel = (*customCollectionUserModel)(nil)

type (
    // CollectionUserModel is an interface to be customized, add more methods here,
    // and implement the added methods in customCollectionUserModel.
    CollectionUserModel interface {
        collectionUserModel
    }

    customCollectionUserModel struct {
        *defaultCollectionUserModel
    }
)

// NewCollectionUserModel returns a model for the database table.
func NewCollectionUserModel(conn *gorm.DB) CollectionUserModel {
	return &customCollectionUserModel{
		defaultCollectionUserModel: newCollectionUserModel(conn),
	}
}
