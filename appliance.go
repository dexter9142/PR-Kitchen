package main

import "sync"

type Appliance struct {
	inUse             int32
	meal              *Meal
	chef              *Chef
	queueWait         int
	prepareMutex      sync.Mutex
	MutexAvailability sync.Mutex
}

func (app *Appliance) getTimeLeft(now int64) int {
	app.MutexAvailability.Lock()
	defer app.MutexAvailability.Unlock()
	if app.inUse == 0 {
		return 0
	}
	return app.meal.getTimeLeft(now) + app.queueWait
}

func (app *Appliance) setAppVal(inUse int32, chef *Chef, meal *Meal) {
	app.MutexAvailability.Lock()
	app.inUse = inUse
	app.meal = meal
	app.chef = chef //setter
	app.MutexAvailability.Unlock()
}

func (app *Appliance) addTimeQueue(value int) {
	app.MutexAvailability.Lock()
	app.queueWait += value
	app.MutexAvailability.Unlock()
}

func (app *Appliance) getAppVal() *Appliance {
	app.MutexAvailability.Lock()
	defer app.MutexAvailability.Unlock() //getter
	return app
}

func (app *Appliance) use(chef *Chef, meal *Meal, now int64) {
	timeForCurrentMeal := meal.getTimeLeft(now)
	app.addTimeQueue(timeForCurrentMeal)
	app.prepareMutex.Lock()
	app.addTimeQueue(-timeForCurrentMeal)
	app.setAppVal(1, chef, meal)

	meal.prepare(chef, now)

	app.setAppVal(0, nil, nil)
	app.prepareMutex.Unlock()
}
