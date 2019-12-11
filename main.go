package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

type user struct {
	ID	string `json:"ID"`
	Name	string `json:"Name"`
	Score	int `json:"Score"`
}

type allUsers []user

var users = allUsers{
	{
		ID:          "1",
		Name:       "Chester",
		Score: 9999999,
	},
}

func homeRoute(w http.ResponseWriter, r* http.Request){
	fmt.Fprintf(w, "Welcome to Chester's math game!")
}

func SaveUsersRoute(w http.ResponseWriter, r* http.Request){

	var newUser user
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintf(w, "user not saved")
	}
	
	json.Unmarshal(reqBody, &newUser)

	//check req body
	if newUser.ID != "" && newUser.Name != "" {
		for _, userIter := range users {
			if userIter.ID == newUser.ID {
				fmt.Fprintf(w, "user exists already")
				w.WriteHeader(400)
				return
			}
		}
		users = append(users, newUser)
	} else {
		fmt.Fprintf(w, "User info missing")
		w.WriteHeader(400)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newUser)
}

func getAllUsersRoute(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(users)
}

func updateScoreRoute(w http.ResponseWriter, r* http.Request){
	userID := mux.Vars(r)["id"]
	var updatedUser user

	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintf(w, "Could not update user")
	}
	json.Unmarshal(reqBody, &updatedUser)

	for i, userIter := range users {
		if userIter.ID == userID {
			userIter.Score = updatedUser.Score
			users = append(users[:i], userIter)

			w.WriteHeader(http.StatusCreated)
			json.NewEncoder(w).Encode(userIter)
		}
	}
}

func getPort() string {
	p := os.Getenv("PORT")
	if p != "" {
	  return ":" + p
	}
	return ":4400"
}

func main() {
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", homeRoute)
	router.HandleFunc("/saveusers", SaveUsersRoute).Methods("POST")
	router.HandleFunc("/getallusers", getAllUsersRoute).Methods("GET")
	router.HandleFunc("/updatescores/{id}", updateScoreRoute).Methods("PATCH")
	port := getPort()
	log.Fatal(http.ListenAndServe(port, router))
}
