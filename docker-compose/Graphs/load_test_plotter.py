import matplotlib.pyplot as plt
import numpy as np

MaxGCPauseMillis_values = [50, 100, 150, 200, 250, 300]
Xmx_values = ["64m", "128m", "256m", "512m", "1g", "2g", "4g"]

colors = ['blue', 'green', 'red', 'cyan', 'yellow', 'purple', 'black']

def plot_values(latency_values, title, filename, ylabel="Latency"):
    plt.figure(figsize=(12, 6))
    for idx, xmx in enumerate(Xmx_values):
        plt.plot(MaxGCPauseMillis_values, latency_values[idx], color=colors[idx], label=f"Xmx {xmx}")
    plt.xticks(MaxGCPauseMillis_values)
    plt.xlabel("MaxGCPauseMillis")
    plt.ylabel(ylabel)
    plt.title(title)
    plt.legend()
    plt.tight_layout()
    plt.savefig(f'../Graphs/LoadTesting/' + filename + '.pdf')
    plt.show()

medians = [[] for _ in Xmx_values]
p90s = [[] for _ in Xmx_values]
p99s = [[] for _ in Xmx_values]

for idx, xmx in enumerate(Xmx_values):
    for max_gc in MaxGCPauseMillis_values:
        file_name = f"../Graphs/LoadTesting/Time_Xmx{xmx}_MaxGCPauseMillis{max_gc}.txt"
        with open(file_name, 'r') as f:
            latencies = [float(line.strip()) for line in f]
            medians[idx].append(np.median(latencies))
            p90s[idx].append(np.percentile(latencies, 90))
            p99s[idx].append(np.percentile(latencies, 99))
print(medians, p90s, p99s)
plot_values(medians, 'Median Latencies', 'median')
plot_values(p90s, 'P90 Latencies', 'p90')
plot_values(p99s, 'P99 Latencies', 'p99')
