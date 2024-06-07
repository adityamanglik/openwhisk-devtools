import numpy as np
import csv

# Open and read the file to find target data
filename = 'locust_stats_history.csv'

# Initialize an empty list to store the data
data = []

with open(filename, mode='r') as file:
    csv_reader = csv.reader(file)
    header = next(csv_reader)  # Skip the header

    for row in csv_reader:
        user_count = int(row[1])
        if user_count >= 999:
            data.append(row)

trans = []

for row in data:
    row = row[4:]
    numeric_row = [float(val) if val != 'N/A' else np.nan for val in row]
    trans.append(numeric_row)

# Convert the list to a NumPy array
data_array = np.array(trans, dtype=np.float64)

# Calculate median of columns 3 and onward
medians = np.nanmedian(data_array, axis=0)

# Print the medians in CSV format
print(','.join(map(str, medians)))

# print(f'Medians have been written to {output_filename}')
