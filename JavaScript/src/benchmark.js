const Benchmark = require('benchmark');

// Importez votre classe NoiseGenerator depuis le fichier approprié
const NoiseGenerator = require('./lib/NoiseGenerator');

// Création d'une instance de NoiseGenerator avec un seed
var seed = Math.random() * 10000
const myNoiseGenerator = new NoiseGenerator(seed);

// Création d'un nouveau benchmark
const suite = new Benchmark.Suite;

// Ajout de la fonction à tester
suite.add('PerlinNoise Generation', function() {
  // Appel de la méthode perlinNoise à tester
  myNoiseGenerator.perlinNoise(1, 1);
});
suite.add('Noise Generation', function() {
    myNoiseGenerator.noise(1, 1);
});

// Écoute des résultats du benchmark
suite
  .on('cycle', function(event) {
    console.log(String(event.target)); // Affiche le résultat du cycle
    console.log(`Temps d'exécution: ${event.target.times.elapsed.toFixed(6)} s`);
  })
  .on('complete', function() {
    console.log('Résultats du benchmark :');
    for (let i = 0; i < this.length; i++) {
      console.log(`${this[i].name}: ${this[i].hz} ops/sec ±${this[i].stats.rme.toFixed(2)}%`);
    }
  });

// Exécution du benchmark
suite.run();
