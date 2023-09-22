#!/bin/bash

# Create or empty the output file
> JSoutput.txt

# Loop 10,000 times
for i in {1..10000}
do
  # Call the command and get the output
  result=$(./curltime "http://128.110.96.69:9090/api/23bc46b1-71f6-4ed5-8c54-816aa4f8c502/hello/world?seed=$i" | tail -n 1)

  # Append the output to the file
  echo $result >> JSoutput.txt

  # Optionally print progress
  echo "Iteration $i done"
done

