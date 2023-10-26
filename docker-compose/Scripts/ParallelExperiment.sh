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

# Loop
for i in $(seq 1 $ITERATIONS); do
    {
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
    timeValue=$(echo "$result" | grep -E 'time_total:' | awk -F': ' '{print $2}' | tr -d ' ')
    echo $timeValue >>$TIME_OUTPUT_FILE

    # Sleep for the desired duration to achieve the rate
    sleep $SLEEP_TIME

    end_time=$SECONDS
    elapsed_time=$(($end_time - $start_time))
    total_elapsed_time=$(($SECONDS))
    } &
done
