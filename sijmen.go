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
	option_ship = true
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

var wishes = []Wish{
	{shipOrWeapon: option_weapon, weapon: TACHYON},
	{shipOrWeapon: option_weapon, weapon: PLASMA},
	{shipOrWeapon: option_weapon, weapon: LASER},
	{shipOrWeapon: option_weapon, weapon: PARTICLE},
	{shipOrWeapon: option_weapon, weapon: PHOTON},
	{shipOrWeapon: option_weapon, weapon: PROTON},
	{shipOrWeapon: option_ship, ship: &DRONE},
	{shipOrWeapon: option_ship, ship: &RHINO},
	{shipOrWeapon: option_ship, ship: &HUMPBACK},
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
		&wishes[5], &wishes[7], &wishes[2], &wishes[8],
	})
	jsonString, _ := json.Marshal(Output{
		Name:         "Sijmen Huizenga",
		Email:        "sijmenhuizenga@gmail.com",
		Transactions: transactions,
	})
	ioutil.WriteFile("output.json", jsonString, os.ModePerm)
	fmt.Printf("Balance: %v\n", balance)

}


type Wish struct {
	shipOrWeapon bool
	weapon       Weapon
	ship         *Ship
}

func freelancer(stars []*Star, wishlist []*Wish) (uint16, []Transaction) {
	var balance uint16
	var transactions []Transaction

	var currentShip = &SCRAPPY
	var myWeapons []Weapon
	var inventory = map[Resource]uint8{
		Ore:         1,
		Water:       1,
		EngineParts: 1,
		Contraband:  0,
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
			ShipPurchase:     "",
		}

		// sell everything
		for resource, amount := range inventory {
			transaction.sell(resource, amount)
			balance += uint16(amount) * uint16(star.getPrice(resource))
			inventory[resource] = 0
		}

		// on last star don't buy anything
		if i == len(stars)-1 {
			transactions = append(transactions, transaction)
			break
		}

		nextStar := stars[i+1]

		shoppingCost, _, shoppingList := star.bestDeal(nextStar, balance, *currentShip)

		if len(wishlist) > 0 {
			nextWishItem := wishlist[0]
			if nextWishItem.shipOrWeapon == option_ship {
				newShipShoppingcost, _, newShipShoppingList := star.bestDeal(nextStar, balance, *nextWishItem.ship)

				if balance >= nextWishItem.ship.Price + newShipShoppingcost {
					shoppingList = newShipShoppingList
					transaction.ShipPurchase = nextWishItem.ship.Name
					balance -= nextWishItem.ship.Price
					currentShip = nextWishItem.ship
					wishlist = wishlist[1:]
				}
			} else {
				// buying weapons and fighting!
				if balance >= shoppingCost+200 && star.hasContract(nextWeapon(nextWishItem.weapon)) {
					// buy weapon on shoppinglist if we can use it immediately
					myWeapons = append(myWeapons, nextWishItem.weapon)
					transaction.WeaponPurchase = []string{weaponNames[nextWishItem.weapon]}
					wishlist = wishlist[1:]
					balance -= 200
				}
			}
		}

		if len(myWeapons) > 0 {
			// if we have weapons, let's try to fight!
			bestcontract := star.bestContractW(nextWeapons(myWeapons))
			if bestcontract != nil {
				transaction.ContractAccepted = bestcontract.CriminalName
				balance += uint16(bestcontract.Bounty)
			}
		}


		// buy the shoppinglist
		for resource, amount := range shoppingList {
			transaction.buy(resource, amount)
			inventory[resource] += amount
			balance -= uint16(amount) * uint16(star.getPrice(resource))
		}

		// debugging
		//fmt.Printf("Balance at %v %v was: %v", i, star.Name, balance)
		//if transaction.ShipPurchase != "" {
		//	fmt.Printf("\n  Bought ship %v", transaction.ShipPurchase)
		//}
		//if len(transaction.WeaponPurchase) != 0 {
		//	fmt.Printf("\n  Bought weapon %v", transaction.WeaponPurchase)
		//}
		//if transaction.ContractAccepted != "" {
		//	fmt.Printf("\n  Contract %v", transaction.ContractAccepted)
		//}
		//println()

		transactions = append(transactions, transaction)
	}
	return balance, transactions
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
