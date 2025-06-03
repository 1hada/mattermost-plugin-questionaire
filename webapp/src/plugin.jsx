import React from 'react';

import {FormattedMessage} from 'react-intl';
import manifest from './manifest';

import RHSView from './components/right_hand_sidebar';
import {
    ChannelHeaderButtonIcon,
} from './components/icons';

export default class DemoPlugin {
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
    }

    uninitialize() {
        //eslint-disable-next-line no-console
        console.log(manifest.id + '::uninitialize()');
    }
}
