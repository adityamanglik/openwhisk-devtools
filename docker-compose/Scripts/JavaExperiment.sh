#!/bin/bash

# Create or empty the output files
> JavaoutputTime.txt
> JavaoutputStartState.txt

# Loop 10,000 times
for i in {1..1000}
do
  # Call the command and get the output
  # result=$(./curltime "http://128.110.96.69:9090/api/23bc46b1-71f6-4ed5-8c54-816aa4f8c502/hello/world?seed=$i")
  result=$(./curltime "http://128.110.96.69:9090/api/23bc46b1-71f6-4ed5-8c54-816aa4f8c502/helloJava/world?seed=$i")

  # Extract the startState value and append to the relevant file
  startState=$(echo "$result" | grep "startState" | awk -F': ' '{print $2}')
  echo $startState >> JavaoutputStartState.txt

  # Extract the time value and append to the relevant file
  timeValue=$(echo "$result" | tail -n 1)
  echo $timeValue >> JavaoutputTime.txt

  # Optionally print progress
  echo "Iteration $i done"
done
