package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"
)

type OauthHandler struct {
	ClientID     string
	ClientSecet  string
	TimeDuration time.Duration
}

func NewOauthHanlder(clientID, ClientSecret string, timeD time.Duration) *OauthHandler {
	return &OauthHandler{ClientID: clientID, ClientSecet: ClientSecret, TimeDuration: timeD}
}

type OAuthAccessResponse struct {
	AccessToken string `json:"access_token"`
}

func (h *OauthHandler) OauthAutentication(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	err := r.ParseForm()
	if err != nil {
		response := APIResponse("failed", http.StatusInternalServerError, "error to parse form", err.Error())
		resp, _ := json.Marshal(response)
		w.WriteHeader(http.StatusBadRequest)
		w.Write(resp)
		return
	}
	code := r.FormValue("code")
	reqURL := fmt.Sprintf("https://github.com/login/oauth/access_token?client_id=%s&client_secret=%s&code=%s", h.ClientID, h.ClientSecet, code)
	req, err := http.NewRequest(http.MethodPost, reqURL, nil)
	if err != nil {
		fmt.Fprintf(os.Stdout, "could not create HTTP request: %v", err)
		w.WriteHeader(http.StatusBadRequest)
	}

	req.Header.Set("accept", "application/json")
	client := &http.Client{
		Timeout: h.TimeDuration,
	}

	res, err := client.Do(req)
	if err != nil {
		response := APIResponse("failed", http.StatusInternalServerError, "error to parse form", err.Error())
		resp, _ := json.Marshal(response)
		w.WriteHeader(http.StatusBadRequest)
		w.Write(resp)
		return
	}

	defer res.Body.Close()
	var response OAuthAccessResponse
	fmt.Println(res.Status)
	fmt.Println("ini auth")

	err = json.NewDecoder(res.Body).Decode(&response)
	if err != nil {
		response := APIResponse("failed", http.StatusInternalServerError, "error to parse form", err.Error())
		resp, _ := json.Marshal(response)
		w.WriteHeader(http.StatusBadRequest)
		w.Write(resp)
		return
	}

	fmt.Println(res.Body)

	w.Header().Set("Location", "/welcome.html?access_token="+response.AccessToken)
	w.WriteHeader(http.StatusFound)

}
