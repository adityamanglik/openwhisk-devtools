import statistics

def parse_line(filread):
    # [int(parse_line(line.strip())) for line in file]
    ret_val = []
    for line in filread:
        l2 = line.strip()
        parts = l2.split(", ")
        if('9501' in parts[0]):
            continue
        if('FAKE' in parts[0]):
            continue
        parts = parts[2]
        ret_val.append(int(parts.split(": ")[1]))
    # print(parts)
    return ret_val

def analyze_file(memory_file, server_file, client_file):
    num_entries = 5
    gc_client_impact = 0
    gc_server_impact = 0
    cycle_count = 0
    
    with open(memory_file, 'r') as file:
        memory_log = parse_line(file)
    # print(memory_log[:10], len(memory_log))
    with open(server_file, 'r') as file:
        server_log = [int(line.strip()) for line in file]
    # print(server_log[:10], len(server_log))
    with open(client_file, 'r') as file:
        client_log = [int(line.strip()) for line in file]
    # print(client_log[:10], len(client_log))
    
    # Check when HeapIdle increases == GC cycle
    index = 1
    while(index < len(memory_log)):
        if (memory_log[index] > memory_log[index - 1]):
            cycle_count += 1
            # For the corresponding index, extract +-10 entries in the server_time
            server_times = server_log[index - num_entries : index + num_entries]
            # print(server_times)
            # Calculate median and STD for the 10 entries
            med = statistics.mean(server_times)
            # Sum up all values that exceed median + STD
            stdd = statistics.stdev(server_times)
            # The sum is the GC impact metric
            for val in server_times:
                if (val > med + stdd):
                    gc_server_impact += val - (med + stdd)
            client_times = client_log[index - num_entries : index + num_entries]
            # Calculate median and STD for the 10 entries
            med = statistics.mean(client_times)
            # Sum up all values that exceed median + STD
            stdd = statistics.stdev(client_times)
            # The sum is the GC impact metric
            for val in client_times:
                if (val > med + stdd):
                    gc_client_impact += val - (med + stdd)
        index += 1
    return (gc_server_impact, gc_client_impact, cycle_count)

# Replace 'data.txt' with the path to your file
ans = analyze_file('/users/am_CU/openwhisk-devtools/docker-compose/Experiments/GCScheduler/Graphs/GCScheduler/Go/10000/memory.txt', '/users/am_CU/openwhisk-devtools/docker-compose/Experiments/GCScheduler/Graphs/GCScheduler/Go/10000/server_time.txt', '/users/am_CU/openwhisk-devtools/docker-compose/Experiments/GCScheduler/Graphs/GCScheduler/Go/10000/client_time.txt')
print("Server: ", ans[0], " Client: ", ans[1], " Cycle Count: ", ans[2])