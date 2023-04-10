package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"teste-api-golang/configs"
	"teste-api-golang/models"

	_ "github.com/go-sql-driver/mysql"
)

func main() {

	db, err := configs.Connect()
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("testedb")

	users := make([]models.User, 0)
	err = db.Raw("select * from users").Scan(&users).Error
	if err != nil {
		fmt.Println(err)
		return
	}

	products := make([]models.Product, 0)
	err = db.Raw("select * from products").Scan(&products).Error
	if err != nil {
		fmt.Println(err)
		return
	}

	// Execute queries aqui...

	http.HandleFunc("/checkout", checkout)
	http.HandleFunc("/", helloHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))

}

func helloHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello, friend!"))
}

func checkout(w http.ResponseWriter, r *http.Request) {
	// Obter os dados do produto a partir da solicitação
	var checkout models.Checkout
	err := json.NewDecoder(r.Body).Decode(&checkout)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	db, err := configs.Connect()
	if err != nil {
		fmt.Println(err)
		return
	}

	// Verificar se o produto está disponível no estoque
	var product models.Product
	err = db.Where("products_id = ?", checkout.ProductID).First(&product).Error
	if err != nil {
		fmt.Println(err)
		return
	}

	if product.Quantity <= 0 {
		fmt.Println("Produto sem estoque")
		return
	}

	// Verificar se o user está apto a comprar o produto
	var user models.User
	err = db.Where("user_id = ?", checkout.UserID).First(&user).Error
	if err != nil {
		fmt.Println(err)
		return
	}

	if user.Amount < product.Price {
		fmt.Println("Saldo do cliente é insuficiente")
		return
	}

	// Atualize o saldo do usuário
	err = db.Exec("UPDATE users SET quantia=quantia-? WHERE user_id=?", product.Price, user.ID).Error
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Descontar a quantidade do produto no estoque
	err = db.Exec("UPDATE products SET quantidade = quantidade - 1 WHERE products_id = ?", product.ID).Error
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Retornar uma resposta de sucesso
	w.WriteHeader(http.StatusOK)

}
