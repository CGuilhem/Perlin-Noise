class Application {
    constructor(title, renderer, noiseGenerator) {
        this.title = title
        this.renderer = renderer
        this.world = new World(noiseGenerator)
        this.startTime = Date.now()
        this.nbFrame = 0
        this.runLoop = this.runLoop.bind(this)
    }

    runLoop() {
        this.updateFPS()
        this.world.loadChunks()
        this.world.draw(this.renderer)
        this.renderer.render()
        requestAnimationFrame(this.runLoop)
    }

    updateFPS() {
        const deltaTime = (Date.now() - this.startTime) / 1000
        if (deltaTime > 1) {
            document.title = `${title} - (${parseInt(this.nbFrame / deltaTime)} FPS)`
            this.nbFrame = 0
            this.startTime = Date.now()
        } else {
            this.nbFrame++
        }
    }
}