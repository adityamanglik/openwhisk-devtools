#!/bin/bash

# Check for required parameters
if [ "$#" -ne 2 ]; then
    echo "Usage: $0 <API_URL> <LANGUAGE>"
    exit 1
fi

API_URL=$1
LANGUAGE=$2
ITERATIONS=5000

# Function to print the progress bar and iteration
print_progress() {
    local current=$1
    local total=$2
    local width=50
    local progress=$(( ($current * $width) / $total ))
    local remaining=$(( $width - $progress ))
    printf "\r["
    printf "%${progress}s" | tr ' ' '#'
    printf "%${remaining}s" ' ' 
    printf "] (%d/%d)" $current $total
}


# Create or empty the output files
TIME_OUTPUT_FILE="${LANGUAGE}OutputTime.txt"
ACTIVATION_ID_OUTPUT_FILE="${LANGUAGE}activation_ids.txt"

if [ "$LANGUAGE" == "JS" ]; then
USED_HEAP_SIZE_FILE="usedHeapSize.txt"
TOTAL_HEAP_SIZE_FILE="totalHeapSize.txt"
HEAP_SIZE_LIMIT_FILE="HeapSizeLimit.txt"

    > $USED_HEAP_SIZE_FILE
    > $TOTAL_HEAP_SIZE_FILE
    > $HEAP_SIZE_LIMIT_FILE
fi

if [ "$LANGUAGE" == "Java" ]; then
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
fi

# Loop 10,000 times
for i in $(seq 1 $ITERATIONS); do
  # Call the command and get the output
  result=$(./curltime "${API_URL}?seed=$i")

  # Check for bad response and raise error if found
  if echo "$result" | grep -q '"status":404' && echo "$result" | grep -q '"message":"Error: Not found."'; then
    echo "Error: Bad response received during iteration $i"
    # 
    exit 1
  fi

  # Extract the activation ID and append to the relevant file
  activationId=$(echo "$result" | grep 'OpenWhisk Activation ID:' | awk -F': ' '{print $2}' | tr -d ' \r')
  echo $activationId >> $ACTIVATION_ID_OUTPUT_FILE

  # Extract the time value and append to the relevant file
  timeValue=$(echo "$result" | grep -E 'time_total:' | awk -F': ' '{print $2}' | tr -d ' ')
  echo $timeValue >> $TIME_OUTPUT_FILE

  if [ "$LANGUAGE" == "Java" ]; then
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
  fi
  
  # If LANGUAGE is JS, extract the memory statistics and append to the relevant files
  if [ "$LANGUAGE" == "JS" ]; then
      usedHeapSizeValue=$(echo "$result" | grep -Eo 'usedHeapSize: [0-9]+' | awk '{print $2}')
      echo $usedHeapSizeValue >> $USED_HEAP_SIZE_FILE

      totalHeapSizeValue=$(echo "$result" | grep -Eo 'totalHeapSize: [0-9]+' | awk '{print $2}')
      echo $totalHeapSizeValue >> $TOTAL_HEAP_SIZE_FILE

      heapSizeLimitValue=$(echo "$result" | grep -Eo 'HeapSizeLimit: [0-9]+' | awk '{print $2}')
      echo $heapSizeLimitValue >> $HEAP_SIZE_LIMIT_FILE
  fi

    # Print progress
  print_progress $i $ITERATIONS
done
