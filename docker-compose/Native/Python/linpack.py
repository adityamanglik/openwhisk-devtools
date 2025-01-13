from http.server import BaseHTTPRequestHandler, HTTPServer
import resource
import json
import random
import time
import psutil
from urllib.parse import urlparse, parse_qs
import numpy as np

PORT = 9900


def main_logic(seed, array_size, req_num):
    """Main logic function that builds linked lists, does nested operations, and sums up float values."""
    random.seed(seed)  # Ensure reproducibility with a given seed
    sum_val = 0
    # Start the timer
    start_time = time.perf_counter()

    # ADD LOGIC HERE ####################################################
    n = array_size  # Use the array_size parameter as the matrix size for LINPACK
    # Estimate the number of operations performed (in FLOPs)
    ops = (2.0 * n) * n * n / 3.0 + (2.0 * n) * n

    # Create an n x n array of random numbers in the range [-0.5, 0.5]
    A = np.random.random_sample((n, n)) - 0.5
    B = A.sum(axis=1)

    # Convert A and B into matrices
    A = np.matrix(A)
    B = np.matrix(B.reshape((n, 1)))

    # Solve the linear system (Ax = B) and measure the latency
    linpack_start = time.perf_counter()
    x = np.linalg.solve(A, B)
    linpack_latency = time.perf_counter() - linpack_start

    # Calculate MFLOPS from the operations count and latency
    mflops = (ops * 1e-6) / linpack_latency

    # For compatibility with the baseline response,
    # assign the computed MFLOPS as the 'sum' value.
    sum_val = mflops
    # END LOGIC HERE ####################################################

    end_time = time.perf_counter()
    duration_seconds = end_time - start_time
    duration_microseconds = duration_seconds * 1_000_000

    # Memory usage (optional; might slow things down if done too frequently)
    process = psutil.Process()
    memory_info = process.memory_info()
    memory_full_info = process.memory_full_info()

    response = {
        "sum": sum_val,
        "executionTime": duration_microseconds,
        "requestNumber": req_num,
        "arraysize": array_size,
        "usedHeapSize": memory_full_info.uss,  # Unique set size
        "totalHeapSize": memory_info.vms       # Virtual memory size
    }
    return response

class RequestHandler(BaseHTTPRequestHandler):
    def do_GET(self):
        parsed_path = urlparse(self.path)
        path = parsed_path.path
        query_components = parse_qs(parsed_path.query)

        # Defaults
        seed = 42
        array_size = 10000
        req_num = 2**53 - 1  # Just as in your original code

        if 'seed' in query_components:
            seed = int(query_components['seed'][0])
        if 'arraysize' in query_components:
            array_size = int(query_components['arraysize'][0])
        if 'requestnumber' in query_components:
            req_num = int(query_components['requestnumber'][0])

        if path.startswith("/Python"):
            response = main_logic(seed, array_size, req_num)
            self.send_response(200)
            self.send_header('Content-type', 'application/json')
            self.end_headers()
            self.wfile.write(bytes(json.dumps(response), "utf8"))
        else:
            self.send_response(404)
            self.end_headers()


if __name__ == "__main__":
    # Optional: limit_memory(512 * 1024 * 1024)  # e.g., 512 MB
    server = HTTPServer(('0.0.0.0', PORT), RequestHandler)
    print("Server running on port", PORT)
    server.serve_forever()
