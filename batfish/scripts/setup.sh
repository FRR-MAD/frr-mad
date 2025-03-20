#!/bin/bash

# Create a virtual environment (if it doesn't already exist)
if [ ! -d "venv" ]; then
    python -m venv venv
fi

# Activate the virtual environment
source venv/bin/activate

# Install dependencies
pip install -r requirements.txt

echo "Setup complete. Virtual environment is ready to use."