package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

// define types
type Customer struct {
	Id        uint16 `json:"id"`
	Name      string `json:"name"`
	CompanyId uint16 `json:"companyId"`
	Company   string `json:"company"`
	Notes     string `json:"notes"`
}

// mock "database" of customers
var customers = map[string]Customer{
	"1": { 1, "Turner Kraus", 1, "Udacity" , "Does not have a TV show" },
	"2": { 2, "Stephen Colbert", 2, "Paramount", ""	},
	"3": { 3, "Jimmy Fallon", 3, "NBC Universal", "" },
	"4": { 4, "Seth Meyers", 3, "NBC Universal", ""	},
}

func index(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./public/index.html")
}

func getCustomers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(customers)
}

func getCustomer(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	customerId := mux.Vars(r)["id"]

	if _, ok := customers[customerId]; ok {
		customer := customers[customerId]
		json.NewEncoder(w).Encode(customer)
	} else {
		w.WriteHeader(http.StatusNotFound)
	}
}

func addCustomer(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var newCustomer Customer
	err := json.NewDecoder(r.Body).Decode(&newCustomer)
	
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	} else {
		w.WriteHeader(http.StatusCreated)
		id := createNewId()
		// ensure the id created matches the id within the struct
		newCustomer.Id = uint16(id)
		customers[strconv.FormatInt(int64(id), 10)] = newCustomer
		json.NewEncoder(w).Encode(newCustomer)
	}
}

func updateCustomer(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	customerId := mux.Vars(r)["id"]
	if customer, ok := customers[customerId]; ok {
		var newCustomer Customer
		err := json.NewDecoder(r.Body).Decode(&newCustomer)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		} else {
			updatedCustomer := replaceCustomerValues(customer, newCustomer)
			customers[customerId] = updatedCustomer
			json.NewEncoder(w).Encode(updatedCustomer)
		}
	} else {
		w.WriteHeader(http.StatusNotFound)
	}
}

func deleteCustomer(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	customerId := mux.Vars(r)["id"]
	if _, ok := customers[customerId]; ok {
		delete(customers, customerId)
		w.WriteHeader((http.StatusOK))
	} else {
		w.WriteHeader(http.StatusNotFound)
	}
}

func main() {
	router := mux.NewRouter();
	router.HandleFunc("/", index)
	router.HandleFunc("/customers", getCustomers).Methods("GET")
	router.HandleFunc("/customers/{id}", getCustomer).Methods("GET")
	router.HandleFunc("/customers", addCustomer).Methods("POST")
	router.HandleFunc("/customers/{id}", updateCustomer).Methods("PUT")
	router.HandleFunc("/customers/{id}", deleteCustomer).Methods("DELETE")

	fmt.Println("Server is starting on port 3000")
	http.ListenAndServe(":3000", router)
}

// utils
// identifies the first available integer id for customers, returns it as a string 
func createNewId() int {
	i := 1
	found := false
	for !found {
		key := strconv.Itoa(i)
		if _, ok := customers[key]; ok {
			i++
		} else {
			found = true
		}
	}
	return i
}

// updates any customer fields that are present in new customer, retains everything else from existing customer
func replaceCustomerValues(existingCustomer, newCustomer Customer) Customer {
	updatedCustomer := existingCustomer

	// editable fields exclude Id
	if newCustomer.Company != "" {
		updatedCustomer.Company = newCustomer.Company
	}
	if newCustomer.CompanyId != 0 {
		updatedCustomer.CompanyId = newCustomer.CompanyId
	}
	if newCustomer.Name != "" {
		updatedCustomer.Name = newCustomer.Name
	}
	if newCustomer.Notes != "" {
		updatedCustomer.Notes = newCustomer.Notes
	}

	return updatedCustomer
}
