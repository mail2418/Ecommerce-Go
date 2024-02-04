package controllers

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/mail2418/ecommerce-project/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func AddAddress() gin.HandlerFunc{
	return func(c *gin.Context){
		user_id := c.Query("id")
		if user_id == ""{
			c.Header("Content-Type","application/json")
			c.JSON(http.StatusNotFound, gin.H{"error":"invalid code"})
			c.Abort()
			return
		}
		address, err := ObjectIDFromHex(user_id)
		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, "internal server error")
		}

		var addresses models.Address
		addresses.Address_ID = primitive.NewObjectID()
		if err = c.BindJSON(&addresses); err != nil {
			c.IndentedJSON(http.StatusNotAcceptable,err.Error())
		}

		ctx, cancel := context.WithTimeout(context.Background(), 100 *time.Second)
		match_filter := bson.D{{Key: "$match", Value: bson.D{primitive.E{Key: "id", Value: address}}}}
		unwind := bson.D{{Key: "$unwind",Value: bson.D{primitive.E{Key: "path", Value: "$address"}}}}
		group := bson.D{{Key: "$group", Value: bson.D{primitive.E{Key: "id", Value: "$address_id"}, {Key: "count", Value: bson.D{primitive.E{Key: "$sum", Value: 1}}}}}}
		pointcursor, err := UserCollection.Aggregate(ctx, mongo.Pipeline{match_filter, unwind, group})

		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, "internal server error")
		}
	}
}

func EditHomeAddress() gin.HandlerFunc{
	return func(c *gin.Context){
		
	}
}

func EditWorkAddress() gin.HandlerFunc{
	return func(c *gin.Context){
		
	}
}

func DeleteAddress() gin.HandlerFunc{
	return func(c *gin.Context){
		user_id := c.Query("id")
		if user_id == ""{
			c.Header("Content-Type","application/json")
			c.JSON(http.StatusNotFound, gin.H{"error":"invalid search index"})
			c.Abort()
			return
		}
		addresses:= make([]models.Address,0)
		new_user_id, err := primitive.ObjectIDFromHex(user_id)
		if err != nil{
			c.IndentedJSON(http.StatusInternalServerError, "internal server error")
		}
		ctx, cancel := context.WithTimeout(context.Background(), 100 *time.Second)
		defer cancel()

		filter := bson.D{primitive.E{Key:"id",Value: new_user_id}}
		update := bson.D{{Key: "$set", Value: bson.D{primitive.E{Key: "address", Value: addresses}}}}
		_, err = UserCollection.UpdateOne(ctx, filter, update)
		if err != nil {
			c.IndentedJSON(http.StatusNotFound, "wrong command")
			return
		}
		defer cancel()
		ctx.Done()
		c.IndentedJSON(http.StatusOK, "successfully delete address")

	}
}