import {connect} from 'react-redux';

import {getCurrentTeam} from 'mattermost-redux/selectors/entities/teams';

import {isEnabled, getQuestionServerAddress, getQuestionPort} from 'selectors';

import RHSView from './rhs_view';

const mapStateToProps = (state) => ({
    enabled: isEnabled(state),
    questionServerAddress: getQuestionServerAddress(state),
    questionPort: getQuestionPort(state),
    team: getCurrentTeam(state),
});

export default connect(mapStateToProps)(RHSView);

/*
To access a custom plugin user setting from Mattermost Redux, you'll need to first ensure the setting is stored in the Mattermost server configuration under the Plugin section. 
You can then access this setting from your plugin's web app code by using the Redux store and appropriate selectors. 
Here's a more detailed breakdown:
1. Storing the Custom Setting:

    Server Configuration:
    Plugin settings are stored within the Mattermost server configuration under the [ Plugins ] section, indexed by plugin IDs.
    Plugin Section:
    Each plugin has its own section within this configuration, where custom settings are defined.
    Example:
    You might define a setting like my_plugin.custom_setting in the [ Plugins ] section for your plugin, according to Mattermost documentation. 

2. Accessing the Setting in Your Web App:

    Register Admin Console Custom Setting:
    If you want administrators to be able to configure this setting, you'll need to register a custom component in the web app to manage it. 
    This is done using the registerAdminConsoleCustomSetting function. 

Redux Integration:
Once the setting is stored in the server configuration, you can access it within your plugin's web app using Redux. 
Selectors:
You'll use Redux selectors to retrieve the custom setting from the store. These selectors will typically map the plugin's settings to a specific key in the Redux state. 
Example:
If your custom setting is my_plugin.custom_setting, you might have a selector like getPluginSetting('my_plugin', 'custom_setting'). 

3. Example using React Redux 
JavaScript

import { connect } from 'react-redux';
import { getPluginSetting } from 'mattermost-redux/selectors/settings'; // Or similar

// Your component
const MyComponent = ({ customSettingValue }) => {
  return (
    <div>
      <p>Custom Setting Value: {customSettingValue}</p>
    </div>
  );
};

// Map the Redux store to props
const mapStateToProps = (state) => {
  const customSettingValue = getPluginSetting(state, 'my_plugin', 'custom_setting');
  return {
    customSettingValue,
  };
};

// Connect the component to the Redux store
export default connect(mapStateToProps, null)(MyComponent);

In summary, to get a custom plugin user setting from Mattermost Redux:

    Store the setting in the Mattermost server configuration under the [ Plugins ] section. 

Register a custom component in your plugin's web app (if needed for administrative configuration). 
Use Redux selectors to retrieve the setting from the store. 
*/