import matplotlib.pyplot as plt
import numpy as np
import sys

SMALL_SIZE = 28
MEDIUM_SIZE = 30
BIGGER_SIZE = 38

plt.rc('font', size=SMALL_SIZE)          # controls default text sizes
plt.rc('axes', titlesize=SMALL_SIZE)     # fontsize of the axes title
plt.rc('axes', labelsize=MEDIUM_SIZE)    # fontsize of the x and y labels
plt.rc('xtick', labelsize=SMALL_SIZE)    # fontsize of the tick labels
plt.rc('ytick', labelsize=SMALL_SIZE)    # fontsize of the tick labels
plt.rc('legend', fontsize=SMALL_SIZE)    # legend fontsize
plt.rc('figure', titlesize=BIGGER_SIZE)  # fontsize of the figure title

def parse_memory_log(memory_file):
    # Extract heapidle and heapalloc
    ret_val = []
    second_container = []
    for idx, line in enumerate(memory_file):
        l2 = line.strip()
        if l2 == "":
            continue
        parts = l2.split(", ")
        if '9502' in parts[1]:
            second_container.append(idx)
        parts = parts[2]
        # val = []
        # for part in parts:
            # val.append(int(part.split(": ")[1]))
        ret_val.append(int(parts.split(": ")[1]))
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

