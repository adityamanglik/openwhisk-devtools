import matplotlib.pyplot as plt
import numpy as np
import os
import sys

def plot_js_memory_stats(used_heap, total_heap, heap_limit, output_file):
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

def plot_line(input_size, data, states, gc1_collections, gc1_collection_times, gc2_collections, gc2_collection_times, output_file):
    x = np.arange(1, len(data) + 1)

    fig, ax1 = plt.figure(figsize=(12, 6)), plt.gca()
    
    ax1.set_xlabel('Iteration')
    ax1.set_ylabel('Response Time (s)', color='blue')
    ax1.plot(x, data, label=output_file+' Response Time', linewidth=2, color='blue')
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


def plot_gc_stats(input_size, gc1_collections, gc1_collection_times, gc2_collections, gc2_collection_times, output_file):
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


if __name__ == '__main__':
        # Check if the right number of arguments is provided
    if len(sys.argv) != 2:
        print("Usage: python script_name.py <size>")
        sys.exit(1)
        
    # Get size from command line
    input_size = int(sys.argv[1])

    output_file = os.path.join(os.getcwd(), "JSOutputTime.txt")
    state_file = os.path.join(os.getcwd(), "JSactivation_ids.txt_startStates.txt")
    used_heap_file = os.path.join(os.getcwd(), "usedHeapSize.txt")
    total_heap_file = os.path.join(os.getcwd(), "totalHeapSize.txt")
    heap_limit_file = os.path.join(os.getcwd(), "HeapSizeLimit.txt")


    with open(output_file, 'r') as f:
        latency_data = [float(line.strip()) for line in f.readlines()]

    with open(state_file, 'r') as f:
        states = [line.strip().split(': ')[1].replace('"', '') for line in f.readlines()]

    # Loading JS memory data
    with open(used_heap_file, 'r') as f:
        used_heap = [float(line.strip()) for line in f.readlines()]

    with open(total_heap_file, 'r') as f:
        total_heap = [float(line.strip()) for line in f.readlines()]

    with open(heap_limit_file, 'r') as f:
        heap_limit = [float(line.strip()) for line in f.readlines()]

    plot_line_orig(input_size, latency_data, states, output_file)
    plot_histogram(input_size, latency_data, states, output_file)
    plot_js_memory_stats(input_size, used_heap, total_heap, heap_limit, output_file)