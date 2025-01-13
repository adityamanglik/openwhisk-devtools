from http.server import BaseHTTPRequestHandler, HTTPServer
import resource
import json
import random
import time
import psutil
from urllib.parse import urlparse, parse_qs

PORT = 9900


def limit_memory(max_memory):
    """Limit the process to use a maximum of `max_memory` bytes."""
    soft, hard = resource.getrlimit(resource.RLIMIT_AS)
    resource.setrlimit(resource.RLIMIT_AS, (max_memory, hard))


class ListNode:
    def __init__(self, value):
        self.value = value
        self.next = None


class LinkedList:
    def __init__(self):
        self.head = None
        self.tail = None

    def push_front(self, value):
        new_node = ListNode(value)
        if self.head is None:
            self.head = new_node
            self.tail = new_node
        else:
            new_node.next = self.head
            self.head = new_node

    def push_back(self, value):
        new_node = ListNode(value)
        if self.tail is None:
            self.head = new_node
            self.tail = new_node
        else:
            self.tail.next = new_node
            self.tail = new_node

    def remove(self, node):
        """Remove the given node from the linked list."""
        if self.head is None:
            return
        # If the node to be removed is the head
        if self.head == node:
            self.head = self.head.next
            if self.head is None:
                self.tail = None
            return
        # Otherwise, walk the list to find the node
        current = self.head
        while current.next and current.next != node:
            current = current.next
        if current.next == node:
            current.next = node.next
            if node.next is None:  # node was the tail
                self.tail = current


def generate_random_normal_builtin(mean, std):
    """
    Faster approach: use Python's built-in random.gauss for normal distribution.
    """
    return random.gauss(mean, std)


def main_logic(seed, array_size, req_num):
    """Main logic function that builds linked lists, does nested operations, and sums up float values."""
    random.seed(seed)  # Ensure reproducibility with a given seed

    # Start the timer
    start_time = time.perf_counter()

    linked_list = LinkedList()

    for i in range(array_size):
        # Generate normal distribution using built-in function
        num = generate_random_normal_builtin(seed, seed)
        linked_list.push_front(num)

        if i % 5 == 0:
            # Build a small nested linked list
            nested_list = LinkedList()
            for j in range(10):
                nested_list.push_back(generate_random_normal_builtin(seed, seed))
            linked_list.push_back(nested_list)

        if i % 5 == 0:
            # push_front + remove head
            temp_num = generate_random_normal_builtin(seed, seed)
            linked_list.push_front(temp_num)
            # Remove the (new) head immediately
            linked_list.remove(linked_list.head)

    # Summation
    sum_val = 0.0
    current = linked_list.head
    while current:
        value = current.value
        # If value is a LinkedList, we iterate its nodes
        if isinstance(value, LinkedList):
            nested_current = value.head
            while nested_current:
                # Accumulate only floats
                if isinstance(nested_current.value, float):
                    sum_val += nested_current.value
                nested_current = nested_current.next
        elif isinstance(value, float):
            sum_val += value
        current = current.next

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
