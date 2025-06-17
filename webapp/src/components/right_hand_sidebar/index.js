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
