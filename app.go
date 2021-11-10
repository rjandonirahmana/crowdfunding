package main

import (
	"fmt"
	auth "funding/auth"
	"funding/handler"
	a "funding/handler/admin"
	handlercampaign "funding/handler/campaign"
	handleruser "funding/handler/user"
	"funding/repository"
	"funding/usecase"
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
	repoaccount := repository.NewRepositoryUser(db)
	repoadmin := repository.NewRepositoryAdmin(db)
	repoCampaign := repository.NewRepositoryCampaign(db)

	serviceadmin := usecase.NewServiceAdmin(repoadmin)
	serviceaccount := usecase.NewService(repoaccount)
	middlware := handler.NewMiddleWare(auth, serviceaccount)
	serviceCampaign := usecase.NewServiceCampaign(repoCampaign)

	handleraccount := handleruser.AccountHandler(serviceaccount, auth)
	handerCampaign := handlercampaign.NewHandlerCampaign(serviceCampaign, serviceaccount)
	handlerAdmin := a.NewAdminHandler(serviceadmin)

	http.HandleFunc("/register", handleraccount.RegisterUser)
	http.HandleFunc("/login", handleraccount.Login)
	http.HandleFunc("/campaigns", handerCampaign.GetCampaigns)
	http.HandleFunc("/campaign", handerCampaign.GetCampaigID)
	http.HandleFunc("/create", middlware.MidllerWare(handerCampaign.CreateCampaign))
	http.HandleFunc("/admin", handlerAdmin.RegisterAdmin)

	fmt.Println("starting web server at http://localhost:8181/")

	err = http.ListenAndServe(":8281", nil)
	if err != nil {
		fmt.Println(err)
	}

}
