#!/bin/bash

# Check for required parameters
if [ "$#" -ne 2 ]; then
    echo "Usage: $0 <API_URL> <LANGUAGE>"
    exit 1
fi

API_URL=$1
LANGUAGE=$2

# Create or empty the output files
TIME_OUTPUT_FILE="${LANGUAGE}OutputTime.txt"
ACTIVATION_ID_OUTPUT_FILE="${LANGUAGE}activation_ids.txt"

> $TIME_OUTPUT_FILE
> $ACTIVATION_ID_OUTPUT_FILE

# Loop 10,000 times
for i in {1..1000}
do
  # Call the command and get the output
  result=$(./curltime "${API_URL}?seed=$i")

   # Extract the activation ID and append to the relevant file
  activationId=$(echo "$result" | grep 'OpenWhisk Activation ID:' | awk -F': ' '{print $2}' | tr -d ' \r')
  echo $activationId >> $ACTIVATION_ID_OUTPUT_FILE

  # Extract the time value and append to the relevant file
  timeValue=$(echo "$result" | grep -E 'time_total:' | awk -F': ' '{print $2}' | tr -d ' ')
  echo $timeValue >> $TIME_OUTPUT_FILE

  # Optionally print progress
  echo "Iteration $i done"
done
