package main

import (
	"fmt"
	"log"
	"net/http"
	"teste-api-golang/rest"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	fmt.Println("Start mini-loja...")
	// Execute queries aqui...
	//Checkout
	http.HandleFunc("/checkout", rest.Checkout)
	//http.HandleFunc("/", rest.HelloHandler)//

	//CRUD = CREATE / READ/ UPDATE/ DELETE
	//User

	http.HandleFunc("/create_user", rest.CreateUser)
	http.HandleFunc("/get_users", rest.GetUsers)
	http.HandleFunc("/get_user", rest.GetUserByID)
	http.HandleFunc("/update_user", rest.UpdateUser)
	http.HandleFunc("/delete_user", rest.DeleteUser)

	log.Fatal(http.ListenAndServe(":8080", nil))

}
