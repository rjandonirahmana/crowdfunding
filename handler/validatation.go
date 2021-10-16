package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"funding/account"
	"funding/admin"
	auth "funding/auth"
	"net/http"
)

type MiddlewaresAuth struct {
	auth         auth.Authentication
	service      account.Service
	ServiceAdmin admin.Service
}

func NewMiddleWare(auth auth.Authentication, service account.Service) *MiddlewaresAuth {
	return &MiddlewaresAuth{auth: auth, service: service}
}

func (m *MiddlewaresAuth) MidllerWare(next http.HandlerFunc) http.HandlerFunc {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		c, err := r.Cookie("token")
		if err != nil {
			fmt.Println("error 1")
			response := APIResponse("failed", http.StatusUnauthorized, "error", nil)
			resp, _ := json.Marshal(response)
			w.Write(resp)
			return
		}

		id, err := m.auth.ValidateToken(c.Value)
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

		ctx := context.WithValue(context.Background(), "CurrentUser", user)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)

	})
}

func (m *MiddlewaresAuth) MidllerWareAdmin(next http.HandlerFunc) http.HandlerFunc {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		w.Header().Set("Content-Type", "application/json")
		c, err := r.Cookie("token_admin")
		if err != nil {
			fmt.Println("error 1")
			response := APIResponse("failed", http.StatusUnauthorized, "error", nil)
			resp, _ := json.Marshal(response)
			w.Write(resp)
			return
		}

		id, err := m.auth.ValidateTokenAdmin(c.Value)
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

		ctx := context.WithValue(context.Background(), "CurrentAdmin", user)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)

	})
}
