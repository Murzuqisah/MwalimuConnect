package functions

type Blockchain struct {
	Blocks []Block
}


func (chain *Blockchain) AddBlock(data []byte) {
	previousBlock := chain.Blocks[len(chain.Blocks)-1]
	newBlock := GenerateBlock(previousBlock, data)
	chain.Blocks = append(chain.Blocks, newBlock)
}
