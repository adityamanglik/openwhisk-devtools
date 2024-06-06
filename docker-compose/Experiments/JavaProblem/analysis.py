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
        if user_count >= 499:
            data.append(row)
            # try:
            #     # Convert the row to float starting from the 4th column (index 3)
            #     numeric_row = [float(val) if val != 'N/A' else np.nan for val in row[3:]]
            #     data.append(numeric_row)
            # except ValueError:
            #     # Skip the row if conversion to float fails
            #     continue
            

trans = []

for row in data:
    # print(row)
    row = row[4:]
    # print(row)
    numeric_row = [float(val) for val in row]
    trans.append(numeric_row)
# print(trans)    
            
# Convert the list to a NumPy array
data_array = np.array(trans, dtype=np.float64)

# print(data_array)

# Calculate median of columns 3 and onward
medians = np.median(data_array, axis=0)
print(medians)
# # Print the medians
# for i, median in enumerate(medians, start=4):  # Start index 4 to match column index in original file
#     print(f'Median of column {i}: {median}')
