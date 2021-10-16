package main

import (
	"fmt"
	"funding/account"
	auth "funding/auth"
	"funding/campaign"
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
	secretKeyAdmin := os.Getenv("SECRET_KEY_ADMIN")
	dbUserName := os.Getenv("DB_USERNAME")
	dbName := os.Getenv("DB_NAME")
	dbPass := os.Getenv("DB_PASSWORD")

	dbString := fmt.Sprintf("host=localhost user=%s password=%s dbname=%s sslmode=disable", dbUserName, dbPass, dbName)
	db, err := sqlx.Connect("postgres", dbString)
	if err != nil {
		log.Fatalln(err)
	}

	auth := auth.NewAuthentication(secretKey, secretKeyAdmin)
	repoaccount := account.NewRepository(db)
	serviceaccount := account.NewService(repoaccount)
	handleraccount := handler.AccountHandler(serviceaccount, auth)
	middlware := handler.NewMiddleWare(auth, serviceaccount)

	repoCampaign := campaign.NewRepository(db)

	serviceCampaign := campaign.NewServiceCampaign(repoCampaign)
	handerCampaign := handler.NewHandlerCampaign(serviceCampaign, serviceaccount)

	http.HandleFunc("/register", handleraccount.RegisterUser)
	http.HandleFunc("/login", handleraccount.Login)
	http.HandleFunc("/campaigns", handerCampaign.GetCampaigns)
	http.HandleFunc("/campaign", handerCampaign.GetCampaigID)
	http.HandleFunc("/create", middlware.MidllerWare(handerCampaign.CreateCampaign))

	fmt.Println("starting web server at http://localhost:8181/")

	http.ListenAndServe(":8181", nil)

}
