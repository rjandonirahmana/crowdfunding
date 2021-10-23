package admin

import (
	"bytes"
	"encoding/json"
	"fmt"
	"funding/admin"
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
	db         = dbTest()
	repo       = admin.NewRepositoryAdmin(db)
	service    = admin.NewServiceAdmin(repo)
	controller = NewAdminHandler(service)
)

func TestCreateAdmin(t *testing.T) {
	testCases := []struct {
		codetest     int
		request      map[string]interface{}
		expectedCode int
		expectedMsg  string
	}{
		{
			codetest:     1,
			request:      map[string]interface{}{"name": "siaparahp", "email": "doniiii123@gmail.com", "password": "1234567890rahp", "confirm_password": "1234567890rahp", "job_id": 1},
			expectedCode: 200,
			expectedMsg:  "sucess",
		},
	}

	for _, test := range testCases {
		reqBody, err := json.Marshal(test.request)
		assert.NoError(t, err)
		req := httptest.NewRequest(http.MethodPost, "/admin", bytes.NewBuffer(reqBody))
		req.Header.Set("Content-Type", "application/json")
		res := httptest.NewRecorder()

		handler := http.HandlerFunc(controller.RegisterAdmin)
		handler.ServeHTTP(res, req)

		assert.Equal(t, test.expectedCode, res.Code)

		fmt.Println(res)

	}
}
