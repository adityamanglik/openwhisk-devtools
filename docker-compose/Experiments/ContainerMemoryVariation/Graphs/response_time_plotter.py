import matplotlib.pyplot as plt
import numpy as np
import sys

SMALL_SIZE = 18
MEDIUM_SIZE = 20
BIGGER_SIZE = 24

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

def plot_latency(client_times, server_times, memory_log, second_container, output_image_file, output_image_file_1):
    # plot all iterations in line graph
    client_times[0] = 1.1*max(client_times[1:])
    print(client_times[0])
    fig, ax1 = plt.subplots(figsize=(18, 7))
    _, med, _, _, _, stdd = calculate_statistics(client_times)
    
    ax1.set_xlabel('Request Number')
    ax1.set_ylabel('Client Time')
    # ax1.set_ylim([med - 5*stdd, med + 5*stdd])
    
    # Plot med + std on y axis
    median = np.median(client_times)
    stdd = np.std(client_times)
    ax1.axhline(y=median, c = 'green', alpha = 0.27, linestyle = '--')
    # ax1.axhline(y=median+stdd, c = 'green', alpha = 0.27, linestyle = '--')
    # ax1.axhline(y=median+2*stdd, c = 'green', alpha = 0.27, linestyle = '--')
    
    # Delineate regions of interest
    # ax1.axvline(x=8, c = 'blue', alpha = 0.27, linestyle = '-')
    ax1.axvline(x=350, c = 'blue', alpha = 0.27, linestyle = '-')
    
    # Plot client times on the primary y-axis
    ax1.plot(range(0, 2), client_times[0:2], color='r', alpha=0.9)
    ax1.plot(range(1, 350), client_times[1:350], color='purple', alpha=0.9, label='Transient')
    ax1.plot(range(350, len(client_times)), client_times[350:], color='blue', alpha=0.9, label='Stable')
    # Plot cold start latency separately
    ax1.plot(0,client_times[0],marker="*", markersize=20, markeredgecolor="black", markerfacecolor="red", label='Cold Start')
    
    # Shade part of plot to clearly delineate
    ax1.axvspan(350, len(client_times), facecolor='b', alpha=0.1)
    
    plt.title('Response Times')
    ax1.legend(loc='upper center')
    # ax1.set_yscale('symlog')
    plt.savefig(output_image_file)
    
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
    ax2.legend(loc='upper right')    
    # print(output_image_file.split('.')[0] + "_1.pdf")
    plt.tight_layout(pad=0)
    plt.savefig(output_image_file_1, bbox_inches='tight', pad_inches=0, format='pdf', dpi=1200)
        
    

def plot_histograms(client_times, server_times, output_image_file):
    # Plotting
    fig, ax1 = plt.subplots(figsize=(10, 6))

    # Plot client times on the primary y-axis
    client_times[0] = 1.1*max(client_times[1:])
    print(client_times[:10])
    ax1.hist(client_times, bins=200, color='r', alpha=0.7, label='Client Response Times')
    ax1.set_xlabel('Time (milliseconds)')
    ax1.set_ylabel('Client Frequency', color='g')

    # Plot cold start marker
    ax1.plot(client_times[0], 3, marker="*", markersize=20, markeredgecolor="black", markerfacecolor="red", label='Cold Start')
    
    # Create a secondary y-axis for server times
    # ax2 = ax1.twinx()
    # ax2.hist(server_times, bins=200, color='b', alpha=0.7, label='Server Execution Times')
    # ax2.set_ylabel('Server Frequency', color='b')

    # Add titles and legends
    plt.title('Histogram of Response Times')
    ax1.legend(loc='upper center')
    # ax2.legend(loc='upper left')
    
    # Calculate statistics
    client_stats = calculate_statistics(client_times)
    # server_stats = calculate_statistics(server_times)
    
    # Add text box for client statistics
    stats_text = f'Client Times\nAverage: {client_stats[0]:.2f}\nMedian: {client_stats[1]:.2f}\nP90: {client_stats[2]:.2f}\nP99: {client_stats[3]:.2f}'
    props = dict(boxstyle='round', facecolor='yellow', alpha=0.3)
    ax1.text(0.75, 0.975, stats_text, transform=ax1.transAxes, fontsize=10,
             verticalalignment='top', bbox=props)
    
    exceeding_gc_threshold = sum(1 for time in client_times if time > (client_stats[1] + client_stats[5]))
    print(f'Number of client_times in GC influence: {exceeding_gc_threshold}, Fraction: {exceeding_gc_threshold/len(client_times)}')

    # Add text box for server statistics
    # stats_text = f'Server Times\nAverage: {server_stats[0]:.2f}\nMedian: {server_stats[1]:.2f}\nP90: {server_stats[2]:.2f}\nP99: {server_stats[3]:.2f}'
    # ax2.text(0.7, 0.92, stats_text, transform=ax2.transAxes, fontsize=10,
            #  verticalalignment='top', horizontalalignment='right', bbox=props)
    
    # Add vertical line to distinguish GC times
    ax1.axvline(x=client_stats[1] + client_stats[5], c = 'blue', alpha = 0.27, linestyle = '-')
    # ax1.axvline(x=client_stats[1] + 2*client_stats[5], c = 'red', alpha = 0.27, linestyle = '--')
    # ax1.axvline(x=client_stats[1] + 3*client_stats[5], c = 'red', alpha = 0.27, linestyle = '--')
    
    # Shade part of plot to clearly delineate
    ax1.axvspan(client_stats[1] + client_stats[5], max(client_times), facecolor='b', alpha=0.1)
    
    # ax1.text((client_stats[1] + client_stats[5]), 30, str(exceeding_gc_threshold/len(client_times)), fontsize=10,
            #  verticalalignment='top', bbox=props)
    # ax1.vline(y, server_stats[1] + server_stats[5], ch, n)
    
    # Find fraction of requests that are in GC region
    

    # Save the plot to the specified file
    plt.savefig(output_image_file, bbox_inches='tight', pad_inches=0, format='pdf', dpi=1200)
    plt.close()
    
