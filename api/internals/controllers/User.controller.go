package controllers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/morelmiles/go-events/config"
	"github.com/morelmiles/go-events/internals/models"
	"golang.org/x/crypto/bcrypt"
)

// Compares the password before saving it
func ComparePassword(password string, hashedPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

// Checks if the user exists
func checkIfUserExists(userId string) bool {
	var user models.User
	config.DB.First(&user, userId)

	return user.ID != 0
}

func GetUsers(w http.ResponseWriter, r *http.Request) {

	var users []models.User
	config.DB.Find(&users)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(&users)
}

// GetUserById - Fetches a list of all users.
func GetUserById(w http.ResponseWriter, r *http.Request) {
	userId := mux.Vars(r)["id"]

	if !checkIfUserExists(userId) {
		json.NewEncoder(w).Encode("user not found!")
		return
	}

	var user models.User
	config.DB.First(&user, userId)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

// CreateUser - Creates a new user
func CreateUser(w http.ResponseWriter, r *http.Request) {
	var user models.User
	var err error

	json.NewDecoder(r.Body).Decode(&user)

	// // Fix this middleware
	// if _, err := middleware.HashPassword(user.Password); err != nil {
	// 	log.Printf("Error: %d", err)
	// 	return
	// }

	newUser := config.DB.Create(&user)
	err = newUser.Error

	if err != nil {
		log.Panic(err)
	} else {
		json.NewEncoder(w).Encode(&user)
	}
}

// UpdateUserById -  Updates a single user by the ID specified
func UpdateUserById(w http.ResponseWriter, r *http.Request) {
	userId := mux.Vars(r)["id"]
	if !checkIfUserExists(userId) {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode("user not found!")
		return
	}

	var user models.User

	config.DB.First(&user, userId)
	json.NewDecoder(r.Body).Decode(&user)
	config.DB.Save(&user)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

// DeleteUserById - Updates a single user by the ID specified.
func DeleteUserById(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	userId := mux.Vars(r)["id"]
	if !checkIfUserExists(userId) {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode("user not found!")
		return
	}

	var user models.User
	config.DB.Delete(&user, userId)
	json.NewEncoder(w).Encode(user)
}

func GetAllEventsByUser(w http.ResponseWriter, r *http.Request) {
	userId := mux.Vars(r)["id"]

	var user models.User
	var properties models.Event

	config.DB.Model(&user).Find(properties).Where("id = ?", userId)
	json.NewEncoder(w).Encode(user)
}
