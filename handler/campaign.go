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
	w.Header().Set("Content-Type", "application/json")
	user := r.Context().Value("user").(account.User)
	var input campaign.CreateCampaignInput

	if r.Method != http.MethodPost {
		response := APIResponse("failed", http.StatusInternalServerError, "failed post bang", nil)
		resp, _ := json.Marshal(response)
		w.Write(resp)
		return
	}

	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		response := APIResponse("failed", http.StatusInternalServerError, "error", err)
		resp, _ := json.Marshal(response)
		w.Write(resp)
		return
	}
	err = h.service.Create(input, user)
	if err != nil {
		response := APIResponse("failed", http.StatusInternalServerError, fmt.Sprintf("error %v", err), err)
		resp, _ := json.Marshal(response)
		w.Write(resp)
		return
	}

	response := APIResponse("sucess", http.StatusOK, "successfully Create Campaign", user)
	resp, _ := json.Marshal(response)
	w.Write(resp)

}

func (h *HandlerCampaign) GetCampaigns(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	campaigns, err := h.service.GetAllCampaigns()
	if err != nil {
		resp := APIResponse("failed to get campaigns", http.StatusBadRequest, "failed", err.Error())
		respBody, _ := json.Marshal(resp)
		w.Write(respBody)
		return
	}
	response := APIResponse("sucess fully get campaigns", http.StatusOK, "success", campaigns)
	resp, _ := json.Marshal(response)
	w.Write(resp)

}
