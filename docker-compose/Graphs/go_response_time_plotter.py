import matplotlib.pyplot as plt
import numpy as np
import os
import sys

ITERATIONS = 5000

def read_data(file_name):
    """Utility function to read data from a file."""
    with open(file_name, 'r') as f:
        return [float(line.strip()) for line in f.readlines() if line.strip()]

def plot_go_memory_stats(input_size, heap_alloc, heap_idle, heap_inuse, heap_objects, heap_released, heap_sys):
    x = np.arange(1, len(heap_alloc) + 1)

    plt.figure(figsize=(12, 6))
    
    plt.plot(x, heap_alloc, label='Heap Alloc Memory', color='blue')
    plt.plot(x, heap_idle, label='Heap Idle Memory', color='green')
    plt.plot(x, heap_inuse, label='Heap In-use Memory', color='cyan')
    plt.plot(x, heap_objects, label='Heap Objects Value', color='purple')
    plt.plot(x, heap_released, label='Heap Released Memory', color='magenta')
    plt.plot(x, heap_sys, label='Heap Sys Memory', color='red')
    
    plt.xlabel('Iteration')
    plt.ylabel('Memory (bytes)')
    plt.grid(True)
    plt.legend()
    plt.title('Go Memory Metrics Over {} Iterations'.format(len(heap_alloc)))
    plt.savefig(f'../Graphs/Go/{input_size}/' + 'memory_stats_plot.pdf')
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
    plt.savefig(f'../Graphs/Go/{input_size}/'+ 'histogram_plot.pdf')
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
    plt.savefig(f'../Graphs/Go/{input_size}/'+ 'line_plot.pdf')
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
    plt.savefig(f'../Graphs/Go/{input_size}/'+ 'combined_line_plot.pdf')


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
    plt.savefig(f'../Graphs/Go/{input_size}/'+ 'gc_stats_plot.pdf')
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
    plt.savefig(f'../Graphs/Go/{input_size}/' + 'combined_Go_metrics_plot.pdf')
    plt.show()



if __name__ == '__main__':
    if len(sys.argv) != 2:
        print("Usage: python script_name.py <size>")
        sys.exit(1)

    os.chdir('/users/am_CU/openwhisk-devtools/docker-compose/Graphs/')

    input_size = sys.argv[1]
    directory_path = f'../Graphs/Go/{input_size}/'
    if not os.path.exists(directory_path):
        os.makedirs(directory_path)

    base_path = os.getcwd()
    base_path = os.path.join(base_path, f'Go/{input_size}/')

    latency_data = read_data(os.path.join(base_path, "GoOutputTime.txt"))
    states = [line.split(': ')[1].replace('"', '').strip() for line in open(os.path.join(base_path, "Goactivation_ids.txt_startStates.txt"))]
    
    heap_alloc = read_data(os.path.join(base_path, "heapAllocMemory.txt"))
    heap_idle = read_data(os.path.join(base_path, "heapIdleMemory.txt"))
    heap_inuse = read_data(os.path.join(base_path, "heapInuseMemory.txt"))
    heap_objects = read_data(os.path.join(base_path, "heapObjects.txt"))
    heap_released = read_data(os.path.join(base_path, "heapReleasedMemory.txt"))
    heap_sys = read_data(os.path.join(base_path, "heapSysMemory.txt"))

    # Truncate datasets to ITERATIONS
    latency_data = latency_data[:ITERATIONS]
    states = states[:ITERATIONS]
    heap_alloc = heap_alloc[:ITERATIONS]
    heap_idle = heap_idle[:ITERATIONS]
    heap_inuse = heap_inuse[:ITERATIONS]
    heap_objects = heap_objects[:ITERATIONS]
    heap_released = heap_released[:ITERATIONS]
    heap_sys = heap_sys[:ITERATIONS]

    # Plotting functions
    plot_line_orig(input_size, latency_data, states)
    plot_histogram(input_size, latency_data, states)
    plot_go_memory_stats(input_size, heap_alloc, heap_idle, heap_inuse, heap_objects, heap_released, heap_sys)