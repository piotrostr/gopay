#!/bin/bash

geth --dev \
  --http \
  --http.api eth,web3,personal,net \
  --http.corsdomain "http://remix.ethereum.org"
