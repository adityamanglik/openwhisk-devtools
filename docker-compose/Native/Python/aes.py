from http.server import BaseHTTPRequestHandler, HTTPServer
import resource
import json
import random
import time
import psutil
import string
import pyaes
from urllib.parse import urlparse, parse_qs

PORT = 9900


def main_logic(seed, array_size, req_num):
    """Main logic function that builds linked lists, does nested operations, and sums up float values."""
    random.seed(seed)  # Ensure reproducibility with a given seed
    sum_val = 0
    # Start the timer
    start_time = time.perf_counter()

    # ADD LOGIC HERE ####################################################

    def generate(length):
        letters = string.ascii_lowercase + string.digits
        return ''.join(random.choice(letters) for i in range(length))

    # Use the baseline's array_size as the length of the message,
    # and req_num as the number of iterations.
    length_of_message = array_size
    num_of_iterations = 10

    message = generate(length_of_message)

    # 128-bit key (16 bytes)
    KEY = b'\xa1\xf6%\x8c\x87}_\xcd\x89dHE8\xbf\xc9,'

    for loops in range(num_of_iterations):
        aes = pyaes.AESModeOfOperationCTR(KEY)
        ciphertext = aes.encrypt(message)
        # Uncomment the following lines to print the ciphertext and plaintext if desired
        # print(ciphertext)
        aes = pyaes.AESModeOfOperationCTR(KEY)
        plaintext = aes.decrypt(ciphertext)
        sum_val += len(plaintext)

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
