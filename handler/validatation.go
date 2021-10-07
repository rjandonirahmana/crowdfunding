package handler

import (
	"encoding/json"
	"funding/account"
	auth "funding/authentikasi"
	"net/http"
)

type MiddlewaresAuth struct {
	auth    auth.Authentication
	service account.Service
}

func (m *MiddlewaresAuth) MidllerWare(http.HandlerFunc) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		cookie, err := r.Cookie("token")
		if err != nil {
			response := APIResponse("failed", http.StatusInternalServerError, "error", err)
			resp, _ := json.Marshal(response)
			w.Write(resp)
			return
		}

		token := cookie.Value

		id, err := m.auth.ValidateToken(token)
		if err != nil {
			response := APIResponse("failed", http.StatusUnauthorized, "error", err)
			resp, _ := json.Marshal(response)
			w.Write(resp)
			return
		}

		user, err := m.service.FindByID(id)
		if err != nil {
			response := APIResponse("failed", http.StatusUnprocessableEntity, "error", err)
			resp, _ := json.Marshal(response)
			w.Write(resp)
			return
		}

		w.Header().Set("user", user.Email)

	}
}
