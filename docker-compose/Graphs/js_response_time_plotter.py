import matplotlib.pyplot as plt
import numpy as np
import os
import sys
from matplotlib.widgets import CheckButtons

ITERATIONS = 5000

def read_data(file_name):
    """Utility function to read data from a file."""
    with open(file_name, 'r') as f:
        return [float(line.strip()) for line in f.readlines()]

def plot_js_memory_stats(input_size, used_heap, total_heap, heap_limit):
    x = np.arange(1, len(used_heap) + 1)

    fig, ax1 = plt.subplots(figsize=(12, 6))

    # Plotting Used Heap and Total Heap on the left y-axis
    ax1.plot(x, used_heap, label='Used Heap Size', color='blue')
    ax1.plot(x, total_heap, label='Total Heap Size', color='darkblue')
    ax1.set_xlabel('Iteration')
    ax1.set_ylabel('Memory (bytes)')
    ax1.grid(True)
    
    # Instantiate a second y-axis that shares the same x-axis for Heap Size Limit
    ax2 = ax1.twinx()
    ax2.plot(x, heap_limit, label='Heap Size Limit', color='red')
    ax2.set_ylabel('Memory (bytes)')

    # Combined legend for both y-axes
    lines, labels = ax1.get_legend_handles_labels()
    lines2, labels2 = ax2.get_legend_handles_labels()
    ax2.legend(lines + lines2, labels + labels2, loc='upper left')

    plt.title('JavaScript Memory Metrics Over {} Iterations'.format(len(used_heap)))
    plt.savefig(f'../Graphs/JS/{input_size}/'+ 'memory_stats_plot.pdf')
    # plt.show()

def plot_histogram(input_size, data, states):
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
    plt.yscale('symlog')
    plt.grid(True, which='both', linestyle='--', linewidth=0.5)
    plt.legend()
    plt.savefig(f'../Graphs/JS/{input_size}/'+ 'histogram_plot.pdf')
    # plt.show()


def plot_line_orig(input_size, data, states):
    x = np.arange(1, len(data) + 1)

    plt.figure(figsize=(12, 6))
    plt.plot(x, data, label='Response Time', linewidth=2, color='blue')
    cold_x = [x[i] for i, state in enumerate(states) if state == "cold"]
    cold_data = [data[i] for i, state in enumerate(states) if state == "cold"]
    plt.scatter(cold_x, cold_data, color='red', label='Cold Activation', s=50)
    plt.title('Response Time Over {} Iterations'.format(len(data)))
    plt.xlabel('Iteration')
    plt.ylabel('Response Time (s)')
    plt.grid(True)
    plt.yscale('symlog')
    plt.legend()
    plt.savefig(f'../Graphs/JS/{input_size}/'+ 'line_plot.pdf')
    # plt.show()

def plot_line(input_size, data, states, gc1_collections, gc1_collection_times, gc2_collections, gc2_collection_times):
    x = np.arange(1, len(data) + 1)

    fig, ax1 = plt.figure(figsize=(12, 6)), plt.gca()
    
    ax1.set_xlabel('Iteration')
    ax1.set_ylabel('Response Time (s)', color='blue')
    ax1.plot(x, data, label='Response Time', linewidth=2, color='blue')
    cold_x = [x[i] for i, state in enumerate(states) if state == "cold"]
    cold_data = [data[i] for i, state in enumerate(states) if state == "cold"]
    ax1.scatter(cold_x, cold_data, color='red', label='Cold Activation', s=50)
    ax1.tick_params(axis='y', labelcolor='blue')

    ax2 = ax1.twinx()  # instantiate a second axes that shares the same x-axis
    ax2.set_ylabel('GC Metrics', color='red')  # we already handled the x-label with ax
    ax2.plot(x, gc1_collections, label='GC1 Collections', linewidth=1, color='darkred')
    ax2.plot(x, gc1_collection_times, label='GC1 Collection Times', linewidth=1, color='red')
    ax2.plot(x, gc2_collections, label='GC2 Collections', linewidth=1, color='lightcoral')
    ax2.plot(x, gc2_collection_times, label='GC2 Collection Times', linewidth=1, color='pink')
    ax2.tick_params(axis='y', labelcolor='red')
    ax2.set_yscale('symlog')

    fig.tight_layout()  # otherwise the right y-label is slightly clipped
    plt.grid(True)
    fig.legend(loc='center right')
    plt.title('Response Time and GC Metrics Over {} Iterations'.format(len(data)))
    plt.savefig(f'../Graphs/JS/{input_size}/'+ 'combined_line_plot.pdf')


