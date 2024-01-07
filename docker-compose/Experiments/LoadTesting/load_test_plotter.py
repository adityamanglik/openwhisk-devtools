import matplotlib.pyplot as plt
import numpy as np
import csv

def extract_values_from_csv(file_path):
    with open(file_path, newline='') as csvfile:
        # Create a CSV reader object
        csv_reader = csv.reader(csvfile)
        last_row = None
        for row in csv_reader:
            last_row = row  # Keep overwriting until the last row is stored
        
        # Extracting the required values
        # Indices are determined by the column positions in your CSV
        p90 = float(last_row[10])
        p99 = float(last_row[13])
        median = float(last_row[6])  # This is assuming 'Total Median Response Time' is the median value you want
        aver = float(last_row[20])
        return p90, p99, median, aver

MaxGCPauseMillis_values = [50, 100, 150, 200, 250, 300]
Xmx_values = ["64m", "128m", "256m", "512m", "1g", "2g", "4g"]

colors = ['blue', 'green', 'red', 'cyan', 'yellow', 'purple', 'black']

def plot_values(latency_values, title, filename, ylabel="Latency"):
    plt.figure(figsize=(12, 6))
    for idx, xmx in enumerate(Xmx_values):
        plt.plot(MaxGCPauseMillis_values, latency_values[idx], color=colors[idx], label=f"Xmx {xmx}")
    plt.xticks(MaxGCPauseMillis_values)
    # plt.yticks([480, 510])
    # plt.yscale('symlog')
    plt.xlabel("MaxGCPauseMillis")
    plt.ylabel(ylabel)
    plt.title(title)
    plt.legend()
    plt.tight_layout()
    plt.savefig(f'../Graphs/LoadTesting/Java/' + filename + '.pdf')
    plt.show()

medians = [[] for _ in Xmx_values]
p90s = [[] for _ in Xmx_values]
p99s = [[] for _ in Xmx_values]
averages = [[] for _ in Xmx_values]
for idx, xmx in enumerate(Xmx_values):
    for max_gc in MaxGCPauseMillis_values:
        file_name = f"../Graphs/LoadTesting/Java/Time_Xmx{xmx}_MaxGCPauseMillis{max_gc}.csv"
        p90, p99, median, aver = extract_values_from_csv(file_name)
        medians[idx].append(median)
        p90s[idx].append(p90)
        p99s[idx].append(p99)
        averages[idx].append(aver)

# print average of each list
for aver in averages:
    print(sum(aver)/len(aver))

plot_values(medians, 'Median Latencies', 'median')
plot_values(p90s, 'P90 Latencies', 'p90')
plot_values(p99s, 'P99 Latencies', 'p99')
plot_values(averages, 'Averages', 'averages')
