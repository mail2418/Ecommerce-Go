package database

import "errors"

var (
	ErrCantFindProduct = errors.New("error cant find product")
	ErrCantDecodeProducts = errors.New("error cant find product")
	ErrUserIdIsNotValid = errors.New("error user id is not valid")
	ErrCantUpdateUser = errors.New("error cant update user")
	ErrCantRemoveItemCart = errors.New("error cant remove item from cart")
	ErrCantGetItem = errors.New("error cant get item from cart")
	ErrCantBuyCartItem = errors.New("error cant update the purchase")
)

func AddProductToCart(){

}

func RemoveCartItem(){

}

func BuyItemFromCart(){

}

func InstantBuyer(){
	
}