package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
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

const (
	option_ship   = true
	option_weapon = false
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

var (
    weapon_TACHYON = Wish{shipOrWeapon: option_weapon, weapon: TACHYON}
    weapon_PLASMA = Wish{shipOrWeapon: option_weapon, weapon: PLASMA}
    weapon_LASER = Wish{shipOrWeapon: option_weapon, weapon: LASER}
    weapon_PARTICLE = Wish{shipOrWeapon: option_weapon, weapon: PARTICLE}
    weapon_PHOTON = Wish{shipOrWeapon: option_weapon, weapon: PHOTON}
    weapon_PROTON = Wish{shipOrWeapon: option_weapon, weapon: PROTON}
	ship_RHINO = Wish{shipOrWeapon: option_ship, ship: &RHINO}
	ship_DRONE = Wish{shipOrWeapon: option_ship, ship: &DRONE}
    ship_HUMPBACK = Wish{shipOrWeapon: option_ship, ship: &HUMPBACK}
)

var wishes = []Wish{
	weapon_TACHYON, weapon_PLASMA, weapon_LASER, weapon_PARTICLE, weapon_PHOTON, weapon_PROTON, ship_RHINO, ship_HUMPBACK,
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

	balance, transactions := freelancer(starsP, []*Wish{
		&weapon_PLASMA, &weapon_TACHYON, &ship_RHINO, &ship_HUMPBACK,
	})
	jsonString, _ := json.Marshal(Output{
		Name:         "Sijmen Huizenga",
		Email:        "sijmenhuizenga@gmail.com",
		Transactions: transactions,
	})
	ioutil.WriteFile("output.json", jsonString, os.ModePerm)
	fmt.Printf("Balance: %v\n", balance)

	//var wg sync.WaitGroup
	//go func() {wishfull(starsP, []*Wish{&weapon_TACHYON}); wg.Done()}()
	//go func() {wishfull(starsP, []*Wish{&weapon_PLASMA}); wg.Done()}()
	//go func() {wishfull(starsP, []*Wish{&weapon_LASER}); wg.Done()}()
	//go func() {wishfull(starsP, []*Wish{&weapon_PARTICLE}); wg.Done()}()
	//go func() {wishfull(starsP, []*Wish{&weapon_PHOTON}); wg.Done()}()
	//go func() {wishfull(starsP, []*Wish{&weapon_PROTON}); wg.Done()}()
	//go func() {wishfull(starsP, []*Wish{&ship_RHINO}); wg.Done()}()
	//go func() {wishfull(starsP, []*Wish{&ship_HUMPBACK}); wg.Done()}()
	//
	//wg.Add(8)
	//wg.Wait()
}

// used for brute-forcing all wish combinations
func wishfull(stars []*Star, wishlist []*Wish) {
	for wishI := range wishes {
		if inW(wishlist, &wishes[wishI]) {
			continue
		}
		nextList := append(wishlist, &wishes[wishI])
		printWishes(nextList)
		balance, transactions := freelancer(stars, nextList)
		if balance > 25000 {
			jsonString, _ := json.Marshal(Output{
				Name:         "Sijmen Huizenga",
				Email:        "sijmenhuizenga@gmail.com",
				Transactions: transactions,
			})
			ioutil.WriteFile("output/"+strconv.Itoa(int(balance))+".json", jsonString, os.ModePerm)
			fmt.Printf("Found great option: %v\n", balance)
			printWishes(nextList)
		}
		if len(nextList) < 4 {
			wishfull(stars, nextList)
		}
	}
}

func printWishes(wishlist []*Wish) {
	for _, w := range wishlist {
		if w.shipOrWeapon == option_weapon {
			fmt.Printf(", %v", w.weapon)
		} else {
			fmt.Printf(", %v", w.ship.Name)
		}
	}
	println()
}

func inW(hay []*Wish, search *Wish) bool {
	for _, w := range hay {
		if w == search {
			return true
		}
	}
	return false
}

type Wish struct {
	shipOrWeapon bool
	weapon       Weapon
	ship         *Ship
}

func freelancer(stars []*Star, wishlist []*Wish) (uint16, []Transaction) {
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
		wishlist:  wishlist,
		myWeapons: []Weapon{},
		myShip:    &SCRAPPY,
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

		bestState, bestLink := findBestNextLink(s, stars, starI, nextLinks, nrOfVisitedStars, 3)

		s = *bestState
		nextLinks = bestLink.next
		starI = int8(int(starI) + int(bestLink.step))
		s.transaction.JumpTo = stars[starI].Name
		println("  next step: ", bestLink.step)
		println("  nr of visited stars: ", nrOfVisitedStars)
		println("  balance: ", s.balance)

		nrOfVisitedStars++
		transactions = append(transactions, s.transaction)
	}
	return s.balance, transactions
}

