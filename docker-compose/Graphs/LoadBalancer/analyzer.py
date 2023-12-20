def parse_line(line):
    parts = line.split(", ")
    return {p.split(": ")[0]: int(p.split(": ")[1]) for p in parts}

def analyze_file(file_path):
    with open(file_path, 'r') as file:
        data = [parse_line(line.strip()) for line in file]
    
    data = data[2000:]

    if not data:
        print("No data found in the file.")
        return

    # Initialize min and max dictionaries
    min_values = {key: float('inf') for key in data[0]}
    max_values = {key: float('-inf') for key in data[0]}

    # Compute min and max for each column
    for entry in data:
        for key in entry:
            min_values[key] = min(min_values[key], entry[key])
            max_values[key] = max(max_values[key], entry[key])

    # Print results
    for key in min_values:
        print(f"{key}: Min = {min_values[key]}, Max = {max_values[key]}")

# Replace 'data.txt' with the path to your file
analyze_file('/users/am_CU/openwhisk-devtools/docker-compose/Graphs/LoadBalancer/Go/100/memory.txt')