package main

import "math/rand"

type ChefList struct {
	chefList      []*Chef
	chefIdCounter int
}

func NewChefList() *ChefList {
	ret := new(ChefList)
	ret.chefIdCounter = 0
	for i := 0; i < chefN; i++ {
		randomChef := chefPersonas[rand.Intn(len(chefPersonas))]
		randomChef.id = ret.chefIdCounter
		ret.chefIdCounter++
		if i == 0 {
			randomChef.rank = 3
		}
		ret.chefList = append(ret.chefList, NewChef(&randomChef))
	}
	return ret
}

func (cl ChefList) start() {
	for _, chef := range cl.chefList {
		go chef.startWorking()
	}
}

var chefPersonas = []Chef{{
	rank:        1,
	proficiency: 1,
	name:        "Jimmy Chef",
	catchPhrase: "Yes",
}, {
	rank:        2,
	proficiency: 2,
	name:        "Wolfgang Puck",
	catchPhrase: "Belissimo",
}, {
	rank:        1,
	proficiency: 3,
	name:        "Jamie Oliver",
	catchPhrase: "Impecable",
}, {
	rank:        3,
	proficiency: 2,
	name:        "James",
	catchPhrase: "This one is trash",
}, {
	rank:        3,
	proficiency: 3,
	name:        "Gordon Ramsay",
	catchPhrase: "WHERE IS THE LAMB SAUCE?",
}}
