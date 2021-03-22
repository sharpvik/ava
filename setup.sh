#!/usr/bin/bash

## Environment Setup
if [ ! -f ".env" ]; then
        python3 make_env.py > .env
fi
