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

	// this is the lookahead-approach.
	// problem is that it causes results with unvisited stars
	lookahead := uint8(10)
	decisionPart := 5

	var visitedStars []uint8
	profit := uint16(0)
	var transactions []*Transaction

	for ; len(visitedStars)+1 < len(stars); {
		// prepare some variables to make visitStar happy
		lastVisited := uint8(0)
		var previousVisitedStars []uint8
		if len(visitedStars) > 0 {
			lastVisited = visitedStars[len(visitedStars)-1]
			previousVisitedStars = visitedStars[:len(visitedStars)-1]
		}

		additionalVisitedStars, _ := visitStar(lastVisited, stars, previousVisitedStars, profit, 0, inventory, uint8(len(visitedStars))+lookahead)
		if len(additionalVisitedStars) == len(visitedStars) {
			break
		}

		// Add the best {decisionPart} as part of our route
		visitedStars = append(visitedStars, additionalVisitedStars[len(visitedStars):len(visitedStars)+decisionPart]...)
		fmt.Printf("total: %v\n", visitedStars)

		profit, transactions, inventory = profitForRoute(visitedStars, stars)
		println(" Profit", profit)
	}

	for i, star := range stars {
		if ! in(uint8(i), visitedStars) {
			println("Not visited ", star.Name)
		}
	}

	fmt.Printf("path: %v\n", visitedStars)
	println("blanace: ", profit)
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
		return visitedStars[:len(visitedStars)-1], 0
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
	if nrOfStars == uint8(len(excludeIndexes)-1) {
		// only the last star is remaining. Return the last star
		return []uint8{nrOfStars-1}
	}

	var out []uint8
	i := uint8(max(int8(startingStar-3), 0))
	end := min(nrOfStars-1, startingStar+3)
	for ; i <= end; i++ {
		if in(i, excludeIndexes) {
			continue
		}
		starmap := starMapBools(nrOfStars, excludeIndexes)
		starmap[i] = true
		if !canReachEndVisitingAllStars(i, starmap, uint8(len(excludeIndexes)+1)) {
			continue
		}
		out = append(out, i)
	}
	return out
}

func starMapBools(nrOfStars uint8, excludeIndexes []uint8) []bool {
	starMap := make([]bool, nrOfStars)
	for _, excludedI := range excludeIndexes {
		starMap[excludedI] = true
	}
	return starMap
}

func canReachEndVisitingAllStars(starI uint8, starmap []bool, nrOfVisistedStars uint8) bool {

	nrOfStars := uint8(len(starmap))

	if nrOfVisistedStars+1 == nrOfStars {
		// starI is the second to last star
		return nrOfStars - starI - 1 <= 3
	}

	i := uint8(max(int8(starI-3), 0))
	end := min(nrOfStars-1, starI+3)
	for ; i <= end; i++ {
		if starmap[i] {
			// never travel into already visited stars
			continue
		}
		starmap[i] = true
		x := canReachEndVisitingAllStars(i, starmap, nrOfVisistedStars+1)
		starmap[i] = false
		if x {
			return true
		}
	}
	fmt.Printf("%v  %v\n", starI, starmap)
	return false
}

//func hasInvalidNeighbour(starI uint8, nrOfStars uint8, excludeIndexes []uint8) bool {
//	// checks if one of starI neighbours is invalid
//
//	i := uint8(max(int8(starI-3), 0))
//	end := min(nrOfStars-1, starI+3)
//	for ; i <= end; i++ {
//		if in(i, excludeIndexes) {
//			// ignore any neighbours that are in the excluded list
//			continue
//		}
//		if !hasNeighbour(i, nrOfStars, append(excludeIndexes, i)) {
//			// a neighbour without any other neighbours is invalid
//			return true
//		}
//		if !canReachFinalStar(i, nrOfStars, append(excludeIndexes, i)) {
//			// a neighbour from where we cannot reach the final star is invalid
//			return true
//		}
//	}
//	return false
//}
//
//func canReachFinalStar(starI uint8, nrOfStars uint8, excludeIndexes []uint8) bool {
//	// go through the sorted excludeIndexes, if the delta between two values (or a combination of consequtive values) is larger than 2 than we can't reach the end.
//	gap := 1
//	for i := starI+1; i < nrOfStars; i++ {
//		if gap >= 3 {
//			return false
//		}
//		if in(i, excludeIndexes) {
//			gap++
//			continue
//		} else {
//			gap = 0
//			excludeIndexes = append(excludeIndexes, i)
//		}
//	}
//	return true
//}
//
//func hasNeighbour(starI uint8, nrOfStars uint8, excludeIndexes []uint8) bool {
//	i := uint8(max(int8(starI-3), 0))
//	end := min(nrOfStars-1, starI+3)
//	for ; i <= end; i++ {
//		if !in(i, excludeIndexes) {
//			return true
//		}
//	}
//	return false
//}


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
