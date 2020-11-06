import json
from typing import List

ships = [
    {"price": 0, "capacity": 10},  # this is the default ship
    {"price": 500, "capacity": 15, "name": "RHINO"},
    # {"price": 1000, "capacity": 20, "name": "DRONE"},
    {"price": 2000, "capacity": 30, "name": "HUMPBACK"},
]

things = [
    "ore",
    "water",
    "engineparts",
    "contraband"
]

def knaphoor(cost: int, profit: int, inventory: List[int], inventoryCapacity: int,
             budget: int, prices: List[int], profits: List[int], continuationI: int = 0) -> (int, int, List[int],):

    if cost > budget:
        return 999999, 0, []

    if inventoryCapacity == 0:
        return cost, profit, inventory

    bestProfit = profit
    bestCost = cost
    bestInventory = inventory
    for i in range(continuationI, len(prices)):
        newCost, newProfit, newInventory = knaphoor(cost + prices[i],
                                                    profit + profits[i],
                                                    [(inventory[j] + 1 if j == i else inventory[j])
                                                     for j in range(len(inventory))],
                                                    inventoryCapacity - 1,
                                                    budget, prices, profits, i)
        if newProfit > bestProfit:
            bestProfit = newProfit
            bestCost = newCost
            bestInventory = newInventory
    return bestCost, bestProfit, bestInventory


def solve():
    transactions = []
    shipI = 0

    def makeTransaction(planet, transaction):
        x = {
            "planet": planet["name"],
            "deltaOre": transaction["ore"],
            "deltaWater": transaction["water"],
            "deltaEngineParts": transaction["engineparts"],
            "deltaContraband": transaction["contraband"],
        }
        if "shipPurchase" in transaction:
            x["shipPurchase"] = transaction["shipPurchase"]
        transactions.append(x)

    with open('starmap.json') as json_file:
        starmap = json.load(json_file)
        starI = 0
        inventory = {thing: 1 for thing in things}
        inventory["contraband"] = 0

        balance = 0

        for currentStar in starmap:
            # sell everything
            t = {thing: -1 * inventory[thing] for thing in things}
            for thing in things:
                balance += inventory[thing] * currentStar[thing + "Price"]
            inventory = {thing: 0 for thing in things}

            # check if on last star
            if starI == len(starmap) - 1:
                makeTransaction(currentStar, t)
                continue
                pass

            nextStar = starmap[starI + 1]

            # if we still have money left to buy a ship, lets do it!
            if shipI < len(ships) - 1 and balance >= ships[shipI + 1]["price"] + 50:
                balance -= ships[shipI + 1]["price"]
                shipI += 1
                t["shipPurchase"] = ships[shipI]["name"]

            capacity = ships[shipI]["capacity"]

            # do not buy things that won't make profit
            # do not buy contraband
            buyableThings = [thing for thing in things
                             if currentStar[thing + "Price"] != 0
                             and (nextStar[thing + "Price"] - currentStar[thing + "Price"]) > 0]
            buyableThings.sort(key=lambda thing: nextStar[thing + "Price"] - currentStar[thing + "Price"])

            prices = [currentStar[thing + "Price"] for thing in buyableThings]
            profits = [nextStar[thing + "Price"] - currentStar[thing + "Price"] for thing in buyableThings]

            # Find the best thing we can buy
            cost, profit, shoppingList = knaphoor(
                cost=0,
                profit=0,
                inventory=[0] * len(prices),
                inventoryCapacity=capacity,
                budget=balance,
                prices=prices,
                profits=profits)

            # buy it
            for i in range(len(buyableThings)):
                # add to transaction
                t[buyableThings[i]] += shoppingList[i]

                # remove from inventory
                inventory[buyableThings[i]] += shoppingList[i]

            # update balance
            balance -= cost

            makeTransaction(currentStar, t)
            starI += 1

    output = {
        "name": "Delftsche Zwervers (Sijmen Huizenga)",
        "email": "sijmenhuizenga@gmail.com",
        "transactions": transactions
    }

    print("Total balance: " + str(balance))

    with open('output.json', 'w') as outfile:
        json.dump(output, outfile)


solve()
