package main

import "time"

func removeFromArr(arr *[]*Meal, ptr *Meal) {
	index := -1
	for i, meal := range *arr {
		if meal == ptr {
			index = i
			break
		}
	}
	if index != -1 {
		*arr = append((*arr)[:index], (*arr)[index+1:]...) //remove from array basically
	}
}

func getTime() int64 {
	return time.Now().UnixNano() / int64(timeUnit)
}

func Div(str string) string {
	return "<div>" + str + "</div>"
}
