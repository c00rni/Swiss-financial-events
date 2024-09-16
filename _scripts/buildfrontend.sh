#!/bin/bash

if [[ ! -e _dist ]]; then
    mkdir _dist
fi

cd _api-auth-front
npm run build

