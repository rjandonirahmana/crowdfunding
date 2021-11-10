package handleruser

import (
	"encoding/json"
	"errors"
	"fmt"
	auth "funding/auth"
	"funding/handler"
	"funding/model"
	"funding/usecase"
	"io"
	"net/http"
	"time"

	"github.com/go-playground/validator/v10"
)

type HanlderUser struct {
	service usecase.ServiceUser
	auth    auth.Authentication
}

func AccountHandler(service usecase.ServiceUser, auth auth.Authentication) *HanlderUser {
	return &HanlderUser{service: service, auth: auth}
}

func (h *HanlderUser) RegisterUser(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodPost {
		response := handler.APIResponse("failed", http.StatusUnprocessableEntity, "cannot continue req", nil)
		resp, _ := json.Marshal(response)
		w.Write(resp)
		return

	}

	var account model.RegisterUserInput

	err := json.NewDecoder(r.Body).Decode(&account)
	if err != nil {
		response := handler.APIResponse("failed", http.StatusInternalServerError, "error to decode json", nil)
		resp, _ := json.Marshal(response)
		w.WriteHeader(500)
		io.WriteString(w, string(resp))
		return

	}
	validate := validator.New()
	err = validate.Struct(&account)
	if err != nil {
		response := handler.APIResponse("failed", http.StatusInternalServerError, "error to continue", err.Error())
		resp, _ := json.Marshal(response)
		w.WriteHeader(500)
		w.Write(resp)
		return
	}

	user, err := h.service.RegisterUser(&account)
	if err != nil {
		response := handler.APIResponse("failed", http.StatusUnprocessableEntity, "error", err.Error())
		resp, _ := json.Marshal(response)
		w.WriteHeader(422)
		w.Write(resp)
		return

	}

	token, err := h.auth.GenerateToken(user.ID)
	if err != nil {
		response := handler.APIResponse("failed", http.StatusUnprocessableEntity, "error to create token", err)
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

	response := handler.APIResponseToken("sucess", http.StatusOK, "successfully create new account", user, token)
	resp, _ := json.Marshal(response)
	w.Write(resp)
}

func (h *HanlderUser) Login(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodPost {
		response := handler.APIResponse("failed", http.StatusInternalServerError, "cant continue if method not post", errors.New("need to method post"))
		resp, _ := json.Marshal(response)
		w.WriteHeader(500)
		w.Write(resp)
		return
	}

	var input model.LoginInput

	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		response := handler.APIResponse("failed", http.StatusInternalServerError, "error", err)
		resp, _ := json.Marshal(response)
		w.WriteHeader(500)
		w.Write(resp)
		return
	}

	user, err := h.service.Login(&input)
	if err != nil {

		response := handler.APIResponse("failed", 422, fmt.Sprintf("%v", err), err)
		resp, _ := json.Marshal(response)
		w.WriteHeader(422)
		w.Write(resp)
		return
	}

	token, err := h.auth.GenerateToken(user.ID)
	if err != nil {
		response := handler.APIResponse("failed", http.StatusInternalServerError, "error", err)
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

	response := handler.APIResponseToken("sucess", http.StatusOK, "successfully Login", user, token)
	resp, _ := json.Marshal(response)
	w.Write(resp)

}
