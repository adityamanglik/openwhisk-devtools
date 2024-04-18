from http.server import BaseHTTPRequestHandler, HTTPServer
import json
import math
import random
import time
import mmap
import os
import struct
from urllib.parse import urlparse, parse_qs

PORT = 9100

def initialize_data_sparse(mem):
    # Define the pattern to be written sparsely
    pattern = b'\x01'  # Example pattern

    # Write the pattern sparsely into the memory
    for i in range(0, len(mem), 1024 * 1024):  # Write 1 byte per megabyte
        mem[i:i+1] = pattern

def allocate_huge_pages(size_gb):
    # Calculate the size in bytes
    size_bytes = size_gb * (1024 ** 3)  # 1 GB = 1024^3 bytes

    # Open a temporary file to back the mmap object
    with open("/dev/zero", "r+b") as f:
        # Create a memory-mapped file using huge pages
        mem = mmap.mmap(f.fileno(), size_bytes, mmap.MAP_PRIVATE | mmap.MAP_ANONYMOUS)

    return mem

def mainLogic():
    start_time = time.perf_counter()
    huge_mem = allocate_huge_pages(20)
    initialize_data_sparse(huge_mem)
    huge_mem.close()
    end_time = time.perf_counter()
    duration_microseconds = (end_time - start_time) * 1_000_000
    return {"state": "finished", "exec_time": duration_microseconds}

class RequestHandler(BaseHTTPRequestHandler):
    def do_GET(self):
        parsed_path = urlparse(self.path)
        path = parsed_path.path
        if path.startswith("/Python"):
            response = mainLogic()
            self.send_response(200)
            self.send_header('Content-type', 'application/json')
            self.end_headers()
            self.wfile.write(bytes(json.dumps(response), "utf8"))
        else:
            self.send_response(404)
            self.end_headers()

if __name__ == "__main__":
    server = HTTPServer(('0.0.0.0', PORT), RequestHandler)
    print("Server running on port", PORT)
    server.serve_forever()
