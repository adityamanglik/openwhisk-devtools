import matplotlib.pyplot as plt
import numpy as np
import sys

def remove_outliers(data, lower_percentile=0, upper_percentile=99.9):
    lower_bound = np.percentile(data, lower_percentile)
    upper_bound = np.percentile(data, upper_percentile)
    return [x for x in data if lower_bound <= x <= upper_bound]

def calculate_statistics(file_path):
    with open(file_path, 'r') as file:
        times = [float(line.strip()) for line in file if line.strip()]

    times = np.array(times)
    average = np.mean(times)
    median = np.median(times)
    p90 = np.percentile(times, 90)
    p99 = np.percentile(times, 99)
 
    return average, median, p90, p99

def plot_histograms(client_file, server_file, output_image_file):
    # Read data from files
    with open(client_file, 'r') as f:
        client_times = [float(line.strip()) for line in f.readlines()]

    with open(server_file, 'r') as f:
        server_times = [float(line.strip()) for line in f.readlines()]

    # Remove outliers
    client_times = remove_outliers(client_times)
    server_times = remove_outliers(server_times)

    # Plotting
    fig, ax1 = plt.subplots(figsize=(10, 6))

    # Plot client times on the primary y-axis
    ax1.hist(client_times, bins=200, color='r', alpha=0.7, label='Client Response Times')
    ax1.set_xlabel('Time (milliseconds)')
    ax1.set_ylabel('Client Frequency', color='g')

    # Create a secondary y-axis for server times
    ax2 = ax1.twinx()
    ax2.hist(server_times, bins=200, color='b', alpha=0.7, label='Server Execution Times')
    ax2.set_ylabel('Server Frequency', color='b')

    # Add titles and legends
    plt.title('Histogram of Response Times')
    ax1.legend(loc='upper right')
    ax2.legend(loc='upper left')
    
    # Calculate statistics
    client_stats = calculate_statistics(client_file)
    server_stats = calculate_statistics(server_file)
    
    # Add text box for client statistics
    stats_text = f'Client Times\nAverage: {client_stats[0]:.2f}\nMedian: {client_stats[1]:.2f}\nP90: {client_stats[2]:.2f}\nP99: {client_stats[3]:.2f}'
    props = dict(boxstyle='round', facecolor='yellow', alpha=0.5)
    ax1.text(0.85, 0.92, stats_text, transform=ax1.transAxes, fontsize=10,
             verticalalignment='top', bbox=props)

    # Add text box for server statistics
    stats_text = f'Server Times\nAverage: {server_stats[0]:.2f}\nMedian: {server_stats[1]:.2f}\nP90: {server_stats[2]:.2f}\nP99: {server_stats[3]:.2f}'
    ax2.text(0.15, 0.92, stats_text, transform=ax2.transAxes, fontsize=10,
             verticalalignment='top', horizontalalignment='right', bbox=props)

    # Save the plot to the specified file
    plt.savefig(output_image_file)
    plt.close()

# Command-line arguments usage
if __name__ == "__main__":
    if len(sys.argv) != 4:
        print("Usage: python script.py <client_time_file> <server_time_file> <output_image_file>")
        sys.exit(1)

    plot_histograms(sys.argv[1], sys.argv[2], sys.argv[3])
