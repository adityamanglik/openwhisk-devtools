import matplotlib.pyplot as plt
import sys

def read_data(file_path):
    heap_alloc = []
    heap_idle = []
    heap_inuse = []

    with open(file_path, 'r') as file:
        for line in file:
            parts = line.strip().split(", ")
            if '9501' in parts[1]:
                continue
            alloc = int(parts[2].split(": ")[1])
            idle = int(parts[3].split(": ")[1])
            inuse = int(parts[4].split(": ")[1])
            heap_alloc.append(alloc)
            heap_idle.append(idle)
            heap_inuse.append(inuse)

    return heap_alloc, heap_idle, heap_inuse

def plot_data(heap_alloc, heap_idle, heap_inuse, output_path):
    plt.figure(figsize=(10, 6))

    plt.plot(heap_alloc, label='HeapAlloc')
    plt.plot(heap_idle, label='HeapIdle')
    # plt.plot(heap_inuse, label='HeapInuse')

    plt.xlabel('Request Number')
    plt.ylabel('Memory (bytes)')
    plt.title('Heap Memory Usage Over Time')
    plt.legend()
    
    plt.savefig(output_path)
    plt.close()
    # plt.show()

# Path to your log file
file_path = sys.argv[1]
output_path = sys.argv[2]

# Read and plot data
heap_alloc, heap_idle, heap_inuse = read_data(file_path)
plot_data(heap_alloc, heap_idle, heap_inuse, output_path)
