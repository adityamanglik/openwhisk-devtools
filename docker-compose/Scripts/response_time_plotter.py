import matplotlib.pyplot as plt
import numpy as np
import sys

def plot_scatter(data, output_file):
    # Generate X values (iteration numbers)
    x = np.arange(1, len(data) + 1)

    # Create the scatter plot
    plt.figure(figsize=(12, 6))
    plt.scatter(x, data, s=2, label=output_file+' Response Time')
    plt.title('Response Time Over {} Iterations'.format(len(data)))
    plt.xlabel('Iteration')
    plt.ylabel('Response Time (s)')
    plt.grid(True)
    plt.legend()

    # Save the plot as a PNG file
    plt.savefig('../Graphs/'+ output_file.replace('.txt', '_scatter_plot.png'))

    # Show the plot
    # plt.show()

def plot_line(data, output_file):
    # # Only use the first 500 data points
    # data = data[:500]

    # Generate X values (iteration numbers)
    x = np.arange(1, len(data) + 1)

    # Create the line plot
    plt.figure(figsize=(12, 6))
    plt.plot(x, data, label=output_file+' Response Time', linewidth=2)
    plt.title('Response Time for First 500 Iterations')
    plt.xlabel('Iteration')
    plt.ylabel('Response Time (s)')
    plt.grid(True)
    plt.legend()

    # Save the plot as a PNG file
    plt.savefig('../Graphs/'+ output_file.replace('.txt', '_line_plot.png'))

    # Show the plot
    # plt.show()

if __name__ == '__main__':
    # Command line arguments
    if len(sys.argv) < 2:
        print("Usage: python script_name.py output_file.txt")
        sys.exit(1)

    output_file = sys.argv[1]
    
    # Read the data from the specified output file
    with open(output_file, 'r') as f:
        lines = f.readlines()

    # Convert string data to floats
    data = [float(line.strip()) for line in lines]

    # Plot both scatter and line plots
    plot_scatter(data, output_file)
    plot_line(data, output_file)
