#!/bin/bash

# Check for required parameters
if [ "$#" -ne 3 ]; then
    echo "Usage: $0 <API_URL> <LANGUAGE> <ITERATIONS>"
    exit 1
fi

API_URL=$1
LANGUAGE=$2
ITERATIONS=$3

# Functions

print_progress() {
    local current=$1
    local total=$2
    local elapsed=$3
    local width=50
    local progress=$(( ($current * $width) / $total ))
    local remaining=$(( $width - $progress ))
    printf "\r["
    printf "%${progress}s" | tr ' ' '#'
    printf "%${remaining}s" ' ' 
    printf "] (%d/%d) Elapsed: %.2fs" $current $total $elapsed
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
    
    HEAP_COMMITTED_MEMORY_FILE="heapCommittedMemory.txt"
    HEAP_INIT_MEMORY_FILE="heapInitMemory.txt"
    HEAP_MAX_MEMORY_FILE="heapMaxMemory.txt"
    HEAP_USED_MEMORY_FILE="heapUsedMemory.txt"

    > $GC1_COLLECTIONS_OUTPUT_FILE
    > $GC1_COLLECTION_TIME_OUTPUT_FILE
    > $GC2_COLLECTIONS_OUTPUT_FILE
    > $GC2_COLLECTION_TIME_OUTPUT_FILE

    > $HEAP_COMMITTED_MEMORY_FILE
    > $HEAP_INIT_MEMORY_FILE
    > $HEAP_MAX_MEMORY_FILE
    > $HEAP_USED_MEMORY_FILE
}

initialize_go_files() {
    HEAP_ALLOC_MEMORY_FILE="heapAllocMemory.txt"
    HEAP_IDLE_MEMORY_FILE="heapIdleMemory.txt"
    HEAP_INUSE_MEMORY_FILE="heapInuseMemory.txt"
    HEAP_OBJECTS_FILE="heapObjects.txt"
    HEAP_RELEASED_MEMORY_FILE="heapReleasedMemory.txt"
    HEAP_SYS_MEMORY_FILE="heapSysMemory.txt"
    SUM_FILE="sum.txt"

    > $HEAP_ALLOC_MEMORY_FILE
    > $HEAP_IDLE_MEMORY_FILE
    > $HEAP_INUSE_MEMORY_FILE
    > $HEAP_OBJECTS_FILE
    > $HEAP_RELEASED_MEMORY_FILE
    > $HEAP_SYS_MEMORY_FILE
    > $SUM_FILE
}

initialize_native_java_files() {
    GC1_COLLECTIONS_OUTPUT_FILE="gc1NativeCollections.txt"
    GC1_COLLECTION_TIME_OUTPUT_FILE="gc1NativeCollectionTime.txt"
    GC2_COLLECTIONS_OUTPUT_FILE="gc2NativeCollections.txt"
    GC2_COLLECTION_TIME_OUTPUT_FILE="gc2NativeCollectionTime.txt"
    
    HEAP_COMMITTED_MEMORY_FILE="nativeHeapCommittedMemory.txt"
    HEAP_INIT_MEMORY_FILE="nativeHeapInitMemory.txt"
    HEAP_MAX_MEMORY_FILE="nativeHeapMaxMemory.txt"
    HEAP_USED_MEMORY_FILE="nativeHeapUsedMemory.txt"

    > $GC1_COLLECTIONS_OUTPUT_FILE
    > $GC1_COLLECTION_TIME_OUTPUT_FILE
    > $GC2_COLLECTIONS_OUTPUT_FILE
    > $GC2_COLLECTION_TIME_OUTPUT_FILE

    > $HEAP_COMMITTED_MEMORY_FILE
    > $HEAP_INIT_MEMORY_FILE
    > $HEAP_MAX_MEMORY_FILE
    > $HEAP_USED_MEMORY_FILE
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
    "NativeJava")
        initialize_native_java_files
        ;;
    "Go")
        initialize_go_files
        ;;
    *)
        echo "Unsupported language: $LANGUAGE"
        exit 1
        ;;
esac

