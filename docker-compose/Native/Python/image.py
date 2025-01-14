from http.server import BaseHTTPRequestHandler, HTTPServer
import resource
import json
import math
import random
import time
import psutil
import gc
import os
from urllib.parse import urlparse, parse_qs
from PIL import Image, ImageOps

PORT = 9900

def limit_memory(max_memory):
    soft, hard = resource.getrlimit(resource.RLIMIT_AS)
    resource.setrlimit(resource.RLIMIT_AS, (max_memory, hard))

class ListNode:
    def __init__(self, value):
        self.value = value
        self.next = None

class ListNode:
    def __init__(self, value):
        self.value = value
        self.next = None

class LinkedList:
    def __init__(self):
        self.head = None
        self.tail = None

    def pushFront(self, value):
        newNode = ListNode(value)
        if self.head is None:
            self.head = newNode
            self.tail = newNode
        else:
            newNode.next = self.head
            self.head = newNode

    def pushBack(self, value):
        newNode = ListNode(value)
        if self.tail is None:
            self.head = newNode
            self.tail = newNode
        else:
            self.tail.next = newNode
            self.tail = newNode

    def remove(self, node):
        # Handle the case where the list is empty
        if self.head is None:
            return

        # If the node to be removed is the head
        if self.head == node:
            self.head = self.head.next
            if self.head is None:  # If the list becomes empty
                self.tail = None
            return

        # If the node to be removed is not the head
        current = self.head
        while current.next is not None and current.next != node:
            current = current.next

        # If the node was found in the list
        if current.next == node:
            current.next = node.next
            if node.next is None:  # If the node is the tail
                self.tail = current

def generateRandomNormal(mean, stdDev):
    u1 = random.random()
    u2 = random.random()
    z0 = math.sqrt(-2 * math.log(u1)) * math.cos(2 * math.pi * u2)
    return z0 * stdDev + mean

def ImageLogic(seed, ARRAY_SIZE, REQ_NUM):
    # Start the timer
    start_time = time.perf_counter()
    
    # file_names = ["Resources/img1.jpg", "Resources/img2.jpg"]
    # selected_file = file_names[random.randint(0, len(file_names) - 1)]
    selected_file = "Resources/img1.jpg"
    if not os.path.exists(selected_file):
        raise FileNotFoundError(f"File {selected_file} not found.")
    
    img = Image.open(selected_file).convert('RGB')
    # Calculate the duration
    end_time = time.perf_counter()
    duration_seconds = end_time - start_time
    print(f"load: {duration_seconds}")
    
    # Resize the image
    img = img.resize((ARRAY_SIZE, ARRAY_SIZE))
    sum_val = sumPixels(img)
        
    # Add random seed to every pixel
    pixels = img.load()
    r1 = random.randint(0, 255)
    r2 = random.randint(0, 255)
    r3 = random.randint(0, 255)
    for y in range(img.height):
        for x in range(img.width):
            r, g, b = pixels[x, y]
            r = clamp(r + r1)
            g = clamp(g + r2)
            b = clamp(b + r3)
            pixels[x, y] = (r, g, b)
            
    # Calculate the duration
    end_time = time.perf_counter()
    duration_seconds = end_time - start_time
    print(f"randomization: {duration_seconds}")
    
    # Calculate the duration
    end_time = time.perf_counter()
    duration_seconds = end_time - start_time
    print(f"resize: {duration_seconds}")

    # Flip horizontally
    img = ImageOps.mirror(img)
    sum_val += sumPixels(img)
    # Calculate the duration
    end_time = time.perf_counter()
    duration_seconds = end_time - start_time
    print(f"mirror: {duration_seconds}")

    # Rotate 90 degrees
    img = img.rotate(90, expand=True)
    sum_val += sumPixels(img)
    # Calculate the duration
    end_time = time.perf_counter()
    duration_seconds = end_time - start_time
    print(f"rotate: {duration_seconds}")

    
    # Calculate the duration
    end_time = time.perf_counter()
    duration_seconds = end_time - start_time

    # Convert duration to microseconds
    duration_microseconds = duration_seconds * 1_000_000
    
     # Collect memory usage statistics
    process = psutil.Process()
    
    memory_info = process.memory_info()
    memory_full_info = process.memory_full_info()

    # Print all available statistics
    # print("memory_info:")
    # for attr in dir(memory_info):
    #     if not attr.startswith('_'):
    #         print(f"{attr}: {getattr(memory_info, attr)}")

    # print("\n----------------------\nmemory_full_info:")
    # for attr in dir(memory_full_info):
    #     if not attr.startswith('_'):
    #         print(f"{attr}: {getattr(memory_full_info, attr)}")
    
    response = {
        "sum": sum_val,
        "executionTime": duration_microseconds,  # Placeholder for execution time calculation
        "requestNumber": REQ_NUM,
        "arraysize": ARRAY_SIZE,
        "usedHeapSize": memory_full_info.uss,  # Placeholder for heap size calculation
        "totalHeapSize": memory_info.vms  # Placeholder for total heap size calculation
    }
    return response

