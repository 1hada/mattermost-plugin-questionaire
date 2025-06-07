import React from 'react';

import {FormattedMessage} from 'react-intl';
import manifest from './manifest';

import RHSView from './components/right_hand_sidebar';
import {
    ChannelHeaderButtonIcon,
} from './components/icons';
import reducer from './reducer';

import {
    fetchPluginSettings,
    getStatus,
} from './actions';

export default class QuestionarePlugin {
    initialize(registry, store) {
        const {toggleRHSPlugin} = registry.registerRightHandSidebarComponent(
            RHSView,
            <FormattedMessage
                id='plugin.name'
                defaultMessage='Questionaire'
            />);

        registry.registerChannelHeaderButtonAction(
            <ChannelHeaderButtonIcon/>,
            () => store.dispatch(toggleRHSPlugin),
            <FormattedMessage
                id='plugin.name'
                defaultMessage='Questionaire'
            />,
        );
        
        //registry.registerAdminConsoleCustomSetting('CustomSetting', CustomSetting);
        registry.registerReducer(reducer);
        // Immediately fetch the current plugin status.
        store.dispatch(fetchPluginSettings());
        store.dispatch(getStatus());
    }

    uninitialize() {
        //eslint-disable-next-line no-console
        console.log(manifest.id + '::uninitialize()');
    }
}
