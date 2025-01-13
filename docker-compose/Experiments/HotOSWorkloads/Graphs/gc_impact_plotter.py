import matplotlib.pyplot as plt

def parse_data(file_path):
    with open(file_path, 'r') as file:
        lines = file.readlines()

    # Initialize data storage
    workload1_server, workload1_client = [], []
    workload2_server, workload2_client = [], []
    arraysizes = []
    current_workload = 1

    for line in lines:
        if 'NOGC' in line:
            current_workload = 2
            continue

        if 'Server:' in line:
            parts = line.split()
            server_time = float(parts[1])
            client_time = float(parts[3])

            if current_workload == 1:
                workload1_server.append(server_time)
                workload1_client.append(client_time)
            else:
                workload2_server.append(server_time)
                workload2_client.append(client_time)
        
        else:
            if current_workload == 1:
                line = line.strip()
                try:
                    val = line
                    arraysizes.append(val)
                except:
                    continue

    return workload1_server, workload1_client, workload2_server, workload2_client, arraysizes

def plot_data(workload1_server, workload1_client, workload2_server, workload2_client, output_path, arraysizes):
    plt.figure(figsize=(10, 6))

    # Plot for workload 1
    plt.plot(arraysizes, workload1_server, label='EM - Server', color='red', linestyle='dashed')
    plt.plot(arraysizes, workload1_client, label='EM - Client', color='red')

    # Plot for workload 2
    plt.plot(arraysizes, workload2_server, label='GC - Server', color='blue', linestyle='dashed')
    plt.plot(arraysizes, workload2_client, label='GC - Client', color='blue')

    plt.xlabel('Arraysize')
    plt.ylabel('GC Impact Metric')
    # plt.yscale('symlog')
    plt.title('GC Impact on Server and Client Latencies for increasing arraysizes')
    plt.legend()
    plt.grid(True)
    plt.tight_layout()
    plt.savefig(output_path)

# Replace 'your_file_path.txt' with the actual path of your file
file_path = '../../analyzer3.log'
output_path = './Go/gc_impact_scale.pdf'
workload1_server, workload1_client, workload2_server, workload2_client, arraysizes = parse_data(file_path)
print(workload1_server, workload1_client, workload2_server, workload2_client, arraysizes)
plot_data(workload1_server, workload1_client, workload2_server, workload2_client, output_path, arraysizes)
