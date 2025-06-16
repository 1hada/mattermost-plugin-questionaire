package main

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/mattermost/mattermost/server/public/plugin"
)

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

    routePath := "/custom_config_settings"
    router.HandleFunc(routePath, p.handleCustomConfigSettings)

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
}

