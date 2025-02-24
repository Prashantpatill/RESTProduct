package main

import (
	"bytes"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
)

//Intitialise a test Instance of the application

var tApp App

func TestMain(m *testing.M) {

	//Setup Db Connection

	err := tApp.Initalise(DbUser, DbPassword, DbTest)
	if err != nil {
		log.Fatalf("Error While Connecting to the Database %s", DbTest)
	}
	log.Printf("Connection Established Sucessfully %s", DbTest)

	createTable()
	m.Run()
}

func createTable() {
	query := ` 
	CREATE TABLE IF NOT EXISTS products(
	id int NOT NULL AUTO_INCREMENT PRIMARY KEY,
	name varchar(255),
	quantity int,
	price float(10,7)
	);`

	_, err := tApp.DB.Exec(query)

	if err != nil {
		log.Fatalf("Error while creating the table %s", err)
	}

}
func addProducts(name string, quantity int, price float64) {
	query := "INSERT INTO products (name, quantity, price) VALUES (?, ?, ?)"
	_, err := tApp.DB.Exec(query, name, quantity, price)
	if err != nil {
		log.Println("Errror is ", err)
	}
}

func clearTable() {
	query := "DELETE from products"
	_, err := tApp.DB.Exec(query)
	if err != nil {
		log.Fatalf("Error created when clearinf table %s", err)
	}

}

func SendRequest(response *http.Request) *httptest.ResponseRecorder {
	recorder := httptest.NewRecorder()
	tApp.Router.ServeHTTP(recorder, response)
	return recorder
}

func checkStatusCode(t *testing.T, expected int, actual int) {
	if expected != actual {
		t.Errorf("Expected Status %v found %v", expected, actual)
	}
}

func TestGetProducts(t *testing.T) {
	//clearTable()

	request, _ := http.NewRequest("GET", "/products", nil)
	response := SendRequest(request)
	checkStatusCode(t, http.StatusAccepted, response.Code)

}

func TestPostProduccts(t *testing.T) {

	clearTable()

	body := []byte(`{"name" : "Mangoes","quantity": 5,"price": 700}`)

	request, _ := http.NewRequest("POST", "/products", bytes.NewBuffer(body))
	request.Header.Set("Content-Type", "application/json")
	response := SendRequest(request)
	checkStatusCode(t, http.StatusCreated, response.Code)

}
