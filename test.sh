#!/bin/sh
set -e

curl -x localhost:8080 --cacert ca.pem --location --request POST https://postman-echo.com/post?a=1 \
	--data-raw "This is expected to be sent back as part of response body."