def clamp(value):
    return max(0, min(255, value))

def sumPixels(img):
    return sum(sum(pixel[:3]) for pixel in img.getdata())

def mainLogic(seed, ARRAY_SIZE, REQ_NUM):
    lst = LinkedList()
    # Start the timer
    start_time = time.perf_counter()

    for i in range(ARRAY_SIZE):
        num = generateRandomNormal(seed, seed)
        lst.pushFront(num)

        if i % 5 == 0:
            nestedList = LinkedList()
            for j in range(10):
                nestedList.pushBack(generateRandomNormal(seed, seed))
            lst.pushBack(nestedList)

        if i % 5 == 0:
            tempNum = generateRandomNormal(seed, seed)
            lst.pushFront(tempNum)
            lst.remove(lst.head)

    sum_val = 0
    current = lst.head
    while current is not None:
        if isinstance(current.value, ListNode):
            nestedCurrent = current.value.head
            while nestedCurrent is not None:
            # Here, we ensure nestedCurrent.value is a float before adding
                if isinstance(nestedCurrent.value, float):
                    sum_val += nestedCurrent.value
                nestedCurrent = nestedCurrent.next
        elif isinstance(current.value, float):  # Ensure current.value is a float
            sum_val += current.value
        current = current.next
    # End the timer
    end_time = time.perf_counter()

    # Calculate the duration
    duration_seconds = end_time - start_time

    # Convert duration to microseconds
    duration_microseconds = duration_seconds * 1_000_000
    
     # Collect memory usage statistics
    process = psutil.Process()
    
    memory_info = process.memory_info()
    memory_full_info = process.memory_full_info()

    # Print all available statistics
    print("memory_info:")
    for attr in dir(memory_info):
        if not attr.startswith('_'):
            print(f"{attr}: {getattr(memory_info, attr)}")

    print("\n----------------------\nmemory_full_info:")
    for attr in dir(memory_full_info):
        if not attr.startswith('_'):
            print(f"{attr}: {getattr(memory_full_info, attr)}")
    
    response = {
        "sum": sum_val,
        "executionTime": duration_microseconds,  # Placeholder for execution time calculation
        "requestNumber": REQ_NUM,
        "arraysize": ARRAY_SIZE,
        "usedHeapSize": memory_full_info.uss,  # Placeholder for heap size calculation
        "totalHeapSize": memory_info.vms  # Placeholder for total heap size calculation
    }
    return response

class RequestHandler(BaseHTTPRequestHandler):
    def do_GET(self):
        parsed_path = urlparse(self.path)
        path = parsed_path.path
        query_components = parse_qs(parsed_path.query)

        seed = 42
        ARRAY_SIZE = 10000
        REQ_NUM = 2**53-1  # Equivalent to Number.MAX_SAFE_INTEGER in JavaScript

        if 'seed' in query_components:
            seed = int(query_components['seed'][0])
        if 'arraysize' in query_components:
            ARRAY_SIZE = int(query_components['arraysize'][0])
        if 'requestnumber' in query_components:
            REQ_NUM = int(query_components['requestnumber'][0])

        if path.startswith("/Python"):
            response = ImageLogic(seed, ARRAY_SIZE, REQ_NUM)
            self.send_response(200)
            self.send_header('Content-type', 'application/json')
            self.end_headers()
            self.wfile.write(bytes(json.dumps(response), "utf8"))
        elif path.startswith("/ImageProcess"):
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
