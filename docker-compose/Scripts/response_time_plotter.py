import matplotlib.pyplot as plt
import numpy as np
import sys

def plot_scatter(data, states, output_file):
    x = np.arange(1, len(data) + 1)
    
    cold_data = [data[i] for i, state in enumerate(states) if state == "cold"]
    warm_data = [data[i] for i, state in enumerate(states) if state == "warm"]
    cold_x = [x[i] for i, state in enumerate(states) if state == "cold"]
    warm_x = [x[i] for i, state in enumerate(states) if state == "warm"]

    plt.figure(figsize=(12, 6))
    plt.scatter(cold_x, cold_data, s=50, color='red', label='Cold Response Time')
    plt.scatter(warm_x, warm_data, s=2, color='blue', label='Warm Response Time')
    plt.title('Response Time Over {} Iterations'.format(len(data)))
    plt.xlabel('Iteration')
    plt.ylabel('Response Time (s)')
    plt.grid(True)
    plt.legend()
    plt.savefig('../Graphs/'+ output_file.replace('.txt', '_scatter_plot.png'))
    # plt.show()

def plot_line(data, states, output_file):
    x = np.arange(1, len(data) + 1)

    plt.figure(figsize=(12, 6))
    plt.plot(x, data, label=output_file+' Response Time', linewidth=2, color='blue')
    cold_x = [x[i] for i, state in enumerate(states) if state == "cold"]
    cold_data = [data[i] for i, state in enumerate(states) if state == "cold"]
    plt.scatter(cold_x, cold_data, color='red', label='Cold Activation', s=50)
    plt.title('Response Time for First 500 Iterations')
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

    plot_scatter(data, states, output_file)
    plot_line(data, states, output_file)
