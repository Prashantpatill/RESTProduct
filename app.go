package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

type App struct {
	Router *mux.Router
	DB     *sql.DB
}

func (app *App) Initalise() error {
	var err error

	dbConnector := fmt.Sprintf("%v:%v@tcp(127.0.0.1:3306)/%v", DbUser, DbPassword, DbName)
	app.DB, err = sql.Open("mysql", dbConnector)

	if err != nil {
		log.Println("Error Encountered during Opening the sconnecion string")
		return err
	}
	app.Router = mux.NewRouter().StrictSlash(true)
	app.handleRoutes()
	return nil
}

func (app *App) Run(address string) {
	log.Fatal(http.ListenAndServe(address, app.Router))

}

func (app *App) handleRoutes() {
	app.Router.HandleFunc("/products", app.getProducts).Methods("GET")
	app.Router.HandleFunc("/products/{id}", app.getProductsById).Methods("GET")
	app.Router.HandleFunc("/products", app.insertProducts).Methods("POST")
	app.Router.HandleFunc("/products/{id}", app.updateDetails).Methods("PUT")
	app.Router.HandleFunc("/products/{id}", app.deleteProducts).Methods("DELETE")
}

func sendResponse(w http.ResponseWriter, status int, payload interface{}) {
	response, _ := json.Marshal(payload)
	w.Header().Set("content-type", "application/json")
	w.WriteHeader(status)
	w.Write(response)
}

func sendError(w http.ResponseWriter, status int, err string) {
	error_message := map[string]string{"error": err}
	sendResponse(w, status, error_message)
}

func (app *App) getProducts(w http.ResponseWriter, r *http.Request) {
	products, err := getProducts(app.DB)

	if err != nil {
		sendError(w, http.StatusInternalServerError, err.Error())
		return
	}
	sendResponse(w, http.StatusAccepted, products)

}

func (app *App) getProductsById(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)

	key, err := strconv.Atoi(vars["id"])

	if err != nil {
		sendError(w, http.StatusInternalServerError, "Invalid product Id")
	}

	p := product{Id: key}

	err = p.getProductsById(app.DB)

	if err != nil {
		switch err {
		case sql.ErrNoRows:
			sendError(w, http.StatusInternalServerError, "No product with Product ID Found")
		default:
			sendError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}
	sendResponse(w, http.StatusAccepted, p)

}

func (app *App) insertProducts(w http.ResponseWriter, r *http.Request) {

	var p product

	err := json.NewDecoder(r.Body).Decode(&p)

	if err != nil {
		sendError(w, http.StatusBadRequest, "check The pauload")
		return
	}

	err = p.insertProducts(app.DB)

	if err != nil {
		sendError(w, http.StatusBadRequest, err.Error())
	}
	sendResponse(w, http.StatusOK, p)

}

func (app *App) updateDetails(w http.ResponseWriter, r *http.Request) {

	muxvar := mux.Vars(r)

	key, err := strconv.Atoi(muxvar["id"])

	if err != nil {
		sendError(w, http.StatusBadRequest, "Invalid Product Id")
		return
	}

	var p product

	err = json.NewDecoder(r.Body).Decode(&p)

	if err != nil {
		sendError(w, http.StatusBadRequest, "Invalid Product Data")
		return
	}

	p.Id = key

	err = p.updateDetails(app.DB)

	if err != nil {
		sendError(w, http.StatusBadRequest, "Invalid Product")
	}
	sendResponse(w, http.StatusOK, p)

}

func (app *App) deleteProducts(w http.ResponseWriter, r *http.Request) {
	muxvar := mux.Vars(r)
	key, err := strconv.Atoi(muxvar["id"])

	if err != nil {
		sendError(w, http.StatusBadRequest, "Invalid ProductId")
	}

	var p product
	p.Id = key

	err = p.deleteProducts(app.DB)

	if err != nil {
		sendError(w, http.StatusBadRequest, "Invalid Product")
	}
	sendResponse(w, http.StatusOK, map[string]string{"message": "Product deleted successfully"})

}
