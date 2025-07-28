package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"encoding/hex"
	"fmt"
	"log"
)

// 挖矿奖励
const subsidy = 10

type TXInput struct {
	Txid      []byte
	Vout      int
	ScriptSig string
}

type TXOutput struct {
	Value        int
	ScriptPubKey string
}

type Transaction struct {
	ID   []byte
	Vin  []TXInput
	Vout []TXOutput
}

// IsCoinbase checks whether the transaction is coinbase
func (tx Transaction) IsCoinbase() bool {
	return len(tx.Vin) == 1 && len(tx.Vin[0].Txid) == 0 && tx.Vin[0].Vout == -1
}

// SetID 为交易设置唯一的标识符。
// 该函数通过将交易的详细信息编码并计算哈希值来生成一个唯一的ID。
// 这个ID对于在整个系统中唯一标识一个交易至关重要。
func (tx *Transaction) SetID() {
	// 创建一个字节缓冲区来存储编码后的交易信息。
	var encoded bytes.Buffer
	// 创建一个32字节的数组来存储哈希值。
	var hash [32]byte

	// 创建一个gob编码器，用于将交易信息编码到缓冲区中。
	enc := gob.NewEncoder(&encoded)
	// 将交易信息编码。
	// 如果编码过程中出现错误，记录错误并终止程序。
	err := enc.Encode(tx)
	if err != nil {
		log.Panic(err)
	}

	// 使用SHA-256算法计算编码后的交易信息的哈希值。
	hash = sha256.Sum256(encoded.Bytes())
	// 将计算出的哈希值赋给交易的ID字段，作为交易的唯一标识符。
	tx.ID = hash[:]
}

func (in *TXInput) CanUnlockOutputWith(unlockingData string) bool {
	return in.ScriptSig == unlockingData
}

func (out *TXOutput) CanBeUnlockedWith(unlockingData string) bool {
	return out.ScriptPubKey == unlockingData
}

// NewCoinbaseTX 铸币交易
func NewCoinbaseTX(to, data string) *Transaction {
	if data == "" {
		data = fmt.Sprintf("Reward to '%s'", to)
	}

	txin := TXInput{[]byte{}, -1, data}
	txout := TXOutput{subsidy, to}
	tx := Transaction{nil, []TXInput{txin}, []TXOutput{txout}}
	tx.SetID()

	return &tx
}

func NewUTXOTransaction(from, to string, amount int, bc *Blockchain) *Transaction {
	var inputs []TXInput
	var outputs []TXOutput

	balance, validOutputs := bc.FindSpendableOutputs(from, amount)

	if balance < amount {
		log.Panic("ERROR: Not enough funds")
	}

	for txid, outs := range validOutputs {
		txID, err := hex.DecodeString(txid)
		if err != nil {
			log.Panic(err)
		}

		for _, outIdx := range outs {
			input := TXInput{txID, outIdx, from}
			inputs = append(inputs, input)
		}
	}

	outputs = append(outputs, TXOutput{amount, to})
	if balance > amount {
		outputs = append(outputs, TXOutput{balance - amount, from})
	}

	tx := Transaction{nil, inputs, outputs}
	tx.SetID()

	return &tx
}
