function generateRandomNormal(mean, stdDev) {
    // Box-Muller transform to generate a random number from a normal distribution
    const u1 = Math.random();
    const u2 = Math.random();
    const z0 = Math.sqrt(-2 * Math.log(u1)) * Math.cos(2 * Math.PI * u2);
    return z0 * stdDev + mean;
  }
  
  function main(params) {
    // 1. Generate a random number X between 1000 - 10000 using a standard normal distribution
    const mean = 55000;
    const stdDev = 25000;
    let X = Math.floor(generateRandomNormal(mean, stdDev));
    if (X < 10000) X = 10000;
    if (X > 100000) X = 100000;
  
    // 2. Create an array
    const randomNumbers = [];
  
    // 3. Write X random numbers into the array
    for (let i = 0; i < X; i++) {
      const randomNum = Math.random();
      randomNumbers.push(randomNum);
    }
  
    // 4. Count the length of the array and return it as word count
    const arrayLength = randomNumbers.length;
  
    return { payload: `Word count is ${arrayLength}` };
  }
  
  // For testing purposes
//   console.log(main({}));
  