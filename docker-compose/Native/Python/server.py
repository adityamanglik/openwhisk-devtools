from http.server import BaseHTTPRequestHandler, HTTPServer
import json
import math
import random
from urllib.parse import urlparse, parse_qs

PORT = 8800

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
        if self.head == node:
            self.head = self.head.next
            if self.head is None:
                self.tail = None
        else:
            current = self.head
            while current.next is not None and current.next != node:
                current = current.next
            if current.next == node:
                current.next = node.next
                if node.next is None:
                    self.tail = current

def generateRandomNormal(mean, stdDev):
    u1 = random.random()
    u2 = random.random()
    z0 = math.sqrt(-2 * math.log(u1)) * math.cos(2 * math.pi * u2)
    return z0 * stdDev + mean

def mainLogic(seed, ARRAY_SIZE, REQ_NUM):
    lst = LinkedList()

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
                sum_val += nestedCurrent.value
                nestedCurrent = nestedCurrent.next
        else:
            sum_val += current.value
        current = current.next

    response = {
        "sum": sum_val,
        "executionTime": 0,  # Placeholder for execution time calculation
        "requestNumber": REQ_NUM,
        "arraysize": ARRAY_SIZE,
        "usedHeapSize": 0,  # Placeholder for heap size calculation
        "totalHeapSize": 0  # Placeholder for total heap size calculation
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
            response = mainLogic(seed, ARRAY_SIZE, REQ_NUM)
            self.send_response(200)
            self.send_header('Content-type', 'application/json')
            self.end_headers()
            self.wfile.write(bytes(json.dumps(response), "utf8"))
        else:
            self.send_response(404)
            self.end_headers()

if __name__ == "__main__":
    server = HTTPServer(('localhost', PORT), RequestHandler)
    print("Server running on port", PORT)
    server.serve_forever()
