package model

import (
	"github.com/dgrijalva/jwt-go"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Inventory struct {
	ID           primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	ItemId       string             `json:"itemid,omitempty"`
	ItemName     string             `json:"name,omitempty"`
	ItemPrice    float64            `json:"price,omitempty"`
	ItemQuantity int                `json:"quantity,omitempty"`
}

type Cart struct {
	ID           primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	ItemId       string             `json:"itemid,omitempty"`
	Username     string             `json:"username,omitempty"`
	ItemPrice    float64            `json:"price,omitempty"`
	ItemQuantity int                `json:"quantity,omitempty"`
}

type Credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type Claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}
