package main

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"fmt"
	"math/big"
	"os"

	"crypto/sha256"
	"log"
	"testando/blockchain"

	"github.com/btcsuite/btcutil/base58"
	"golang.org/x/crypto/ripemd160"

	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

type wallet struct {
	PrivateKey string
	PublicKey  string
	Address    string
}

func generate_p2phk(publicKey []byte) string {
	step1 := sha256.Sum256(publicKey[:])
	ripemd160Hash := ripemd160.New()
	ripemd160Hash.Write(step1[:])
	ripemd_checksum := ripemd160Hash.Sum(nil)
	ripemd_output := []byte{0x00}
	ripemd_output = append(ripemd_output, ripemd_checksum[:]...)
	step2 := sha256.Sum256(ripemd_output[:])
	step3 := sha256.Sum256(step2[:])
	step4 := ripemd_output
	step4 = append(step4, step3[:4]...)
	step5 := base58.Encode(step4)
	return step5
}

func Wallet() wallet {
	curve := elliptic.P256()

	// carregando 8 bytes aleatórios com alta entropia
	seed := make([]byte, 8)
	rand.Read(seed)
	//

	// gerando uma chave privada a partir de um seed númerico de 2**64 bits
	privateKeyHash := sha256.Sum256(seed)
	//

	// gerando chave pública usando a curva eliptica secp256k1
	privateKey := new(ecdsa.PrivateKey)
	privateKey.PublicKey.Curve = curve
	privateKey.D = new(big.Int).SetBytes(privateKeyHash[:])
	privateKey.PublicKey.X, privateKey.PublicKey.Y = curve.ScalarBaseMult(privateKey.D.Bytes())
	publicKeyBytes := elliptic.Marshal(curve, privateKey.PublicKey.X, privateKey.PublicKey.Y)
	//

	return wallet{
		fmt.Sprintf("%x", privateKeyHash[:]),
		fmt.Sprintf("%x", publicKeyBytes),
		generate_p2phk(publicKeyBytes[:]),
	}
}

func (w *wallet) GetPrivateKey() string {
	return w.PrivateKey
}

func (w *wallet) GetPublicKey() string {
	return w.PublicKey
}

func (w *wallet) GetAddress() string {
	return w.Address
}

func main() {
	var wallet wallet
	var addr string
	_, err := os.Stat("wallets/wallet.db")
	if err != nil {
		if os.IsNotExist(err) {
			os.Mkdir("wallets", 0755)
			wallet = Wallet()
			fmt.Println("Gerando sua primeira wallet na nossa aplicação :\n")
			privKey := wallet.GetPrivateKey()
			publicKey := wallet.GetPublicKey()
			addr = wallet.GetAddress()
			db, err := sql.Open("sqlite3", "wallets/wallet.db")
			if err != nil {
				fmt.Println(err)
			}
			defer db.Close()
			_, err = db.Exec(`
			CREATE TABLE IF NOT EXISTS wallet (
				id INTEGER PRIMARY KEY,
				privateKey TEXT,
				publicKey TEXT
			)
			`)
			if err != nil {
				log.Fatal(err)
			}
			_, err = db.Exec(`
			CREATE TABLE IF NOT EXISTS addresses (
				id INTEGER PRIMARY KEY,
				address TEXT
			)
			`)
			if err != nil {
				log.Fatal(err)
			}
			_, err = db.Exec(`INSERT INTO wallet(privateKey, publicKey) VALUES (?,?)`, privKey, publicKey)
			if err != nil {
				fmt.Println(err)
			}
			_, err = db.Exec(`INSERT INTO addresses(address) VALUES (?)`, addr)
			if err != nil {
				fmt.Println(err)
			}
		}
	} else {
		fmt.Println("Wallet loaded!")
		db, err := sql.Open("sqlite3", "wallets/wallet.db")
		if err != nil {
			fmt.Println(err)
		}
		defer db.Close()
		rows, err := db.Query("SELECT privateKey, publicKey FROM wallet")
		if err != nil {
			log.Fatal(err)
		}
		defer rows.Close()
		for rows.Next() {
			rows.Scan(&wallet.PrivateKey, &wallet.PublicKey)
		}
		rows, err = db.Query("SELECT address FROM addresses")
		if err != nil {
			log.Fatal(err)
		}
		defer rows.Close()
		for rows.Next() {
			rows.Scan(&wallet.Address)
		}
	}

	fmt.Println("-------")
	fmt.Printf("Address ( P2PHK ) : %s\n", wallet.Address)
	Blockchain := blockchain.LoadBlockchain()
	Block, err := Blockchain.GetLatestBlock()
	if err != nil {
		log.Fatal(err)
	}
	Block2 := blockchain.Block{}
	Block2.MineBlock(
		[]blockchain.Transaction{
			blockchain.Transaction{
				"LGCoin",
				wallet.Address,
				float64(Blockchain.MiningReward),
			},
		},
		Block.HashBlock(),
		Blockchain.Difficulty,
	)
	Blockchain.AddBlock(Block2)
	fmt.Printf("Balance : %.8f LG Coins\n", Blockchain.GetBalanceByAddress(wallet.Address))
	fmt.Println("-------")
	if err != nil {
		log.Fatal("Erro!")
	}

}
