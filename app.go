package main

import (
	"fmt"
	"funding/account"
	auth "funding/authentikasi"
	"funding/handler"
	"log"
	"net/http"
	"os"

	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	secretKey := os.Getenv("SECRET_KEY")
	dbUserName := os.Getenv("DB_USERNAME")
	dbName := os.Getenv("DB_NAME")
	dbPass := os.Getenv("DB_PASSWORD")

	dbString := fmt.Sprintf("host=localhost port=5432 user=%s password=%s dbname=%s sslmode=disable", dbUserName, dbPass, dbName)
	db, err := sqlx.Connect("postgres", dbString)
	if err != nil {
		log.Fatalln(err)
	}

	auth := auth.NewAuthentication([]byte(secretKey))
	repoaccount := account.NewRepository(db)
	serviceaccount := account.NewService(repoaccount)
	handleraccount := handler.AccountHandler(serviceaccount, auth)

	// repoCampaign := campaign.NewRepository(db)
	// serviceCampaign := campaign.NewServiceCampaign(repoCampaign)
	// handerCampaign := handler.NewHandlerCampaign(serviceCampaign, serviceaccount)

	http.HandleFunc("/register", handleraccount.RegisterUser)
	http.HandleFunc("/login", handleraccount.Login)

	fmt.Println("starting web server at http://localhost:8181/")

	http.ListenAndServe(":8181", nil)

}
