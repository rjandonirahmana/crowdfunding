package admin

import (
	"encoding/json"
	"funding/auth"
	"funding/handler"
	"funding/model"
	"funding/usecase"
	"net/http"
	"time"

	"github.com/go-playground/validator/v10"
)

type adminhandler struct {
	auth    auth.Authentication
	service usecase.ServiceAdmin
}

func NewAdminHandler(service usecase.ServiceAdmin) *adminhandler {
	return &adminhandler{service: service}
}

func (h *adminhandler) RegisterAdmin(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var input model.InputAdmin
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		response := handler.APIResponse("failed", http.StatusInternalServerError, "error parse json", nil)
		resp, _ := json.Marshal(response)
		w.Write(resp)
		return
	}

	validate := validator.New()
	err = validate.Struct(&input)
	if err != nil {
		response := handler.APIResponse("failed", http.StatusInternalServerError, "error validation input", err.Error())
		resp, _ := json.Marshal(response)
		w.Write(resp)
		return
	}

	admin, err := h.service.Register(input)
	if err != nil {
		response := handler.APIResponse("failed", http.StatusInternalServerError, "error to create account admin", err.Error())
		resp, _ := json.Marshal(response)
		w.Write(resp)
		return
	}

	token, err := h.auth.GenerateTokenAdmin(admin.ID)
	if err != nil {
		response := handler.APIResponse("failed", http.StatusInternalServerError, "error to create token", err.Error())
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

	resonse := handler.APIResponseToken("success", http.StatusOK, "done", admin, token)
	resp, _ := json.Marshal(resonse)
	w.Write(resp)

}