def plot_gc_stats(input_size, gc1_collections, gc1_collection_times, gc2_collections, gc2_collection_times):
    x = np.arange(1, len(gc1_collections) + 1)

    fig, ax1 = plt.subplots(figsize=(12, 6))

    # Plotting GC1 and GC2 Collections on the left y-axis
    ax1.plot(x, gc1_collections, label='GC1 Collections', color='darkred')
    ax1.plot(x, gc2_collections, label='GC2 Collections', color='lightcoral')
    ax1.set_xlabel('Iteration')
    ax1.set_ylabel('GC Object Collections')
    ax1.grid(True)
    
    # Instantiate a second y-axis that shares the same x-axis
    ax2 = ax1.twinx()
    ax2.plot(x, gc1_collection_times, label='GC1 Collection Time', color='red')
    ax2.plot(x, gc2_collection_times, label='GC2 Collection Time', color='pink')
    ax2.set_ylabel('GC Collection Time')

    # Combined legend for both y-axes
    lines, labels = ax1.get_legend_handles_labels()
    lines2, labels2 = ax2.get_legend_handles_labels()
    ax2.legend(lines + lines2, labels + labels2, loc='upper left')

    plt.title('GC Metrics Over {} Iterations'.format(len(gc1_collections)))
    plt.savefig(f'../Graphs/JS/{input_size}/'+ 'gc_stats_plot.pdf')
    # plt.show()

def plot_js_metrics(input_size, data, states, used_heap, total_heap, heap_limit):
    x = np.arange(1, len(data) + 1)
    fig, ax1 = plt.subplots()  # Removed the figsize argument to let matplotlib auto-size
    
    # Response Time
    l1, = ax1.plot(x, data, label='Response Time', linewidth=2, color='blue', marker='o')
    cold_x = [x[i] for i, state in enumerate(states) if state == "cold"]
    cold_data = [data[i] for i, state in enumerate(states) if state == "cold"]
    ax1.scatter(cold_x, cold_data, color='red', label='Cold Activation', s=50)
    ax1.set_xlabel('Iteration')
    ax1.set_ylabel('Response Time (s)', color='blue')
    ax1.tick_params(axis='y', labelcolor='blue')
    ax1.grid(True)

    # Memory statistics
    ax2 = ax1.twinx()
    l2, = ax2.plot(x, used_heap, label='Used Heap Size', color='darkgreen')
    l3, = ax2.plot(x, total_heap, label='Total Heap Size', color='green')
    l4, = ax2.plot(x, heap_limit, label='Heap Size Limit', color='red')
    ax2.set_ylabel('Memory (bytes)', color='green')
    ax2.tick_params(axis='y', labelcolor='green')
    
    fig.tight_layout()
    
    plt.title('JavaScript Response Time and Memory Metrics Over {} Iterations'.format(len(data)))
    plt.savefig(f'../Graphs/JS/{input_size}/' + 'combined_js_metrics_plot.pdf')
    plt.show()



if __name__ == '__main__':
    # Check the number of arguments
    if len(sys.argv) != 2:
        print("Usage: python script_name.py <size>")
        sys.exit(1)

    os.chdir('/users/am_CU/openwhisk-devtools/docker-compose/Graphs/')

    # Get size from the command line and set up directories
    input_size = sys.argv[1]
    directory_path = f'../Graphs/JS/{input_size}/'
    if not os.path.exists(directory_path):
        os.makedirs(directory_path)

    # Base file path
    base_path = os.getcwd()
    base_path = os.path.join(base_path, f'JS/{input_size}/')

    # Reading data
    latency_data = read_data(os.path.join(base_path, "JSOutputTime.txt"))
    states = [line.split(': ')[1].replace('"', '').strip() for line in open(os.path.join(base_path, "JSactivation_ids.txt_startStates.txt"))]
    used_heap = read_data(os.path.join(base_path, "usedHeapSize.txt"))
    total_heap = read_data(os.path.join(base_path, "totalHeapSize.txt"))
    heap_limit = read_data(os.path.join(base_path, "HeapSizeLimit.txt"))

    latency_data = latency_data[:ITERATIONS]
    states = states[:ITERATIONS]
    used_heap = used_heap[:ITERATIONS]
    total_heap = total_heap[:ITERATIONS]
    heap_limit = heap_limit[:ITERATIONS]

    # Plotting functions
    plot_line_orig(input_size, latency_data, states)
    plot_histogram(input_size, latency_data, states)
    plot_js_memory_stats(input_size, used_heap, total_heap, heap_limit)
    plot_js_metrics(input_size, latency_data, states, used_heap, total_heap, heap_limit)