package main

//import statements
import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

//user object structure
type User struct {
	Userid         int
	Username       string
	Fname          string
	Lname          string
	Email          string
	Password       string
	IsAdmin        bool
	CreatedDate    time.Time
	LastModDate    time.Time
	LastLoggedDate time.Time
}

// creates a db connection
func dbConn() (db *sql.DB) {
	dbDriver := "mysql"
	dbUser := "root"
	dbPass := "admin"
	dbName := "ecnapplication"
	db, err := sql.Open(dbDriver, dbUser+":"+dbPass+"@/"+dbName+"?parseTime=true")
	if err != nil {
		panic(err.Error())
	}
	return db
}

// returns all users
func getUsers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	db := dbConn()
	result, err := db.Query("SELECT * FROM user")
	if err != nil {
		panic(err.Error())
	}
	defer result.Close()
	var users []User
	for result.Next() {
		var user User
		err := result.Scan(&user.Userid, &user.Username, &user.Email, &user.Fname, &user.Lname, &user.Password, &user.CreatedDate, &user.LastModDate, &user.LastLoggedDate, &user.IsAdmin)
		if err != nil {
			panic(err.Error())
		}
		users = append(users, user)
	}
	json.NewEncoder(w).Encode(users)
}

// returns a single user from user id
func getUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	db := dbConn()
	params := mux.Vars(r)
	userid := params["userid"]
	result, err := db.Query("SELECT * FROM user WHERE user_id=?", userid)
	if err != nil {
		panic(err.Error())
	}
	var users []User
	defer result.Close()
	for result.Next() {
		var user User
		err := result.Scan(&user.Userid, &user.Username, &user.Email, &user.Fname, &user.Lname, &user.Password, &user.CreatedDate, &user.LastModDate, &user.LastLoggedDate, &user.IsAdmin)
		if err != nil {
			panic(err.Error())
		}
		users = append(users, user)
	}
	json.NewEncoder(w).Encode(users)
}

//creates a new user
func createUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	db := dbConn()

	insForm, err := db.Prepare("INSERT INTO user(username, email, fname, lname, password, createddate, lastmoddate, lastloggeddate) VALUES(?,?,?,?,?,?,?,?)")
	if err != nil {
		panic(err.Error())
	}

	//reads jason object
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err.Error())
	}
	keyVal := make(map[string]string)
	json.Unmarshal(body, &keyVal)
	username := keyVal["username"]
	email := keyVal["email"]
	fname := keyVal["fname"]
	lname := keyVal["lname"]
	password := keyVal["password"]
	dt := time.Now()
	_, err = insForm.Exec(username, email, fname, lname, password, dt, dt, dt)
	defer db.Close()
	if err != nil {
		panic(err.Error())
	}
	fmt.Fprintf(w, "New user was created")

}

//update a user
func updateUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	db := dbConn()
	params := mux.Vars(r)
	userid := params["userid"]

	insForm, err := db.Prepare("UPDATE user SET username=?, email=?, fname=?, lname=?, password=?, lastmoddate=? WHERE user_id=?")
	if err != nil {
		panic(err.Error())
	}

	//reads jason object
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err.Error())
	}
	keyVal := make(map[string]string)
	json.Unmarshal(body, &keyVal)
	username := keyVal["username"]
	email := keyVal["email"]
	fname := keyVal["fname"]
	lname := keyVal["lname"]
	password := keyVal["password"]
	dt := time.Now()
	_, err = insForm.Exec(username, email, fname, lname, password, dt, userid)

	defer db.Close()
	if err != nil {
		panic(err.Error())
	}
	fmt.Fprintf(w, "User "+userid+" was updated")
}

//deletes user from userid
func deleteUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	db := dbConn()
	params := mux.Vars(r)
	userid := params["userid"]
	delForm, err := db.Prepare("DELETE FROM user WHERE user_id=?")
	if err != nil {
		panic(err.Error())
	}
	delForm.Exec(userid)
	defer db.Close()

	if err != nil {
		panic(err.Error())
	}
	fmt.Fprintf(w, "User "+userid+" was deleted")
}

//user login and update last logged time
func userAuth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	db := dbConn()
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err.Error())
	}
	keyVal := make(map[string]string)
	json.Unmarshal(body, &keyVal)
	username := keyVal["username"]
	password := keyVal["password"]

	//check if user credentials valid
	result, err := db.Query("SELECT * FROM user WHERE username=? AND password=?", username, password)
	if err != nil {
		panic(err.Error())
	}
	var users []User
	defer result.Close()
	for result.Next() {
		var user User
		err := result.Scan(&user.Userid, &user.Username, &user.Email, &user.Fname, &user.Lname, &user.Password, &user.CreatedDate, &user.LastModDate, &user.LastLoggedDate, &user.IsAdmin)
		if err != nil {
			panic(err.Error())
		}
		users = append(users, user)
	}

	//update last logged time
	insForm, err := db.Prepare("UPDATE user SET lastloggeddate=? WHERE username=? AND password=?")
	if err != nil {
		panic(err.Error())
	}
	dt := time.Now()
	_, err = insForm.Exec(dt, username, password)

	defer db.Close()
	if err != nil {
		panic(err.Error())
	}

	json.NewEncoder(w).Encode(users)

}

func main() {
	router := mux.NewRouter().StrictSlash(true)

	router.HandleFunc("/api/users", getUsers).Methods("GET")
	router.HandleFunc("/api/users/{userid}", getUser).Methods("GET")
	router.HandleFunc("/api/users", createUser).Methods("POST")
	router.HandleFunc("/api/users/{userid}", updateUser).Methods("PUT")
	router.HandleFunc("/api/users/{userid}", deleteUser).Methods("DELETE")
	router.HandleFunc("/api/login", userAuth).Methods("GET")
	log.Fatal(http.ListenAndServe(":8080", router))
}
