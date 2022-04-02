package api

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type HttpResponse struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

type AuthUser struct {
	Uid       primitive.ObjectID `json:"-" bson:"_id,omitempty"`
	Email     string             `json:"email" bson:"email" binding:"required,email"`
	Password  string             `json:"password,omitempty" bson:"password" binding:"required"`
	Token     string             `json:"token" bson:"token"`
	LoginTime time.Time          `json:"-" bson:"lt,omitempty"`
}
