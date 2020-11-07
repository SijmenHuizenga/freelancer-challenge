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

var ships = []Ship{
	{Price: 0, Capacity: 10, Name: "SCRAPPY"},
	//{Price: 500, Capacity: 15, Name: "RHINO"},
	{Price: 1000, Capacity: 20, Name: "DRONE"},
	{Price: 2000, Capacity: 30, Name: "HUMPBACK"},
}

var weaponNames = map[Weapon]string {
	0: "TACHYON",
	1: "PLASMA",
	2: "LASER",
	3: "PARTICLE",
	4: "PHOTON",
	5: "PROTON",
}


const (
	TACHYON = 0
	PLASMA = 1
	LASER = 2
	PARTICLE = 3
	PHOTON = 4
	PROTON = 5
    totalWeaponCount = 6
)

func weaponType(name string) Weapon {
	switch name {
		case "TACHYON": return TACHYON
		case "PLASMA": return PLASMA
		case "LASER": return LASER
		case "PARTICLE": return PARTICLE
		case "PHOTON": return PHOTON
		case "PROTON": return PROTON
		case "totalWeaponCount": return totalWeaponCount
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

	transactions := freelancer(stars)
	jsonString, _ := json.Marshal(Output {
		Name: "Sijmen Huizenga",
		Email: "sijmenhuizenga@gmail.com",
		Transactions: transactions,
	})
	ioutil.WriteFile("output.json", jsonString, os.ModePerm)
}

func freelancer(stars []Star) []Transaction {
	var balance uint16
	var transactions []Transaction
	var currentShip = 0

	var weaponShoppinglist = []Weapon {
		PARTICLE, PROTON, // of omdraaien
		LASER,
		PLASMA,
		PHOTON, TACHYON, // of gewoon niet
	}
	var myWeapons []Weapon

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

		// buying weapons and fighting!
		var acceptedContract *Contract = nil
		if len(weaponShoppinglist) > 0 && balance > 200 && star.hasContract(nextWeapon(weaponShoppinglist[0])){
			buying := weaponShoppinglist[0]
			// buy weapon on shoppinglist if we can use it immediately
			myWeapons = append(myWeapons, weaponShoppinglist[0])
			weaponShoppinglist = weaponShoppinglist[1:]
			acceptedContract = star.bestContract(nextWeapon(buying))
			transaction.WeaponPurchase = []string{weaponNames[buying]}
			balance -= 200
		} else if len(myWeapons) > 0 {
			// if I have weapons, let's try to fight!
			acceptedContract = star.bestContractW(nextWeapons(myWeapons))
		}
		if acceptedContract != nil {
			transaction.ContractAccepted = acceptedContract.CriminalName
			balance += uint16(acceptedContract.Bounty)
		}

		// debugging
		fmt.Printf("Balance at %v %v was: %v", i, star.Name, balance)
		if transaction.ShipPurchase != "" {
			fmt.Printf("\n  Bought ship %v", transaction.ShipPurchase)
		}
		if len(transaction.WeaponPurchase) != 0 {
			fmt.Printf("\n  Bought weapon %v", transaction.WeaponPurchase)
		}
		println()

		transactions = append(transactions, transaction)
	}
	println("Total balance: ", balance)
	return transactions
}


func nextWeapon(w Weapon) Weapon{
	if w == totalWeaponCount-1 {
		return 0
	}
	return w+1
}
func nextWeapons(w []Weapon) []Weapon{
	out := make([]Weapon, len(w))
	for _, w := range w {
		out = append(out, nextWeapon(w))
	}
	return out
}