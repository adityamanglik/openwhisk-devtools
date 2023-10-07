#!/bin/bash

# Check for required parameters
if [ "$#" -ne 2 ]; then
    echo "Usage: $0 <API_URL> <LANGUAGE>"
    exit 1
fi

API_URL=$1
LANGUAGE=$2
ITERATIONS=5000

# Functions

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

initialize_js_files() {
    USED_HEAP_SIZE_FILE="usedHeapSize.txt"
    TOTAL_HEAP_SIZE_FILE="totalHeapSize.txt"
    HEAP_SIZE_LIMIT_FILE="HeapSizeLimit.txt"

    > $USED_HEAP_SIZE_FILE
    > $TOTAL_HEAP_SIZE_FILE
    > $HEAP_SIZE_LIMIT_FILE
}

initialize_java_files() {
    GC1_COLLECTIONS_OUTPUT_FILE="gc1Collections.txt"
    GC1_COLLECTION_TIME_OUTPUT_FILE="gc1CollectionTime.txt"
    GC2_COLLECTIONS_OUTPUT_FILE="gc2Collections.txt"
    GC2_COLLECTION_TIME_OUTPUT_FILE="gc2CollectionTime.txt"

    > $GC1_COLLECTIONS_OUTPUT_FILE
    > $GC1_COLLECTION_TIME_OUTPUT_FILE
    > $GC2_COLLECTIONS_OUTPUT_FILE
    > $GC2_COLLECTION_TIME_OUTPUT_FILE
}

# Main script

# Create or empty the common output files
TIME_OUTPUT_FILE="${LANGUAGE}OutputTime.txt"
ACTIVATION_ID_OUTPUT_FILE="${LANGUAGE}activation_ids.txt"

> $TIME_OUTPUT_FILE
> $ACTIVATION_ID_OUTPUT_FILE

case "$LANGUAGE" in
    "JS")
        initialize_js_files
        ;;
    "Java")
        initialize_java_files
        ;;
    *)
        echo "Unsupported language: $LANGUAGE"
        exit 1
        ;;
esac

# Loop 
for i in $(seq 1 $ITERATIONS); do
  retry_counter=0
  max_retries=5
  success=0

  while [ $retry_counter -lt $max_retries ]; do
    result=$(./curltime "${API_URL}?seed=$i")

    # If good response, break out of inner loop
    if ! (echo "$result" | grep -q '"status":404' && echo "$result" | grep -q '"message":"Error: Not found."'); then
      success=1
      break
    fi

    echo "Attempt $(($retry_counter + 1)) failed for iteration $i. Resetting API and retrying..."
    ssh $OW_SERVER_NODE "cd $OW_DIRECTORY/Scripts/; source action_API_setup.sh"
    retry_counter=$(($retry_counter + 1))
  done

  # If after 5 retries we still haven't received a good response, quit or decide to continue
  if [ $success -eq 0 ]; then
    echo "\nError: Bad response received during iteration $i after $max_retries attempts."
    exit 1
  fi

  # Common extraction
  activationId=$(echo "$result" | grep 'OpenWhisk Activation ID:' | awk -F': ' '{print $2}' | tr -d ' \r')
  echo $activationId >> $ACTIVATION_ID_OUTPUT_FILE
  timeValue=$(echo "$result" | grep -E 'time_total:' | awk -F': ' '{print $2}' | tr -d ' ')
  echo $timeValue >> $TIME_OUTPUT_FILE

  # Language-specific extraction
  case "$LANGUAGE" in
      "JS")
          usedHeapSizeValue=$(echo "$result" | grep -Eo 'usedHeapSize: [0-9]+' | awk '{print $2}')
          totalHeapSizeValue=$(echo "$result" | grep -Eo 'totalHeapSize: [0-9]+' | awk '{print $2}')
          heapSizeLimitValue=$(echo "$result" | grep -Eo 'HeapSizeLimit: [0-9]+' | awk '{print $2}')

          echo $usedHeapSizeValue >> $USED_HEAP_SIZE_FILE
          echo $totalHeapSizeValue >> $TOTAL_HEAP_SIZE_FILE
          echo $heapSizeLimitValue >> $HEAP_SIZE_LIMIT_FILE
          ;;

      "Java")
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
          ;;
  esac

  # Print progress
  print_progress $i $ITERATIONS
done

echo # Print a newline after completion
