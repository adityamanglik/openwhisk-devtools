import matplotlib.pyplot as plt
import numpy as np
import os
import sys

def plot_heap_stats(input_size, heap_committed_memory, heap_init_memory, heap_max_memory, heap_used_memory):
    x = np.arange(1, len(heap_committed_memory) + 1)

    fig, ax1 = plt.subplots(figsize=(12, 6))

    # Plotting heap committed and heap used memory on the left y-axis
    ax1.plot(x, heap_committed_memory, label='Heap Committed Memory', color='darkgreen')
    ax1.plot(x, heap_used_memory, label='Heap Used Memory', color='lightgreen')
    ax1.set_xlabel('Iteration')
    ax1.set_ylabel('Heap Memory (Bytes)')
    ax1.grid(True)
    
    # Instantiate a second y-axis that shares the same x-axis
    ax2 = ax1.twinx()
    ax2.plot(x, heap_init_memory, label='Heap Init Memory', color='green')
    ax2.plot(x, heap_max_memory, label='Heap Max Memory', color='lime')
    ax2.set_ylabel('Heap Memory (Bytes)')

    # Combined legend for both y-axes
    lines, labels = ax1.get_legend_handles_labels()
    lines2, labels2 = ax2.get_legend_handles_labels()
    ax2.legend(lines + lines2, labels + labels2, loc='upper left')

    plt.title('Heap Memory Metrics Over {} Iterations'.format(len(heap_committed_memory)))
    plt.savefig(f'../Graphs/Java/{input_size}/'+ 'heap_stats_plot.pdf')
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
    plt.savefig(f'../Graphs/Java/{input_size}/'+ 'histogram_plot.pdf')
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
    plt.yscale('symlog')
    plt.grid(True)
    plt.legend()
    plt.savefig(f'../Graphs/Java/{input_size}/'+ 'line_plot.pdf')
    # plt.show()

def plot_line(input_size, data, states, gc1_collections, gc1_collection_times, gc2_collections, gc2_collection_times, heap_committed_memory, heap_init_memory, heap_max_memory, heap_used_memory):
    x = np.arange(1, len(data) + 1)

    fig, ax1 = plt.subplots(figsize=(12, 6))
    
    ax1.set_xlabel('Iteration')
    ax1.set_ylabel('Response Time (s)', color='blue')
    ax1.plot(x, data, label='Response Time', linewidth=2, color='blue', marker='o')
    cold_x = [x[i] for i, state in enumerate(states) if state == "cold"]
    cold_data = [data[i] for i, state in enumerate(states) if state == "cold"]
    ax1.scatter(cold_x, cold_data, color='red', label='Cold Activation', s=50)
    ax1.tick_params(axis='y', labelcolor='blue')

    ax2 = ax1.twinx()
    ax2.set_ylabel('GC Metrics', color='red')
    ax2.plot(x, gc1_collections, label='GC1 Collections', linewidth=1, color='darkred', marker='^')
    ax2.plot(x, gc1_collection_times, label='GC1 Collection Times', linewidth=1, color='red', marker='v')
    ax2.plot(x, gc2_collections, label='GC2 Collections', linewidth=1, color='lightcoral', marker='<')
    ax2.plot(x, gc2_collection_times, label='GC2 Collection Times', linewidth=1, color='pink', marker='>')
    ax2.tick_params(axis='y', labelcolor='red')
    ax2.set_yscale('symlog')

    ax3 = ax1.twinx()
    ax3.spines['right'].set_position(('outward', 60))
    ax3.set_ylabel('Heap Memory (Committed & Used)')
    ax3.plot(x, heap_committed_memory, label='Heap Committed Memory', color='darkgreen', linestyle='--')
    ax3.plot(x, heap_used_memory, label='Heap Used Memory', color='lightgreen', linestyle=':')
    ax3.tick_params(axis='y', labelcolor='green')

    ax4 = ax1.twinx()
    ax4.spines['right'].set_position(('outward', 120))
    ax4.set_ylabel('Heap Memory (Initial & Max)')
    ax4.plot(x, heap_init_memory, label='Heap Init Memory', color='green', linestyle='-')
    ax4.plot(x, heap_max_memory, label='Heap Max Memory', color='lime', linestyle='-.')
    ax4.tick_params(axis='y', labelcolor='darkgreen')

    fig.tight_layout()
    plt.grid(True)
    fig.legend(loc='center right')
    plt.title('Response Time, GC, and Heap Memory Metrics Over {} Iterations'.format(len(data)))
    plt.savefig(f'../Graphs/Java/{input_size}/'+ 'combined_line_plot.pdf')
    # plt.show()



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
    plt.savefig(f'../Graphs/Java/{input_size}/'+ 'gc_stats_plot.pdf')
    # plt.show()


