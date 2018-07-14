package main

import (
	"math/big"
	"bytes"
	"encoding/binary"
	"log"
	"fmt"
	"math"
	"crypto/sha256"
	"runtime"
)

const targetBits = 24

type ProofOfWork struct {
	block  *Block
	target *big.Int
}

func NewProofOfWork(b *Block) *ProofOfWork {
	target := big.NewInt(1)
	target.Lsh(target, uint(256-targetBits))

	pow := &ProofOfWork{b, target}
	return pow
}

// IntToHex converts an int64 to a byte array
func IntToHex(num int64) []byte {
	buff := new(bytes.Buffer)
	err := binary.Write(buff, binary.BigEndian, num)
	if err != nil {
		log.Panic(err)
	}

	return buff.Bytes()
}

func (pow *ProofOfWork) prepareData(nonce int) []byte {
	data := bytes.Join(
		[][]byte{
			pow.block.PrevBlockHash,
			pow.block.Data,
			IntToHex(pow.block.Timestamp),
			IntToHex(int64(targetBits)),
			IntToHex(int64(nonce)),
		},
		[]byte{},
	)

	return data
}

const maxNonce = math.MaxInt64

type work struct {
	FromNonce int
	ToNonce   int
}

type result struct {
	Nonce int
	Hash  []byte
}

func (pow *ProofOfWork) Run() (int, []byte) {

	resultChan := make(chan result)
	workChan := make(chan work, 1)
	for workers := 0; workers <= runtime.NumCPU(); workers++ {
		go func(resultChan chan result, workChan chan work) {
			for batch := range workChan {
				var hashInt big.Int
				var hash [32]byte
				for nonce := batch.FromNonce; nonce < batch.ToNonce; nonce ++ {
					data := pow.prepareData(nonce)
					hash = sha256.Sum256(data)
					//fmt.Printf("\r%x", hash)
					hashInt.SetBytes(hash[:])

					if hashInt.Cmp(pow.target) == -1 {
						resultChan <- result{nonce, hash[:]}
						return
					} else {
						nonce++
					}

				}

			}
		}(resultChan, workChan)

	}
	fmt.Printf("Mining the block containing \"%s\"\n", pow.block.Data)
	minNonce := 0
	for nonce := 1000; nonce < maxNonce; nonce += 1000 {
		select {
		case r := <-resultChan:
			return r.Nonce, r.Hash
		default:
			batch := work{minNonce, nonce}
			workChan <- batch
			minNonce = nonce
		}
	}

	fmt.Println("\n\n")
	return 0, []byte{}
}

func (pow *ProofOfWork) Validate() bool {
	var hashInt big.Int

	data := pow.prepareData(pow.block.Nonce)
	hash := sha256.Sum256(data)
	hashInt.SetBytes(hash[:])

	return hashInt.Cmp(pow.target) == -1
}
