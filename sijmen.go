package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	. "math"
	"os"
	"time"
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
	Capacity uint8
	Name     string
}

var ships = []*Ship{
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

	var stars []*Star
	json.Unmarshal(byteValue, &stars)
	transactions := freelancer(stars)
	jsonString, _ := json.Marshal(Output {
		Name: "Delftsche Zwervers (Sijmen Huizenga)",
		Email: "sijmenhuizenga@gmail.com",
		Transactions: transactions,
	})
	ioutil.WriteFile("output.json", jsonString, os.ModePerm)
}

func freelancer(stars []*Star) []*Transaction {
	var inventory = map[Resource]uint8{
		Ore:         1,
		Water:       1,
		EngineParts: 1,
		Contraband:  0,
	}

	visistedStars, _ := visitStar(0, stars, []uint8{}, 0, 0, inventory, 150)
	balance, transactions, _ := profitForRoute(visistedStars, stars)

	println("blanace: ", balance)
	return transactions
}

func visitStar(currentStar uint8, stars []*Star, visitedStars []uint8, balance uint16, ship uint8, inventory map[Resource]uint8, maxNrVisitedStars uint8) (newVisitedStars []uint8, newProfit uint16) {

	// sell everything
	for resource, amount := range inventory {
		balance += uint16(amount) * uint16(stars[currentStar].getPrice(resource))
		inventory[resource] = 0
	}

	visitedStars = append(visitedStars, currentStar)

	// on last star don't buy anything
	if len(visitedStars) == len(stars) || uint8(len(visitedStars)) == maxNrVisitedStars {
		solutionsPerSecond()
		return visitedStars, balance
	}

	neighbours := findJumpableStars(currentStar, uint8(len(stars)), visitedStars)

	if len(neighbours) == 0 {
		// nowhere to go but still stars to visit. This is an invalid solution.
		return visitedStars, 0
	}

	bestBalance := balance
	bestVisitedStars := visitedStars

	for _, neighbourI := range neighbours {
		cost, profit, shoppingList := stars[currentStar].bestDeal(stars[neighbourI], balance, ships[ship].Capacity)

		inventoryAtNextStar := make(map[Resource]uint8)
		for k,v := range inventory {
			inventoryAtNextStar[k] = v
		}

		for resource, amount := range shoppingList {
			inventory[resource] += amount
		}

		newVisitedStars, newBalance := visitStar(neighbourI, stars, visitedStars, balance+profit-cost, ship, inventoryAtNextStar, maxNrVisitedStars)
		if newBalance > bestBalance {
			bestBalance = newBalance
			bestVisitedStars = newVisitedStars
		}
	}

	return bestVisitedStars, bestBalance
}

var counter = float64(0)
var lastSecond = time.Now().UnixNano()

func solutionsPerSecond() bool {
	counter++
	if counter > 5000 {
		deltaS := float64(time.Now().UnixNano() - lastSecond)
		fmt.Printf("%v solutions/second\n", Round((counter / deltaS) * 1000000000))
		lastSecond = time.Now().UnixNano()
		counter = 0
		return true
	}
	return false
}

func findJumpableStars(startingStar uint8, nrOfStars uint8, excludeIndexes []uint8) []uint8 {
	var out []uint8
	i := uint8(max(int8(startingStar-3), 0))
	end := min(nrOfStars-1, startingStar+3)
	for ; i <= end; i++ {
		if in(i, excludeIndexes) {
			continue
		}
		out = append(out, i)
	}
	return out
}

func max(a int8, b int8) int8{
	if a > b {
		return a
	}
	return b
}
func min(a uint8, b uint8) uint8{
	if a < b {
		return a
	}
	return b
}
func min16(a uint16, b uint16) uint16{
	if a < b {
		return a
	}
	return b
}

func in(a uint8, list []uint8) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

func profitForRoute(route []uint8, stars []*Star) (uint16, []*Transaction, map[Resource]uint8){
	var balance uint16
	var transactions []*Transaction
	var currentShip = 0

	var inventory = map[Resource]uint8 {
		Ore: 1,
		Water: 1,
		EngineParts: 1,
		Contraband: 0,
	}

	for routeI, starI := range route {
		currentStar := stars[starI]
		transaction := Transaction{
			Planet:           stars[starI].Name,
			DeltaOre:         0,
			DeltaWater:       0,
			DeltaEngineParts: 0,
			JumpTo:           "",
			WeaponPurchase:   []string{},
			ContractAccepted: "",
			ShipPurchase:     "",
		}


		// sell everything
		for resource, amount := range inventory {
			transaction.sell(resource, amount)
			balance += uint16(amount) * uint16(currentStar.getPrice(resource))
			inventory[resource] = 0
		}

		// on last star don't buy anything
		if routeI == len(route) - 1 {
			transactions = append(transactions, &transaction)
			break
		}

		nextStar := stars[route[routeI+1]]
		transaction.JumpTo = nextStar.Name


		// buy a ship if we can
		//if currentShip < len(ships) - 1 && balance >= ships[currentShip + 1].Price + 50 {
		//	currentShip ++
		//	balance -= ships[currentShip].Price
		//	transaction.ShipPurchase = ships[currentShip].Name
		//}

		_, _, shoppingList := currentStar.bestDeal(nextStar, balance, ships[currentShip].Capacity)

		// buy the shoppinglist
		for resource, amount := range shoppingList {
			transaction.buy(resource, amount)
			inventory[resource] += amount
			balance -= uint16(amount) * uint16(currentStar.getPrice(resource))
		}

		transactions = append(transactions, &transaction)
	}
	return balance, transactions, inventory
}
