import subprocess
import json
import csv

def run_curl_command():
    # Run the curl command and capture the output
    output = subprocess.check_output(['./client'])
    # Decode the byte string to UTF-8
    output_str = output.decode('utf-8')
    # Parse the JSON response
    data = json.loads(output_str)
    return data

def main():
    # Number of times to execute the curl command
    N = 1000
    # Initialize a list to store the results
    results = []

    # Execute the curl command N times
    for _ in range(N):
        result = run_curl_command()
        results.append(result)

    # Write the results to a CSV file
    with open('curl_results.csv', 'w', newline='') as csvfile:
        fieldnames = ['state', 'exec_time', 'request_number', 'thp_status', 'nr_thp']
        writer = csv.DictWriter(csvfile, fieldnames=fieldnames)

        writer.writeheader()
        for result in results:
            writer.writerow(result)

if __name__ == "__main__":
    main()