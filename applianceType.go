package main

import (
	"math"
	"strconv"
	"sync"
)

type ApplianceList struct {
	numOfApparat int
	list         []*Appliance
	listMutex    sync.Mutex
}

func (al *ApplianceList) getApplianceAndWait(now int64) (*Appliance, int) {

	al.listMutex.Lock()
	appa := al.list[0].getAppVal()
	minWait := math.MaxInt32

	for _, loopAppa := range al.list {
		timeLeft := loopAppa.getTimeLeft(now)
		if timeLeft == 0 {
			minWait = 0
			appa = loopAppa.getAppVal()
			break
		}
		if minWait > timeLeft {
			minWait = timeLeft
			appa = loopAppa.getAppVal()
		}
	}
	al.listMutex.Unlock()

	return appa, minWait
}

func newApparat(numOfApparat int) *ApplianceList {
	ret := new(ApplianceList)
	ret.numOfApparat = numOfApparat
	for i := 0; i < numOfApparat; i++ {
		ret.list = append(ret.list, new(Appliance))
	}
	return ret
}

func (al *ApplianceList) getTimeLeft(now int64) int {
	minWait := math.MaxInt32
	for i, _ := range al.list {
		timeLeft := al.list[i].getTimeLeft(now)
		if timeLeft == 0 {
			return 0
		}
		if minWait > timeLeft {
			minWait = timeLeft
		}
	}
	return minWait
}

func (al *ApplianceList) getStatus() string {
	ret := ""
	for i, Apparat := range al.list {
		identification := "Id:" + strconv.Itoa(i)
		if Apparat.inUse == 1 {
			identification += " Used by chef id:"
			if Apparat.chef != nil {
				identification += strconv.Itoa(Apparat.chef.id)
			}
		} else {
			identification += " Free"
		}
		ret += Div(identification)
	}
	return ret
}
