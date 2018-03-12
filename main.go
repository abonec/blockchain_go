package main

import (
	"bytes"
	"strconv"
	"crypto/sha256"
	"time"
	"fmt"
)

type Block struct {
	Timestamp     int64
	Data          []byte
	Hash          []byte
	PrevBlockHash []byte
}

type Blockchain struct {
	blocks []*Block
}

func (b *Block) SetHash() {
	timestamp := []byte(strconv.FormatInt(b.Timestamp, 10))
	headers := bytes.Join([][]byte{b.PrevBlockHash, b.Data, timestamp}, []byte{})
	hash := sha256.Sum224(headers)
	b.Hash = hash[:]
}

func (chain *Blockchain) AddBlock(data string) {
	prevBlock := chain.blocks[len(chain.blocks)-1]
	block := NewBlock(data, prevBlock.Hash)
	chain.blocks = append(chain.blocks, block)
}

func NewBlock(data string, prevBlockHash []byte) *Block {
	block := &Block{Timestamp: time.Now().Unix(), Data: []byte(data), PrevBlockHash: prevBlockHash}
	block.SetHash()
	return block
}

func NewGenesisBlock() *Block {
	return NewBlock("Genesis block", []byte{})
}

func NewBlockchain() *Blockchain {
	return &Blockchain{blocks: []*Block{NewGenesisBlock()}}
}

func main() {
	chain := NewBlockchain()

	chain.AddBlock("Send 1 BTC to somebody")
	chain.AddBlock("Send 100 BTC to someone")

	for _, block := range chain.blocks {
		fmt.Printf("PrevHash: %x\n", block.PrevBlockHash)
		fmt.Printf("Data: %s\n", block.Data)
		fmt.Printf("Hash: %x\n", block.Hash)
		fmt.Printf("=================\n")
	}
}
