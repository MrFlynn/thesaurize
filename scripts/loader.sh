#!/bin/bash

set -uo pipefail

if [[ -z "$REDIS_URL" ]]; then
    >&2 echo "Variable REDIS_URL not set. Exiting..."
    exit 1
fi

# Check if gcc is installed.
command -v gcc git
if [[ $? -eq 1 ]]; then
    apt-get update
    apt-get install -y gcc git
fi

# Install package and run.
pip install \
    https://github.com/MrFlynn/thesaurize/releases/download/v1.11.13/thesaurize_loader-0.2.4-py3-none-any.whl

thesaurize-loader \
    --file=https://www.openoffice.org/lingucomponent/MyThes-1.zip \
    --connection="redis://${REDIS_URL}:6379"