#!/bin/bash

# Build the Python package
python -m build

# Run tests
python -m pytest tests/