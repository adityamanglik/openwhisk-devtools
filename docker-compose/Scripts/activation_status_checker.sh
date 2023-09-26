#!/bin/bash

# Define the OpenWhisk CLI command with the configuration file
WSK_CLI="WSK_CONFIG_FILE=./.wskprops ./openwhisk-src/bin/wsk -i"

# Assuming you have a file named activation_ids.txt with one activation ID per line
while read -r activation_id; do
  # Fetch the logs for the activation ID
  logs=$($WSK_CLI activation logs "$activation_id")

  # Check the logs for cold/warm status
  if echo "$logs" | grep -q "starting up"; then
    echo "$activation_id: cold"
  elif echo "$logs" | grep -q "already running"; then
    echo "$activation_id: warm"
  else
    echo "$activation_id: unknown"
  fi
done < activation_ids.txt