if __name__ == '__main__':
    # Check if the right number of arguments is provided
    if len(sys.argv) != 2:
        print("Usage: python script_name.py <size>")
        sys.exit(1)

    # Get size from command line
    input_size = int(sys.argv[1])

    directory_path = f'../Graphs/Java/{input_size}/'
    if not os.path.exists(directory_path):
        os.makedirs(directory_path)

    # Read files
    output_file = os.path.join(os.getcwd(), "JavaOutputTime.txt")
    state_file = os.path.join(os.getcwd(), "Javaactivation_ids.txt_startStates.txt")
    gc1_collections_file = os.path.join(os.getcwd(), "gc1Collections.txt")
    gc1_collection_time_file = os.path.join(os.getcwd(), "gc1CollectionTime.txt")
    gc2_collections_file = os.path.join(os.getcwd(), "gc2Collections.txt")
    gc2_collection_time_file = os.path.join(os.getcwd(), "gc2CollectionTime.txt")
    heap_committed_memory_file = os.path.join(os.getcwd(), "heapCommittedMemory.txt")
    heap_init_memory_file = os.path.join(os.getcwd(), "heapInitMemory.txt")
    heap_max_memory_file = os.path.join(os.getcwd(), "heapMaxMemory.txt")
    heap_used_memory_file = os.path.join(os.getcwd(), "heapUsedMemory.txt")

    with open(output_file, 'r') as f:
        latency_data = [float(line.strip()) for line in f.readlines()]

    with open(state_file, 'r') as f:
        states = [line.strip().split(': ')[1].replace('"', '') for line in f.readlines()]

    # Loading GC1 data
    with open(gc1_collections_file, 'r') as f:
        gc1_collections = [float(line.strip()) for line in f.readlines()]

    with open(gc1_collection_time_file, 'r') as f:
        gc1_collection_times = [float(line.strip()) for line in f.readlines()]

    # Loading GC2 data
    with open(gc2_collections_file, 'r') as f:
        gc2_collections = [float(line.strip()) for line in f.readlines()]

    with open(gc2_collection_time_file, 'r') as f:
        gc2_collection_times = [float(line.strip()) for line in f.readlines()]

    # Loading heap data
    with open(heap_committed_memory_file, 'r') as f:
        heap_committed_memory = [float(line.strip()) for line in f.readlines()]

    with open(heap_init_memory_file, 'r') as f:
        heap_init_memory = [float(line.strip()) for line in f.readlines()]

    with open(heap_max_memory_file, 'r') as f:
        heap_max_memory = [float(line.strip()) for line in f.readlines()]

    with open(heap_used_memory_file, 'r') as f:
        heap_used_memory = [float(line.strip()) for line in f.readlines()]

    plot_gc_stats(input_size, gc1_collections, gc1_collection_times, gc2_collections, gc2_collection_times)
    plot_line_orig(input_size, latency_data, states)
    plot_line(input_size, latency_data, states, gc1_collections, gc1_collection_times, gc2_collections, gc2_collection_times, heap_committed_memory, heap_init_memory, heap_max_memory, heap_used_memory)  # Modified the function arguments
    plot_histogram(input_size, latency_data, states)
    # Now plot heap stats
    plot_heap_stats(input_size, heap_committed_memory, heap_init_memory, heap_max_memory, heap_used_memory)
