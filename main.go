package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/jackc/pgx"
	"log"
	"net/http"
	"os"
	"time"
)

//create struct for database
type user struct {
	Id        int    `json:"id"`
	Name      string `json:"name"`
	Lastname  string `json:"lastname"`
	Age       string `json:"age"`
	Birthdate string `json:"birthdate"`
}

var conn *pgx.Conn
var err1 error
var userr user

func main() {
	//create mux routing
	r := mux.NewRouter()
	//create requests
	r.HandleFunc("/user", addUser).Methods("POST")
	r.HandleFunc("/user", getUser).Methods("GET")
	r.HandleFunc("/user", deleteUser).Methods("DELETE")
	r.HandleFunc("/user", putUser).Methods("PUT")
	//create server
	srv := &http.Server{
		Handler:      r,
		Addr:         "127.0.0.1:5000",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	//run server and print error if server is down
	log.Fatal(srv.ListenAndServe())
}

func addUser(writer http.ResponseWriter, request *http.Request) {
	//create connection to database
	conn, err1 = pgx.Connect(context.Background(), "postgres://postgres:admin@Localhost:5432/users")
	if err1 != nil {
		fmt.Fprintf(os.Stderr, "Unable to connection to database: %v\n", err1)
		os.Exit(1)
	}
	//create new decoder
	decoder := json.NewDecoder(request.Body)
	//decode user from request
	decoder.Decode(&userr)
	conn.Exec(context.Background(), "INSERT INTO users (id, name, lastname, age, birthdate) VALUES ($1,$2,$3,$4,$5)", &userr.Id, &userr.Name, &userr.Lastname, &userr.Age, &userr.Birthdate)
	//close connection
	conn.IsClosed()
	writer.WriteHeader(http.StatusOK)
	return
}

func putUser(writer http.ResponseWriter, request *http.Request) {
	//create connection to database
	conn, err1 = pgx.Connect(context.Background(), "postgres://postgres:admin@Localhost:5432/users")
	if err1 != nil {
		fmt.Fprintf(os.Stderr, "Unable to connection to database: %v\n", err1)
		os.Exit(1)
	}
	//create new decoder
	decoder := json.NewDecoder(request.Body)
	//decode user from request
	decoder.Decode(&userr)
	if userr.Name != "" {
		conn.Exec(context.Background(), "update users set name=$2 where id=$1", &userr.Id, &userr.Name)
	}
	if userr.Lastname != "" {
		conn.Exec(context.Background(), "update users set lastname=$2 where id=$1", &userr.Id, &userr.Lastname)
	}
	if userr.Age != "" {
		conn.Exec(context.Background(), "update users set age=$2 where id=$1", &userr.Id, &userr.Age)
	}
	if userr.Birthdate != "" {
		conn.Exec(context.Background(), "update users set birthdate=$2 where id=$1", &userr.Id, &userr.Birthdate)
	}
	//close connection
	conn.IsClosed()
	writer.WriteHeader(http.StatusOK)
	return

}

func deleteUser(writer http.ResponseWriter, request *http.Request) {
	//create connection to database
	conn, err1 = pgx.Connect(context.Background(), "postgres://postgres:admin@Localhost:5432/users")
	if err1 != nil {
		fmt.Fprintf(os.Stderr, "Unable to connection to database: %v\n", err1)
		os.Exit(1)
	}
	//create new decoder
	decoder := json.NewDecoder(request.Body)
	//decode user from request
	decoder.Decode(&userr)
	conn.Exec(context.Background(), "delete from users where id=$1", userr.Id)
	//close connection
	conn.IsClosed()
	writer.WriteHeader(http.StatusOK)
	return
}

func getUser(writer http.ResponseWriter, request *http.Request) {
	//create connection to database
	conn, err1 = pgx.Connect(context.Background(), "postgres://postgres:admin@Localhost:5432/users")
	var users []user
	if err1 != nil {
		fmt.Fprintf(os.Stderr, "Unable to connection to database: %v\n", err1)
		os.Exit(1)
	}
	//writing to rows result of sql query
	rows, _ := conn.Query(context.Background(), "SELECT * FROM users")
	//writing the received strings to an array users
	for rows.Next() {
		err := rows.Scan(&userr.Id, &userr.Name, &userr.Lastname, &userr.Age, &userr.Birthdate)
		if err != nil {
			log.Fatalf("Unable to scan the row. %v \n", err)
		}
		users = append(users, userr)
	}
	//converting array to json and writing it to writer
	json.NewEncoder(writer).Encode(users)
	//close connection
	conn.IsClosed()
	return
}
