package main

import (
	"os"
	"time"
)

var restaurantAddress = "http://localhost"

const restaurantPort = ":7500"
const kitchenPort = ":8000"

const chefN = 3
const ovenN = 3
const stoveN = 2
const orderListMaxSize = 3

const timeUnit = 100 * time.Millisecond

var kitchen Kitchen

var menu = []Menu{pizza, salad, zeama, sswmlc, idwmm, waffles, aubergine, lasagna, burger, gyros}
var apparatToId = map[string]int{"": 0, "oven": 1, "stove": 2}
var idToApparat = map[int]string{0: "nil", 1: "oven", 2: "stove"}

var pizza = Menu{1, "pizza", 20, 2, "oven"}
var salad = Menu{2, "salad", 10, 1, ""}
var zeama = Menu{3, "zeama", 7, 1, "stove"}
var sswmlc = Menu{4, "Scallop Sashimi with Meyer Lemon Confit", 32, 3, ""}
var idwmm = Menu{5, "Island Duck with Mulberry Mustard", 35, 3, "oven"}
var waffles = Menu{6, "Waffles", 10, 1, "stove"}
var aubergine = Menu{7, "Aubergine", 20, 2, ""}
var lasagna = Menu{8, "Lasagna", 30, 2, "oven"}
var burger = Menu{9, "Burger", 15, 1, "oven"}
var gyros = Menu{10, "Gyros", 15, 1, ""}

func main() {
	if args := os.Args; len(args) > 1 {
		restaurantAddress = args[1]
	}
	kitchen.start()
}
