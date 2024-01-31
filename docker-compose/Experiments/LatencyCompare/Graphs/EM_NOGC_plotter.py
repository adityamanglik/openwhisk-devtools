import matplotlib.pyplot as plt
import numpy as np
import sys

def parse_memory_log(memory_file):
    # Extract heapidle and heapalloc
    ret_val = []
    second_container = []
    for idx, line in enumerate(memory_file):
        l2 = line.strip()
        if l2 == "":
            continue
        parts = l2.split(", ")
        if '9501' in parts[1]:
            second_container.append(idx)
        parts = parts[2:4]
        val = []
        for part in parts:
            val.append(int(part.split(": ")[1]))
        ret_val.append(val)
    # print(parts)
    return ret_val, second_container

def remove_outliers(data, lower_percentile=0, upper_percentile=99.99):
    lower_bound = np.percentile(data, lower_percentile)
    upper_bound = np.percentile(data, upper_percentile)
    answer = [x for x in data if lower_bound <= x <= upper_bound]
    outliers = [x for x in data if lower_bound >= x or x >= upper_bound]
    for x in outliers:
        print("Removed outlier value from plotting: ", x)
    return answer

def calculate_statistics(times):
    times = np.array(times)
    average = np.mean(times)
    median = np.median(times)
    p90 = np.percentile(times, 90)
    p99 = np.percentile(times, 99)
    summed = np.sum(times)
    stdd = np.std(times)
 
    return average, median, p90, p99, summed, stdd

def plot_latency(client_times, NOGC_client_times, output_image_file):
    client_times = client_times[2000:2500]
    NOGC_client_times = NOGC_client_times[2000:2500]
    # server_times = server_times[:5000]
    # memory_log = memory_log[:5000]
    
    # plot all iterations in line graph
    # Plotting
    fig, ax1 = plt.subplots(figsize=(10, 6))
    _, med, _, _, _, stdd = calculate_statistics(client_times)
    stats_text = f'EM Median: {med:.2f} Std: {stdd:.2f}\n'
    print(med, stdd)
    _, med, _, _, _, stdd = calculate_statistics(NOGC_client_times)
    stats_text += f'NOGC Median: {med:.2f} Std: {stdd:.2f}'
    print(med, stdd)
    # Plot client times on the primary y-axis
    
    ax1.plot(NOGC_client_times, color='b', alpha=0.6, label='NOGC Response Times')
    ax1.plot(client_times, color='r', alpha=0.95, label='EM Response Times')
    ax1.set_xlabel('Request Number')
    ax1.set_ylabel('Client Time', color='r')
    # ax1.set_ylim([med - 5*stdd, med + 5*stdd])
    
    # Plot med + std on y axis
    # median = np.median(client_times)
    # stdd = np.std(client_times)
    # ax1.axhline(y=median, c = 'green', alpha = 0.27, linestyle = '--')
    # ax1.axhline(y=median+stdd, c = 'green', alpha = 0.27, linestyle = '--')

    # Create a secondary y-axis for server times
    # ax2 = ax1.twinx()
    # ax2.plot(server_times, color='b', alpha=0.7, label='Server Execution Times')
    # ax2.set_ylabel('Server Time', color='b')
        
    # GC_iterations = []
    # for idx in range(1, len(memory_log)):
    #     # heapalloc, heapidle
    #     # mark iterations with HeapIdle increase or HeapAlloc decrease as GC calls
    #     if memory_log[idx][0] < memory_log[idx - 1][0]:
    #         # print("HeapAlloc")
    #         # print(memory_log[idx][0], memory_log[idx - 1][0])
    #         GC_iterations.append(idx)
    #     # elif memory_log[idx][1] > memory_log[idx - 1][1]:
    #         # print("HeapIdle")
    #         # print(memory_log[idx][1], memory_log[idx - 1][1])
    #         # GC_iterations.append(idx)
    #     idx += 1
    # # Mark GC cycle iterations with green vertical lines
    # for iter in GC_iterations:
    #     ax1.axvline(x=iter, c = 'green', alpha = 0.27, linestyle = '--')
        
    # # Plot second container calls
    # # if second_container != []:
    # #     for req in second_container:
    # #         ax1.scatter(x=req, y = 3500, c = 'blue', alpha = 0.27, marker = '*')
    
    # Add median and stdd in text box
    props = dict(boxstyle='round', facecolor='yellow', alpha=0.5)
    ax1.text(0.65, 0.92, stats_text, transform=ax1.transAxes, fontsize=10,
             verticalalignment='top', bbox=props)
        
    # # Add titles and legends
    plt.title('Response Times')
    ax1.legend(loc='upper left')
    # ax2.legend(loc='upper right')    
    plt.savefig(output_image_file)

