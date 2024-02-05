package database

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/mail2418/ecommerce-project/controllers"
	"github.com/mail2418/ecommerce-project/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	ErrCantFindProduct = errors.New("error cant find product")
	ErrCantDecodeProducts = errors.New("error cant find product")
	ErrUserIdIsNotValid = errors.New("error user id is not valid")
	ErrCantUpdateUser = errors.New("error cant update user")
	ErrCantRemoveItemCart = errors.New("error cant remove item from cart")
	ErrCantGetItem = errors.New("error cant get item from cart")
	ErrCantBuyCartItem = errors.New("error cant update the purchase")
)

func AddProductToCart(ctx context.Context, prodCollection, useColletion *mongo.Collection, productID primitive.ObjectID, userID string) error {
	searchfromdb, err := prodCollection.Find(ctx, bson.M{"id":productID})
	if err != nil {
		log.Println(err)
		return ErrCantFindProduct
	}
	var productCart []models.ProductUser
	err = searchfromdb.All(ctx, &productCart)
	if err != nil {
		log.Println(err)
		return ErrCantDecodeProducts
	}
	id,err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		log.Println(err)
		return ErrUserIdIsNotValid
	}
	filter := bson.D{primitive.E{Key: "id",Value: id}}
	update := bson.D{{Key: "$push", Value: bson.D{primitive.E{Key: "usercart", Value: bson.D{{Key: "$each", Value: productCart}}}}}}
	_,err = controllers.UserCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		return ErrCantUpdateUser
	}
	return nil

}

func RemoveCartItem(ctx context.Context, prodCollection, userCollection *mongo.Collection, productID primitive.ObjectID, userID string) error {
	id, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		log.Println(err)
		return ErrUserIdIsNotValid
	}
	filter := bson.D(primitive.E{Key: "id", Value: id}) 
	update := bson.M{"$pull":bson.M{"usercart":bson.M{"id":productID}}}
	_, err := controllers.UserCollection.UpdateMany(ctx, filter, update)
	if err != nil {
		return ErrCantRemoveItemCart
	}
	return nil
}

func BuyItemFromCart(ctx context.Context, userCollection *mongo.Collection, userID string) error {
	//  fetch cart from user
	//  finnd cart total
	// create order with items
	// add order to user collection
	// add items in cart to order list
	// empty cart
	id,err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		log.Println(err)
		return ErrUserIdIsNotValid
	}
	var getcartitems models.User
	var ordercart models.Order

	ordercart.Order_ID = primitive.NewObjectID()
	ordercart.Ordered_At = time.Now()
	ordercart.Order_Cart = make([]models.ProductUser, 0)
	ordercart.Payment_Method.COD = true

	unwind := bson.D{{Key: "$unwind", Value: bson.D{primitive.E{Key: "path", Value: "$usercart"}}}}
	grouping := bson.D{{Key: "$group", Value: bson.D{primitive.E{Key: "id", Value: "$id"}, {Key: "total", Value: bson.D{primitive.E{Key: "$sum", Value: "$usercart.price"}}}}}}
	currentResults, err := controllers.UserCollection.Aggregate(ctx, mongo.Pipeline{unwind, grouping})
	ctx.Done()
	if err != nil {
		panic(err)
	}
	var getusercart []bson.M
	if err = currentResults.All(ctx, &getusercart); err != nil {
		panic(err)
	}
	var total_price int32

	for _,user_item := range getusercart{
		price := user_item["total"]
		total_price = price.(int32)
	}
	ordercart.Price = int64(total_price)
	filter := bson.D{primitive.E{Key: "id", Value: id}}
	update := bson.D{{Key: "$push", Value: bson.D{primitive.E{Key: "orders",Value: ordercart}}}}
	_, err = controllers.UserCollection.UpdateMany(ctx, filter, update)
	if err != nil {
		panic(err)
	}
	err = controllers.UserCollection.FindOne(ctx, bson.D{primitive.E{Key: "id", Value: id}}).Decode(&getcartitems)
	if err != nil {
		panic(err)
	}
	filter2 := bson.D{primitive.E{Key: "id", Value: id}}
	update2 := bson.M{"$push":bson.M{"orders.$[].order_list":bson.M{"$each":getcartitems.User_Cart}}}
	_, err = controllers.UserCollection.UpdateOne(ctx, filter2, update2)
	if err != nil {
		log.Println(err)
	}
	usercart_empty := make([]models.ProductUser,0)
	filter3 := bson.D{primitive.E{Key: "id", Value: id}}
	update3 := bson.D{{Key: "$set", Value:bson.D{primitive.E{Key: "usercart", Value: usercart_empty}}}}
	_, err = controllers.UserCollection.UpdateOne(ctx, filter3, update3)
	if err != nil {
		return ErrCantBuyCartItem
	}
	return nil
}

func InstantBuyer(){
	
}