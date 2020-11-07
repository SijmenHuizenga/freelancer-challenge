package main

import "log"

type Output struct {
	Name string `json:"name"`
	Email string `json:"email"`
	Transactions []Transaction `json:"transactions"`
}

type Transaction struct {
	Planet           string   `json:"planet"`
	DeltaOre         int8   `json:"deltaOre"`
	DeltaWater       int8   `json:"deltaWater"`
	DeltaEngineParts int8   `json:"deltaEngineParts"`
	DeltaContraband  int8   `json:"deltaContraband"`
	JumpTo           string   `json:"jumpTo,omitempty"`
	WeaponPurchase   []string `json:"weaponPurchase,omitempty"`
	ContractAccepted string   `json:"contractAccepted,omitempty"`
	ShipPurchase     string   `json:"shipPurchase,omitempty"`
}

func (t *Transaction) sell(resource Resource, amount uint8) {
	switch resource {
	case Water:
		t.DeltaWater -= int8(amount)
		return
	case Ore:
		t.DeltaOre -= int8(amount)
		return
	case EngineParts:
		t.DeltaEngineParts -= int8(amount)
		return
	case Contraband:
		t.DeltaContraband -= int8(amount)
		return
	}
	log.Fatal("The following resource does not exist: ", resource)
}

func (t *Transaction) buy(resource Resource, amount uint8) {
	switch resource {
	case Water:
		t.DeltaWater += int8(amount)
		return
	case Ore:
		t.DeltaOre += int8(amount)
		return
	case EngineParts:
		t.DeltaEngineParts += int8(amount)
		return
	case Contraband:
		t.DeltaContraband += int8(amount)
		return
	}
	log.Fatal("The following resource does not exist: ", resource)
}
