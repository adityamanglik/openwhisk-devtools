import numpy as np
import csv

def extract_values_from_csv(file_path):
    peak_requests_per_sec = 0.0
    p90, p99, median, aver = 0.0, 0.0, 0.0, 0.0

    with open(file_path, newline='') as csvfile:
        csv_reader = csv.reader(csvfile)
        next(csv_reader)  # Skip the header row
        row = []
        for row in csv_reader:
            if row[4] != 'N/A':  # Check if 'Requests/s' is available
                current_requests_per_sec = float(row[4])
                if current_requests_per_sec > peak_requests_per_sec:
                    peak_requests_per_sec = current_requests_per_sec
        p90 = float(row[10])
        p99 = float(row[13])
        median = float(row[6])
        aver = float(row[20])

    return peak_requests_per_sec, p90, p99, median, aver

def calculate_statistics(file_path):
    with open(file_path, 'r') as file:
        times = [int(line.strip()) for line in file if line.strip()]

    times = np.array(times)
    average = np.mean(times)
    median = np.median(times)
    p90 = np.percentile(times, 90)
    p99 = np.percentile(times, 99)

    return average, median, p90, p99

# Java data extraction
file_name = f"/users/am_CU/openwhisk-devtools/docker-compose/Graphs/LoadTesting/Java/LoadLatencyCurve.csv"
peak_requests_per_sec, p90, p99, median, aver = extract_values_from_csv(file_name)
print("Java -->")
print(f"Peak Requests/Sec: {peak_requests_per_sec}\n"
      f"Average Response Time: {aver} ms\n"
      f"Median Response Time: {median} ms\n"
      f"P90 Response Time: {p90} ms\n"
      f"P99 Response Time: {p99} ms\n")

# Path to the FunctionTime.txt file
function_time_file = f"/users/am_CU/openwhisk-devtools/docker-compose/Graphs/LoadTesting/Java/FunctionTime.txt"

# Calculate and print the statistics
average, median, p90, p99 = calculate_statistics(function_time_file)
print("Function Time Statistics:")
print(f"Average Response Time: {average} ns")
print(f"Median Response Time: {median} ns")
print(f"P90 Response Time: {p90} ns")
print(f"P99 Response Time: {p99} ns")

# Go data extraction
file_name = f"/users/am_CU/openwhisk-devtools/docker-compose/Graphs/LoadTesting/Go/LoadLatencyCurve.csv"
peak_requests_per_sec, p90, p99, median, aver = extract_values_from_csv(file_name)
print("\nGo -->")
print(f"Peak Requests/Sec: {peak_requests_per_sec}\n"
      f"Average Response Time: {aver} ms\n"
      f"Median Response Time: {median} ms\n"
      f"P90 Response Time: {p90} ms\n"
      f"P99 Response Time: {p99} ms\n")

# Path to the FunctionTime.txt file
function_time_file = f"/users/am_CU/openwhisk-devtools/docker-compose/Graphs/LoadTesting/Go/FunctionTime.txt"

# Calculate and print the statistics
average, median, p90, p99 = calculate_statistics(function_time_file)
print("Function Time Statistics:")
print(f"Average Response Time: {average} ns")
print(f"Median Response Time: {median} ns")
print(f"P90 Response Time: {p90} ns")
print(f"P99 Response Time: {p99} ns")