def plot_latency(client_times, server_times, memory_log, output_image_file, output_image_file_1):
    # plot all iterations in line graph
    fig, ax1 = plt.subplots(figsize=(15, 6))
    client_times = [x//1000 for x in client_times]
    _, med, _, _, _, stdd = calculate_statistics(client_times)
    # Plot client times on the primary y-axis
    ax1.plot(client_times, color='r', alpha=0.9, label='Client Response Times')
    ax1.set_xlabel('Request Number')
    ax1.set_ylabel('Client Time (ms)', color='r')
    # ax1.set_ylim([med - 5*stdd, med + 5*stdd])
    
    # Plot med + std on y axis
    median = np.median(client_times)
    stdd = np.std(client_times)
    ax1.axhline(y=median, c = 'green', alpha = 0.27, linestyle = '--')
    # ax1.axhline(y=median+stdd, c = 'green', alpha = 0.27, linestyle = '--')
    
    plt.title('Response Times')
    # ax1.legend(loc='upper left')
    plt.savefig(output_image_file, bbox_inches='tight', pad_inches=0, format='pdf', dpi=1200)
    
    ax2 = ax1.twinx()
    ax2.plot(memory_log, color='b', alpha=0.4, label='HeapAlloc')
    ax2.set_ylabel('Allocated heap memory', color='b')
   
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
    # if second_container != []:
    #     for req in second_container:
    #         ax1.scatter(x=req, y = 3500, c = 'blue', alpha = 0.27, marker = '*')
        
    # Add titles and legends
    # plt.title('Response Times vs Heap memory allocation')
    # ax1.legend(loc='upper left')
    # ax2.legend(loc='upper right')    
    # print(output_image_file.split('.')[0] + "_1.pdf")
    plt.savefig(output_image_file_1, bbox_inches='tight', pad_inches=0, format='pdf', dpi=1200)
        
    

def plot_histograms(client_times, server_times, output_image_file):
    # Plotting
    fig, ax1 = plt.subplots(figsize=(10, 6))

    # Plot client times on the primary y-axis
    ax1.hist(client_times, bins=200, color='r', alpha=0.7, label='Client Response Times')
    ax1.set_xlabel('Time (milliseconds)')
    ax1.set_ylabel('Client Frequency', color='g')

    # Create a secondary y-axis for server times
    # ax2 = ax1.twinx()
    # ax2.hist(server_times, bins=200, color='b', alpha=0.7, label='Server Execution Times')
    # ax2.set_ylabel('Server Frequency', color='b')

    # Add titles and legends
    plt.title('Histogram of Response Times')
    ax1.legend(loc='upper right')
    # ax2.legend(loc='upper left')
    
    # Calculate statistics
    client_stats = calculate_statistics(client_times)
    # server_stats = calculate_statistics(server_times)
    
    # Add text box for client statistics
    stats_text = f'Client Times\nAverage: {client_stats[0]:.2f}\nMedian: {client_stats[1]:.2f}\nP90: {client_stats[2]:.2f}\nP99: {client_stats[3]:.2f}'
    props = dict(boxstyle='round', facecolor='yellow', alpha=0.3)
    ax1.text(0.8, 0.92, stats_text, transform=ax1.transAxes, fontsize=10,
             verticalalignment='top', bbox=props)
    
    exceeding_gc_threshold = sum(1 for time in client_times if time > (client_stats[1] + client_stats[5]))
    print(f'Number of client_times in GC influence: {exceeding_gc_threshold}, Fraction: {exceeding_gc_threshold/len(client_times)}')

    # Add text box for server statistics
    # stats_text = f'Server Times\nAverage: {server_stats[0]:.2f}\nMedian: {server_stats[1]:.2f}\nP90: {server_stats[2]:.2f}\nP99: {server_stats[3]:.2f}'
    # ax2.text(0.7, 0.92, stats_text, transform=ax2.transAxes, fontsize=10,
            #  verticalalignment='top', horizontalalignment='right', bbox=props)
    
    # Add vertical line to distinguish GC times
    ax1.axvline(x=client_stats[1] + client_stats[5], c = 'red', alpha = 0.27, linestyle = '--')
    # ax1.text((client_stats[1] + client_stats[5]), 30, str(exceeding_gc_threshold/len(client_times)), fontsize=10,
            #  verticalalignment='top', bbox=props)
    # ax1.vline(y, server_stats[1] + server_stats[5], ch, n)
    
    # Find fraction of requests that are in GC region
    

    # Save the plot to the specified file
    plt.savefig(output_image_file, bbox_inches='tight', pad_inches=0, format='pdf', dpi=1200)
    plt.close()
    
def plot_hdr_histograms(client_times, output_file):
    # Define the percentiles we are interested in
    percentiles = [50, 90, 95, 99, 99.9, 99.99, 99.999]

    # Calculate the response times at each percentile
    percentile_values = [np.percentile(client_times, p) for p in percentiles]
    percentiles = [str(x) for x in percentiles]
    # Create the plot
    plt.figure(figsize=(10, 6))
    plt.plot(percentiles, percentile_values, marker='o', label='Baseline')

    # Add the expected service level line
    # expected_service_level = median + 3  # Example value for demonstration
    # plt.axhline(y=expected_service_level, color='orange', linestyle='--', label='Expected Service Level')

    # Set the plot labels and title
    plt.xlabel('Percentile')
    plt.ylabel('Response Time (ms)')
    plt.title('Response Time by Percentile Distribution')

    # Set the x-axis to a logarithmic scale
    # plt.xscale('symlog')
    # plt.xticks(percentiles, labels=[f"{p}%" for p in percentiles])

    # Add grid and legend
    plt.grid(True)
    plt.legend(loc='upper left')

    # Save the plot to the specified file
    plt.savefig(output_file, bbox_inches='tight', pad_inches=0, format='pdf', dpi=1200)
    plt.close()
    

# Command-line arguments usage
if __name__ == "__main__":
    # if len(sys.argv) != 6:
        # print("Usage: python script.py <client_time_file> <server_time_file> <memory_file> <dist_image_file> <latency_image_file>")
        # sys.exit(1)
    with open(sys.argv[1], 'r') as f:
        client_times = [float(line.strip().split(', ')[1]) for line in f.readlines()]
    # print(client_times[:10])
    with open(sys.argv[2], 'r') as f:
        server_times = [float(line.strip().split(',')[1]) for line in f.readlines()]
    # print(server_times[:10])
    with open(sys.argv[3], 'r') as f:
        memory_log = [float(line.strip().split(',')[1]) for line in f.readlines()]
    # print(len(memory_log))
    # Skip warm up
    # TODO: Pass warm up and actual request numbers from go file
    # memory_log = memory_log[len(memory_log)//2:]
    # print(memory_log[:10])
    # print(second_container)
    client_times = client_times[:100]
    client_times = [x*8 for x in client_times]
    plot_histograms(client_times, server_times, sys.argv[4])
    plot_latency(client_times, server_times, memory_log, sys.argv[5], sys.argv[6])
    plot_hdr_histograms(client_times, sys.argv[7])
