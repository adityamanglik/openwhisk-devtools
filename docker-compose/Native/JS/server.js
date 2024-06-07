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
const fs = require('fs');
const path = require('path');
const { Buffer } = require('buffer'); // To handle binary data
const serverPort = 8800;

const server = http.createServer((req, res) => {
    if (req.url.startsWith('/JS')) {
        return jsonHandler(req, res);
    } else if (req.url.startsWith('/ImageProcess')) {
        return imageProcessHandler(req, res);
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
    const nestedCount = 10;
    
    let lst = new LinkedList();

    for (let i = 0; i < ARRAY_SIZE; i++) {
        const num = generateRandomNormal(seed, seed);
        lst.pushFront(num);
        
        // Stress GC with nested list
        if (i % 5 === 0) {
            let nestedList = new LinkedList();
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

async function imageProcessHandler(req, res) {
    const queryObject = url.parse(req.url, true).query;
    let seed = 42; // default seed value
    let ARRAY_SIZE = 1000; // default array size value

    if (queryObject.seed) {
        seed = parseInt(queryObject.seed);
    }

    if (queryObject.arraysize) {
        ARRAY_SIZE = parseInt(queryObject.arraysize);
    }

    if (queryObject.requestnumber) {
        REQ_NUM = parseInt(queryObject.requestnumber);
    }

    try {
        const jsonResponse = await imageLogic(seed, ARRAY_SIZE, REQ_NUM);
        res.writeHead(200, {'Content-Type': 'application/json'});
        res.end(JSON.stringify(jsonResponse));
    } catch (err) {
        res.writeHead(500, {'Content-Type': 'application/json'});
        res.end(JSON.stringify({ error: err.message }));
    }
}

async function imageLogic(seed, ARRAY_SIZE, REQ_NUM) {
    // console.log(`In ImageLogic`);
    const start = Date.now();

    const fileNames = ["Resources/img1.jpg", "Resources/img2.jpg"];
    const selectedFile = fileNames[Math.floor(Math.random() * fileNames.length)];
    const img = fs.readFileSync(path.join(__dirname, selectedFile));

    const imgBuffer = Buffer.from(img).toString('base64');
    const imgData = Buffer.from(imgBuffer, 'base64');

    // Process image (this example assumes a grayscale image for simplicity)
    let sum = 0;
    for (let i = 0; i < imgData.length; i++) {
        imgData[i] = clamp(imgData[i] + Math.floor(Math.random() * seed));
        sum += imgData[i];
    }

    // Resize (simple example, not true resize)
    const resizedData = imgData.slice(0, ARRAY_SIZE * ARRAY_SIZE);

    sum += sumPixels(resizedData);

    // Flip horizontally (simple example)
    const flippedData = flipHorizontally(resizedData, ARRAY_SIZE);

    sum += sumPixels(flippedData);

    // Rotate 90 degrees (simple example)
    const rotatedData = rotate(flippedData, ARRAY_SIZE, 90);

    sum += sumPixels(rotatedData);

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
    // console.log(`Response: ${JSON.stringify(response)}`);
    return response;
}

function resize(data, size) {
    return data.slice(0, size * size); // This is a simplification and not a true resize
}

function sumPixels(data) {
    return data.reduce((acc, val) => acc + val, 0);
}

function flipHorizontally(data, size) {
    const flippedData = new Uint8Array(data.length);
    for (let i = 0; i < size; i++) {
        for (let j = 0; j < size; j++) {
            flippedData[i * size + j] = data[i * size + (size - j - 1)];
        }
    }
    return flippedData;
}

function rotate(data, size, angle) {
    const rotatedData = new Uint8Array(data.length);
    if (angle === 90) {
        for (let i = 0; i < size; i++) {
            for (let j = 0; j < size; j++) {
                rotatedData[j * size + (size - i - 1)] = data[i * size + j];
            }
        }
    }
    return rotatedData;
}

function clamp(value) {
    return Math.max(0, Math.min(255, value));
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
