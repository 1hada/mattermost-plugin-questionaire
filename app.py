from flask import Flask, request, jsonify
from flask_cors import CORS

app = Flask(__name__)
CORS(app, origins=["http://localhost:8065"])

# Simulated external data source
DATA = {
    "question_uuid": "187-5309",
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

HTML_TABLE = """<table border="1">
  <tr>
    <th>Name</th>
    <th>Age</th>
    <th>City</th>
  </tr>
  <tr>
    <td>Alice</td>
    <td>25</td>
    <td>New York</td>
  </tr>
  <tr>
    <td>Bob</td>
    <td>30</td>
    <td>Los Angeles</td>
  </tr>
  <tr>
    <td>Charlie</td>
    <td>28</td>
    <td>Chicago</td>
  </tr>
</table>"""

@app.route('/buttons', methods=['GET'])
def get_buttons():
    return jsonify({
        "question_uuid": DATA["question_uuid"],
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
    is_correct = selected_id == correct_id
    question_uuid = data.get('question_uuid')
    print(f"{data = }") # FOR TESTING

    return jsonify({
        "status": "received",
        "response_stats": HTML_TABLE if is_correct else f"<div>Wow these are some cool stats for {question_uuid}.</div>",
    }), 200

if __name__ == '__main__':
    app.run(debug=True)