# Loop 
for i in $(seq 1 $ITERATIONS); do
  start_time=$SECONDS
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
    ssh $OW_SERVER_NODE "cd $OW_DIRECTORY/; source Scripts/action_API_setup.sh"
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
          # Extract the heapCommittedMemory value and append to the relevant file
          heapCommittedMemoryValue=$(echo "$result" | grep -Eo '"heapCommittedMemory: ": [0-9]+' | awk '{print $3}')
          echo $heapCommittedMemoryValue >> heapCommittedMemory.txt

          # Extract the heapInitMemory value and append to the relevant file
          heapInitMemoryValue=$(echo "$result" | grep -Eo '"heapInitMemory: ": [0-9]+' | awk '{print $3}')
          echo $heapInitMemoryValue >> heapInitMemory.txt

          # Extract the heapMaxMemory value and append to the relevant file
          heapMaxMemoryValue=$(echo "$result" | grep -Eo '"heapMaxMemory: ": [0-9]+' | awk '{print $3}')
          echo $heapMaxMemoryValue >> heapMaxMemory.txt

          # Extract the heapUsedMemory value and append to the relevant file
          heapUsedMemoryValue=$(echo "$result" | grep -Eo '"heapUsedMemory: ": [0-9]+' | awk '{print $3}')
          echo $heapUsedMemoryValue >> heapUsedMemory.txt

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
          
      "Go")
        heapAllocMemoryValue=$(echo "$result" | grep -Eo '"heapAllocMemory": [0-9]+' | awk '{print $2}')
        echo $heapAllocMemoryValue >> $HEAP_ALLOC_MEMORY_FILE

        heapIdleMemoryValue=$(echo "$result" | grep -Eo '"heapIdleMemory": [0-9]+' | awk '{print $2}')
        echo $heapIdleMemoryValue >> $HEAP_IDLE_MEMORY_FILE

        heapInuseMemoryValue=$(echo "$result" | grep -Eo '"heapInuseMemory": [0-9]+' | awk '{print $2}')
        echo $heapInuseMemoryValue >> $HEAP_INUSE_MEMORY_FILE

        heapObjectsValue=$(echo "$result" | grep -Eo '"heapObjects": [0-9]+' | awk '{print $2}')
        echo $heapObjectsValue >> $HEAP_OBJECTS_FILE

        heapReleasedMemoryValue=$(echo "$result" | grep -Eo '"heapReleasedMemory": [0-9]+' | awk '{print $2}')
        echo $heapReleasedMemoryValue >> $HEAP_RELEASED_MEMORY_FILE

        heapSysMemoryValue=$(echo "$result" | grep -Eo '"heapSysMemory": [0-9]+' | awk '{print $2}')
        echo $heapSysMemoryValue >> $HEAP_SYS_MEMORY_FILE

        sumValue=$(echo "$result" | grep -Eo '"sum": [0-9]+' | awk '{print $2}')
        echo $sumValue >> $SUM_FILE
        ;;
      "NativeJava")
        
# Extraction and writing to the files
        heapCommittedMemoryValue=$(echo "$result" | sed -n 's/.*"heapCommittedMemory: ":[ \t]*\([0-9]*\).*/\1/p')
        echo $heapCommittedMemoryValue >> nativeHeapCommittedMemory.txt

        heapInitMemoryValue=$(echo "$result" | sed -n 's/.*"heapInitMemory: ":[ \t]*\([0-9]*\).*/\1/p')
        echo $heapInitMemoryValue >> nativeHeapInitMemory.txt

        heapMaxMemoryValue=$(echo "$result" | sed -n 's/.*"heapMaxMemory: ":[ \t]*\([0-9]*\).*/\1/p')
        echo $heapMaxMemoryValue >> nativeHeapMaxMemory.txt

        heapUsedMemoryValue=$(echo "$result" | sed -n 's/.*"heapUsedMemory: ":[ \t]*\([0-9]*\).*/\1/p')
        echo $heapUsedMemoryValue >> nativeHeapUsedMemory.txt

        gc1CollectionsValue=$(echo "$result" | awk -F'"gc1CollectionCount":' '{print $2}' | awk -F, '{print $1}')
        echo $gc1CollectionsValue >> gc1NativeCollections.txt

        gc1CollectionTimeValue=$(echo "$result" | awk -F'"gc1CollectionTime":' '{print $2}' | awk -F, '{print $1}')
        echo $gc1CollectionTimeValue >> gc1NativeCollectionTime.txt

        gc2CollectionsValue=$(echo "$result" | awk -F'"gc2CollectionCount":' '{print $2}' | awk -F, '{print $1}')
        echo $gc2CollectionsValue >> gc2NativeCollections.txt

        gc2CollectionTimeValue=$(echo "$result" | awk -F'"gc2CollectionTime":' '{print $2}' | awk -F, '{print $1}')
        echo $gc2CollectionTimeValue >> gc2NativeCollectionTime.txt
        ;;
  esac

  end_time=$SECONDS
  elapsed_time=$(($end_time-$start_time))
  total_elapsed_time=$(($SECONDS))

  # Print progress
  print_progress $i $ITERATIONS $total_elapsed_time
done

echo # Print a newline after completion
