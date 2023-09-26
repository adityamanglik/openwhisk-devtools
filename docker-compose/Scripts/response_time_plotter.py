import matplotlib.pyplot as plt
import numpy as np
import sys


def plot_histogram(data, states, output_file):
    plt.figure(figsize=(12, 6))

    # Create histogram bins
    n, bins, patches = plt.hist(data, bins=50, color='blue', label='Warm Response Time', alpha=0.7)
    
    # Get the counts of cold starts in each bin
    cold_counts, _ = np.histogram([data[i] for i, state in enumerate(states) if state == "cold"], bins=bins)
    
    # Mark the cold starts at the top of the histogram bar with a number in red
    for count, rect, cold_count in zip(n, patches, cold_counts):
        height = count
        plt.text(rect.get_x() + rect.get_width() / 2, height + 5, str(int(cold_count)), ha='center', va='bottom', color='red')
    
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

if __name__ == '__main__':
    if len(sys.argv) < 3:
        print("Usage: python script_name.py latency_output_file.txt state_file.txt")
        sys.exit(1)

    output_file = sys.argv[1]
    state_file = sys.argv[2]

    with open(output_file, 'r') as f:
        lines = f.readlines()

    with open(state_file, 'r') as f:
        states = [line.strip().split(': ')[1].replace('"', '') for line in f.readlines()]

    data = [float(line.strip()) for line in lines]

    plot_histogram(data, states, output_file)
    plot_line(data, states, output_file)
