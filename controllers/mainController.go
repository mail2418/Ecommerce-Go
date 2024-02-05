package controllers

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/mail2418/ecommerce-project/database"
	"github.com/mail2418/ecommerce-project/models"
	generate "github.com/mail2418/ecommerce-project/tokens"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

var UserCollection *mongo.Collection = database.UserData(database.Client, "Users")
var ProductCollection *mongo.Collection = database.ProductData(database.Client, "Products")
var Validate = validator.New()

// Sign up
func SignUp() gin.HandlerFunc{
	return func(c *gin.Context){
		ctx, cancel := context.WithTimeout(context.Background(), 100 * time.Second)
		defer cancel()

		var user models.User
		if err := c.BindJSON(&user); err!= nil{
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		validationErr := Validate.Struct(user)
		if validationErr != nil{
			c.JSON(http.StatusBadRequest, gin.H{"error":validationErr})
		}
		count, err := UserCollection.CountDocuments(ctx, bson.M{"email":user.Email})
		if err!= nil{
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":err.Error(),
			})
			return
		}
		if count > 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error":"user already sign in"})
		}
		count,err = UserCollection.CountDocuments(ctx, bson.M{"phone":user.Phone})

		defer cancel()

		if err != nil {
			log.Panic(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error":err})
			return
		}
		if count >0 {
			c.JSON(http.StatusBadRequest,gin.H{"error":"phone already in use"})
		}

		password := HashPassword(*&user.Password)
		user.Password = password
		user.Created_At, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		user.Updated_At, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		user.ID = primitive.NewObjectID()
		user.User_ID = user.ID.Hex()

		token, refreshToken, _ := generate.TokenGenerator(*&user.Email, *&user.First_Name, *&user.Last_Name, user.User_ID)
		user.Token = token
		user.Refresh_Token = refreshToken
		user.User_Cart = make([]models.ProductUser,0)
		user.Address_Details = make([]models.Address,0)
		user.Order_Status = make([]models.Order,0)

		_, inserter := UserCollection.InsertOne(ctx, user)
		if inserter != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error":"user can't created"})
			return
		}
		defer cancel()
		c.JSON(http.StatusCreated,gin.H{"success":"user successfully created"})
	}
}
// Hash Password
func HashPassword(password string) string {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password),14)
	if err != nil {
		log.Panic(err)
	}
	return string(bytes)
}
// Verify Password
func VerifyPassword(userPassword string, givenPassword string) (bool,string) {
	err := bcrypt.CompareHashAndPassword([]byte(givenPassword), []byte(userPassword))
	valid := true
	msg := ""
	if err != nil {
		msg = "Login or Password is incorrect"
		valid = false
	}
	return valid, msg
}

// Login
func Login() gin.HandlerFunc{
	return func(c *gin.Context){
		var ctx,cancel = context.WithTimeout(context.Background(), 100* time.Second)
		defer cancel()
		var user models.User
		var foundUser models.User
		if err := c.BindJSON(&user); err != nil{
			c.JSON(http.StatusBadRequest, gin.H{"error":err})
			return
		}

		err := UserCollection.FindOne(ctx, bson.M{"email":user.Email}).Decode(&foundUser)
		defer cancel()

		if err != nil{
			c.JSON(http.StatusInternalServerError, gin.H{"error":"login or password incorrect"})
			return
		}

		passwordIsValid, msg := VerifyPassword(*&user.Password, *&foundUser.Password)
		defer cancel()

		if !passwordIsValid{
			c.JSON(http.StatusInternalServerError, gin.H{"error":msg})
			fmt.Println(msg)
			return
		}
		token, refreshToken, _ := generate.TokenGenerator(*&foundUser.Email, *&foundUser.First_Name, *&foundUser.Last_Name, foundUser.User_ID)
		defer cancel()
		generate.UpdateAllTokens(token, refreshToken, foundUser.User_ID)

		c.JSON(http.StatusOK, "Successfully login")
	}
}

// Product View Admin
func ProductViewerAdmin() gin.HandlerFunc{
	return func (c *gin.Context){
		ctx,cancel := context.WithTimeout(context.Background(), 100 *time.Second)
		var products models.Product
		defer cancel()
		if err := c.BindJSON(&products); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error":err.Error()})
			return
		}
		products.Product_ID = primitive.NewObjectID()
		_, err := ProductCollection.InsertOne(ctx, products)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error":"not inserted"})
			return
		}
		defer cancel()
		c.JSON(http.StatusOK, "success view product")
	}
}

// Search Product
func SearchProduct() gin.HandlerFunc{
	return func(c *gin.Context){
		var productList []models.Product
		var ctx, cancel = context.WithTimeout(context.Background(), 100 * time.Second)
		defer cancel()
		cursor, err := ProductCollection.Find(ctx, bson.D{{}})

		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, "something went error, please try after some time")
			return
		}
		err = cursor.All(ctx, &productList)
		if err != nil{
			log.Println(err)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		defer cursor.Close(ctx)
		if err := cursor.Err(); err != nil {
			log.Println(err)
			c.IndentedJSON(400, "invalid")
			return
		}
		defer cancel()
		c.IndentedJSON(http.StatusOK, productList)
	}
}

// Search Product by Query
func SearchProductByQuery() gin.HandlerFunc{
	return func(c *gin.Context){
		var searchProducts []models.Product
		queryParam := c.Query("name")

		if queryParam == ""{
			log.Println("query is empty")
			c.Header("Content-Type","application/json")
			c.JSON(http.StatusNotFound, gin.H{"error":"invalid search index"})
			c.Abort()
			return
		}
		ctx, cancel := context.WithTimeout(context.Background(), 100 *time.Second)
		defer cancel()

		searchQueryDB, err := ProductCollection.Find(ctx, bson.M{"product_name":bson.M{"$regex":queryParam}})
		if err != nil {
			c.IndentedJSON(http.StatusNotFound, "something went wrong while fetching datas")
			return
		}
		err = searchQueryDB.All(ctx, &searchProducts)
		if err != nil {
			log.Println(err)
			c.IndentedJSON(400, "invalid")
			return
		}
		defer searchQueryDB.Close(ctx)

		if err := searchQueryDB.Err(); err != nil {
			log.Println(err)
			c.IndentedJSON(400, "invalid request")
			return
		}
		defer cancel()
		c.IndentedJSON(200, searchProducts)
	}
}