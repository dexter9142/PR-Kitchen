package main

import (
	"math"
	"strconv"
	"sync"
	"time"
)

type Menu struct {
	id              int
	name            string
	preparationTime int
	complexity      int
	chefingApparat  string
}

type Meal struct {
	prepared          int32
	inUse             int32
	timeRequired      int
	complexity        int
	appliance         int
	preparingTime     int64
	foodId            int
	cookId            int
	parent            *Order
	MutexAvailability sync.Mutex
}

type Delivery struct {
	OrderId        int            `json:"order_id"`
	TableId        int            `json:"table_id"`
	Items          []int          `json:"items"`
	Priority       int            `json:"priority"`
	MaxWait        int            `json:"max_wait"`
	PickUpTime     int64          `json:"pick_up_time"`
	CookingTime    int            `json:"cooking_time"`
	CookingDetails []MealDelivery `json:"cooking_details"`
}

type MealDelivery struct {
	FoodId int `json:"food_id"`
	CookId int `json:"cook_id"`
}

type Order struct {
	id          int
	tableId     int
	items       []int
	mealCounter int32
	priority    int
	pickUpTime  int64
	maxWait     int
	mealList    []*Meal
}

type OrderList struct {
	deliveryMutex sync.Mutex
	mealMutex     sync.Mutex
	ovenList      []*Meal
	stoveList     []*Meal
	nilList       []*Meal
	orderArr      []*Order
	maxLen        int
}

type OrderProp struct {
	Id         int   `json:"order_id"`
	TableId    int   `json:"table_id"`
	WaiterId   int   `json:"waiter_id"`
	Items      []int `json:"items"`
	Priority   int   `json:"priority"`
	MaxWait    int   `json:"max_wait"`
	PickUpTime int64 `json:"pick_up_time"`
}

func newDelivery(order *Order) *Delivery {
	ret := new(Delivery)
	ret.OrderId = order.id
	ret.TableId = order.tableId
	ret.Items = order.items
	ret.Priority = order.priority
	ret.MaxWait = order.maxWait
	ret.PickUpTime = order.pickUpTime
	ret.CookingTime = int(getTime() - order.pickUpTime)
	var cookingDetails []MealDelivery
	for _, meal := range order.mealList {
		cookingDetails = append(cookingDetails, MealDelivery{meal.foodId, meal.cookId})
	}
	ret.CookingDetails = cookingDetails
	return ret
}

func (m *Meal) getTimeLeft(now int64) int {

	if m.inUse == 1 {
		elapsed := int(now - m.preparingTime)
		return m.timeRequired - elapsed
	}
	elapsed := int(now - m.parent.pickUpTime)
	limit := m.parent.maxWait
	priority := m.parent.priority
	return limit - elapsed - m.timeRequired - priority
}

func (m *Meal) get() *Meal {
	m.MutexAvailability.Lock()
	defer m.MutexAvailability.Unlock()
	return m
}

func (m *Meal) set(meal *Meal) {
	m.MutexAvailability.Lock()
	defer m.MutexAvailability.Unlock()
	m.parent = meal.parent
	m.inUse = meal.inUse
	m.prepared = meal.prepared
}

func (m *Meal) getBusyMeal() *Meal {
	m.inUse = 1
	return m
}

func (m *Meal) prepare(cook *Chef, now int64) {
	writeMeal := m.get()
	if writeMeal.prepared == 1 {
		return
	}

	writeMeal.inUse = 1
	writeMeal.preparingTime = now
	writeMeal.cookId = cook.id
	m.set(writeMeal)
	time.Sleep(time.Duration(m.timeRequired) * timeUnit)
	writeMeal.inUse = 1
	writeMeal.prepared = 1
	writeMeal.parent.mealCounter -= 1
	m.set(writeMeal)
}
func newMeal(parent *Order, id int) *Meal {
	food := menu[id]
	return &Meal{
		prepared:          0,
		inUse:             0,
		timeRequired:      food.preparationTime,
		complexity:        food.complexity,
		appliance:         apparatToId[food.chefingApparat],
		preparingTime:     0,
		foodId:            id,
		cookId:            -1,
		parent:            parent,
		MutexAvailability: sync.Mutex{},
	}
}

func parseOrder(postOrder *OrderProp) *Order {
	ret := new(Order)
	ret.id = postOrder.Id
	ret.tableId = postOrder.TableId
	ret.items = postOrder.Items
	ret.mealCounter = 0
	ret.priority = postOrder.Priority
	ret.pickUpTime = postOrder.PickUpTime
	ret.maxWait = postOrder.MaxWait
	for _, id := range postOrder.Items {
		ret.mealCounter += 1
		meal := newMeal(ret, id)
		ret.mealList = append(ret.mealList, meal)
	}
	return ret
}

