import json

transactions = []

ships = [
    {"price": 0, "capacity": 10},  # this is the default ship
    {"price": 500, "capacity": 15, "name": "RHINO"},
    # {"price": 1000, "capacity": 20, "name": "DRONE"},
    {"price": 2000, "capacity": 30, "name": "HUMPBACK"},
]

shipI = 0

things = [
    "ore",
    "water",
    "engineparts",
    "contraband"
]

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
    print("Transaction", x)
    transactions.append(x)


with open('starmap.json') as json_file:
    starmap = json.load(json_file)
    i = 0
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
        if i == len(starmap) - 1:
            makeTransaction(currentStar, t)
            continue
            pass

        nextStar = starmap[i + 1]

        # if we still have money left to buy a ship, lets do it!
        if shipI < len(ships) - 1 and balance >= ships[shipI + 1]["price"] + 50:
            balance -= ships[shipI + 1]["price"]
            shipI += 1
            t["shipPurchase"] = ships[shipI]["name"]

        capacity = ships[shipI]["capacity"]

        # Calculate how much we will gain per piece of each thing
        winstMargins = {thing: nextStar[thing + "Price"] - currentStar[thing + "Price"]
                        for thing in things if currentStar[thing + "Price"] != 0}
        winstMarginsKeys = list(winstMargins.keys())
        winstMarginsVals = list(winstMargins.values())

        # best shoppinglist starts out with all 0
        bestShoppinglist = {thing: 0 for thing in things if currentStar[thing + "Price"] != 0}
        bestWinst = 0

        def magic(margin_counter, commulative_shoppinglist):

            if margin_counter == len(winstMargins):
                # shoppinglist complete!
                global balance
                global bestWinst
                global bestShoppinglist

                nr_of_items_bought = sum(commulative_shoppinglist.values())
                if nr_of_items_bought > capacity:
                    # too many things in our shopping cart
                    return

                cost = 0
                for key in winstMarginsKeys:
                    cost += currentStar[key + "Price"] * commulative_shoppinglist[key]
                if cost > balance:
                    # this setup is too expensive
                    return

                winst = 0
                for key in winstMarginsKeys:
                    winst += winstMargins[key] * commulative_shoppinglist[key]
                if winst < bestWinst:
                    return

                # this is the best shoppinglist (until now)
                bestShoppinglist = commulative_shoppinglist
                bestWinst = winst

                return
            for k in range(capacity+1):
                magic(margin_counter+1, {**commulative_shoppinglist, winstMarginsKeys[margin_counter]: k})


        magic(0, {thing: 0 for thing in things if currentStar[thing + "Price"] != 0})

        for thing in bestShoppinglist:
            t[thing] += bestShoppinglist[thing]
            inventory[thing] += bestShoppinglist[thing]
            balance -= bestShoppinglist[thing] * currentStar[thing+"Price"]

        # print("Going to buy:")
        # print(bestShoppinglist)

        # bestInvestment = max(winstMargins.items(), key=operator.itemgetter(1))[0]
        #
        # # See how much we can buy of that thing
        # howMuchBuy = int(min(balance / currentStar[bestInvestment + "Price"], ships[shipI]["capacity"]))
        # print(currentStar["name"] + "\t investing in " + bestInvestment + " bying " + str(howMuchBuy))
        #
        # # buy that thing
        # t[bestInvestment] = t[bestInvestment] + howMuchBuy
        # balance -= howMuchBuy * currentStar[bestInvestment + "Price"]
        # inventory[bestInvestment] += howMuchBuy

        makeTransaction(currentStar, t)
        i += 1

output = {
    "name": "Delftsche Zwervers (Sijmen Huizenga)",
    "email": "sijmenhuizenga@gmail.com",
    "transactions": transactions
}

print("Total balance: " + str(balance))

with open('output.json', 'w') as outfile:
    json.dump(output, outfile)

# bewaren tot de eerste piek
# koopt geen spullen op
# je koopt iets niet als het in de volgende goedkoper te kopen zijn

# REGELS
# koop niet als de volgende planeet als de volgende planeet het voor een lagere heeft.
#   goedkoper op andere planeet
#   alleen bij een hogere prijs

# verkopen op een piek
#   de prijs bij de volgende planeet gelijk dan verkopen. ivm tijdelijke bagageprijs

# Als de prijs daalt moet je altijd je spullen verkopen
# je moet altijd verkopen als de prijs daalt
# Nooit een goed kopen wat je dan verkoopt
# Als de volgende planeet gelijke prijs of lager heeft dan


# als je in een dal zit, dan moet je altijd kopen.
