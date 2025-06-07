import React from 'react';
import PropTypes from 'prop-types';
import {Client4} from 'mattermost-redux/client';

export default class RHSView extends React.PureComponent {
    static propTypes = {
        team: PropTypes.object.isRequired,
    }

    constructor(props) {
        super(props);
        this.state = {
            mainText: '',
            sideNote: '',
            buttons: [],
            correctButtonId: null,
            currentUser: null,
            selectedButtonId: null,
            statusMessage: '',
            iconColor: 'black',
            iconSymbol: 'icon-information-outline',
            questionServer: '',
        };
    }

    componentDidMount() {
        // Get the User settings
        this.loadSettings().then(() => {
            fetch(`http://${this.state.questionServer}/buttons`)
                .then((res) => res.json())
                .then((data) => {
                    this.setState({
                        mainText: data.main_text,
                        sideNote: data.side_note,
                        buttons: data.buttons,
                        correctButtonId: data.correct_button,
                    });
                })
                .catch((err) => console.error('Failed to fetch buttons:', err));
        });
        
        Client4.getMe().then((user) => {
            this.setState({
                currentUser: {
                    id: user.id,
                    username: user.username,
                },
            });
        }).catch((err) => {
            console.error('Failed to get current user:', err);
        });
    }
    

    loadSettings() {
        return fetch('/plugins/com.mattermost.questionare/custom_config_settings')
            .then((response) => response.json())
            .then((config) => {
                this.setState({
                    questionServer: `${config.QuestionServerAddress}:${config.QuestionPort}`,
                });
            })
            .catch((err) => console.error('Failed to load settings:', err));
    }

    handleButtonClick = (buttonId) => {
        fetch(`http://${this.state.questionServer}/button-click`, {
            method: 'POST',
            headers: {'Content-Type': 'application/json'},
            body: JSON.stringify({
                id: buttonId,
                user_id: this.state.currentUser?.id,
                username: this.state.currentUser?.username,
                correct_button: this.state.correctButtonId,
            }),
        })
        .then((res) => res.json())
        .then((data) => {
            const isCorrect = data.is_correct;
            this.setState({
                selectedButtonId: buttonId,
                statusMessage: isCorrect ? '✅ Correct answer!' : '❌ Wrong answer!',
                iconColor: isCorrect ? 'green' : 'red',
                iconSymbol: isCorrect ? 'icon-check' : 'icon-close',
            });
        })
        .catch((err) => {
            console.error('Error sending button click:', err);
        });
    }

    renderButtons = () => {
        return this.state.buttons.map((btn) => (
            <button
                key={btn.id}
                onClick={() => this.handleButtonClick(btn.id)}
                style={style.button}
            >
                {btn.label}
            </button>
        ));
    }

    render() {
        return (
            <div style={style.rhs}>
                <div style={style.mainText}>{this.state.mainText}</div>
                <div style={style.sideNote}>{this.state.sideNote}</div>
                <br/>
                <div>{this.renderButtons()}</div>
                <br/><br/>
                <div style={{marginTop: '20px'}}>
                    <div style={{color: this.state.iconColor, fontWeight: 'bold'}}>
                        <i className={`icon ${this.state.iconSymbol}`} style={{marginRight: '5px'}} />
                        {this.state.statusMessage || 'Choose an option above'}
                    </div>
                </div>
            </div>
        );
    }
}

const style = {
    rhs: {
        padding: '10px',
    },
    button: {
        display: 'block',
        margin: '10px 0',
        padding: '10px 20px',
        fontSize: '14px',
        cursor: 'pointer',
    },
    mainText: {
        fontSize: '16px',
        fontWeight: 'bold',
        marginBottom: '10px',
    },
    sideNote: {
        fontSize: '14px',
        color: 'gray',
        marginBottom: '10px',
    },
};
