import matplotlib.pyplot as plt
import sys

def read_data(file_path):
    heap_alloc = []
    heap_idle = []
    heap_inuse = []

    with open(file_path, 'r') as file:
        for line in file:
            parts = line.strip().split(", ")
            if len(parts) != 3:
                continue
            alloc, idle, inuse = [int(part.split(": ")[1]) for part in parts]
            heap_alloc.append(alloc)
            heap_idle.append(idle)
            heap_inuse.append(inuse)

    return heap_alloc, heap_idle, heap_inuse

def plot_data(heap_alloc, heap_idle, heap_inuse, output_path):
    plt.figure(figsize=(10, 6))

    plt.plot(heap_alloc, label='HeapAlloc')
    plt.plot(heap_idle, label='HeapIdle')
    plt.plot(heap_inuse, label='HeapInuse')

    plt.xlabel('Sample Number')
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
