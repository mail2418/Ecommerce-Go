package main

import (
	// "log"
	// "os"

	// "github.com/gin-gonic/gin"
	// "github.com/mail2418/ecommerce-project/controllers"
	// "github.com/mail2418/ecommerce-project/database"
	// "github.com/mail2418/ecommerce-project/middleware"
	// "github.com/mail2418/ecommerce-project/routes"
)

// func main(){
// 	port := os.Getenv("PORT")
// 	if port == "" {
// 		port = "8080"
// 	}
// 	app := controllers.NewApplication(database.ProductData(database.Client, "Products"), database.UserData(database.Client, "Users"))
// 	router := gin.New()
// 	router.Use(gin.Logger())

// 	routes.UserRoutes(router)
// 	router.Use(middleware.Authentication())

// 	router.GET("/addtocart", app.AddToCart())
// 	router.GET("/removeitem", app.RemoveItem())
// 	router.GET("/cartcheckout", app.BuyFromCart())
// 	router.GET("/instantbuy", app.InstantBuy())

// 	log.Fatal(router.Run(":"+port))
// }