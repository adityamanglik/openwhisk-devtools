import matplotlib.pyplot as plt
import csv
import numpy as np
import os

def extract_percentiles(file_path):
    try:
        with open(file_path, 'r', newline='') as file:
            csv_reader = csv.reader(file)
            last_line = None
            for row in csv_reader:
                last_line = row
            # print(last_line)
            if last_line is None:
                return None

            headers = ["Type","Name","Request Count","Failure Count","Median Response Time","Average Response Time","Min Response Time","Max Response Time","Average Content Size","Requests/s","Failures/s","50%","66%","75%","80%","90%","95%","98%","99%","99.9%","99.99%","100%"]
            percentiles_labels = ["50%", "95%", "99%", "99.9%", "99.99%", "100%"]
            header_indices = [headers.index(header) for header in percentiles_labels]
            # print(header_indices)
            percentiles = {header: last_line[index] for header, index in zip(percentiles_labels, header_indices)}
            return percentiles

    except FileNotFoundError:
        print(f"File not found: {file_path}")
        return None
    except Exception as e:
        print(f"An error occurred: {e}")
        return None

def plot_sla(output_file):
    plt.figure(figsize=(10, 6))
    # Extract files
    # List of files: files_4g = ['12000rpm.csv', '1200rpm.csv', '120rpm.csv', '1rpm.csv', '2rpm.csv', '3000rpm.csv', '30rpm.csv', '6000rpm.csv', '60rpm.csv', '6rpm.csv']
    # Retrieve all csv files from both directories
    files_4g = [f for f in os.listdir('./4g') if f.endswith('.csv')]
    files_4g = ['1rpm.csv', '1200rpm.csv', '12000rpm.csv']#, '120rpm.csv', '2rpm.csv', '3000rpm.csv', '30rpm.csv', '6000rpm.csv',  '6rpm.csv']
    
    files_128m = [f for f in os.listdir('./128m') if f.endswith('.csv')]
    files_128m = ['1rpm.csv', '1200rpm.csv', '12000rpm.csv']
    # Sorting files to ensure matching pairs are processed together if required
    # files_4g.sort()
    # files_4g = files_4g[:3]
    # files_128m.sort()
    # files_128m = files_128m[:3]
    percentiles_labels = ["50%", "95%", "99%", "99.9%", "99.99%", "100%"]
    
    # Get a colormap from matplotlib
    colormap = plt.cm.get_cmap('plasma', len(files_4g) + len(files_128m))  # Using 'viridis' but you can choose any other
    colors = [colormap(i) for i in range(len(files_4g))]
    
    for i, (f4g, f128m) in enumerate(zip(files_4g, files_128m)):
        file_path_4g = os.path.join('./4g', f4g)
        file_path_128m = os.path.join('./128m', f128m)
        # Start plotting
        percentiles_4g = extract_percentiles(file_path_4g)
        print(file_path_4g, percentiles_4g)
        percentiles_128m = extract_percentiles(file_path_128m)
        # print(file_path_128m, percentiles_128m)

        if percentiles_4g is None or percentiles_128m is None:
            print("Could not load all data.")
            return

        
        percentile_values_4g = [float(percentiles_4g[p]) for p in percentiles_labels]
        percentile_values_128m = [float(percentiles_128m[p]) for p in percentiles_labels]

        plt.scatter(percentiles_labels, percentile_values_4g, s=50, edgecolor='black', linewidth=0.5, marker='o', color=colors[i], label=f'4G-{f4g[:-4]}')
        plt.scatter(percentiles_labels, percentile_values_128m, s=100, edgecolor='black', linewidth=0.5, marker='*', color=colors[i], label=f'128M-{f128m[:-4]}')

    plt.xlabel('Percentile')
    plt.ylabel('Response Time (ms)')
    plt.title('Response Time Percentile Distribution')
    # plt.xticks(rotation=45)
    # plt.yscale('symlog')
    plt.grid(True)
    plt.legend(loc='upper left', ncols=3)
    plt.savefig(output_file)
    plt.close()

def main():
    plot_sla("./SLA_Memory_Impact.pdf")

if __name__ == "__main__":
    main()
