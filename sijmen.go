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

	// this is the lookahead-approach.
	// problem is that it causes results with unvisited stars
	//lookahead := uint8(10)
	//decisionPart := 5
	//
	//var visitedStars []uint8
	//lastVisited := uint8(0)
	//profit := uint16(0)
	//var transactions []*Transaction
	//
	//for ; len(visitedStars)+1 < len(stars); {
	//	// visit the next {lookahead} stars
	//	additionalVisitedStars, _ := visitStar(lastVisited, stars, visitedStars, profit, 0, inventory, uint8(len(visitedStars))+lookahead)
	//	if len(additionalVisitedStars) == len(visitedStars)+1 {
	//		// no new stars visited. We are finished
	//		break
	//	}
	//
	//	fmt.Printf(    "additional: %v \n", additionalVisitedStars)
	//	// Add the best {decisionPart} as part of our route
	//	visitedStars = append(visitedStars, additionalVisitedStars[len(visitedStars):len(visitedStars)+decisionPart]...)
	//	lastVisited = visitedStars[len(visitedStars)-1]
	//	fmt.Printf("total: %v\n", visitedStars)
	//
	//	profit, transactions, inventory = profitForRoute(visitedStars, stars)
	//	println(" Profit", profit)
	//
	//	// remove last to account for the fact that visitStar() will start at the last star
	//	visitedStars = visitedStars[:len(visitedStars)-1]
	//
	//	fmt.Printf(" Inventory %v \n", inventory)
	//}

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

	//fmt.Printf("%v  %v\n", bestVisitedStars, bestBalance)
	return bestVisitedStars, bestBalance
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
