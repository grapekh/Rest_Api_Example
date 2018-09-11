package main

/* Copyright (c) 2018 by Howard I Grapek <howiegrapek@yahoo.com>
 * All rights reserved.
 *
 * License: 
 * Redistribution and use in source and binary forms, with or without
 * modification, are permitted provided that the following conditions are met:
 *
 *   - Redistributions of source code must retain the above copyright notice, this
 *     list of conditions and the following disclaimer.
 *
 *   - Redistributions in binary form must reproduce the above copyright notice,
 *     this list of conditions and the following disclaimer in the documentation
 *     and/or other materials provided with the distribution.
 *
 * THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS"
 * AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE
 * IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE
 * DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE
 * FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL
 * DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR
 * SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER
 * CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY,
 * OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
 * OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
 *
 */

/*
 * DESCRIPTION: 
 * Simple RestAPI Interface with Json data. 
 * Tested with Google Chrome's Postman Application
 *
 * This has several kinds of API Actions: 
 *      list all people (GET /people), 
 *      Display s person via ID (GET /people/{id})),
 *      Delete a person via ID (DELETE /people/{id})), 
 *      Create a person record via ID (POST /people/{id}))
 *
 * Dependencies: Gorilla/MUX
 *  go get github.com/gorilla/mux
 *
 * This program listens to port 8888 on localhost. 
 * Go for it.  
 */

import (
    "log"
    "fmt"
    "net/http"
    "encoding/json"
    "github.com/gorilla/mux"
)

// This is the people data... 
// Inside each of the tags there is an omitempty parameter. This means that if the address
// is null, it will be excluded from the JSON data instead of showing up as an empty string or value.
// For this example, id1 and id4 will have full person with address, and id2 is only name, no address. 

type Person struct {
    ID        string   `json:"id,omitempty"`
    Firstname string   `json:"firstname,omitempty"`
    Lastname  string   `json:"lastname,omitempty"`
    Address   *Address `json:"address,omitempty"`
}

type Address struct {
    City  string `json:"city,omitempty"`
    State string `json:"state,omitempty"`
}

// Because there is no database (just testing), just create a public variable that is global.
// it is just a slice of Person and contain all the data used herein
var people []Person

// Return the full person variable to the frontend. 

func GetPerson(w http.ResponseWriter, req *http.Request) {
    params := mux.Vars(req)
    for _, item := range people {
        if item.ID == params["id"] {
            json.NewEncoder(w).Encode(item)
            return
        }
    }
    json.NewEncoder(w).Encode(&Person{})
}

// Return all people in the people slice as json to the frontend. 

func GetPeople(w http.ResponseWriter, req *http.Request) {
    json.NewEncoder(w).Encode(people)
}

// Receive JSON data to work with in the request.
// Decode the JSON data that was passed in and store it in a Person object.
// We assign the new object an id based on what mux found and then we append it to our global slice. 
// In the end, our global array will be returned and it should include everything including our newly 
// added piece of data.

func CreatePerson(w http.ResponseWriter, req *http.Request) {
    params := mux.Vars(req)
    var person Person
    _ = json.NewDecoder(req.Body).Decode(&person)
    person.ID = params["id"]
    people = append(people, person)
    json.NewEncoder(w).Encode(people)
}

// Loop through the data just like GetPerson stuff... 
// The difference is that instead of printing the data, we need to remove it.
// When the id to be deleted has been found, we can recreate our slice with all 
// data excluding that found at the index.

func DeletePerson(w http.ResponseWriter, req *http.Request) {
    params := mux.Vars(req)
    for index, item := range people {
        if item.ID == params["id"] {
            people = append(people[:index], people[index+1:]...)
            break
        }
    }
    json.NewEncoder(w).Encode(people)
}

func main() {
    fmt.Println("Starting Test Restful API on port 8888 = http://localhost:8888")
    fmt.Println("to see ALL PEOPLE: Try with: localhost:8888/people")
    fmt.Println("To see single person (ID 2): Try with: localhost:8888/people/2")

    router := mux.NewRouter()

    // Set up some people for testing. 
    people = append(people, Person{ID: "1", Firstname: "Joe", Lastname: "Dirt", Address: &Address{City: "Los Angeles", State: "CA"}})
    people = append(people, Person{ID: "2", Firstname: "Roger", Lastname: "Rabbit"})
    people = append(people, Person{ID: "4", Firstname: "Mickey", Lastname: "Mouse", Address: &Address{City: "Orlando", State: "FL"}})

    router.HandleFunc("/people",      GetPeople).Methods("GET")
    router.HandleFunc("/people/{id}", GetPerson).Methods("GET")
    router.HandleFunc("/people/{id}", CreatePerson).Methods("POST")
    router.HandleFunc("/people/{id}", DeletePerson).Methods("DELETE")

    log.Fatal(http.ListenAndServe(":8888", router))
}
