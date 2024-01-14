import statistics
import sys

def parse_memory_log(memory_file):
    # [int(parse_line(line.strip())) for line in file]
    ret_val = []
    for line in memory_file:
        l2 = line.strip()
        parts = l2.split(", ")
        parts = parts[3]
        ret_val.append(int(parts.split(": ")[1]))
    # print(parts)
    return ret_val

def calculate_gc_impact(memory_log, server_log, client_log):
    num_entries = 5
    gc_client_impact = []
    gc_server_impact = []
    cycle_count = 0
    
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
                    gc_server_impact.append(val - (med + stdd))
            client_times = client_log[index - num_entries : index + num_entries]
            # Calculate median and STD for the 10 entries
            med = statistics.mean(client_times)
            # Sum up all values that exceed median + STD
            stdd = statistics.stdev(client_times)
            # The sum is the GC impact metric
            for val in client_times:
                if (val > med + stdd):
                    gc_client_impact.append(val - (med + stdd))
        index += 1
    # Take average of values
    ans = []
    ans.append(sum(gc_server_impact)/len(gc_server_impact))
    ans.append(sum(gc_client_impact)/len(gc_client_impact))
    ans.append(cycle_count)
    return ans

if __name__ == "__main__":
    with open(sys.argv[1], 'r') as file:
        memory_log = parse_memory_log(file)
    # print(memory_log[:10], len(memory_log))
    with open(sys.argv[2], 'r') as file:
        server_log = [int(line.strip().split(',')[1]) for line in file]
    # print(server_log[:10], len(server_log))
    with open(sys.argv[3], 'r') as file:
        client_log = [int(line.strip().split(',')[1]) for line in file]
    # print(client_log[:10], len(client_log))
        
    ans = calculate_gc_impact(memory_log, server_log, client_log)
    print("Server: ", ans[0], " Client: ", ans[1], " Cycle Count: ", ans[2])