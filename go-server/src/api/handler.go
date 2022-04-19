package api

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"go-server/src/db"
	"math/rand"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func SignUpHandler(c *gin.Context) {
	var err error
	var user AuthUser
	if err := c.ShouldBindJSON(&user); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, HttpResponse{Code: http.StatusBadRequest, Msg: err.Error()})
		return
	}

	user.Token = genToken()
	user.Password = getHash(user.Password)
	user.LoginTime = time.Now()

	_, err = db.Mongo.Collection(string(db.USER_AUTH)).InsertOne(db.Ctx, user)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusOK, HttpResponse{Code: 10001, Msg: err.Error()})
		return
	}

	c.JSON(http.StatusOK, HttpResponse{Code: http.StatusOK, Data: user.Token})
}

func SignInHandler(c *gin.Context) {
	var user AuthUser
	if err := c.ShouldBindJSON(&user); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, HttpResponse{Code: http.StatusBadRequest, Msg: err.Error()})
		return
	}

	token := genToken()
	filter := bson.M{"email": user.Email, "password": getHash(user.Password)}
	update := bson.M{"$set": bson.M{"token": token, "lt": time.Now()}}
	option := options.FindOneAndUpdate()
	option.SetReturnDocument(options.After)

	err := db.Mongo.Collection(string(db.USER_AUTH)).FindOneAndUpdate(db.Ctx, filter, update, option).Decode(&user)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, HttpResponse{Code: http.StatusBadRequest, Msg: err.Error()})
		return
	}

	c.JSON(http.StatusOK, HttpResponse{Code: http.StatusOK, Data: token})
}

func getHash(key string) string {
	hash := md5.New()
	hash.Write([]byte(key))
	return hex.EncodeToString(hash.Sum(nil))
}

func genToken() string {
	b := make([]byte, 8)
	rand.Read(b)
	return fmt.Sprintf("%x", b)
}
