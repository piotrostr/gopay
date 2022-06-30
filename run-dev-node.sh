#!/bin/bash

# geth --dev \
#   --http \
#   --http.api eth,web3,personal,net \
#   --http.corsdomain "http://remix.ethereum.org"

ganache-cli --account "0x$PRIVATE_KEY,10000000000000000000"
