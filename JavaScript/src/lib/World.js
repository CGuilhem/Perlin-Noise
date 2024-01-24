class World {
    constructor(noiseGenerator) {
        this.chunkManager = new ChunkManager(noiseGenerator)
        this.renderDistance = 10
        this.loadDistance = 2
    }

    draw(renderer) {
        const chunks = this.chunkManager.chunks
        for (var iChunk in chunks) {
            const chunk = chunks[iChunk]
            const chunkElement = chunk.element
            chunkElement.draw(renderer)
        }
    }

    loadChunks() {
        var isBuffered = false
        for (var i = 0; i < this.loadDistance; i++) {
            const minX = 0
            const minZ = 0
            const maxX = i
            const maxZ = i
            for (var x = minX; x < maxX; x++) {
                for (var z = minZ; z < maxZ; z++) {
                    this.chunkManager.load(x * CHUNK_SIZE, z * CHUNK_SIZE)
                    isBuffered = this.chunkManager.addToBuffer(x * CHUNK_SIZE, z * CHUNK_SIZE)
                }
            }
            if(isBuffered){
                break
            }
        }
        if(!isBuffered){
            this.loadDistance++
        }
        if(this.loadDistance > this.renderDistance){
            this.loadDistance = 2
        }
    }
}