def plot_NOGC_histograms(client_times, output_file):
    print('SLA Plot')
    # Discard cold start value
    client_times = client_times[350:]
    
    # Define the percentiles we are interested in
    percentiles = [50, 90, 95, 99, 99.9, 99.99, 99.999]

    # Calculate the response times at each percentile
    percentile_values = [np.percentile(client_times, p) for p in percentiles]
    print('Percentiles: ', percentile_values)
    percentiles_2 = [str(x) for x in percentiles]
    # Create the plot
    plt.figure(figsize=(10, 6))
    plt.plot(percentiles_2, percentile_values, marker='o', color='red', label='With GC')

    # Add the expected service level line
    # expected_service_level = median + 3  # Example value for demonstration
    # plt.axhline(y=expected_service_level, color='orange', linestyle='--', label='Expected Service Level')

    # Remove values from GC
    # Calculate statistics
    client_stats = calculate_statistics(client_times)
    # server_stats = calculate_statistics(server_times)
    
    # Add text box for client statistics
    stats_text = f'Client Times\nAverage: {client_stats[0]:.2f}\nMedian: {client_stats[1]:.2f}\nP90: {client_stats[2]:.2f}\nP99: {client_stats[3]:.2f}\nSTD: {client_stats[5]:.2f}'
    print(stats_text)
    
    exceeding_gc_threshold = sum(1 for time in client_times if time > (client_stats[1] + 1.5*client_stats[5]))
    print(f'Number of client_times in GC influence: {exceeding_gc_threshold}, Fraction: {exceeding_gc_threshold/len(client_times)}')
    
    filtered_client_times = [time for time in client_times if time <= (client_stats[1] + 1.5*client_stats[5])]
    percentile_values = [np.percentile(filtered_client_times, p) for p in percentiles]
    percentiles = [str(x) for x in percentiles]
    plt.plot(percentiles_2, percentile_values, marker='o', color='blue', label='Without GC')
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
    output_file = "./Graphs/Go/10000/sla_plot_2.pdf"
    plt.savefig(output_file, bbox_inches='tight', pad_inches=0, format='pdf', dpi=1200)
    plt.close()
    
def plot_hdr_histograms(latencies, memory_sizes):
    # Define the percentiles we are interested in
    percentiles = [50, 90, 95, 99, 99.9, 99.99, 99.999]
    # Create the plot
    plt.figure(figsize=(10, 6))
    for client_times, memory in zip(latencies, memory_sizes):
        # Calculate the response times at each percentile
        percentile_values = [np.percentile(client_times, p) for p in percentiles]
        percentiles_print = [str(x) for x in percentiles]
        plt.plot(percentiles_print, percentile_values, marker='o', label=f'{memory}')

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
    plt.savefig("naivesolution4.pdf", bbox_inches='tight', pad_inches=0, format='pdf', dpi=1200)
    plt.close()

# Command-line arguments usage
if __name__ == "__main__":
    # if len(sys.argv) != 6:
        # print("Usage: python script.py <client_time_file> <server_time_file> <memory_file> <dist_image_file> <latency_image_file>")
        # sys.exit(1)
        
    memory_sizes=["128m", "256m", "512m", "1024m", "2048m"]
    latencies = []
    for mem in memory_sizes:
        read_me = f'times_{mem}.txt'
        with open(read_me, 'r') as f:
            client_times = [float(line.strip().split(', ')[1]) for line in f.readlines()]
            print(client_times[:10])
            latencies.append(client_times)
    
    plot_hdr_histograms(latencies, memory_sizes)
    