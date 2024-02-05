package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID              primitive.ObjectID `json:"id" bson:"id"`
	User_ID         string             `json:"user_id"`
	First_Name      string             `json:"first_name" validate:"required,min=2,max=30"`
	Last_Name       string             `json:"last_name" validate:"required,min=2,max=30"`
	Password        string             `json:"password" validate:"required,min=6"`
	Email           string             `json:"email" validate:"required,email"`
	Phone           string             `json:"phone" validate:"required"`
	Token           string             `json:"token"`
	Refresh_Token  string             `json:"refresh_token"`
	User_Cart       []ProductUser      `json:"user_cart" bson:"usercart"`
	Address_Details []Address          `json:"address_details" bson:"address"`
	Order_Status    []Order              `json:"order_status"`
	Created_At      time.Time          `json:"created_at"`
	Updated_At      time.Time          `json:"updated_at"`
}

type Product struct {
	Product_ID   primitive.ObjectID `bson:"id"`
	Product_Name string             `json:"product_name"`
	Price        uint64             `json:"price"`
	Rating       uint8              `json:"rating"`
	Image        string             `json:"image"`
}

type ProductUser struct {
	Product_ID   primitive.ObjectID `bson:"id"`
	Product_Name string             `json:"product_name" bson:"product_name"`
	Price        uint64             `json:"price" bson:"price"`
	Rating       uint8              `json:"rating" bson:"rating"`
	Image        string             `json:"image" bson:"image"`
}

type Address struct {
	Address_ID primitive.ObjectID `bson:"id"`
	House      string             `json:"house" bson:"house"`
	Street     string             `json:"street" bson:"street"`
	City       string             `json:"city" bson:"city"`
	Post_Code  string             `json:"post_code" bson:"post_code"`
}

type Order struct {
	Order_ID       primitive.ObjectID `bson:"id"`
	Order_Cart     []ProductUser      `json:"order_cart" bson:"order_cart"`
	Ordered_At     time.Time          `json:"ordered_at" bson:"ordered_at"`
	Price          int64             `json:"price" bson:"price"`
	Discount       uint16             `json:"discount" bson:"discount"`
	Payment_Method Payment            `json:"payment_method" bson:"payment_method"`
}

type Payment struct {
	Digital bool
	COD     bool
}
