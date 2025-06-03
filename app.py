from flask import Flask, request, jsonify
from flask_cors import CORS

app = Flask(__name__)
CORS(app, origins=["http://localhost:8065"])

# Simulated external data source
DATA = {
    "main_text": "This is text that should be added before the buttons.",
    "side_note": "This is text that should be added between the main-text and buttons and be grayed out.",
    "buttons": [
        {"id": "button1", "label": "Approve"},
        {"id": "button2", "label": "Reject"},
        {"id": "button3", "label": "Maybe"},
        {"id": "button4", "label": "Escalate"},
    ],
    "correct_button": "button1",  # This will be dynamically set in real app
}

@app.route('/buttons', methods=['GET'])
def get_buttons():
    return jsonify({
        "main_text": DATA["main_text"],
        "side_note": DATA["side_note"],
        "buttons": DATA["buttons"],
        "correct_button": DATA["correct_button"],
    })

@app.route('/button-click', methods=['POST'])
def handle_click():
    data = request.get_json()
    selected_id = data.get('id')
    correct_id = data.get('correct_button')
    print(f"{data = }") # FOR TESTING

    is_correct = selected_id == correct_id
    return jsonify({
        "status": "received",
        "button": selected_id,
        "is_correct": is_correct,
    }), 200

if __name__ == '__main__':
    app.run(debug=True)
