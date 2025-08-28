from flask import Flaks 
from datetime import datetime
import os

app = Flask(__name__)


@app.route('/')
def home():
    return f"""
    <h1> Hello from flask! </h1>
    """

if __name__ == '__main__':
    app.run(debug=True,host='0.0.0.0')