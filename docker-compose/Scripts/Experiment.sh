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
  result=$(./curltime "${API_URL}?seed=$i")

  # Check for bad response and raise error if found
  if echo "$result" | grep -q '"status":404' && echo "$result" | grep -q '"message":"Error: Not found."'; then
    echo "\nError: Bad response received during iteration $i"
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
          # ... your Java-specific extraction logic here ...
          ;;
  esac

  # Print progress
  print_progress $i $ITERATIONS
done

echo # Print a newline after completion
