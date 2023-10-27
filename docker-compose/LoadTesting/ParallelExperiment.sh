#!/bin/bash

# Check for required parameters
if [ "$#" -ne 4 ]; then
    echo "Usage: $0 <API_URL> <LANGUAGE> <ITERATIONS> <RATE>"
    exit 1
fi

API_URL=$1
LANGUAGE=$2
ITERATIONS=$3
RATE=$4
SLEEP_TIME=$(awk "BEGIN {print 1.0/$RATE}")

# Functions

# Create or empty the common output files
TIME_OUTPUT_FILE="${LANGUAGE}OutputTime.txt"

>$TIME_OUTPUT_FILE

# Traffic Injection Loop
for i in $(seq 1 $ITERATIONS); do
    {
    start_time=$SECONDS

    result=$(./curltime "${API_URL}?seed=$i")

    # Common extraction
    timeValue=$(echo "$result" | grep -E 'time_total:' | awk -F': ' '{print $2}' | tr -d ' ')
    echo $timeValue >>$TIME_OUTPUT_FILE

    end_time=$SECONDS
    elapsed_time=$(($end_time - $start_time))
    total_elapsed_time=$(($SECONDS))
    } &
    # Sleep for the desired duration to achieve the rate
    sleep $SLEEP_TIME
done

# Wait for all background processes to finish before script exits
wait