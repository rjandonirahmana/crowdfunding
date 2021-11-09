package handlercampaign

// import (
// 	"bytes"
// 	"context"
// 	"encoding/json"
// 	"fmt"
// 	"funding/account"
// 	"funding/auth"
// 	"funding/campaign"
// 	"funding/handler"
// 	"log"
// 	"net/http"
// 	"net/http/httptest"
// 	"os"
// 	"testing"

// 	"github.com/jmoiron/sqlx"
// 	"github.com/joho/godotenv"
// 	_ "github.com/lib/pq"
// 	"github.com/stretchr/testify/assert"
// )

// func dbTest() *sqlx.DB {
// 	err := godotenv.Load("../../.env")
// 	if err != nil {
// 		log.Fatal("Error loading .env file")
// 	}

// 	// secretKey := os.Getenv("SECRET_KEY")
// 	// secretKeyAdmin := os.Getenv("SECRET_KEY_ADMIN")
// 	dbUserName := os.Getenv("DB_USERNAME")
// 	dbName := os.Getenv("DB_NAME")
// 	dbPass := os.Getenv("DB_PASSWORD")

// 	dbString := fmt.Sprintf("host=localhost user=%s password=%s dbname=%s sslmode=disable", dbUserName, dbPass, dbName)
// 	db, err := sqlx.Connect("postgres", dbString)
// 	if err != nil {
// 		log.Fatalln(err)
// 	}

// 	return db
// }

// var (
// 	secretKey       = os.Getenv("SECRET_KEY")
// 	secretKeyAdmin  = os.Getenv("SECRET_KEY_ADMIN")
// 	authentication  = auth.NewAuthentication(secretKey, secretKeyAdmin)
// 	db              = dbTest()
// 	repoAccount     = account.NewRepository(db)
// 	serviceAccount  = account.NewService(repoAccount)
// 	repo            = campaign.NewRepository(db)
// 	service         = campaign.NewServiceCampaign(repo)
// 	MiddlewaresAuth = handler.NewMiddleWare(authentication, serviceAccount)
// 	controller      = NewHandlerCampaign(service, serviceAccount)
// )

// func TestCreateCampaign(t *testing.T) {
// 	token1, err := authentication.GenerateToken(2)
// 	assert.NoError(t, err)
// 	assert.Nil(t, err)

// 	testCases := []struct {
// 		name         string
// 		request      map[string]interface{}
// 		token        string
// 		userID       int
// 		expectedcode int
// 	}{
// 		{
// 			name:         "test1",
// 			request:      map[string]interface{}{"name": "coba", "short_description": "coba", "description": "cobaja", "goal_amount": 1000},
// 			userID:       10,
// 			token:        token1,
// 			expectedcode: 200,
// 		}, {
// 			name:         "test2",
// 			request:      map[string]interface{}{"name": "coba", "short_description": "coba", "description": "cobaja", "goal_amount": "1000"},
// 			userID:       10,
// 			token:        token1,
// 			expectedcode: 422,
// 		}, {
// 			name:         "test2",
// 			request:      map[string]interface{}{"name": "", "short_description": "", "description": "", "goal_amount": 1000},
// 			token:        token1,
// 			userID:       3,
// 			expectedcode: 422,
// 		},
// 	}

// 	for _, test := range testCases {

// 		user, err := serviceAccount.FindByID(test.userID)
// 		assert.NoError(t, err)
// 		reqBody, err := json.Marshal(test.request)
// 		assert.NoError(t, err)

// 		// cookie := http.Cookie{
// 		// 	Name:  "token",
// 		// 	Value: test.token,
// 		// }

// 		req := httptest.NewRequest(http.MethodPost, "/create", bytes.NewBuffer(reqBody))
// 		req.Header.Set("Content-Type", "application/json")
// 		res := httptest.NewRecorder()

// 		// http.SetCookie(res, &cookie)
// 		ctx := context.WithValue(context.Background(), "CurrentUser", user)
// 		req = req.WithContext(ctx)
// 		handler := http.HandlerFunc(controller.CreateCampaign)
// 		handler.ServeHTTP(res, req)
// 		fmt.Println(res)

// 		assert.Equal(t, test.expectedcode, res.Code)
// 	}
// }

// func TestGetCampaign(t *testing.T) {
// 	testCases := []struct {
// 		name         string
// 		request      string
// 		expectedcode int
// 	}{
// 		{
// 			name:         "test1",
// 			request:      "36",
// 			expectedcode: 200,
// 		}, {
// 			name:         "test2",
// 			request:      "0",
// 			expectedcode: 422,
// 		}, {
// 			name:         "test3",
// 			request:      ",,10",
// 			expectedcode: 422,
// 		}, {
// 			name:         "test4",
// 			request:      "100",
// 			expectedcode: 417,
// 		},
// 	}

// 	for _, test := range testCases {
// 		method := http.MethodGet
// 		if test.name == "test4" {
// 			method = http.MethodPost
// 		}
// 		req := httptest.NewRequest(method, "/campaign", nil)
// 		req.Header.Set("Content-Type", "application/json")
// 		q := req.URL.Query()
// 		q.Add("id", test.request)

// 		req.URL.RawQuery = q.Encode()
// 		res := httptest.NewRecorder()

// 		handler := http.HandlerFunc(controller.GetCampaigID)
// 		handler.ServeHTTP(res, req)

// 		fmt.Println(res)

// 		assert.Equal(t, test.expectedcode, res.Code)
// 	}
// }

// func TestGetCampaigns(t *testing.T) {
// 	testCases := []struct {
// 		name         string
// 		expectedcode int
// 	}{
// 		{
// 			name:         "test1",
// 			expectedcode: 200,
// 		},
// 	}

// 	for _, test := range testCases {
// 		req := httptest.NewRequest(http.MethodGet, "/campaigns", nil)
// 		req.Header.Set("Content-Type", "application/json")
// 		res := httptest.NewRecorder()

// 		fmt.Println(test.name)

// 		handler := http.HandlerFunc(controller.GetCampaigns)
// 		handler.ServeHTTP(res, req)

// 		assert.Equal(t, test.expectedcode, res.Code)

// 		fmt.Println(res)
// 	}
// }
