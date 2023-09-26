#!/bin/bash

# Check if input file is provided
if [ "$#" -ne 1 ]; then
    echo "Usage: $0 <FILE_WITH_ACTIVATION_IDS>"
    exit 1
fi

INPUT_FILE=$1
OUTPUT_FILE="$1_startStates.txt"

# Empty or create the output file
> $OUTPUT_FILE

# Loop through each activation ID in the input file
while IFS= read -r activation_id; do
    # Fetch activation details
    activation_output=$(WSK_CONFIG_FILE=/users/am_CU/openwhisk-devtools/docker-compose/.wskprops /users/am_CU/openwhisk-devtools/docker-compose/openwhisk-src/bin/wsk -i activation get "$activation_id")
    
    # Extract startState (if exists)
    startState=$(echo "$activation_output" | grep "startState" | awk -F': ' '{print $2}' | tr -d ',' | tr -d ' ')

    # Check if startState is found and write to output file
    if [[ ! -z "$startState" ]]; then
        echo "$activation_id: $startState" >> $OUTPUT_FILE
    else
        echo "$activation_id: startState not found" >> $OUTPUT_FILE
    fi

done < "$INPUT_FILE"
