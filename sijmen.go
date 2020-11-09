package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
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

var (
	SCRAPPY  = Ship{Price: 0, Capacity: 10, Name: "SCRAPPY"}
	RHINO    = Ship{Price: 500, Capacity: 15, Name: "RHINO"}
	DRONE    = Ship{Price: 1000, Capacity: 20, Name: "DRONE"}
	HUMPBACK = Ship{Price: 2000, Capacity: 30, Name: "HUMPBACK"}
)

var SHIPS = []*Ship{&SCRAPPY, &HUMPBACK}

var weaponNames = map[Weapon]string{
	0: "TACHYON",
	1: "PLASMA",
	2: "LASER",
	3: "PARTICLE",
	4: "PHOTON",
	5: "PROTON",
}

const (
	TACHYON          = Weapon(0)
	PLASMA           = Weapon(1)
	LASER            = Weapon(2)
	PARTICLE         = Weapon(3)
	PHOTON           = Weapon(4)
	PROTON           = Weapon(5)
	totalWeaponCount = Weapon(6)
)

func weaponType(name string) Weapon {
	switch name {
	case "TACHYON":
		return TACHYON
	case "PLASMA":
		return PLASMA
	case "LASER":
		return LASER
	case "PARTICLE":
		return PARTICLE
	case "PHOTON":
		return PHOTON
	case "PROTON":
		return PROTON
	case "totalWeaponCount":
		return totalWeaponCount
	}
	log.Fatal("Weapon doesnt exist")
	return 0
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

	for starI := range stars {
		for contractI := range stars[starI].Contracts {
			stars[starI].Contracts[contractI].WeaponType = weaponType(stars[starI].Contracts[contractI].WeaponName)
		}
	}

	var starsP []*Star
	for i := range stars {
		starsP = append(starsP, &stars[i])
	}

	balance, transactions := freelancer(starsP)
	jsonString, _ := json.Marshal(Output{
		Name:         "Sijmen Huizenga",
		Email:        "sijmenhuizenga@gmail.com",
		Transactions: transactions,
	})
	ioutil.WriteFile("output.json", jsonString, os.ModePerm)
	fmt.Printf("Balance: %v\n", balance)

}

func freelancer(stars []*Star) (uint16, []Transaction) {
	var transactions []Transaction

	var s = State{
		transaction: Transaction{},
		balance:     0,
		inventory: map[Resource]uint8{
			Ore:         1,
			Water:       1,
			EngineParts: 1,
			Contraband:  0,
		},
		myWeapons: []Weapon{},
		myShip:    0,
		worth: 21,
	}
	var nextLinks = STARMAP
	var nrOfVisitedStars = 0
	var starI = int8(0)

	for true {
		s.transaction = Transaction{
			Planet:           stars[starI].Name,
			DeltaOre:         0,
			DeltaWater:       0,
			DeltaEngineParts: 0,
			JumpTo:           "",
			WeaponPurchase:   []string{},
			ContractAccepted: "",
			ShipPurchase:     "",
		}
		println("Visiting", stars[starI].Name, starI)

		// sell everything
		for resource, amount := range s.inventory {
			s.transaction.sell(resource, amount)
			s.balance += uint16(amount) * uint16(stars[starI].getPrice(resource))
			s.inventory[resource] = 0
		}

		// on last star don't buy anything
		if int(starI) == len(stars)-1 {
			transactions = append(transactions, s.transaction)
			break
		}

		if nextLinks == nil || len(*nextLinks) == 0 {
			log.Fatal("No next links. That can't be. You were supposed to be infinite")
		}

		bestState, bestLink, _ := findBestNextLink(s, stars, starI, nextLinks, nrOfVisitedStars, 4)

		s = *bestState
		nextLinks = bestLink.next
		starI = int8(int(starI) + int(bestLink.step))
		s.transaction.JumpTo = stars[starI].Name
		println("  next step: ", bestLink.step)
		println("  nr of visited stars: ", nrOfVisitedStars)
		println("  balance: ", s.balance)
		if s.transaction.ShipPurchase != "" {
			println( "  bought ship: ", s.transaction.ShipPurchase)
		}
		if len(s.transaction.WeaponPurchase) > 0{
			println( "  bought weapon: ", s.transaction.WeaponPurchase[0])
		}

		nrOfVisitedStars++
		transactions = append(transactions, s.transaction)
		//if nrOfVisitedStars > 4 {
		//	break
		//}
	}
	return s.balance, transactions
}

