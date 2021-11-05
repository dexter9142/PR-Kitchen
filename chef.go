package main

import (
	"fmt"
	"math/rand"
	"strconv"
	"sync"
	"sync/atomic"
	"time"
)

var chefStatus = [...]string{" Waiting.", " Preparing", " Delivering"}

type Chef struct {
	id            int
	rank          int
	proficiency   int
	name          string
	catchPhrase   string
	isWorking     int32
	statusId      int
	orderId       int
	foodId        int
	applianceType int
	timeStarted   int64
	timeRequired  int
}

func NewChef(chef *Chef) *Chef {
	ret := new(Chef)

	ret.id = chef.id
	ret.rank = chef.rank
	ret.proficiency = chef.proficiency
	ret.name = chef.name
	ret.catchPhrase = chef.catchPhrase
	ret.isWorking = 0
	ret.statusId = 0
	ret.orderId = 0
	ret.foodId = 0
	ret.applianceType = 0
	ret.timeStarted = 0
	ret.timeRequired = 0

	return ret
}

func (c *Chef) startWorking() {
	c.isWorking = 1
	for c.isWorking == 1 {
		didATask := false
		var wg sync.WaitGroup
		stuffToDo := rand.Intn(c.proficiency) + 1
		for i := 0; i < stuffToDo; i++ {
			wg.Add(1)
			go func() {
				meal := kitchen.orderList.getMeal(c)
				if meal != nil {
					didATask = true
					now := getTime()
					c.statusId = 1
					c.orderId = meal.parent.id
					c.foodId = meal.foodId
					c.timeStarted = now
					c.timeRequired = meal.timeRequired
					switch meal.appliance {
					case 0:
						c.applianceType = 0
						meal.prepare(c, now)
					case 1:
						c.applianceType = 1
						appliance, waitAppliance := kitchen.ovens.getApplianceAndWait(now)
						c.timeRequired += waitAppliance
						appliance.use(c, meal, now)
					case 2:
						c.applianceType = 2
						appliance, waitAppliance := kitchen.stoves.getApplianceAndWait(now)
						c.timeRequired += waitAppliance
						appliance.use(c, meal, now)
					}
				}
				wg.Done()
			}()
		}
		wg.Wait()

		delivery := kitchen.orderList.getDelivery()
		if delivery != nil {
			success := false
			for success == false {
				didATask = true
				c.statusId = 2
				success = kitchen.kitchenWeb.deliver(delivery)
				if success == false {
					fmt.Println("OH NO")
				}
			}
		}
		if !didATask {
			c.statusId = 0
			time.Sleep(timeUnit)
		}
	}
}

func (c *Chef) stopWorking() {
	atomic.StoreInt32(&c.isWorking, 0)
}

func (c *Chef) getStatus() string {
	ret := "Chef " + c.name + " id:" + strconv.Itoa(c.id) + chefStatus[c.statusId] + " "
	if c.statusId != 0 {
		ret += menu[c.foodId].name + " for order id:" + strconv.Itoa(c.orderId)
		if c.applianceType != 0 {
			ret += " using " + idToApparat[c.applianceType]
		}
		ret += " time left:" + strconv.Itoa(c.timeRequired-int(getTime()-c.timeStarted))
	}

	return ret
}
