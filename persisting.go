package main

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"github.com/boltdb/bolt"
)

func (b *Block) Serialize() []byte {
	var result bytes.Buffer
	encoder := gob.NewEncoder(&result)

	err := encoder.Encode(b)
	if err != nil {
		fmt.Println("error while encoding", err)
	}
	return result.Bytes()
}

func DeserializeBlock(b []byte) *Block {
	var block Block

	decoder := gob.NewDecoder(bytes.NewReader(b))
	err := decoder.Decode(&block)
	if err != nil {
		fmt.Println("error while decoding", err)
	}
	return &block
}

type BlockchainIterator struct {
	currentHash []byte
	db          *bolt.DB
}

func (bc *Blockchain) Iterator() *BlockchainIterator {
	return &BlockchainIterator{bc.tip, bc.db}
}

func (i *BlockchainIterator) Next() *Block {
	var block *Block

	err := i.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))
		data := b.Get(i.currentHash)
		block = DeserializeBlock(data)
		return nil
	})
	warning(err, "error while retrieving block from the db")

	i.currentHash = block.PrevBlockHash
	return block
}