func (order Order) isReady() bool {
	if order.mealCounter > 0 {
		return false
	}
	return true
}

func getPriority(meal *Meal, timeLeft int) float64 {
	maxWait := meal.parent.maxWait
	orderPriority := meal.parent.priority
	numMeals := int(meal.parent.mealCounter)
	return math.Tanh(float64(numMeals+meal.complexity+orderPriority)/3)*2 - math.Tanh(1-float64(timeLeft)/float64(maxWait))
} //TanH has 3 values, 1, 0, -1

func NewOrderList() *OrderList {
	ret := new(OrderList)
	ret.deliveryMutex = sync.Mutex{}
	ret.mealMutex = sync.Mutex{}
	ret.ovenList = []*Meal{}
	ret.stoveList = []*Meal{}
	ret.nilList = []*Meal{}
	ret.orderArr = []*Order{}
	ret.maxLen = orderListMaxSize
	return ret
}

func (orderList *OrderList) addOrder(order *Order) bool {
	orderList.deliveryMutex.Lock()
	defer orderList.deliveryMutex.Unlock()
	if len(orderList.orderArr) >= orderList.maxLen {
		return false
	}
	orderList.orderArr = append(orderList.orderArr, order)
	for _, meal := range order.mealList {
		applianceId := meal.appliance
		switch applianceId {
		case 0:
			orderList.nilList = append(orderList.nilList, meal)
		case 1:
			orderList.ovenList = append(orderList.ovenList, meal)
		case 2:
			orderList.stoveList = append(orderList.stoveList, meal)
		}
	}
	return true
}

func (orderList *OrderList) getDelivery() *Delivery {
	//avoid geting the same delivery
	orderList.deliveryMutex.Lock()
	defer orderList.deliveryMutex.Unlock()

	for i, order := range orderList.orderArr {
		if order.isReady() {
			for _, meal := range order.mealList {
				applianceId := meal.appliance
				switch applianceId {
				case 0:
					removeFromArr(&orderList.nilList, meal)
				case 1:
					removeFromArr(&orderList.ovenList, meal)
				case 2:
					removeFromArr(&orderList.stoveList, meal)
				}
			}
			orderList.orderArr = append(orderList.orderArr[:i], orderList.orderArr[i+1:]...)
			return newDelivery(order)
		}
	}
	return nil
}

func (orderList *OrderList) getMeal(cook *Chef) *Meal {
	orderList.mealMutex.Lock()
	defer orderList.mealMutex.Unlock()

	now := getTime()
	var priority float64 = 0
	var ret *Meal
	ovenTimeLeft := kitchen.ovens.getTimeLeft(now)
	for _, meal := range orderList.ovenList {
		whichMeal := meal.get()
		if whichMeal.prepared == 0 && whichMeal.inUse == 0 && whichMeal.complexity <= cook.rank {
			localPriority := getPriority(whichMeal, whichMeal.getTimeLeft(now)+ovenTimeLeft)
			if priority < localPriority {
				priority = localPriority
				ret = whichMeal
			}
		}
	}
	stoveTimeLeft := kitchen.stoves.getTimeLeft(now)
	for _, meal := range orderList.stoveList {
		whichMeal := meal.get()
		if whichMeal.prepared == 0 && whichMeal.inUse == 0 && whichMeal.complexity <= cook.rank {
			localPriority := getPriority(whichMeal, whichMeal.getTimeLeft(now)+stoveTimeLeft)
			if priority < localPriority {
				priority = localPriority
				ret = whichMeal
			}
		}
	}
	for _, meal := range orderList.nilList {
		whichMeal := meal.get()
		if whichMeal.prepared == 0 && whichMeal.inUse == 0 && whichMeal.complexity <= cook.rank {
			localPriority := getPriority(whichMeal, whichMeal.getTimeLeft(now))
			if priority < localPriority {
				priority = localPriority
				ret = whichMeal
			}
		}
	}

	if ret != nil {
		return ret.get()
	}

	return ret
}

func (orderList *OrderList) getStatus() string {
	var ret string

	now := getTime()
	for _, order := range orderList.orderArr {

		ret += Div("Order id:" + strconv.Itoa(order.id) + " Meals to prepare:" + strconv.Itoa(int(order.mealCounter)) + "/" + strconv.Itoa(len(order.items)) +
			" Time passed:" + strconv.Itoa(int(now-order.pickUpTime)) + " Max wait:" + strconv.Itoa(order.maxWait))
	}
	return ret
}
