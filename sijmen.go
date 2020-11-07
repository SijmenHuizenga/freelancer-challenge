package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

type Resource string

const (
	Ore         = "ore"
	Water       = "water"
	EngineParts = "engineparts"
	Contraband  = "contraband"
)
var RESOURCES = []Resource{Ore, Water, EngineParts, Contraband}

type Ship struct {
	Price    uint16
	Capacity uint16
	Name     string
}

var ships = []Ship{
	{Price: 0, Capacity: 10, Name: "SCRAPPY"},
	{Price: 500, Capacity: 15, Name: "RHINO"},
	//{Price: 1000, Capacity: 20, Name: "DRONE"},
	{Price: 2000, Capacity: 30, Name: "HUMPBACK"},
}

func main() {
	jsonFile, err := os.Open("starmap.json")
	if err != nil {
		fmt.Println(err)
	}
	defer jsonFile.Close()
	byteValue, _ := ioutil.ReadAll(jsonFile)

	var stars []Star
	json.Unmarshal(byteValue, &stars)
	transactions := freelancer(stars)
	jsonString, _ := json.Marshal(transactions)
	ioutil.WriteFile("output.json", jsonString, os.ModePerm)
}

func freelancer(stars []Star) []Transaction {
	var balance uint16
	var transactions []Transaction
	var currentShip = 0

	var inventory = map[Resource]uint8 {
		Ore: 1,
		Water: 1,
		EngineParts: 1,
		Contraband: 0,
	}

	for i, star := range stars {
		transaction := Transaction{
			Planet:           star.Name,
			DeltaOre:         0,
			DeltaWater:       0,
			DeltaEngineParts: 0,
			JumpTo:           "",
			WeaponPurchase:   []string{},
			ContractAccepted: "",
			ShipPurchase: 	  "",
		}

		// sell everything
		for resource, amount := range inventory {
			transaction.sell(resource, amount)
			balance += uint16(amount) * uint16(star.getPrice(resource))
			inventory[resource] = 0
		}

		// on last star don't buy anything
		if i == len(stars) - 1 {
			transactions = append(transactions, transaction)
			break
		}

		nextStar := stars[i+1]


		// buy a ship if we can
		if currentShip < len(ships) - 1 && balance >= ships[currentShip + 1].Price + 50 {
			currentShip ++
			balance -= ships[currentShip].Price
			transaction.ShipPurchase = ships[currentShip].Name
		}

		_, _, shoppingList := star.bestDeal(&nextStar, balance, ships[currentShip])

		// buy the shoppinglist
		for resource, amount := range shoppingList {
			transaction.buy(resource, amount)
			inventory[resource] += amount
			balance -= uint16(amount) * uint16(star.getPrice(resource))
		}

		transactions = append(transactions, transaction)
	}
	println("Total balance: ", balance)
	return transactions
}
