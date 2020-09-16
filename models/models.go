package models

import (
	"time"
)

//Create Struct
type Enrollee struct {
	ID          *string     `json:"id,omitempty" bson:"id,omitempty"`
	Name        *string     `json:"name" bson:"name,omitempty" validate:"required,notblank"`
	IsActive    *bool       `json:"is_active" bson:"is_active,omitempty" validate:"required"`
	BirthDate   *string     `json:"birth_date" bson:"birth_date,omitempty" validate:"required,notblank,datetime=2006-01-02"`
	PhoneNumber *string     `json:"phone_number" bson:"phone_number,omitempty"`
	CreatedAt   time.Time   `json:"created_at" bson:"created_at,omitempty"`
	UpdatedAt   time.Time   `json:"updated_at" bson:"updated_at,omitempty"`
	Dependents  []Dependent `json:"dependents" bson:"dependents" validate:"dive"`
}

type Dependent struct {
	ID        *string   `json:"id,omitempty" bson:"id,omitempty"`
	Name      *string   `json:"name" bson:"name,omitempty" validate:"required,notblank"`
	BirthDate *string   `json:"birth_date" bson:"birth_date,omitempty" validate:"required,notblank,datetime=2006-01-02"`
	CreatedAt time.Time `json:"created_at" bson:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at" bson:"updated_at,omitempty"`
}
