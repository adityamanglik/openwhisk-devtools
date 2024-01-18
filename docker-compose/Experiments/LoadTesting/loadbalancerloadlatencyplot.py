import matplotlib.pyplot as plt
import numpy as np
import csv
from collections import defaultdict

def extract_and_aggregate_values(file_path):
    # Initialize dictionaries to store latency values for each user count
    median_values, p95_values, p99_values, p999_values, p9999_values = defaultdict(list), defaultdict(list), defaultdict(list), defaultdict(list), defaultdict(list)

    with open(file_path, newline='') as csvfile:
        csv_reader = csv.reader(csvfile)
        for row in csv_reader:
            if row and row[0] != "Timestamp":  # Skip the header row
                user_count = int(row[1])
                if user_count == 0:
                    continue
                rps = float(row[4])
                if rps < (1.2*user_count):
                    continue
                # Append the latencies to corresponding user count keys
                median_values[user_count].append(float(row[6]))
                p95_values[user_count].append(float(row[11]))
                p99_values[user_count].append(float(row[13]))
                p999_values[user_count].append(float(row[14]))
                p9999_values[user_count].append(float(row[15]))

    # Calculate median for each user count
    for user_count in median_values:
        median_values[user_count] = np.median(median_values[user_count])
        p95_values[user_count] = np.median(p95_values[user_count])
        p99_values[user_count] = np.median(p99_values[user_count])
        p999_values[user_count] = np.median(p999_values[user_count])
        p9999_values[user_count] = np.median(p9999_values[user_count])

    return median_values, p95_values, p99_values, p999_values, p9999_values

def plot_values(values_dict, NOGC_values_dict, title):
    plt.figure(figsize=(10, 6))
    for user_count, latencies in sorted(values_dict.items()):
        NOGC_latencies = NOGC_values_dict[user_count]
        plt.plot(user_count, latencies, marker='o', color='r', label='EM')
        plt.plot(user_count, NOGC_latencies, marker='o', color='b', label='NOGC')
    plt.xlabel("Number of Users")
    plt.ylabel("Latency (ms)")
    plt.title(title)
    plt.legend()
    plt.grid(True)
    plt.tight_layout()
    plt.savefig(f'./Graphs/LoadTesting/Go/' + title + '.pdf')

# Replace 'your_file_path.csv' with the actual path of your CSV file
file_path = './EM_locust_stats_history.csv'
NOGC_file_path = './locust_stats_history.csv'

median_values, p95_values, p99_values, p999_values, p9999_values = extract_and_aggregate_values(file_path)
NOGC_median_values, NOGC_p95_values, NOGC_p99_values, NOGC_p999_values, NOGC_p9999_values = extract_and_aggregate_values(NOGC_file_path)

# Plotting the latencies
plot_values(median_values, NOGC_median_values, "Median")
plot_values(p95_values, NOGC_p95_values, "p95")
plot_values(p99_values, NOGC_p99_values, "p99")
plot_values(p999_values, NOGC_p999_values, "p999")
plot_values(p9999_values, NOGC_p9999_values, "p9999")
