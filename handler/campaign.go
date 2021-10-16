package handler

import (
	"encoding/json"
	"fmt"
	"funding/account"
	"funding/campaign"
	"net/http"
	"strconv"
)

type HandlerCampaign struct {
	service campaign.ServiceCampaign
	account account.Service
}

func NewHandlerCampaign(service campaign.ServiceCampaign, account account.Service) *HandlerCampaign {
	return &HandlerCampaign{service: service, account: account}
}

func (h *HandlerCampaign) CreateCampaign(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("CurrentUser").(account.User)
	w.Header().Set("Content-Type", "application/json")
	var input campaign.CreateCampaignInput

	if r.Method != http.MethodPost {
		response := APIResponse("failed", http.StatusInternalServerError, "failed post bang", nil)
		resp, _ := json.Marshal(response)
		w.Write(resp)
		return
	}

	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		response := APIResponse("failed", http.StatusInternalServerError, "error", err.Error())
		resp, _ := json.Marshal(response)
		w.Write(resp)
		return
	}
	campaign, err := h.service.Create(input, user)
	if err != nil {
		response := APIResponse("failed", http.StatusInternalServerError, fmt.Sprintf("error %v", err), campaign)
		resp, _ := json.Marshal(response)
		w.Write(resp)
		return
	}

	response := APIResponse("sucess", http.StatusOK, "successfully Create Campaign", campaign)
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

func (h *HandlerCampaign) GetCampaigID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodGet {
		response := APIResponse("failed", http.StatusInternalServerError, "method request no allowed", nil)
		resp, _ := json.Marshal(response)
		w.Write(resp)
		return
	}

	fmt.Printf("GET params were:%s", r.URL.Query().Get("id"))
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		response := APIResponse("failed", http.StatusInternalServerError, "FAILED QUERRY", nil)
		resp, _ := json.Marshal(response)
		w.Write(resp)
		return

	}

	campaign, err := h.service.GetCampaignID(uint(id))
	if err != nil {
		response := APIResponse("failed", http.StatusInternalServerError, "FAILED QUERRY", err.Error())
		resp, _ := json.Marshal(response)
		w.Write(resp)
		return
	}

	response := APIResponse("sucessfully get campaign", http.StatusOK, "success", campaign)
	resp, _ := json.Marshal(response)
	w.Write(resp)
}
