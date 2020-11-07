package main

import (
	"log"
)

type Star struct {
	Name             string     `json:"name"`
	OrePrice         uint8      `json:"orePrice"`
	WaterPrice       uint8      `json:"waterPrice"`
	EnginePartsPrice uint8      `json:"enginePartsPrice"`
	ContrabandPrice  uint8      `json:"contrabandPrice"`
	Faction          string     `json:"faction"`
	Contracts        []Contract `json:"contracts"`

	_aggregatedContracts map[Weapon]Bounty
}

type Contract struct {
	CriminalName string `json:"criminalName"`
	WeaponName   string `json:"weaponType"`
	WeaponType   Weapon
	Bounty       Bounty `json:"bounty"`
}

type Weapon uint8
type Bounty int

func (star *Star) getPrice(r Resource) uint8 {
	switch r {
	case Water:
		return star.WaterPrice
	case Ore:
		return star.OrePrice
	case EngineParts:
		return star.EnginePartsPrice
	case Contraband:
		return star.ContrabandPrice
	}
	log.Fatal("The following resource does not exist: ", r)
	return 0
}

type BuyableResource struct {
	resource Resource
	cost     uint8
	profit   uint8
}

func (star *Star) bestDeal(goingToStar *Star, budget uint16, ship Ship) (bestCost uint16, bestProfit uint16, bestInventory map[Resource]uint8) {
	capacity := ship.Capacity
	buyableResources := star.bestDealFindBuyableResources(goingToStar)
	cost, profit, inventory := knaphoor(0, 0, make([]uint8, len(buyableResources)), capacity, budget, buyableResources, 0)
	shoppingList := map[Resource]uint8 {}
	for i, buyableResource := range buyableResources {
		shoppingList[buyableResource.resource] = inventory[i]
	}
	return cost, profit, shoppingList
}

func (star *Star) contractsAggregated() map[Weapon]Bounty {
	if star._aggregatedContracts != nil {
		return star._aggregatedContracts
	}
	var out = map[Weapon]Bounty{}
	for _, contract := range star.Contracts {
		if val, ok := out[contract.WeaponType]; ok {
			out[contract.WeaponType] = val + contract.Bounty
		} else {
			out[contract.WeaponType] = contract.Bounty
		}
	}
	star._aggregatedContracts = out
	return out
}

func (star *Star) hasContract(by Weapon) bool {
	_, ok := star.contractsAggregated()[by]
	return ok
}

func (star *Star) bestContract(by Weapon) *Contract {
	var bestContract *Contract = nil
	for contractI := range star.Contracts{
		if star.Contracts[contractI].WeaponType != by {
			continue
		}
		if bestContract == nil || star.Contracts[contractI].Bounty > bestContract.Bounty {
			bestContract = &star.Contracts[contractI]
		}
	}
	return bestContract
}

func (star *Star) bestContractW(weapons []Weapon) *Contract {
	var bestContract *Contract = nil
	for _, contract := range star.Contracts{
		if !in(contract.WeaponType, weapons) {
			continue
		}
		if bestContract == nil || contract.Bounty > bestContract.Bounty {
			bestContract = &contract
		}
	}
	return bestContract
}

func in(weaponType Weapon, weapons []Weapon) bool {
	for _, w := range weapons {
		if w == weaponType {
			return true
		}
	}
	return false
}

func knaphoor(cost uint16, profit uint16, inventory []uint8, inventoryCapacity uint16, budget uint16,
	resources []BuyableResource, continuationI int) (bestCost uint16, bestProfit uint16, bestInventory []uint8) {

	if cost > budget {
		return 999, 0, inventory
	}

	if inventoryCapacity == 0 {
		return cost, profit, inventory
	}

	bestCost = cost
	bestProfit = profit
	bestInventory = inventory

	for i := continuationI; i < len(resources); i++ {

		newInv := make([]uint8, len(inventory))
		copy(newInv, inventory)
		newInv[i]++

		newCost, newProfit, newInventory := knaphoor(
			cost + uint16(resources[i].cost),
			profit + uint16(resources[i].profit),
			newInv,
			inventoryCapacity - 1,
			budget, resources, i)

		if newProfit > bestProfit {
			bestProfit = newProfit
			bestCost = newCost
			bestInventory = newInventory
		}
	}

	return bestCost, bestProfit, bestInventory
}


func (star *Star) bestDealFindBuyableResources(goingToStar *Star) []BuyableResource {
	var buyableThings []BuyableResource
	for _, resource := range RESOURCES {
		if resource == Contraband && star.Faction == "LIBERTY_POLICE" {
			continue
		}
		profit := int16(goingToStar.getPrice(resource)) - int16(star.getPrice(resource))
		if profit <= 0 {
			continue
		}
		newResource := BuyableResource{
			resource: resource,
			profit:   uint8(profit),
			cost:     star.getPrice(resource),
		}
		i := 0
		for i < len(buyableThings) {
			if buyableThings[i].profit < newResource.profit {
				break
			}
			i++
		}
		buyableThings = insert(buyableThings, i, newResource)
	}
	return buyableThings
}

func insert(a []BuyableResource, index int, value BuyableResource) []BuyableResource {
	if len(a) == index {
		return append(a, value)
	}
	a = append(a[:index+1], a[index:]...)
	a[index] = value
	return a
}
