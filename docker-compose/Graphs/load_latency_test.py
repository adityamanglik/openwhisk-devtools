import matplotlib.pyplot as plt
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

# Java data extraction
file_name = f"/users/am_CU/openwhisk-devtools/docker-compose/Graphs/LoadTesting/Java/LoadLatencyCurve.csv"
peak_requests_per_sec, p90, p99, median, aver = extract_values_from_csv(file_name)
print("Java -->")
print(f"Peak Requests/Sec: {peak_requests_per_sec}\n"
      f"Average Response Time: {aver} ms\n"
      f"Median Response Time: {median} ms\n"
      f"P90 Response Time: {p90} ms\n"
      f"P99 Response Time: {p99} ms\n")

# Go data extraction
file_name = f"/users/am_CU/openwhisk-devtools/docker-compose/Graphs/LoadTesting/Go/LoadLatencyCurve.csv"
peak_requests_per_sec, p90, p99, median, aver = extract_values_from_csv(file_name)
print("Go -->")
print(f"Peak Requests/Sec: {peak_requests_per_sec}\n"
      f"Average Response Time: {aver} ms\n"
      f"Median Response Time: {median} ms\n"
      f"P90 Response Time: {p90} ms\n"
      f"P99 Response Time: {p99} ms\n")

