#!/bin/bash

# Define MySQL credentials
MYSQL_USER="root"
MYSQL_PASS="NewPassword"
MYSQL_DB="k6"

# Path to your k6 output JSON file
JSON_FILE="./output.json"

# Extract the test name and the result data (JSON) from the output file
TEST_NAME=$(jq -r '.meta.name' $JSON_FILE)
RESULT_DATA=$(cat $JSON_FILE)

# Insert the extracted data into MySQL
mysql -u$MYSQL_USER -p$MYSQL_PASS $MYSQL_DB -e "INSERT INTO test_results (test_name, result_data) VALUES ('$TEST_NAME', '$RESULT_DATA');"

echo "Data inserted successfully!"

k6 run --out experimental-prometheus-rw=http://localhost:9090 loadTest.js