func findBestNextLink(s State, stars []*Star, starI int8, nextLinks *[]*Link, nrOfVisitedStars int, lookahead uint8) (*State, *Link, uint16) {
	//nrOfVisitedStars is including the current star
	var bestState *State = nil
	var bestLink *Link = nil
	var bestLookaheadBalance = uint16(0)

	// nrOfUnvisitedStars: current jump is not included
	nrOfUnvisitedStars := uint8(len(stars) - nrOfVisitedStars)

	shouldLookahead := nrOfUnvisitedStars > 2 && lookahead > 0

	for nextLinkI, nextLink := range *nextLinks {
		nextStarI := starI + nextLink.step
		if nextStarI < 0 || nextStarI >= int8(len(stars)) {
			continue
		}
		if nextStarI == int8(len(stars)-1) && nrOfUnvisitedStars == 1 {
			// never visit the last star EXCEPT when we are at the second-to-last star
			continue
		}

		if nextLink.linksToRoot > nrOfUnvisitedStars-1 {
			continue
		}

		nextStar := stars[nextStarI]
		newState := visit(stars[starI], nextStar, State{
			transaction: CopyTransaction(&s.transaction),
			balance:     s.balance,
			inventory:   CopyMap(s.inventory),
			myWeapons:   s.myWeapons,
			myShip:      s.myShip,
			worth:       s.worth,
		})

		if shouldLookahead {
			_, _, lookaheadBestLookaheadBalance := findBestNextLink(State{
				transaction: CopyTransaction(&newState.transaction),
				balance:     newState.balance,
				inventory:   CopyMap(newState.inventory),
				myWeapons:   newState.myWeapons,
				myShip:      newState.myShip,
				worth:       newState.worth,
			}, stars, nextStarI, nextLink.next, nrOfVisitedStars+1, lookahead-1)
			if bestState == nil || lookaheadBestLookaheadBalance > bestLookaheadBalance {
				bestLink = (*nextLinks)[nextLinkI]
				bestState = &newState
				bestLookaheadBalance = lookaheadBestLookaheadBalance
			}
		} else {
			if bestState == nil || newState.balance > bestState.balance {
				bestLink = (*nextLinks)[nextLinkI]
				bestState = &newState
			}
		}
	}
	if bestState == nil || bestLink == nil {
		log.Fatal("No route found. Impossible!")
	}

	if shouldLookahead {
		return bestState, bestLink, bestLookaheadBalance
	} else {
		return bestState, bestLink, bestState.balance
	}
}

type State struct {
	transaction Transaction
	balance     uint16
	worth       uint16
	inventory   map[Resource]uint8
	myWeapons   []Weapon
	myShip      uint8
}

func visit(currentStar *Star, nextStar *Star, s State) State {
	// just sell everything we have
	for resource, amount := range s.inventory {
		s.transaction.sell(resource, amount)
		s.balance += uint16(amount) * uint16(currentStar.getPrice(resource))
		s.worth -= uint16(amount) * uint16(currentStar.getPrice(resource))
		s.inventory[resource] = 0
	}
	shoppingCost, _, shoppingList := currentStar.bestDeal(nextStar, s.balance, *SHIPS[s.myShip])

	// still ships to buy
	if s.myShip < uint8(len(SHIPS)) -1 {
		newShipShoppingcost, _, newShipShoppingList := currentStar.bestDeal(nextStar, s.balance, *SHIPS[s.myShip+1])
		if s.balance >= SHIPS[s.myShip+1].Price+newShipShoppingcost {
			shoppingList = newShipShoppingList
			s.transaction.ShipPurchase = SHIPS[s.myShip+1].Name
			s.balance -= SHIPS[s.myShip+1].Price
			s.myShip++
		}
	}

	// contracts are orderd from best to worst reward
	weaponsIcanBeat := nextWeapons(s.myWeapons)
	for _, contract := range currentStar.Contracts {
		// if we have the weapon to beat this bitch, fight it!
		if in(contract.WeaponType, weaponsIcanBeat) {
			s.transaction.ContractAccepted = contract.CriminalName
			s.balance += uint16(contract.Bounty)
			s.worth += uint16(contract.Bounty)
			break
		}

		// If we can buy the weapon to beat it, buy it and kill the sucker
		if s.balance >= shoppingCost+200 {
			toBuy := previousWeapon(contract.WeaponType)
			// exclude some of the not-so-powerfull weapons
			if toBuy == PARTICLE || toBuy == TACHYON || toBuy == PHOTON || toBuy == PROTON {
				continue
			}
			// buy weapon on shoppinglist if we can use it immediately
			s.myWeapons = append(s.myWeapons, toBuy)
			s.transaction.WeaponPurchase = []string{weaponNames[toBuy]}
			s.transaction.ContractAccepted = contract.CriminalName
			s.balance -= 200
			s.balance += uint16(contract.Bounty)
			s.worth += uint16(contract.Bounty)
			break
		}
	}



	// buy the shoppinglist
	for resource, amount := range shoppingList {
		s.transaction.buy(resource, amount)
		s.inventory[resource] += amount
		s.balance -= uint16(amount) * uint16(currentStar.getPrice(resource))
		s.worth += uint16(amount) * uint16(currentStar.getPrice(resource))
	}


	return s
}



func previousWeapon(w Weapon) Weapon {
	if w == 0 {
		return totalWeaponCount-1
	}
	return w - 1
}
func nextWeapon(w Weapon) Weapon {
	if w == totalWeaponCount-1 {
		return 0
	}
	return w + 1
}
func nextWeapons(w []Weapon) []Weapon {
	var out []Weapon
	for _, w := range w {
		out = append(out, nextWeapon(w))
	}
	return out
}
