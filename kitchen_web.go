package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

type KitchenWeb struct {
	kitchenServer   http.Server
	kitchenHandler  KitchenHandler
	kitchenClient   http.Client
	connectionError error
}

func (kw *KitchenWeb) start() {
	kw.kitchenServer.Addr = kitchenPort
	kw.kitchenServer.Handler = &kw.kitchenHandler

	fmt.Println(time.Now())
	if err := kw.kitchenServer.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}

func (kw *KitchenWeb) deliver(delivery *Delivery) bool {

	requestBody, marshallErr := json.Marshal(delivery)
	if marshallErr != nil {
		log.Fatal(marshallErr)
	}

	request, newRequestError := http.NewRequest(http.MethodPost, restaurantAddress+restaurantPort+"/delivery", bytes.NewBuffer(requestBody))
	if newRequestError != nil {
		fmt.Println("Could not create new request. Error:", newRequestError)
		log.Fatal(newRequestError)
	} else {
		response, doError := kw.kitchenClient.Do(request)
		if doError != nil {
			fmt.Println("ERROR Sending request. ERR:", doError)
			log.Fatal(doError)
		}
		var responseBody = make([]byte, response.ContentLength)
		response.Body.Read(responseBody)
		if string(responseBody) != "OK" {
			return false
		}
		return true
	}
	return true
}

func (kw *KitchenWeb) establishConnection() bool {
	if kitchen.connected == true {
		return false
	}
	request, _ := http.NewRequest(http.MethodConnect, restaurantAddress+restaurantPort+"/", bytes.NewBuffer([]byte{}))
	response, err := kw.kitchenClient.Do(request)
	if err != nil {
		kw.connectionError = err
		return false
	}
	var responseBody = make([]byte, response.ContentLength)
	response.Body.Read(responseBody)
	if string(responseBody) != "OK" {
		return false
	}
	return true
}
