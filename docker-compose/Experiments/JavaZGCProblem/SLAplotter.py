import sys
import numpy as np
import matplotlib.pyplot as plt
def plot_hdr_histograms(client_times, output_file):
    # Define the percentiles we are interested in
    percentiles = [50, 90, 95, 99, 99.9, 99.99, 99.999]

    # Calculate the response times at each percentile
    percentile_values = [np.percentile(client_times, p) for p in percentiles]
    percentiles = [str(x) for x in percentiles]
    # Create the plot
    plt.figure(figsize=(10, 6))
    plt.plot(percentiles, percentile_values, marker='o', label='Baseline')

    # Add the expected service level line
    # expected_service_level = median + 3  # Example value for demonstration
    # plt.axhline(y=expected_service_level, color='orange', linestyle='--', label='Expected Service Level')

    # Set the plot labels and title
    plt.xlabel('Percentile')
    plt.ylabel('Response Time (ms)')
    plt.title('Response Time by Percentile Distribution')

    # Set the x-axis to a logarithmic scale
    # plt.xscale('symlog')
    # plt.xticks(percentiles, labels=[f"{p}%" for p in percentiles])

    # Add grid and legend
    plt.grid(True)
    plt.legend(loc='upper left')

    # Save the plot to the specified file
    plt.savefig(output_file)
    plt.close()
    
if __name__ == "__main__":
    # if len(sys.argv) != 6:
        # print("Usage: python script.py <client_time_file> <server_time_file> <memory_file> <dist_image_file> <latency_image_file>")
        # sys.exit(1)
    with open(sys.argv[1], 'r') as f:
        client_times = [float(line.strip().split(', ')[1]) for line in f.readlines()]
    plot_hdr_histograms(client_times, sys.argv[2])