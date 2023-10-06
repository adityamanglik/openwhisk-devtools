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
GC1_COLLECTIONS_OUTPUT_FILE="gc1Collections.txt"
GC1_COLLECTION_TIME_OUTPUT_FILE="gc1CollectionTime.txt"
GC2_COLLECTIONS_OUTPUT_FILE="gc2Collections.txt"
GC2_COLLECTION_TIME_OUTPUT_FILE="gc2CollectionTime.txt"

> $TIME_OUTPUT_FILE
> $ACTIVATION_ID_OUTPUT_FILE
> $GC1_COLLECTIONS_OUTPUT_FILE
> $GC1_COLLECTION_TIME_OUTPUT_FILE
> $GC2_COLLECTIONS_OUTPUT_FILE
> $GC2_COLLECTION_TIME_OUTPUT_FILE

# Loop 10,000 times
for i in {1..5000}
do
  # Call the command and get the output
  result=$(./curltime "${API_URL}?seed=$i")

  # Check for bad response and raise error if found
  if echo "$result" | grep -q '"status":404' && echo "$result" | grep -q '"message":"Error: Not found."'; then
    echo "Error: Bad response received during iteration $i"
    exit 1
  fi

  # Extract the activation ID and append to the relevant file
  activationId=$(echo "$result" | grep 'OpenWhisk Activation ID:' | awk -F': ' '{print $2}' | tr -d ' \r')
  echo $activationId >> $ACTIVATION_ID_OUTPUT_FILE

  # Extract the time value and append to the relevant file
  timeValue=$(echo "$result" | grep -E 'time_total:' | awk -F': ' '{print $2}' | tr -d ' ')
  echo $timeValue >> $TIME_OUTPUT_FILE

  # Extract the gc1CollectionCount value and append to the relevant file
  gc1CollectionsValue=$(echo "$result" | grep -Eo '"gc1CollectionCount": [0-9]+' | awk '{print $2}')
  echo $gc1CollectionsValue >> $GC1_COLLECTIONS_OUTPUT_FILE

  # Extract the gc1CollectionTime value and append to the relevant file
  gc1CollectionTimeValue=$(echo "$result" | grep -Eo '"gc1CollectionTime": [0-9]+' | awk '{print $2}')
  echo $gc1CollectionTimeValue >> $GC1_COLLECTION_TIME_OUTPUT_FILE

  # Extract the gc2CollectionCount value and append to the relevant file
  gc2CollectionsValue=$(echo "$result" | grep -Eo '"gc2CollectionCount": [0-9]+' | awk '{print $2}')
  echo $gc2CollectionsValue >> $GC2_COLLECTIONS_OUTPUT_FILE

  # Extract the gc2CollectionTime value and append to the relevant file
  gc2CollectionTimeValue=$(echo "$result" | grep -Eo '"gc2CollectionTime": [0-9]+' | awk '{print $2}')
  echo $gc2CollectionTimeValue >> $GC2_COLLECTION_TIME_OUTPUT_FILE

  # Optionally print progress
  echo "Iteration $i done"
done