func findBestNextLink(s State, stars []*Star, starI int8, nextLinks *[]*Link, nrOfVisitedStars int, lookahead uint8) (*State, *Link) {
	//nrOfVisitedStars is including the current star
	var bestState *State = nil
	var bestLink *Link = nil
	var bestLookaheadBalance = uint16(0)
	for _, nextLink := range *nextLinks {
		nextStarI := starI + nextLink.step
		if nextStarI < 0 || nextStarI >= int8(len(stars)) {
			continue
		}
		// nrOfUnvisitedStars: current jump is not included
		nrOfUnvisitedStars := uint8(len(stars) - nrOfVisitedStars)
		if nextStarI == int8(len(stars)-1) && nrOfUnvisitedStars == 1 {
			// never visit the last star EXCEPT when we are at the second-to-last star
			continue
		}

		if nextLink.linksToRoot > nrOfUnvisitedStars-1 {
			continue
		}

		nextStar := stars[nextStarI]
		newState := visit(stars[starI], nextStar, State{
			transaction: s.transaction,
			balance:     s.balance,
			inventory:   CopyMap(s.inventory),
			wishlist:    s.wishlist,
			myWeapons:   s.myWeapons,
			myShip:      s.myShip,
		})

		if nrOfUnvisitedStars > 2 && lookahead > 0 {
			lookaheadBestState, _ := findBestNextLink(State{
				transaction: newState.transaction,
				balance:     newState.balance,
				inventory:   CopyMap(newState.inventory),
				wishlist:    newState.wishlist,
				myWeapons:   newState.myWeapons,
				myShip:      newState.myShip,
			}, stars, nextStarI, nextLink.next, nrOfVisitedStars+1, lookahead-1)
			if bestState == nil || lookaheadBestState.balance > bestLookaheadBalance {
				bestLink = nextLink
				bestState = &newState
				bestLookaheadBalance = lookaheadBestState.balance
			}
		} else {
			if bestState == nil || newState.balance > bestState.balance {
				bestLink = nextLink
				bestState = &newState
			}
		}
	}
	if bestState == nil || bestLink == nil {
		log.Fatal("No route found. Impossible!")
	}

	return bestState, bestLink
}

type State struct {
	transaction Transaction
	balance     uint16
	inventory   map[Resource]uint8
	wishlist    []*Wish
	myWeapons   []Weapon
	myShip      *Ship
}

func visit(currentStar *Star, nextStar *Star, s State) State {
	// return transaction, balance, wishlist, myWeapons, myShip
	shoppingCost, _, shoppingList := currentStar.bestDeal(nextStar, s.balance, *s.myShip)

	if len(s.wishlist) > 0 {
		nextWishItem := s.wishlist[0]
		if nextWishItem.shipOrWeapon == option_ship {
			newShipShoppingcost, _, newShipShoppingList := currentStar.bestDeal(nextStar, s.balance, *nextWishItem.ship)

			if s.balance >= nextWishItem.ship.Price+newShipShoppingcost {
				shoppingList = newShipShoppingList
				s.transaction.ShipPurchase = nextWishItem.ship.Name
				s.balance -= nextWishItem.ship.Price
				s.myShip = nextWishItem.ship
				s.wishlist = s.wishlist[1:]
			}
		} else {
			// buying weapons and fighting!
			if s.balance >= shoppingCost+200 && currentStar.hasContract(nextWeapon(nextWishItem.weapon)) {
				// buy weapon on shoppinglist if we can use it immediately
				s.myWeapons = append(s.myWeapons, nextWishItem.weapon)
				s.transaction.WeaponPurchase = []string{weaponNames[nextWishItem.weapon]}
				s.wishlist = s.wishlist[1:]
				s.balance -= 200
			}
		}
	}

	if len(s.myWeapons) > 0 {
		// if we have weapons, let's try to fight!
		bestcontract := currentStar.bestContractW(nextWeapons(s.myWeapons))
		if bestcontract != nil {
			s.transaction.ContractAccepted = bestcontract.CriminalName
			s.balance += uint16(bestcontract.Bounty)
		}
	}

	// buy the shoppinglist
	for resource, amount := range shoppingList {
		s.transaction.buy(resource, amount)
		s.inventory[resource] += amount
		s.balance -= uint16(amount) * uint16(currentStar.getPrice(resource))
	}

	return s
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
