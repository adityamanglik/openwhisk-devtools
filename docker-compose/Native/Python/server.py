from http.server import BaseHTTPRequestHandler, HTTPServer
import json
import random
import os
import signal
import time
import sys

# MARKER_FOR_SIZE_UPDATE
ARRAY_SIZE = 10000

class RequestHandler(BaseHTTPRequestHandler):
    def do_GET(self):
        if self.path == "/GoNative":
            self.handle_request()
        else:
            self.send_error(404, "File not found")

    def handle_request(self):
        seed = 42  # default seed value
        query = self.path.split("?")
        if len(query) > 1:
            params = query[1].split("&")
            for p in params:
                if p.startswith("seed="):
                    seed = int(p.split("=")[1])

        response, execution_time = self.main_logic(seed)

        self.send_response(200)
        self.send_header('Content-type', 'application/json')
        self.end_headers()
        self.wfile.write(json.dumps(response).encode())

        # Uncomment the line below to print execution time to console
        # print(f"Request processed in {execution_time} ns")

    @staticmethod
    def main_logic(seed):
        start = time.time_ns()
        random.seed(seed)
        arr = [random.randint(0, 100000) for _ in range(ARRAY_SIZE)]
        sum_val = sum(arr)
        execution_time = time.time_ns() - start

        response = {
            "sum": sum_val,
            "executionTime": execution_time  # raw execution time in nanoseconds
        }
        return response, execution_time

def graceful_shutdown(signum, frame):
    print('Shutting down server...')
    sys.exit(0)

def run(server_class=HTTPServer, handler_class=RequestHandler, port=9900):
    server_address = ('', port)
    httpd = server_class(server_address, handler_class)
    signal.signal(signal.SIGINT, graceful_shutdown)
    signal.signal(signal.SIGTERM, graceful_shutdown)
    print(f'Server listening on http://localhost:{port}')
    httpd.serve_forever()

if __name__ == '__main__':
    run(port=9900)
