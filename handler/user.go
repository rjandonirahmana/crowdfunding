package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"funding/account"
	auth "funding/authentikasi"
	"net/http"
	"time"
)

type HanlderUser struct {
	service account.Service
	auth    auth.Authentication
}

func AccountHandler(service account.Service, auth auth.Authentication) *HanlderUser {
	return &HanlderUser{service: service, auth: auth}
}

func (h *HanlderUser) RegisterUser(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodPost {
		response := APIResponse("failed", http.StatusUnprocessableEntity, "cannot continue req", nil)
		resp, _ := json.Marshal(response)
		w.Write(resp)
		return

	}
	var account account.RegisterUserInput

	err := json.NewDecoder(r.Body).Decode(&account)
	if err != nil {
		response := APIResponse("failed", http.StatusInternalServerError, "error", nil)
		resp, _ := json.Marshal(response)
		w.Write(resp)
		return

	}

	if account.Password != account.ConfirmPassword {
		response := APIResponse("failed", http.StatusUnprocessableEntity, "your password and confirm passw doesnt match", nil)
		resp, _ := json.Marshal(response)
		w.Write(resp)
		return
	}

	user, err := h.service.RegisterUser(account)
	if err != nil {
		response := APIResponse("failed", http.StatusUnprocessableEntity, "error", err)
		resp, _ := json.Marshal(response)
		w.Write(resp)
		return

	}

	token, err := h.auth.GenerateToken(user.ID)
	if err != nil {
		response := APIResponse("failed", http.StatusUnprocessableEntity, "error to create token", err)
		resp, _ := json.Marshal(response)
		w.Write(resp)
		return
	}

	cookie := http.Cookie{
		Name:    "token",
		Value:   token,
		Expires: time.Now().Add(time.Hour * 12),
	}

	http.SetCookie(w, &cookie)

	response := APIResponseToken("sucess", http.StatusOK, "successfully create new account", user, token)
	resp, _ := json.Marshal(response)
	w.Write(resp)
}

func (h *HanlderUser) Login(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodPost {
		fmt.Println("a")
		response := APIResponse("failed", http.StatusInternalServerError, "cant continue if method not post", errors.New("need to method post"))
		resp, _ := json.Marshal(response)
		w.Write(resp)
		return
	}

	var input account.LoginInput

	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		response := APIResponse("failed", http.StatusInternalServerError, "error", err)
		resp, _ := json.Marshal(response)
		w.Write(resp)
		return
	}

	user, err := h.service.Login(input)
	if err != nil {

		response := APIResponse("failed", http.StatusInternalServerError, fmt.Sprintf("%v", err), err)
		resp, _ := json.Marshal(response)
		w.Write(resp)
		return
	}

	token, err := h.auth.GenerateToken(user.ID)
	if err != nil {
		response := APIResponse("failed", http.StatusInternalServerError, "error", err)
		resp, _ := json.Marshal(response)
		w.Write(resp)
		return
	}

	cookie := http.Cookie{
		Name:    "token",
		Value:   token,
		Expires: time.Now().Add(time.Hour * 12),
	}

	http.SetCookie(w, &cookie)

	response := APIResponseToken("sucess", http.StatusOK, "successfully Login", user, token)
	resp, _ := json.Marshal(response)
	w.Write(resp)

}
