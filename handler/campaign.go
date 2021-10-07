package handler

import (
	"encoding/json"
	"fmt"
	"funding/account"
	"funding/campaign"
	"net/http"
)

type HandlerCampaign struct {
	service campaign.ServiceCampaign
	account account.Service
}

func NewHandlerCampaign(service campaign.ServiceCampaign, account account.Service) *HandlerCampaign {
	return &HandlerCampaign{service: service, account: account}
}

func (h *HandlerCampaign) CreateCampaign(w http.ResponseWriter, r *http.Request) {
	var input campaign.CreateCampaignInput

	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		response := APIResponse("failed", http.StatusInternalServerError, "error", err)
		resp, _ := json.Marshal(response)
		w.Write(resp)
	}

	email := w.Header().Get("user")
	user, err := h.account.FindByEmail(email)
	if err != nil {
		response := APIResponse("failed", http.StatusInternalServerError, fmt.Sprintf("error %v", err), err)
		resp, _ := json.Marshal(response)
		w.Write(resp)
	}
	err = h.service.Create(input, user)
	if err != nil {
		response := APIResponse("failed", http.StatusInternalServerError, fmt.Sprintf("error %v", err), err)
		resp, _ := json.Marshal(response)
		w.Write(resp)
	}

	response := APIResponse("sucess", http.StatusOK, "successfully Create Campaign", user)
	resp, _ := json.Marshal(response)
	w.Write(resp)

}
