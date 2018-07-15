package main

import (
	"bytes"
	"strconv"
	"crypto/sha256"
	"time"
	"github.com/boltdb/bolt"
)

type Block struct {
	Timestamp     int64
	Data          []byte
	Hash          []byte
	Nonce         int
	PrevBlockHash []byte
}

type Blockchain struct {
	tip []byte
	db  *bolt.DB
}

func (b *Block) SetHash() {
	timestamp := []byte(strconv.FormatInt(b.Timestamp, 10))
	headers := bytes.Join([][]byte{b.PrevBlockHash, b.Data, timestamp}, []byte{})
	hash := sha256.Sum224(headers)
	b.Hash = hash[:]
}

func (chain *Blockchain) AddBlock(data string) {
	var lastHash []byte

	err := chain.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))
		lastHash = b.Get([]byte("l"))
		return nil
	})
	warning(err, "error while reading tip from db")

	newBlock := NewBlock(data, lastHash)

	err = chain.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))
		err := b.Put(newBlock.Hash, newBlock.Serialize())
		warning(err, "error while saving new block in db")
		if err != nil {
			return err
		}
		err = b.Put([]byte("l"), newBlock.Hash)
		warning(err, "error while saving tip in db")
		if err != nil {
			return err
		}
		chain.tip = newBlock.Hash

		return nil
	})
}

func NewBlock(data string, prevBlockHash []byte) *Block {
	block := &Block{Timestamp: time.Now().Unix(), Data: []byte(data), PrevBlockHash: prevBlockHash}
	pow := NewProofOfWork(block)
	nonce, hash := pow.Run()

	block.Hash = hash[:]
	block.Nonce = nonce
	return block
}

func NewGenesisBlock() *Block {
	return NewBlock("Genesis block", []byte{})
}

const dbFile = "./db.bolt"
const blocksBucket = "blocksBucket"

func NewBlockchain() *Blockchain {
	var tip []byte
	db, err := bolt.Open(dbFile, 0600, nil)
	warning(err, "error while opening db file "+dbFile)
	err = db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))
		if b == nil {
			genesis := NewGenesisBlock()
			b, err := tx.CreateBucket([]byte(blocksBucket))
			warning(err, "error while creating bucket")
			if err != nil {
				return err
			}
			err = b.Put(genesis.Hash, genesis.Serialize())
			warning(err, "error while saving genesis block")
			if err != nil {
				return err
			}
			err = b.Put([]byte("l"), genesis.Hash)
			warning(err, "error while saving tip of the blockchain")
			if err != nil {
				return err
			}
			tip = genesis.Hash
		} else {
			tip = b.Get([]byte("l"))
		}
		return nil
	})
	return &Blockchain{tip, db}
}

func main() {
	chain := NewBlockchain()
	defer chain.db.Close()

	cli := CLI{chain}
	cli.Run()
}
