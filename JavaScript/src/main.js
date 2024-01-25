console.time("Execution Time");
const PNG = require('pngjs').PNG;
const fs = require('fs');

function smoothstep(w) {
    if (w <= 0.0) {
        return 0.0;
    }
    if (w >= 1.0) {
        return 1.0;
    }
    return w * w * (3.0 - 2.0 * w);
}

function interpolate(a0, a1, w) {
    return a0 + (a1 - a0) * smoothstep(w);
}

function dotGridGradient(ix, iy, x, y, gradients) {
    if (iy >= 0 && iy < gradients.length && ix >= 0 && ix < gradients[0].length) {
        const dx = x - ix;
        const dy = y - iy;

        if (gradients[iy].length > 0 && gradients[iy][ix].length >= 2) {
            return dx * gradients[iy][ix][0] + dy * gradients[iy][ix][1];
        }
    }

    return 0.0;
}

function perlin(x, y, gradients) {
    const x0 = Math.floor(x);
    const x1 = x0 + 1;
    const y0 = Math.floor(y);
    const y1 = y0 + 1;

    const sx = x - x0;
    const sy = y - y0;

    let n0, n1, ix0, ix1, value;

    n0 = dotGridGradient(x0, y0, x, y, gradients);
    n1 = dotGridGradient(x1, y0, x, y, gradients);
    ix0 = interpolate(n0, n1, sx);

    n0 = dotGridGradient(x0, y1, x, y, gradients);
    n1 = dotGridGradient(x1, y1, x, y, gradients);
    ix1 = interpolate(n0, n1, sx);

    value = interpolate(ix0, ix1, sy);

    return value;
}

const gridSize = 256;
const width = 800;
const height = 800;

const gradients = new Array(gridSize).fill(null).map(() =>
    new Array(gridSize).fill(null).map(() => {
        const angle = Math.random() * 2.0 * Math.PI;
        return [Math.cos(angle), Math.sin(angle)];
    })
);

const imgData = new Uint8Array(width * height * 4);

const png = new PNG({
    width: width,
    height: height,
    colorType: 6, // RGBA
});

// Remplir le buffer de l'image avec des valeurs de bruit de Perlin
for (let i = 0; i < width; i++) {
    for (let j = 0; j < height; j++) {
        let x = i / width;
        let y = j / height;
        const frequency = 5.0;
        x *= frequency;
        y *= frequency;

        let total = 0.0;
        let amplitude = 1.0;
        const octaves = 10;
        const persistence = 0.85;

        for (let o = 0; o < octaves; o++) {
            total += perlin(x, y, gradients) * amplitude;
            x *= 2;
            y *= 2;
            amplitude *= persistence;
        }

        const value = (total + 1.0) / 2.0;

        let col;
        if (value < 0.2) {
            col = [0, 0, 255, 255]; // Bleu pour l'eau
        } else if (value < 0.3) {
            col = [194, 178, 128, 255]; // Beige pour le sable
        } else if (value < 0.5) {
            col = [34, 139, 34, 255]; // Vert pour la forêt
        } else if (value < 0.7) {
            col = [139, 69, 19, 255]; // Marron pour la montagne
        } else {
            col = [255, 250, 250, 255]; // Blanc pour la neige
        }

        const baseIndex = (j * width + i) * 4;
        imgData[baseIndex] = col[0];
        imgData[baseIndex + 1] = col[1];
        imgData[baseIndex + 2] = col[2];
        imgData[baseIndex + 3] = col[3];
    }
}

// Convertir les données brutes dans le format PNG
for (let i = 0; i < imgData.length; i++) {
    png.data[i] = imgData[i];
}

// Écrire le fichier PNG
png.pack()
    .pipe(fs.createWriteStream('perlin_noise.png'))
    .on('finish', () => console.log('Image saved.'));

console.timeEnd("Execution Time");

/******* BENCHMARK *******/
const Benchmark = require('benchmark');

const suite = new Benchmark.Suite();

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
