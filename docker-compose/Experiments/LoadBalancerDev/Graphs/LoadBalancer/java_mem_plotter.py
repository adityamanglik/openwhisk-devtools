import matplotlib.pyplot as plt
import sys

def read_data(file_path):
    heap_used = []
    heap_committed = []
    heap_max = []

    with open(file_path, 'r') as file:
        for line in file:
            parts = line.strip().split(", ")
            if len(parts) != 3:
                continue
            used, committed, max_memory = [int(part.split(": ")[1]) for part in parts]
            heap_used.append(used)
            heap_committed.append(committed)
            heap_max.append(max_memory)

    return heap_used, heap_committed, heap_max

def plot_data(heap_used, heap_committed, heap_max, output_path):
    plt.figure(figsize=(10, 6))

    plt.plot(heap_used, label='HeapUsedMemory')
    plt.plot(heap_committed, label='HeapCommittedMemory')
    plt.plot(heap_max, label='HeapMaxMemory')

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
heap_used, heap_committed, heap_max = read_data(file_path)
plot_data(heap_used, heap_committed, heap_max, output_path)
