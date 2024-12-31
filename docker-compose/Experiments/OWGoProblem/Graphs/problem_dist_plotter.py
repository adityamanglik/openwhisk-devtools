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

def plot_histograms(server_times, server_times_2, output_image_file):
    # Read data from files    
    server_stats = calculate_statistics(server_times)
    
    # Add text box for server statistics
    stats_text = f'Server Times\nAverage: {server_stats[0]:.2f}\nMedian: {server_stats[1]:.2f}\n P90: {server_stats[2]:.2f}\nP99: {server_stats[3]:.2f}\nSummed: {server_stats[4]:.2f}'
    print(stats_text)
    
    # Read data from files    
    server_stats = calculate_statistics(server_times_2)
    
    # Add text box for server statistics
    stats_text = f'Server Times\nAverage: {server_stats[0]:.2f}\nMedian: {server_stats[1]:.2f}\n P90: {server_stats[2]:.2f}\nP99: {server_stats[3]:.2f}\nSummed: {server_stats[4]:.2f}'
    print(stats_text)
    
    # Plotting
    fig, ax1 = plt.subplots(figsize=(10, 6))

    # Plot client times on the primary y-axis
    ax1.hist(server_times, bins=200, color='r', alpha=0.7, label='NO GC Response Times')
    ax1.hist(server_times_2, bins=200, color='b', alpha=0.7, label='GC Response Times')
    ax1.set_xlabel('Time (milliseconds)')
    ax1.set_ylabel('Frequency')

    # Add titles and legends
    plt.title('Response Times')
    ax1.legend(loc='upper right')
    
    # Save the plot to the specified file
    plt.savefig(output_image_file)
    plt.close()
    
def plot_mix_hist(server_time, gc_indices, output_image_file):
    # Create a histogram with bins
    hist, bins = np.histogram(server_time, bins=200)

    # Create a mask for the highlighted indices
    mask = np.zeros_like(hist, dtype=bool)
    for idx in gc_indices:
        # Find which bin the index falls into and set that bin in the mask to True
        bin_idx = np.digitize([server_time[idx]], bins)[0] - 1
        mask[min(bin_idx, len(hist)-1)] = True

    # Plot the original histogram
    plt.bar(bins[:-1], hist, width=np.diff(bins), color='red', align='edge', label='General Data')

    # Plot the highlighted part of the histogram
    plt.bar(bins[:-1], hist * mask, width=np.diff(bins), color='blue', align='edge', label='Highlighted Indices')

    # Plot client times on the primary y-axis
    # ax1.hist(server_times, bins=200, color='r', alpha=0.7, label='Server Response Times')
    plt.xlabel('Time (milliseconds)')
    plt.ylabel('Frequency', color='g')

    # Add titles and legends
    plt.title('Histogram of Response Times')
    plt.legend(loc='upper right')
    
    # Save the plot to the specified file
    plt.savefig(output_image_file)
    plt.close()
    
def gc_requests(memory_log):
    GC_iterations = []
    for idx in range(1, len(memory_log)):
        # heapalloc, heapidle
        # mark iterations with HeapIdle increase or HeapAlloc decrease as GC calls
        if memory_log[idx][0] < memory_log[idx - 1][0]:
            # print("HeapAlloc")
            # print(memory_log[idx][0], memory_log[idx - 1][0])
            GC_iterations.append(idx)
        # elif memory_log[idx][1] > memory_log[idx - 1][1]:
            # print("HeapIdle")
            # print(memory_log[idx][1], memory_log[idx - 1][1])
            # GC_iterations.append(idx)
        idx += 1
    return GC_iterations

# Command-line arguments usage
if __name__ == "__main__":
    # if len(sys.argv) != 6:
    #     print("Usage: python script.py <client_time_file> <server_time_file> <memory_file> <dist_image_file> <latency_image_file>")
    #     sys.exit(1)
    # with open(sys.argv[1], 'r') as f:
        # client_times = [float(line.strip().split(', ')[1]) for line in f.readlines()]
    # print(client_times[:10])
    with open('./Case1/server_time.txt', 'r') as f:
        server_times = [int(line.strip().split(',')[1]) for line in f.readlines()]
    print(server_times[:10])
    server_times = server_times[100:]
    with open('./Case3/memory.txt', 'r') as f:
        memory_log, second_container = parse_memory_log(f)
    memory_log = memory_log[100:]
    with open('./Case3/server_time.txt', 'r') as f:
        server_times_2 = [int(line.strip().split(',')[1]) for line in f.readlines()]
    server_times_2 = server_times_2[100:]

    gc_indic = gc_requests(memory_log)
    print(gc_indic)
    plot_mix_hist(server_times_2, gc_indic, './problem_dist.pdf')
    # print(second_container)
    plot_histograms(server_times, server_times_2, './case_dist.pdf')
