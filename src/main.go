package main

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"fmt"
	"math/big"

	"crypto/sha256"
	"log"
	"testando/blockchain"

	"github.com/btcsuite/btcutil/base58"
	"golang.org/x/crypto/ripemd160"
)

type wallet struct {
	PrivateKey string
	PublicKey  []byte
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
		string(privateKeyHash[:]),
		publicKeyBytes,
		generate_p2phk(publicKeyBytes[:]),
	}
}

func (w *wallet) GetPrivateKey() string {
	fmt.Printf("Private Key ( 64 bits ) : %x\n", w.PrivateKey)
	return w.PrivateKey
}

func (w *wallet) GetPublicKey() []byte {
	fmt.Printf("Public Key : %x\n", w.PublicKey[:])
	return w.PublicKey
}

func (w *wallet) GetAddress() string {
	fmt.Printf("Address ( P2PHK ) : %s\n", w.Address)
	return w.Address
}

func main() {
	wallet := Wallet()
	fmt.Println("Gerando sua primeira wallet na nossa aplicação :\n")
	wallet.GetPrivateKey()
	wallet.GetPublicKey()
	addr := wallet.GetAddress()
	fmt.Println("-------")
	Blockchain := blockchain.Blockchain{
		[]blockchain.Block{
			blockchain.GenesisBlock(),
		},
		4,
		[]blockchain.Transaction{},
		50,
	}
	Block, err := Blockchain.GetLatestBlock()
	Block2 := blockchain.Block{}
	Block2.MineBlock(
		[]blockchain.Transaction{
			blockchain.Transaction{
				"LGCoin",
				addr,
				float64(Blockchain.MiningReward),
			},
		},
		Block.PreviousHash,
		Blockchain.Difficulty,
	)
	fmt.Printf("Hash found : %s\n", Block2.Hash)
	if err != nil {
		log.Fatal("Erro!")
	}
}
