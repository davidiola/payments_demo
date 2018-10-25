package main

import (
	"context"
	"errors"
	"firebase.google.com/go"
	"firebase.google.com/go/db"
	"fmt"
	"log"
)

type Transaction struct {
	CardId            int
	TransactionAmount int
}

type Card struct {
	CardId  int
	Balance int
	Limit   int
}

type TestJson struct {
	field int
}

func checkTransaction(t Transaction, client *db.Client, ctx context.Context) bool {

	//to check a transaction we retrieve the card based upon its cardID, check if balance + transactionAmount <= Limit, if so
	//return success and subtract amount, otherwise return false and leave balance as is

	//retrieve card from database based upon cardID
	retrieveCardLocation := fmt.Sprintf("/%d", t.CardId)
	var currCard Card
	if err := client.NewRef(retrieveCardLocation).Get(ctx, &currCard); err != nil {
		log.Fatal(err)
	}

	makeTransaction := func(tn db.TransactionNode) (interface{}, error) {
		var currCard Card
		if err := tn.Unmarshal(&currCard); err != nil {
			return nil, err
		}
		if currCard.Balance+t.TransactionAmount <= currCard.Limit {
			currCard.Balance = currCard.Balance + t.TransactionAmount
			return currCard, nil
		} else {
			return currCard, errors.New("Transaction unsuccessful")
		}
	}

	if err := client.NewRef(retrieveCardLocation).Transaction(ctx, makeTransaction); err != nil {
		log.Fatal(err)
		return false
	}

	return true

}

func main() {

	ctx := context.Background()
	//modify this to reflect your own firebase database
	config := &firebase.Config{
		DatabaseURL: "https://yourdatabase.firebaseio.com/",
	}
	app, err := firebase.NewApp(ctx, config)
	if err != nil {
		fmt.Println(err)
		log.Fatal(err)
	}

	client, err := app.Database(ctx)
	if err != nil {
		fmt.Println(err)
		log.Fatal(err)
	}

	newCard := Card{
		CardId:  1234,
		Balance: 0,
		Limit:   1000,
	}

	t := Transaction{
		CardId:            1234,
		TransactionAmount: 500,
	}

	//store new card into database if it doesn't exist

	checkCard := Card{}
	locationInDatabase := fmt.Sprintf("/%d", newCard.CardId) //store card in cardID location of database
	if err := client.NewRef(locationInDatabase).Get(ctx, &checkCard); err != nil {
		log.Fatal(err)
	}
	if (Card{}) == checkCard { //card doesn't exist yet, create it in database
		if err := client.NewRef(locationInDatabase).Set(ctx, newCard); err != nil {
			fmt.Println(err)
			log.Fatal(err)
		}
	}

	success := checkTransaction(t, client, ctx)
	if !success {
		fmt.Println("Transaction unsuccessful")
	} else {
		fmt.Println("Successful transaction")
	}

}
