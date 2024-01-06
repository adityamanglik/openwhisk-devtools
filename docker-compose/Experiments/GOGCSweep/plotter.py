import pandas as pd
import matplotlib.pyplot as plt
import os

# Directory containing the CSV files
directory = './Data'

# Array of sizes and GOGC values
sizes = [100, 10000, 1000000]
GOGC = [-1, 1, 10, 50, 100, 200, 400, 800]

# Function to read data from CSV
def read_data(size, gogc):
    filename = f"{directory}/{size}_{gogc}_latencies.csv"
    if os.path.exists(filename):
        return pd.read_csv(filename)
    else:
        return None

# Plotting function
def plot_data(data, column, title):
    plt.figure()
    plt.plot(data['ArraySize'], data[column], marker='o')
    plt.title(title)
    plt.xlabel('Array Size')
    plt.ylabel(column)
    plt.grid(True)
    plt.savefig(f"./Plots/{title}.png")
    plt.close()

# Iterate over sizes and GOGC values, read and plot data
for size in sizes:
    for gc in GOGC:
        data = read_data(size, gc)
        if data is not None:
            for column in data.columns[1:]:  # Skip the first column (ArraySize)
                plot_title = f"{column} - Size {size} - GOGC {gc}"
                plot_data(data, column, plot_title)
