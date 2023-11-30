import matplotlib.pyplot as plt
import sys

def plot_histograms(client_file, server_file, output_image_file):
    # Read data from files
    with open(client_file, 'r') as f:
        client_times = [float(line.strip()) for line in f.readlines()]

    with open(server_file, 'r') as f:
        server_times = [float(line.strip()) for line in f.readlines()]

    # Plotting
    plt.figure(figsize=(10, 6))
    plt.hist(client_times, bins=50, label='Client Response Times')
    plt.hist(server_times, bins=50, label='Server Execution Times')
    plt.xscale('symlog')
    plt.xlabel('Time (milliseconds)')
    plt.ylabel('Frequency')
    plt.title('Histogram of Response Times')
    plt.legend()

    # Save the plot to the specified file
    plt.savefig(output_image_file)
    plt.close()

# Command-line arguments usage
if __name__ == "__main__":
    if len(sys.argv) != 4:
        print("Usage: python script.py <client_time_file> <server_time_file> <output_image_file>")
        sys.exit(1)

    plot_histograms(sys.argv[1], sys.argv[2], sys.argv[3])
