package controllers

import (
	"context"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/mail2418/ecommerce-project/database"
	"github.com/mail2418/ecommerce-project/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)


type Application struct{
	prodCollection *mongo.Collection
	userCollection *mongo.Collection
}

func NewApplication(prodCollection, userCollection *mongo.Collection) *Application{
	return &Application{
		prodCollection: prodCollection,
		userCollection: userCollection,
	}
}

func (app *Application) AddToCart() gin.HandlerFunc{
	return func (c *gin.Context){
		productQueryID := c.Query("id")
		if productQueryID == ""{
			log.Println("product id is empty")
			c.AbortWithError(http.StatusBadRequest, errors.New("product id is empty"))
			return
		}
		userQueryID := c.Query("userID")
		if userQueryID == ""{
			log.Println("user id is empty")
			c.AbortWithError(http.StatusBadRequest, errors.New("user id is empty"))
			return
		}
		productID, err := primitive.ObjectIDFromHex(productQueryID)

		if err != nil {
			log.Println(err)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		var ctx, cancel = context.WithTimeout(context.Background(), 5 * time.Second)
		defer cancel()

		err =  database.AddProductToCart(ctx, app.prodCollection, app.userCollection, productID, userQueryID)
		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, err)
		}
		c.IndentedJSON(200, "successfully added to the cart")
	}
}

func (app *Application)RemoveItem() gin.HandlerFunc{
	return func(c *gin.Context){
		productQueryID := c.Query("id")
		if productQueryID == ""{
			log.Println("product id is empty")
			c.AbortWithError(http.StatusBadRequest, errors.New("product id is empty"))
			return
		}
		userQueryID := c.Query("userID")
		if userQueryID == ""{
			log.Println("user id is empty")
			c.AbortWithError(http.StatusBadRequest, errors.New("user id is empty"))
			return
		}
		productID, err := primitive.ObjectIDFromHex(productQueryID)

		if err != nil {
			log.Println(err)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		var ctx, cancel = context.WithTimeout(context.Background(), 5 * time.Second)
		defer cancel()

		err = database.RemoveCartItem(ctx, app.prodCollection, app.userCollection, productID, userQueryID)
		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, err)
			return
		}
		c.IndentedJSON(http.StatusOK, "successfully removed item from cart")
	}
}

func (app *Application) GetItemFromCart() gin.HandlerFunc{
	return func(c *gin.Context){
		user_id := c.Query("id")
		if user_id == ""{
			c.Header("Content-Type","application/json")
			c.JSON(http.StatusNotFound, gin.H{"error":"invalid id"})
			c.Abort()
			return
		}
		new_user_id, _ := primitive.ObjectIDFromHex(user_id)
		ctx, cancel := context.WithTimeout(context.Background(), 100 * time.Second)
		defer cancel()
		var filledcart models.User
		err := UserCollection.FindOne(ctx, bson.D{primitive.E{Key: "id", Value: new_user_id}}).Decode(&filledcart)
		if err != nil{
			log.Println(err)
			c.IndentedJSON(http.StatusInternalServerError, "not found")
			return
		}
		filter_match := bson.D{{Key: "$match",Value: bson.D{primitive.E{Key: "id", Value: new_user_id}}}}
		unwind := bson.D{{Key: "$unwind",Value: bson.D{primitive.E{Key: "path", Value: "$usercart"}}}}
		grouping := bson.D{{Key: "$group",Value: bson.D{primitive.E{Key: "id", Value: "$id"}, {Key: "total",Value: bson.D{primitive.E{Key: "$sum", Value: new_user_id}}}}}}
		pointcursor, err := UserCollection.Aggregate(ctx, mongo.Pipeline{filter_match, unwind, grouping})
		if err != nil{
			log.Println(err)
		}
		var listing []bson.M
		if err = pointcursor.All(ctx, &listing); err != nil{
			log.Println(err)
			c.AbortWithStatus(http.StatusInternalServerError)
		}
		for _,json := range listing{
			c.IndentedJSON(http.StatusOK, json["total"])
			c.IndentedJSON(http.StatusOK, filledcart.User_Cart)
		}
	}
}

func (app *Application) BuyFromCart() gin.HandlerFunc{
	return func(c *gin.Context){
		userQueryID := c.Query("id")
		if userQueryID == ""{
			log.Panic("user id is empty")
			c.AbortWithError(http.StatusBadRequest, errors.New("user id is empty"))
			return
		}
		ctx, cancel := context.WithTimeout(context.Background(), 100 * time.Second)
		defer cancel()
		err := database.BuyItemFromCart(ctx, app.userCollection, userQueryID)

		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, err)
		}
		c.IndentedJSON(http.StatusOK,"successfully buy item")
		
	}
}

func (app *Application)InstantBuy() gin.HandlerFunc{
	return func(c *gin.Context){
		productQueryID := c.Query("id")
		if productQueryID == ""{
			log.Println("product id is empty")
			c.AbortWithError(http.StatusBadRequest, errors.New("product id is empty"))
			return
		}
		userQueryID := c.Query("userID")
		if userQueryID == ""{
			log.Println("user id is empty")
			c.AbortWithError(http.StatusBadRequest, errors.New("user id is empty"))
			return
		}
		productID, err := primitive.ObjectIDFromHex(productQueryID)

		if err != nil {
			log.Println(err)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		var ctx, cancel = context.WithTimeout(context.Background(), 5 * time.Second)
		defer cancel()

		err = database.InstantBuyer(ctx, app.prodCollection, app.userCollection, productID, userQueryID)
		if err != nil{
			c.IndentedJSON(http.StatusInternalServerError, err)
		}
		c.IndentedJSON(http.StatusOK, "successfully placed to cart")
	}
}