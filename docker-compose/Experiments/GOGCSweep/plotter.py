import pandas as pd
import matplotlib.pyplot as plt
import os

# Directory containing the CSV files
directory = './Data'

# Array of sizes and GOGC values
sizes = [100, 1000, 10000, 50000]
GOGC = [1, 10, 50, 100, 200, 400, 800, -1]
GOGC_plot = [1, 10, 50, 100, 200, 400, 800, 1000]
column_list = ["ArraySize", "ClientP50", "ClientP99", "ClientP999", "ClientP9999", "ServerP50", "ServerP99", "ServerP999", "ServerP9999"]

# Function to read data from all CSV files for a given size
def read_data(size):
    gogc_array = []
    for gogc in GOGC:
        if gogc == -1:
            filename = f"./Data/{size}_DISABLED_latencies.csv"
        else:
            filename = f"./Data/{size}_{gogc}_latencies.csv"
        if os.path.exists(filename):
            data = pd.read_csv(filename, header=None).iloc[1:].to_numpy()
        else:
            print("File not found: ", filename)
        gogc_array.append(data)
    return gogc_array

# Plotting function
def plot_data(data, response_type):
    plot_title = f"{response_type}"
    
    plt.figure()
    for val in data:
        plt.plot(GOGC_plot, val, marker='o')
    plt.title(f"{response_type}")
    plt.xlabel('GOGC')
    plt.ylabel('Latency (microsec)')
    plt.legend(sizes)
    # plt.grid(True)
    # Save plot in the plot directory
    plt.savefig(f"./Plots/{response_type}.png")
    plt.close()

# Iterate over array sizes, read and plot data for each latency metric
all_data = []
for size in sizes:
    data = read_data(size)
    all_data.append(data)
# print(all_data)
# Pick one column from all data and plot it
response_type = 1
for idx, response_type in enumerate(column_list):
    if idx == 0:
        continue
    plotting_data = []
    for array in all_data:  # Skip the arraysize column
        data_series = []
        
        for column in array:  # Skip the arraysize column
            data_series.append(column[0][idx])
        plotting_data.append(data_series)
    print(plotting_data)
    plot_data(plotting_data, response_type)