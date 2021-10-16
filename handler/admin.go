package handler

import (
	"encoding/json"
	"funding/admin"
	"funding/auth"
	"net/http"
	"time"

	"github.com/go-playground/validator/v10"
)

type handler struct {
	auth    auth.Authentication
	service admin.Service
}

func NewAdminHandler(service admin.Service) *handler {
	return &handler{service: service}
}

func (h *handler) RegisterAdmin(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var input admin.InputAdmin
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		response := APIResponse("failed", http.StatusInternalServerError, "error parse json", nil)
		resp, _ := json.Marshal(response)
		w.Write(resp)
		return
	}

	validate := validator.New()
	err = validate.Struct(&input)
	if err != nil {
		response := APIResponse("failed", http.StatusInternalServerError, "error validation input", err.Error())
		resp, _ := json.Marshal(response)
		w.Write(resp)
		return
	}

	admin, err := h.service.Register(input)
	if err != nil {
		response := APIResponse("failed", http.StatusInternalServerError, "error to create account admin", err.Error())
		resp, _ := json.Marshal(response)
		w.Write(resp)
		return
	}

	token, err := h.auth.GenerateTokenAdmin(admin.ID)
	if err != nil {
		response := APIResponse("failed", http.StatusInternalServerError, "error to create token", err.Error())
		resp, _ := json.Marshal(response)
		w.Write(resp)
		return
	}

	cookie := http.Cookie{
		Name:    "token_admin",
		Value:   token,
		Expires: time.Now().Add(time.Hour * 10),
	}
	http.SetCookie(w, &cookie)

	resonse := APIResponseToken("success", http.StatusOK, "done", admin, token)
	resp, _ := json.Marshal(resonse)
	w.Write(resp)

}
