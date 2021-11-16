package main

import (
	elastic "funding/elasticsearch"
	"time"

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

	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {

	el, err := elastic.NewCreateIndex([]string{"http://localhost:9200"})
	if err != nil {
		log.Fatal(err)
	}

	err = el.CreateIndex("campaign")
	if err != nil {
		fmt.Println(err)
	}

	err = godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	fs := http.FileServer(http.Dir("public"))
	http.Handle("/", fs)

	secretKey := os.Getenv("SECRET_KEY")
	secretKeyAdmin := os.Getenv("SECRET_KEY_ADMIN")
	dbUserName := os.Getenv("DB_USERNAME")
	dbName := os.Getenv("DB_NAME")
	dbPass := os.Getenv("DB_PASSWORD")
	clientID := os.Getenv("GITHUB_CLIENTID")
	clientSecet := os.Getenv("GITHUB_CLIENTSECRET")

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
	oauthHanlder := handler.NewOauthHanlder(clientID, clientSecet, time.Second*15)
	handlerAdmin := a.NewAdminHandler(serviceadmin)

	r := mux.NewRouter()

	r.HandleFunc("/register", handleraccount.RegisterUser).Methods("POST")
	r.HandleFunc("/login", handleraccount.Login).Methods("POST")
	r.HandleFunc("/campaigns", handerCampaign.GetCampaigns).Methods("GET")
	r.HandleFunc("/campaign", handerCampaign.GetCampaigID).Methods("GET")
	r.HandleFunc("/create", middlware.MidllerWare(handerCampaign.CreateCampaign)).Methods("POST")
	r.HandleFunc("/admin", handlerAdmin.RegisterAdmin).Methods("POST")
	r.HandleFunc("/auth", oauthHanlder.OauthAutentication)

	fmt.Println("starting web server at http://localhost:8181/")

	server := &http.Server{
		Handler: r,
		Addr:    "127.0.0.1:8181",
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Fatal(server.ListenAndServe())

}
