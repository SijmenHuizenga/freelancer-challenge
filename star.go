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
}

type Contract struct {
	CriminalName string `json:"criminalName"`
	WeaponType   string `json:"weaponType"`
	Bounty       string `json:"bounty"`
}

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

func (star *Star) bestDeal(goingToStar *Star, budget uint16, capacity uint8) (uint16, uint16, map[Resource]uint8) {
	// the naive solution
	bestResource := Resource("")
	bestProfit := uint16(0)
	bestCost := uint16(0)
	for _, resource := range RESOURCES {
		if resource == Contraband && star.Faction == "LIBERTY_POLICE" {
			continue
		}
		profit := int16(goingToStar.getPrice(resource)) - int16(star.getPrice(resource))
		if profit <= 0 {
			continue
		}

		if uint16(profit) > bestProfit {
			bestProfit = uint16(profit)
			bestResource = resource
			bestCost = uint16(star.getPrice(resource))
		}
	}
	if bestProfit == 0 {
		return 999, 0, nil
	}
	shoppingList := map[Resource]uint8 {}
	shoppingList[bestResource] = uint8(min16(uint16(capacity), budget / bestCost))
	return uint16(shoppingList[bestResource]) * bestCost, bestProfit, shoppingList

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
