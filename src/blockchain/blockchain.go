package blockchain

import (
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
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
	block.Transactions = append(block.Transactions, Transactions...)
	block.PreviousHash = PreviousHash
	block.Nonce = 0
	difficultyTarget := strings.Repeat("0", int(Difficulty))
	hashBlock := block.HashBlock()
	//fmt.Println("Mining started...")
	//startMining := time.Now()
	var nonce uint64
	for nonce = 1; hashBlock[:Difficulty] != difficultyTarget; nonce++ {
		block.Nonce = nonce
		hashBlock = block.HashBlock()
	}
	//endMining := time.Now()
	//timeOfMining := endMining.Sub(startMining)
	//fmt.Printf("Total time of mining : %s\n", timeOfMining)
	block.Hash = hashBlock
}

func (blockchain *Blockchain) IsValidChain() bool {
	lastBlock, err := blockchain.GetLatestBlock()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(lastBlock)
	return true
}

func (blockchain *Blockchain) AddBlock(block Block) error {
	if len(blockchain.Chain) == 0 {
		blockchain.Chain = append(blockchain.Chain, block)
		return nil
	}

	lastBlock, err := blockchain.GetLatestBlock()
	if err != nil {
		log.Fatal(err)
	}
	if lastBlock.HashBlock() != block.PreviousHash {
		return errors.New("Previous hash block invalid!")
	}

	if block.Hash[:blockchain.Difficulty] != strings.Repeat("0", int(blockchain.Difficulty)) {
		return errors.New("Block hash invalid!")
	}

	file, err := os.Create(fmt.Sprintf("blocks/blk%05d.db", len(blockchain.Chain)))
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	encoder := json.NewEncoder(file)
	err = encoder.Encode(block)
	if err != nil {
		log.Fatal(err)
	}
	blockchain.Chain = append(blockchain.Chain, block)
	return nil
}

func GenesisBlock() Block {
	block := Block{
		time.Now().Unix(),
		[]Transaction{},
		"0",
		0,
		"0",
	}
	_, err := os.Stat("blocks")
	if err != nil {
		if os.IsNotExist(err) {
			os.Mkdir("blocks", 0755)
			file, err := os.Create(fmt.Sprintf("blocks/blk00000.db"))
			if err != nil {
				log.Fatal(err)
			}
			defer file.Close()
			encoder := json.NewEncoder(file)
			err = encoder.Encode(block)
			if err != nil {
				log.Fatal(err)
			}
		}
	}

	return block
}

func (blockchain *Blockchain) GetBlockByIndex(index uint64) (Block, error) {
	var i uint64
	if uint64(len(blockchain.Chain)) > index {
		return (blockchain.Chain[i]), nil
	}
	log.Fatal("Error: Block index isnt on blockchain!")
	return Block{}, nil
}

func (blockchain *Blockchain) GetBalanceByAddress(address string) float64 {
	// lembrar de adicionar verificação de transações de um endereço para si mesmo
	balance := float64(0)
	for _, block := range blockchain.Chain {
		for _, transaction := range block.Transactions {
			if transaction.Sender == address {
				balance = balance - transaction.Amount
			} else if transaction.Recipient == address {
				balance = balance + transaction.Amount
			}
		}
	}
	return balance
}

func LoadBlockchain() Blockchain {
	_, err := os.Stat("blocks")
	if err != nil {
		if os.IsNotExist(err) {
			return Blockchain{
				[]Block{
					GenesisBlock(),
				},
				4,
				[]Transaction{},
				50,
			}
		}
	}
	loadedBlockchain := Blockchain{
		[]Block{},
		4,
		[]Transaction{},
		50,
	}
	blocks, err := os.ReadDir("blocks")
	currentBlock := Block{}
	var fileBlock *os.File
	var decoder *json.Decoder

	if err != nil {
		log.Fatal(err)
	}

	for _, entry := range blocks {
		if !entry.IsDir() {
			fileBlock, err = os.Open(fmt.Sprintf("blocks/%s", entry.Name()))
			if err != nil {
				log.Fatal(err)
			}
			decoder = json.NewDecoder(fileBlock)
			err = decoder.Decode(&currentBlock)
			if err != nil {
				log.Fatal(err)
			}
			loadedBlockchain.Chain = append(loadedBlockchain.Chain, currentBlock)

			currentBlock = Block{}
			fileBlock.Close()
		}
	}
	return loadedBlockchain
}
