function generateRandomNormal(mean, stdDev) {
  // Box-Muller transform to generate a random number from a normal distribution
  const u1 = Math.random();
  const u2 = Math.random();
  const z0 = Math.sqrt(-2 * Math.log(u1)) * Math.cos(2 * Math.PI * u2);
  return z0 * stdDev + mean;
}

class ListNode {
  constructor(value) {
      this.value = value;
      this.next = null;
  }
}

class LinkedList {
  constructor() {
      this.head = null;
      this.tail = null;
  }

  pushFront(value) {
      const newNode = new ListNode(value);
      if (this.head === null) {
          this.head = newNode;
          this.tail = newNode;
      } else {
          newNode.next = this.head;
          this.head = newNode;
      }
  }

  pushBack(value) {
      const newNode = new ListNode(value);
      if (this.tail === null) {
          this.head = newNode;
          this.tail = newNode;
      } else {
          this.tail.next = newNode;
          this.tail = newNode;
      }
  }

  // JavaScript does not have direct memory management like Go, so we simulate removal
  remove(node) {
      if (this.head === node) {
          this.head = this.head.next;
          if (this.head === null) {
              this.tail = null;
          }
      } else {
          let current = this.head;
          while (current.next !== null && current.next !== node) {
              current = current.next;
          }
          if (current.next === node) {
              current.next = node.next;
              if (node.next === null) {
                  this.tail = current;
              }
          }
      }
  }
}


const http = require('http');
const url = require('url');
const serverPort = 8800;

const server = http.createServer((req, res) => {
    if (req.url.startsWith('/JS')) {
        return jsonHandler(req, res);
    }
    res.writeHead(404);
    res.end();
});

function jsonHandler(req, res) {
    const queryObject = url.parse(req.url, true).query;
    let seed = 42; // default seed value
    let ARRAY_SIZE = 10000; // default array size value
    let REQ_NUM = Number.MAX_SAFE_INTEGER; // default request number

    if (queryObject.seed) {
        seed = parseInt(queryObject.seed);
    }

    if (queryObject.arraysize) {
        ARRAY_SIZE = parseInt(queryObject.arraysize);
    }

    if (queryObject.requestnumber) {
        REQ_NUM = parseInt(queryObject.requestnumber);
    }

    const jsonResponse = mainLogic(seed, ARRAY_SIZE, REQ_NUM);
    res.writeHead(200, {'Content-Type': 'application/json'});
    res.end(JSON.stringify(jsonResponse));
}

function mainLogic(seed, ARRAY_SIZE, REQ_NUM) {
    const start = Date.now();

    let lst = new LinkedList();

    for (let i = 0; i < ARRAY_SIZE; i++) {
        const num = generateRandomNormal(seed, seed);
        lst.pushFront(num);
        
        // Stress GC with nested list
        if (i % 5 === 0) {
            let nestedList = new LinkedList();
            const nestedCount = generateRandomNormal(seed, seed);
            for (let j = 0; j < nestedCount; j++) {
                nestedList.pushBack(generateRandomNormal(seed, seed));
            }
            lst.pushBack(nestedList);
        }

        // Immediate removal after insertion to stress GC
        if (i % 5 === 0) {
            const tempNum = generateRandomNormal(seed, seed);
            lst.pushFront(tempNum);
            lst.remove(lst.head); // Removing the recently added node
        }
    }

    // Sum values and return result
    let sum = 0;
    let current = lst.head;
    while (current !== null) {
        if (typeof current.value === 'number') {
            sum += current.value;
        } else if (current.value instanceof LinkedList) {
            let nestedCurrent = current.value.head;
            while (nestedCurrent !== null) {
                sum += nestedCurrent.value;
                nestedCurrent = nestedCurrent.next;
            }
        }
        current = current.next;
    }

    const executionTime = Date.now() - start;

    // Calculate size of heap
    // const { performance } = require('perf_hooks');
    const usedHeapSize = process.memoryUsage().heapUsed;
    const totalHeapSize = process.memoryUsage().heapTotal;
    // Again, Node.js doesn't provide a direct equivalent to jsHeapSizeLimit, but you can get the resident set size
    // const residentSetSize = process.memoryUsage().rss;

    const response = {
        sum: sum,
        executionTime: executionTime,
        requestNumber: REQ_NUM,
        arraysize: ARRAY_SIZE,
        usedHeapSize: usedHeapSize,
        totalHeapSize: totalHeapSize
    };
    return response;
}

server.listen(serverPort, () => {
    console.log(`Server listening on http://localhost:${serverPort}`);
});

// Graceful shutdown logic
process.on('SIGINT', () => {
    console.log('Shutting down server...');
    server.close(() => {
        console.log('Server shut down gracefully.');
        process.exit(0);
    });
});
