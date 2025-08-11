#!/bin/bash
curl -H "X-CMC_PRO_API_KEY: c1349632-6a87-419a-a00e-d0a93213101a" \
  "https://pro-api.coinmarketcap.com/v1/cryptocurrency/listings/latest?limit=10&convert=USD" | jq
