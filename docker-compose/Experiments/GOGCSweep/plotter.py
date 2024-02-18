import pandas as pd
import numpy as np
import matplotlib.pyplot as plt
import os

# Directory containing the CSV files
directory = './Data'

# Array of sizes and GOGC values
sizes = [1000, 10000, 100000]
GOGC = [1, 10, 100, 500, 999]
GOGC_plot = [1, 10, 100, 500, 999]
column_list = ["ArraySize", "totalExecutionTime", "ClientAvg", "ClientP50", "ClientP99", "ClientP999", "ClientP9999", "ServerAvg", "ServerP50", "ServerP99", "ServerP999", "ServerP9999"]

# Function to read data from all CSV files for a given array size
def read_data(size):
    gogc_array = []
    for gogc in GOGC:
        if gogc == -1:
            filename = f"./Data/{size}_DISABLED_latencies.csv"
        else:
            filename = f"./Data/{size}_{gogc}_latencies.csv"
        if os.path.exists(filename):
            data = pd.read_csv(filename, header=None).iloc[1:].values
        else:
            print("File not found: ", filename)
        # print(data)
        # print('BREAK')
        gogc_array.append(data)
    return gogc_array

# Plotting function
def plot_data(data, response_type):
    plot_title = f"{response_type}"
    
    plt.figure()
    for val in data:
        val = [int(x) for x in val]
        print(GOGC_plot, val, response_type)
        plt.plot(GOGC_plot, val, marker='o')
        
    plt.title(f"{response_type}")
    plt.xlabel('GOGC')
    plt.ylabel('Latency (microsec)')
    plt.legend(sizes)
    plt.yscale('log')
    # plt.grid(True)
    # Save plot in the plot directory
    plt.savefig(f"./Plots/{response_type}.png")
    plt.close()

# Iterate over array sizes, read and plot data for each latency metric
all_data = []
for size in sizes:
    data = read_data(size)
    # print(data)
    # print('BREAK2')
    data = np.vstack(data) 
    # print(data)
    # print('BREAK3')
    data = data.T
    # print(data)
    # print('BREAK4')
    all_data.append(data)
print(all_data)

# Pick one column from all data and plot it
response_type = 1
for idx, response_type in enumerate(column_list):
    if idx < 2:
        continue
    plotting_data = []
    for array in all_data:  # Skip the arraysize column
        # print(array[idx])
        # print('REAK')
        plotting_data.append(array[idx])
    print(plotting_data)
    plot_data(plotting_data, response_type)