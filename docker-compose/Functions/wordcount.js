function generateRandomNormal(mean, stdDev) {
  // Box-Muller transform to generate a random number from a normal distribution
  const u1 = Math.random();
  const u2 = Math.random();
  const z0 = Math.sqrt(-2 * Math.log(u1)) * Math.cos(2 * Math.PI * u2);
  return z0 * stdDev + mean;
}

function main(params) {
  // // 1. Generate a random number X between 1000 - 10000 using a standard normal distribution
  // const mean = 55000*params.seed;
  // const stdDev = 25000;
  // let X = Math.floor(generateRandomNormal(mean, stdDev));
  // if (X < 10000) X = 10000;
  // if (X > 100000) X = 100000;

  // If seed is not provided in params, default to 42
  const seedValue = params.seed ? params.seed : 42;

  // 2. Create an array
  const randomNumbers = [];

  // 3. Write X random numbers into the array
  for (let i = 0; i < 1000000; i++) {
    const randomNum = Math.random() * seedValue;
    randomNumbers.push(randomNum);
  }

  // Calculate the sum of the array
  const arraySum = randomNumbers.reduce((acc, val) => acc + val, 0);

  return { payload: `Seed: ${seedValue}, The sum of the array values is ${arraySum}` };
}

// For testing purposes
//   console.log(main({}));
