package user

import (
	"bytes"
	"encoding/json"
	"fmt"
	"funding/account"
	"funding/auth"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
)

func dbTest() *sqlx.DB {
	err := godotenv.Load("../../.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// secretKey := os.Getenv("SECRET_KEY")
	// secretKeyAdmin := os.Getenv("SECRET_KEY_ADMIN")
	dbUserName := os.Getenv("DB_USERNAME")
	dbName := os.Getenv("DB_NAME")
	dbPass := os.Getenv("DB_PASSWORD")

	dbString := fmt.Sprintf("host=localhost user=%s password=%s dbname=%s sslmode=disable", dbUserName, dbPass, dbName)
	db, err := sqlx.Connect("postgres", dbString)
	if err != nil {
		log.Fatalln(err)
	}

	return db
}

var (
	secretKey      = os.Getenv("SECRET_KEY")
	secretKeyAdmin = os.Getenv("SECRET_KEY_ADMIN")
	authentication = auth.NewAuthentication(secretKey, secretKeyAdmin)
	db             = dbTest()
	repo           = account.NewRepository(db)
	service        = account.NewService(repo)
	controller     = AccountHandler(service, authentication)
)

func TestRegisterUser(t *testing.T) {
	testCases := []struct {
		name       string
		request    map[string]interface{}
		expectCode int
		expectMsg  string
	}{
		{
			name:       "test1",
			request:    map[string]interface{}{"name": "ateng", "email": "ateg@gmail.com", "password": "12345678", "confirm_password": "12345678", "occupation": "IT"},
			expectCode: 200,
			expectMsg:  "success",
		}, {
			name:       "test2",
			request:    map[string]interface{}{"name": "ateng", "email": "ateng@", "password": "12345678", "confirm_password": "12345678", "occupation": "IT"},
			expectCode: 500,
			expectMsg:  "failed",
		}, {
			name:       "test1",
			request:    map[string]interface{}{"name": "ateng", "email": "ateg@gmail.com", "password": "12345678", "confirm_password": "12345678", "occupation": "IT"},
			expectCode: 422,
			expectMsg:  "success",
		},
	}

	for _, test := range testCases {
		reqBody, err := json.Marshal(test.request)
		assert.Nil(t, err)

		req := httptest.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(reqBody))
		req.Header.Set("Content-Type", "application/json")
		res := httptest.NewRecorder()

		hanlder := http.HandlerFunc(controller.RegisterUser)
		hanlder.ServeHTTP(res, req)

		fmt.Println(res)
		assert.Equal(t, test.expectCode, res.Code)
	}
}

// func TestLogin(t *testing.T) {

// }
