#!/bin/bash

set -euo pipefail

if [[ -z "$REDIS_URL" ]]; then
    >&2 echo "Variable REDIS_URL not set. Exiting..."
    exit 1
fi

# Install package and run.
pip install thesaurize-loader
thesaurize-loader \
    --file=https://www.openoffice.org/lingucomponent/MyThes-1.zip \
    --connection="redis://${REDIS_URL}:6379"