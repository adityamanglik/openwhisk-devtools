import matplotlib.pyplot as plt
import csv
import numpy as np

def extract_percentiles(file_path):
    try:
        with open(file_path, 'r', newline='') as file:
            csv_reader = csv.reader(file)
            last_line = None
            for row in csv_reader:
                last_line = row
            print(last_line)
            if last_line is None:
                return None

            headers = ["Type","Name","Request Count","Failure Count","Median Response Time","Average Response Time","Min Response Time","Max Response Time","Average Content Size","Requests/s","Failures/s","50%","66%","75%","80%","90%","95%","98%","99%","99.9%","99.99%","100%"]
            percentiles_labels = ["50%", "95%", "99%", "99.9%", "99.99%", "100%"]
            header_indices = [headers.index(header) for header in percentiles_labels]

            percentiles = {header: last_line[index] for header, index in zip(headers, header_indices)}
            return percentiles

    except FileNotFoundError:
        print(f"File not found: {file_path}")
        return None
    except Exception as e:
        print(f"An error occurred: {e}")
        return None

def plot_sla(output_file):
    percentiles_4g = extract_percentiles('./4g/1.csv')
    print(percentiles_4g)
    percentiles_128m = extract_percentiles('./128m/1.csv')

    if percentiles_4g is None or percentiles_128m is None:
        print("Could not load all data.")
        return

    percentiles_labels = ["50%", "95%", "99%", "99.9%", "99.99%", "100%"]
    percentile_values_4g = [float(percentiles_4g[p]) for p in percentiles_labels]
    percentile_values_128m = [float(percentiles_128m[p]) for p in percentiles_labels]

    plt.figure(figsize=(10, 6))
    plt.plot(percentiles_labels, percentile_values_4g, marker='o', label='4G Network')
    plt.plot(percentiles_labels, percentile_values_128m, marker='o', label='128M Network')

    plt.xlabel('Percentile')
    plt.ylabel('Response Time (ms)')
    plt.title('Response Time by Percentile Distribution for 4G and 128M Networks')
    plt.xticks(rotation=45)
    plt.grid(True)
    plt.legend(loc='upper left')
    plt.savefig(output_file)
    plt.close()

def main():
    plot_sla("./SLA_Memory_Impact.pdf")

if __name__ == "__main__":
    main()
