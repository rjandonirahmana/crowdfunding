package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"funding/account"
	auth "funding/authentikasi"
	"net/http"
	"strings"
)

type MiddlewaresAuth struct {
	auth    auth.Authentication
	service account.Service
}

func NewMiddleWare(auth auth.Authentication, service account.Service) *MiddlewaresAuth {
	return &MiddlewaresAuth{auth: auth, service: service}
}

func (m *MiddlewaresAuth) MidllerWare(next http.HandlerFunc) http.HandlerFunc {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		authorizationHeader := r.Header.Get("Authorization")
		if !strings.Contains(authorizationHeader, "Bearer") {
			fmt.Println("error 1")
			response := APIResponse("failed", http.StatusUnauthorized, "error", nil)
			resp, _ := json.Marshal(response)
			w.Write(resp)
			return
		}

		var token string
		tokenString := strings.Split(authorizationHeader, " ")
		if len(tokenString) == 2 {
			token = tokenString[1]
		}

		id, err := m.auth.ValidateToken(token)
		if err != nil {
			fmt.Println("error 2")
			response := APIResponse("failed", http.StatusUnauthorized, "error", err)
			resp, _ := json.Marshal(response)
			w.Write(resp)
			return
		}

		user, err := m.service.FindByID(int(id))
		if err != nil {
			fmt.Println("error 3")
			response := APIResponse("failed", http.StatusUnprocessableEntity, "error", err)
			resp, _ := json.Marshal(response)
			w.Write(resp)
			return
		}

		ctx := context.WithValue(context.Background(), "user", user)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)

	})
}
