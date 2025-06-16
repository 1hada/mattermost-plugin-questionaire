package main

import (
	"encoding/json"
	"net/http"
	"bytes"
	"fmt"
	"io"
	"time"

	"github.com/gorilla/mux"
	"github.com/mattermost/mattermost/server/public/plugin"
)

// Button represents a single button in the JSON response from /buttons
type Button struct {
	ID   string `json:"id"`
	Text string `json:"label"`
}

// GetButtonsResponse represents the JSON response structure from the /buttons endpoint
type GetButtonsResponse struct {
	QuestionUUID string   `json:"question_uuid"`
	MainText     string   `json:"main_text"`
	SideNote     string   `json:"side_note"`
	Buttons      []Button `json:"buttons"`
	CorrectButton string `json:"correct_button"`
}

// ButtonClickRequest represents the JSON payload for the /button-click endpoint
type ButtonClickRequest struct {
	ID           string `json:"id"`
	UserID           string `json:"user_id"`
	UserName           string `json:"username"`
	CorrectButton string `json:"correct_button"`
	QuestionUUID string `json:"question_uuid"`
}

// ButtonClickResponse represents the JSON response structure from the /button-click endpoint
type ButtonClickResponse struct {
	Status      string `json:"status"`
	ResponseStats string `json:"response_stats"`
}

// ServeHTTP allows the plugin to implement the http.Handler interface. Requests destined for the
// /plugins/{id} path will be routed to the plugin.
//
// The Mattermost-User-Id header will be present if (and only if) the request is by an
// authenticated user.
//
// This demo implementation sends back whether or not the plugin hooks are currently enabled. It
// is used by the web app to recover from a network reconnection and synchronize the state of the
// plugin's hooks.
func (p *Plugin) ServeHTTP(c *plugin.Context, w http.ResponseWriter, r *http.Request) {

	p.router.ServeHTTP(w, r)
}


func (p *Plugin) initializeAPI() {
	router := mux.NewRouter()

    router.HandleFunc("/custom_config_settings", p.handleCustomConfigSettings)
    router.HandleFunc("/buttons", p.fetchButtons)
    router.HandleFunc("/button-click", p.sendButtonClick)

	p.router = router

}


func (p *Plugin) handleCustomConfigSettings(w http.ResponseWriter, r *http.Request) {
	p.API.LogInfo("FOR TESTING handleCustomConfigSettings we are requesting the QuestionServerAddress :" + p.getConfiguration().QuestionServerAddress + " and QuestionPort : " + p.getConfiguration().QuestionPort)

	// Basic auth check
	config := map[string]string{
		"QuestionServerAddress": p.getConfiguration().QuestionServerAddress,
		"QuestionPort": p.getConfiguration().QuestionPort,
	}
	responseJSON, _ := json.Marshal(config)

	w.Header().Set("Content-Type", "application/json")
	if _, err := w.Write(responseJSON); err != nil {
		p.API.LogError("Failed to write custom config settings", "err", err.Error())
	}
	p.API.LogInfo("FOR TESTING handleCustomConfigSettings DONE")
}



// fetchButtons makes a GET request to the /buttons endpoint and parses the response.
func (p *Plugin)  fetchButtons(w http.ResponseWriter, r *http.Request)  {
	p.API.LogInfo("FOR TESTING fetchButtons we are requesting the QuestionServerAddress :" + p.getConfiguration().QuestionServerAddress + " and QuestionPort : " + p.getConfiguration().QuestionPort)
	url := fmt.Sprintf("%s:%s/buttons", p.getConfiguration().QuestionServerAddress, p.getConfiguration().QuestionPort)

	// Create an HTTP client with a timeout
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		p.API.LogError("FOR TESTING failed to make GET request :", "err", err.Error())
		return
	}
	defer resp.Body.Close() // Ensure the response body is closed

	// Check if the HTTP status code indicates success
	if resp.StatusCode != http.StatusOK {
		p.API.LogInfo("FOR TESTING received non-OK status code")
		return
	}

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		p.API.LogError("FOR TESTING failed to read response body :", "err", err.Error())
		return
	}

	// Unmarshal the JSON response into the struct
	var buttonsResponse GetButtonsResponse
	if err = json.Unmarshal(body, &buttonsResponse); err != nil {
		p.API.LogError("FOR TESTING failed to unmarshal JSON response :", "err", err.Error())
		return
	}

	// 2. Marshal the extracted clickData into JSON for the request
	jsonPayload, err := json.Marshal(buttonsResponse)
	if err != nil {
		p.API.LogError("FOR TESTING Failed to marshal JSON payload for Flask", "err", err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if _, err = w.Write(jsonPayload); err != nil {
		p.API.LogError("FOR TESTING Failed to write custom config settings", "err", err.Error())
	}
}

// sendButtonClick is an HTTP handler that processes a button click POST request.
// It reads the click data from the request body, forwards it to the Flask server,
func (p *Plugin) sendButtonClick(w http.ResponseWriter, r *http.Request) {
	p.API.LogInfo("FOR TESTING sendButtonClick we are requesting the QuestionServerAddress :" + p.getConfiguration().QuestionServerAddress + " and QuestionPort : " + p.getConfiguration().QuestionPort)
	url := fmt.Sprintf("%s:%s/button-click", p.getConfiguration().QuestionServerAddress, p.getConfiguration().QuestionPort)

	// Ensure the request method is POST
	if r.Method != http.MethodPost {
		p.API.LogInfo("FOR TESTING Only POST requests are allowed")
		return
	}

	// Set content type for the response back to the client
	w.Header().Set("Content-Type", "application/json")

	// 1. Decode the incoming JSON payload from the client
	var clickData ButtonClickRequest
	err := json.NewDecoder(r.Body).Decode(&clickData)
	if err != nil {
		p.API.LogError("FOR TESTING Failed to decode request body", "err", err.Error())
		return
	}
	defer r.Body.Close()

	// 2. Marshal the extracted clickData into JSON for the request
	jsonPayload, err := json.Marshal(clickData)
	if err != nil {
		p.API.LogError("FOR TESTING Failed to marshal JSON payload for Flask", "err", err.Error())
		return
	}

	// 3. Create and send the POST request to the server
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(jsonPayload))
	if err != nil {
		p.API.LogError("FOR TESTING Failed to create request", "err", err.Error())
		return
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		p.API.LogError("FOR TESTING Failed to send request", "err", err.Error())

		return
	}
	defer resp.Body.Close()

	// 4. Read the response from the Flask server
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		p.API.LogError("FOR TESTING Failed to read response", "err", err.Error())
		return
	}

	// 5. Unmarshal response
	var clickResponse ButtonClickResponse
	if err = json.Unmarshal(body, &clickResponse); err != nil {
		p.API.LogError("FOR TESTING Failed to unmarshal response JSON", "err", err.Error())
		return
	}

	// 6. Set the HTTP status code and encode Flask's response to the client
	w.WriteHeader(resp.StatusCode) // Use Flask's original status code
	if err = json.NewEncoder(w).Encode(clickResponse); err != nil {
		p.API.LogError("FOR TESTING Error encoding response for client", "err", err.Error())
	}
}
