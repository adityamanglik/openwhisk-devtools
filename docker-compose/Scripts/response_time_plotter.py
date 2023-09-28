import matplotlib.pyplot as plt
import numpy as np
import sys


import numpy as np
import matplotlib.pyplot as plt

def plot_histogram(data, states, output_file):
    plt.figure(figsize=(12, 6))

    # Compute mean and std
    mean_val = np.mean(data)
    std_val = np.std(data)

    # Create histogram bins
    n, bins, patches = plt.hist(data, bins=50, color='blue', label='Warm Response Time', alpha=0.7)
    
    # Get the counts of cold starts in each bin
    cold_counts, _ = np.histogram([data[i] for i, state in enumerate(states) if state == "cold"], bins=bins)
    
    # Mark the cold starts at the top of the histogram bar with a number in red
    for count, rect, cold_count in zip(n, patches, cold_counts):
        height = count
        plt.text(rect.get_x() + rect.get_width() / 2, height + 5, str(int(cold_count)), ha='center', va='bottom', color='red')

    # Add text label for mean and std
    plt.text(0.85, 0.85, f"Mean: {mean_val:.2f}\nStd: {std_val:.2f}", transform=plt.gca().transAxes, ha="right", va="top",
             bbox=dict(boxstyle="round", facecolor="white", edgecolor="black"))

    # Add text label for mean and std
    plt.text(0.85, 0.85, f"Cold starts: {len([x for x in states if x == 'cold'])}\nTotal: {len(states)}", transform=plt.gca().transAxes, ha="right", va="top",
             bbox=dict(boxstyle="round", facecolor="white", edgecolor="black"))

    plt.title('Distribution of Response Time Over {} Iterations'.format(len(data)))
    plt.xlabel('Response Time (s)')
    plt.ylabel('Number of Activations')
    plt.grid(True, which='both', linestyle='--', linewidth=0.5)
    plt.legend()
    plt.savefig('../Graphs/'+ output_file.replace('.txt', '_histogram_plot.png'))
    # plt.show()


def plot_line(data, states, output_file):
    x = np.arange(1, len(data) + 1)

    plt.figure(figsize=(12, 6))
    plt.plot(x, data, label=output_file+' Response Time', linewidth=2, color='blue')
    cold_x = [x[i] for i, state in enumerate(states) if state == "cold"]
    cold_data = [data[i] for i, state in enumerate(states) if state == "cold"]
    plt.scatter(cold_x, cold_data, color='red', label='Cold Activation', s=50)
    plt.title('Response Time Over {} Iterations'.format(len(data)))
    plt.xlabel('Iteration')
    plt.ylabel('Response Time (s)')
    plt.grid(True)
    plt.legend()
    plt.savefig('../Graphs/'+ output_file.replace('.txt', '_line_plot.png'))
    # plt.show()

def plot_gc_line(data, output_file, ylabel):
    x = np.arange(1, len(data) + 1)
    plt.figure(figsize=(12, 6))
    plt.plot(x, data, label=output_file, linewidth=2, color='blue')
    plt.title(f'{ylabel} Over {len(data)} Iterations')
    plt.xlabel('Iteration')
    plt.ylabel(ylabel)
    plt.grid(True)
    plt.legend()
    plt.savefig('../Graphs/'+ output_file.replace('.txt', '_line_plot.png'))

if __name__ == '__main__':
    if len(sys.argv) < 6:
        print("Usage: python script_name.py latency_output_file.txt state_file.txt gcCollections.txt gcCollectionTime.txt gcTotalCollectors.txt")
        sys.exit(1)

    output_file = sys.argv[1]
    state_file = sys.argv[2]
    gc_collections_file = sys.argv[3]
    gc_collection_time_file = sys.argv[4]
    gc_total_collectors_file = sys.argv[5]

    with open(output_file, 'r') as f:
        latency_data = [float(line.strip()) for line in f.readlines()]

    with open(state_file, 'r') as f:
        states = [line.strip().split(': ')[1].replace('"', '') for line in f.readlines()]

    # Plotting the latency data
    plot_histogram(latency_data, states, output_file)
    plot_line(latency_data, states, output_file)

    # Loading and plotting the GC Collections
    with open(gc_collections_file, 'r') as f:
        gc_collections = [float(line.strip()) for line in f.readlines()]

    plot_gc_line(gc_collections, "gcCollections", "GC Collections")

    # Loading and plotting the GC Collection Time
    with open(gc_collection_time_file, 'r') as f:
        gc_collection_times = [float(line.strip()) for line in f.readlines()]

    plot_gc_line(gc_collection_times, "gcCollectionTime", "GC Collection Time (ms)")

    # Loading and plotting the Total GC Collectors
    with open(gc_total_collectors_file, 'r') as f:
        gc_total_collectors = [float(line.strip()) for line in f.readlines()]

    plot_gc_line(gc_total_collectors, "gcTotalCollectors", "Total GC Collectors")