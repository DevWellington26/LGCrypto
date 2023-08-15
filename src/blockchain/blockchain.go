package blockchain

import (
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"
)

type Blockchain struct {
	Chain                []Block       `json:"chain"`
	Difficulty           uint64        `json:"difficulty"`
	PenddingTransactions []Transaction `json:"transactions"`
	MiningReward         uint64        `json:"mining_reward"`
}

type Block struct {
	Timestamp    int64         `json:"timestamp"`
	Transactions []Transaction `json:"transactions"`
	PreviousHash string        `json:"previous_hash"`
	Nonce        uint64        `json:"nonce"`
	Hash         string        `json:"hash"`
}

type Transaction struct {
	Sender    string  `json:"sender"`
	Recipient string  `json:"recipient"`
	Amount    float64 `json:"amount"`
}

func (block *Block) HashBlock() string {
	jsonData, err := json.Marshal(block)
	if err != nil {
		fmt.Println("Error : ", err)
	}
	return fmt.Sprintf("%x", sha256.Sum256(jsonData))
}

func (blockchain *Blockchain) GetAllBlockchain() {

}

func (blockchain *Blockchain) GetLatestBlock() (Block, error) {
	if len(blockchain.Chain) == 0 {
		return Block{}, errors.New("Blockchain empty!")
	}
	return blockchain.Chain[len(blockchain.Chain)-1], nil
}

func (block *Block) MineBlock(Transactions []Transaction, PreviousHash string, Difficulty uint64) {
	block.Timestamp = time.Now().Unix()
	_ = append(block.Transactions, Transactions...)
	block.PreviousHash = PreviousHash
	block.Nonce = 0
	difficultyTarget := strings.Repeat("0", int(Difficulty))
	hashBlock := block.HashBlock()
	fmt.Println("Mining started...")
	startMining := time.Now()
	var nonce uint64
	for nonce = 1; hashBlock[:Difficulty] != difficultyTarget; nonce++ {
		block.Nonce = nonce
		hashBlock = block.HashBlock()
	}
	endMining := time.Now()
	timeOfMining := endMining.Sub(startMining)
	fmt.Printf("Total time of mining : %s\n", timeOfMining)
	block.Hash = hashBlock
}

func GenesisBlock() Block {

	return Block{
		time.Now().Unix(),
		[]Transaction{},
		"0",
		0,
		"0",
	}
}

func (blockchain *Blockchain) GetBlockByIndex(index uint64) (Block, error) {
	var i uint64
	if uint64(len(blockchain.Chain)) > index {
		return (blockchain.Chain[i]), nil
	}
	log.Fatal("Error: Block index isnt on blockchain!")
	return Block{}, nil
}