def EM_NOGC_plot_histograms(clienttimes, server_times, output_image_file):
    # TODO: Modify function to plot distributions from both EM and NOGC on same plot
    # Read data from files    
    client_stats = calculate_statistics(client_times)
    server_stats = calculate_statistics(server_times)
    
    # Add text box for client statistics
    stats_text = f'Client Times\nAverage: {client_stats[0]:.2f}\nMedian: {client_stats[1]:.2f}\n STD: {client_stats[5]:.2f}\nP90: {client_stats[2]:.2f}\nP99: {client_stats[3]:.2f}\nSummed: {client_stats[4]:.2f}'
    print(stats_text)
    # Add text box for server statistics
    stats_text = f'Server Times\nAverage: {server_stats[0]:.2f}\nMedian: {server_stats[1]:.2f}\n STD: {client_stats[5]:.2f}\nP90: {server_stats[2]:.2f}\nP99: {server_stats[3]:.2f}\nSummed: {server_stats[4]:.2f}'
    print(stats_text)
    
    # Remove outliers
    print("Removing outliers from client: ")
    # client_times = remove_outliers(client_times)
    print("Removing outliers from server: ")
    # server_times = remove_outliers(server_times)

    # Plotting
    fig, ax1 = plt.subplots(figsize=(10, 6))

    # Plot client times on the primary y-axis
    ax1.hist(client_times, bins=200, color='r', alpha=0.7, label='Client Response Times')
    ax1.set_xlabel('Time (milliseconds)')
    ax1.set_ylabel('Client Frequency', color='g')

    # Create a secondary y-axis for server times
    ax2 = ax1.twinx()
    ax2.hist(server_times, bins=200, color='b', alpha=0.7, label='Server Execution Times')
    ax2.set_ylabel('Server Frequency', color='b')

    # Add titles and legends
    plt.title('Histogram of Response Times')
    ax1.legend(loc='upper right')
    ax2.legend(loc='upper left')
    
    # Calculate statistics
    client_stats = calculate_statistics(client_times)
    server_stats = calculate_statistics(server_times)
    
    # Add text box for client statistics
    stats_text = f'Client Times\nAverage: {client_stats[0]:.2f}\nMedian: {client_stats[1]:.2f}\nP90: {client_stats[2]:.2f}\nP99: {client_stats[3]:.2f}'
    props = dict(boxstyle='round', facecolor='yellow', alpha=0.5)
    ax1.text(0.85, 0.92, stats_text, transform=ax1.transAxes, fontsize=10,
             verticalalignment='top', bbox=props)

    # Add text box for server statistics
    stats_text = f'Server Times\nAverage: {server_stats[0]:.2f}\nMedian: {server_stats[1]:.2f}\nP90: {server_stats[2]:.2f}\nP99: {server_stats[3]:.2f}'
    ax2.text(0.15, 0.92, stats_text, transform=ax2.transAxes, fontsize=10,
             verticalalignment='top', horizontalalignment='right', bbox=props)

    # Save the plot to the specified file
    plt.savefig(output_image_file)
    plt.close()

# Command-line arguments usage
if __name__ == "__main__":
    # if len(sys.argv) != 6:
    #     print("Usage: python script.py <client_time_file> <server_time_file> <memory_file> <dist_image_file> <latency_image_file>")
    #     sys.exit(1)
    with open("../EM.txt", 'r') as f:
        client_times = [float(line.strip().split(', ')[1]) for line in f.readlines()]
    # print(client_times[:10])
    with open("../NOGC.txt", 'r') as f:
        NOGC_client_times = [float(line.strip().split(', ')[1]) for line in f.readlines()]
        
    # Discard first 100 values as they are unstable
    # client_times = client_times[100:]
    # server_times = server_times[100:]
    # NOGC_client_times = NOGC_client_times[100:]
    # NOGC_server_times = NOGC_server_times[100:]
    # print(memory_log[:10])
    # print(second_container)
    # plot_histograms(client_times, server_times, sys.argv[4])
    plot_latency(client_times, NOGC_client_times, "../EM_NOGC.pdf")
