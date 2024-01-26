const Benchmark = require('benchmark');

const suite = new Benchmark.Suite();

const gradients = new Array(256).fill(null).map(() =>
    new Array(256).fill(null).map(() => {
        const angle = Math.random() * 2.0 * Math.PI;
        return [Math.cos(angle), Math.sin(angle)];
    })
);

let perlinTotalTime = 0;

suite.add('perlin', function () {
    const startTime = performance.now();
    perlin(0.5, 0.5, gradients);
    const endTime = performance.now();
    perlinTotalTime += endTime - startTime;
  })
  .on('cycle', function (event) {
    console.log(String(event.target));
  })
  .on('complete', function () {
    console.log('Benchmark completed.');
    console.log('Total time for perlin:', (perlinTotalTime / 1000).toFixed(3), 'seconds');
  })
  .run({ 'async': true });
