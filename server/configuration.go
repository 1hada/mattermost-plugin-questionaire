package main

import (
	"encoding/json"

	"github.com/pkg/errors"

	"github.com/mattermost/mattermost/server/public/model"
	"github.com/mattermost/mattermost/server/public/pluginapi"
)

// configuration captures the plugin's external configuration as exposed in the Mattermost server
// configuration, as well as values computed from the configuration. Any public fields will be
// deserialized from the Mattermost server configuration in OnConfigurationChange.
//
// As plugins are inherently concurrent (hooks being called asynchronously), and the plugin
// configuration can change at any time, access to the configuration must be synchronized. The
// strategy used in this plugin is to guard a pointer to the configuration, and clone the entire
// struct whenever it changes. You may replace this with whatever strategy you choose.
type configuration struct {
	// An address which handles question and answer requests.  
	QuestionServerAddress string

	// The port handling the question and answer requests.
	QuestionPort string

	// disabled tracks whether or not the plugin has been disabled after activation. It always starts enabled.
	disabled bool
}

// Clone deep copies the configuration. 
func (c *configuration) Clone() *configuration {
	return &configuration{
		QuestionServerAddress:   c.QuestionServerAddress,
		QuestionPort:            c.QuestionPort,
		disabled:                c.disabled,
	}
}

// getConfiguration retrieves the active configuration under lock, making it safe to use
// concurrently. The active configuration may change underneath the client of this method, but
// the struct returned by this API call is considered immutable.
func (p *Plugin) getConfiguration() *configuration {
	p.configurationLock.RLock()
	defer p.configurationLock.RUnlock()

	if p.configuration == nil {
		return &configuration{}
	}

	return p.configuration
}

// setConfiguration replaces the active configuration under lock.
//
// Do not call setConfiguration while holding the configurationLock, as sync.Mutex is not
// reentrant. In particular, avoid using the plugin API entirely, as this may in turn trigger a
// hook back into the plugin. If that hook attempts to acquire this lock, a deadlock may occur.
//
// This method panics if setConfiguration is called with the existing configuration. This almost
// certainly means that the configuration was modified without being cloned and may result in
// an unsafe access.
func (p *Plugin) setConfiguration(configuration *configuration) {
	p.configurationLock.Lock()
	defer p.configurationLock.Unlock()

	if configuration != nil && p.configuration == configuration {
		panic("setConfiguration called with the existing configuration")
	}

	p.configuration = configuration
}


// OnConfigurationChange is invoked when configuration changes may have been made.
//
func (p *Plugin) OnConfigurationChange() error {
	p.API.LogInfo("FOR TESTING log OnConfigurationChange START")

	if p.client == nil {
		p.client = pluginapi.NewClient(p.API, p.Driver)
	}

	configuration := p.getConfiguration().Clone()

	// Load the public configuration fields from the Mattermost server configuration.
	if loadConfigErr := p.API.LoadPluginConfiguration(configuration); loadConfigErr != nil {
		return errors.Wrap(loadConfigErr, "failed to load plugin configuration")
	}

	p.setConfiguration(configuration)
	p.API.LogInfo("FOR TESTING log OnConfigurationChange END")

	return nil
}


// ConfigurationWillBeSaved is invoked before saving the configuration to the
// backing store.
// An error can be returned to reject the operation. Additionally, a new
// config object can be returned to be stored in place of the provided one.
// Minimum server version: 8.0
//
// This demo implementation logs a message to the demo channel whenever config
// is going to be saved.
// If the Username config option is set to "invalid" an error will be
// returned, resulting in the config not getting saved.
// If the Username config option is set to "replaceme" the config value will be
// replaced with "replaced".
func (p *Plugin) ConfigurationWillBeSaved(newCfg *model.Config) (*model.Config, error) {
	p.API.LogInfo("FOR TESTING ConfigurationWillBeSaved")

	cfg := p.getConfiguration()
	if cfg.disabled {
		return nil, nil
	}
	
	configData := newCfg.PluginSettings.Plugins[manifest.Id]
	js, err := json.Marshal(configData)
	if err != nil {
		p.API.LogError(
			"Failed to marshal config data ConfigurationWillBeSaved",
			"error", err.Error(),
		)
		return nil, nil
	}

	if err := json.Unmarshal(js, &cfg); err != nil {
		p.API.LogError(
			"Failed to unmarshal config data ConfigurationWillBeSaved",
			"error", err.Error(),
		)
		return nil, nil
	}

	p.API.LogInfo("NewCfg returned.")
	return newCfg, nil
}

// setEnabled wraps setConfiguration to configure if the plugin is enabled.
func (p *Plugin) setEnabled(enabled bool) {
	p.API.LogInfo("FOR TESTING setEnabled")

	var configuration = p.getConfiguration().Clone()
	configuration.disabled = !enabled

	p.setConfiguration(configuration)
}
