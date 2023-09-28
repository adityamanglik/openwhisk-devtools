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
GC_COLLECTIONS_OUTPUT_FILE="gcCollections.txt"
GC_COLLECTION_TIME_OUTPUT_FILE="gcCollectionTime.txt"
GC_TOTAL_COLLECTORS_OUTPUT_FILE="gcTotalCollectors.txt"

> $TIME_OUTPUT_FILE
> $ACTIVATION_ID_OUTPUT_FILE
> $GC_COLLECTIONS_OUTPUT_FILE
> $GC_COLLECTION_TIME_OUTPUT_FILE
> $GC_TOTAL_COLLECTORS_OUTPUT_FILE

# Loop 10,000 times
for i in {1..1000}
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

  # Extract the gcTotalCollections value and append to the relevant file
  gcCollectionsValue=$(echo "$result" | grep -Eo '"gcTotalCollectionCount": [0-9]+' | awk '{print $2}')
  echo $gcCollectionsValue >> $GC_COLLECTIONS_OUTPUT_FILE

  # Extract the gcTotalCollectionTime value and append to the relevant file
  gcCollectionTimeValue=$(echo "$result" | grep -Eo '"gcTotalCollectionTime": [0-9]+' | awk '{print $2}')
  echo $gcCollectionTimeValue >> $GC_COLLECTION_TIME_OUTPUT_FILE

  # Extract the gcTotalCollectors value and append to the relevant file
  gcTotalCollectorsValue=$(echo "$result" | grep -Eo '"gcTotalCollectors": [0-9]+' | awk '{print $2}')
  echo $gcTotalCollectorsValue >> $GC_TOTAL_COLLECTORS_OUTPUT_FILE

  # Optionally print progress
  echo "Iteration $i done"
done
