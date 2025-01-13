from http.server import BaseHTTPRequestHandler, HTTPServer
import json
import random
import time
from urllib.parse import urlparse, parse_qs

PORT = 9900

def main_logic(seed, array_size, req_num):
    pass

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
