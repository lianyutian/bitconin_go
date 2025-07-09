package main

import (
	"time"
)

type Block struct {
	// 区块产生时间
	Timestamp int64
	// 区块内容
	Data []byte
	// 前一个区块 Hash
	PrevBlockHash []byte
	// 本身区块 Hash
	Hash  []byte
	Nonce int
}

func NewBlock(data string, prevBlockHash []byte) *Block {
	block := &Block{time.Now().Unix(), []byte(data), prevBlockHash, []byte{}, 0}
	pow := NewProofOfWork(block)
	nonce, hash := pow.Run()

	block.Hash = hash[:]
	block.Nonce = nonce
	return block
}
