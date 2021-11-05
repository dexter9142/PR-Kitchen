package main

import "time"

type Kitchen struct {
	kitchenWeb KitchenWeb
	orderList  *OrderList
	ovens      *ApplianceList
	stoves     *ApplianceList
	cookList   *ChefList
	connected  bool
}

func (k *Kitchen) start() {
	k.cookList = NewChefList()
	k.orderList = NewOrderList()
	k.ovens = newApparat(ovenN)
	k.stoves = newApparat(stoveN)

	go k.tryConnectRestaurant()
	k.kitchenWeb.start()
}

func (k *Kitchen) tryConnectRestaurant() {
	k.connected = false
	for !k.connected {
		if k.kitchenWeb.establishConnection() {
			k.connectionSuccessful()
			break
		} else {
			time.Sleep(timeUnit)
		}
	}
}
func (k *Kitchen) connectionSuccessful() {
	if k.connected {
		return
	}
	k.connected = true
	k.cookList.start()
}

func (k *Kitchen) getStatus() string {
	ret := "Cooks:"
	for _, cook := range k.cookList.chefList {
		ret += Div(cook.getStatus())
	}
	ret += "Ovens:"
	ret += k.ovens.getStatus()
	ret += "Stoves:"
	ret += k.stoves.getStatus()
	ret += "OrderList:"
	ret += k.orderList.getStatus()

	return ret
}
