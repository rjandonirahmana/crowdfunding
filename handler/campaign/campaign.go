package handlercampaign

import (
	"encoding/json"
	"fmt"
	"funding/account"
	"funding/campaign"
	"funding/handler"
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
		response := handler.APIResponse("failed", http.StatusInternalServerError, "failed post bang", nil)
		resp, _ := json.Marshal(response)
		w.Write(resp)
		return
	}

	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		response := handler.APIResponse("failed", http.StatusUnprocessableEntity, "error", err.Error())
		resp, _ := json.Marshal(response)
		w.WriteHeader(422)
		w.Write(resp)
		return
	}
	campaign, err := h.service.Create(input, user)
	if err != nil {
		response := handler.APIResponse("failed", 422, fmt.Sprintf("error %v", err), campaign)
		resp, _ := json.Marshal(response)
		w.WriteHeader(422)
		w.Write(resp)
		return
	}

	response := handler.APIResponse("sucess", http.StatusOK, "successfully Create Campaign", campaign)
	resp, _ := json.Marshal(response)
	w.Write(resp)

}

func (h *HandlerCampaign) GetCampaigns(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	campaigns, err := h.service.GetAllCampaigns()
	if err != nil {
		resp := handler.APIResponse("failed to get campaigns", http.StatusBadRequest, "failed", err.Error())
		respBody, _ := json.Marshal(resp)
		w.WriteHeader(400)
		w.Write(respBody)
		return
	}
	response := handler.APIResponse("sucess fully get campaigns", http.StatusOK, "success", campaigns)
	resp, _ := json.Marshal(response)
	w.Write(resp)

}

func (h *HandlerCampaign) GetCampaigID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodGet {
		response := handler.APIResponse("failed", http.StatusInternalServerError, "method request no allowed", nil)
		resp, _ := json.Marshal(response)
		w.WriteHeader(http.StatusExpectationFailed)
		w.Write(resp)
		return
	}

	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		response := handler.APIResponse("failed", http.StatusUnprocessableEntity, "FAILED QUERRY", nil)
		resp, _ := json.Marshal(response)
		w.WriteHeader(422)
		w.Write(resp)
		return

	}

	campaign, err := h.service.GetCampaignID(uint(id))
	if err != nil {
		response := handler.APIResponse("failed", http.StatusInternalServerError, "FAILED QUERRY", err.Error())
		resp, _ := json.Marshal(response)
		w.WriteHeader(422)
		w.Write(resp)
		return
	}

	response := handler.APIResponse("sucessfully get campaign", http.StatusOK, "success", campaign)
	resp, _ := json.Marshal(response)
	w.Write(resp)
}
