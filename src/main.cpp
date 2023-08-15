#include <iostream>
#include <openssl/rand.h>
#include <openssl/sha.h>
#include <openssl/ripemd.h>
#include <secp256k1.h>
#include <secp256k1_ecdh.h>
#include <stdint.h>

class Wallet{
    private:
        uint8_t privateKey[65];
        uint8_t publicKey[65];
        secp256k1_pubkey pubkey; // utilizado apenas para gerar a chave bruta
        secp256k1_context* ctx;
        unsigned char keypair[96];
        char address[25]; // P2PHK Address
        
    public:
        Wallet(void);
        ~Wallet(void);
        void createAddress();
        uint8_t* getPrivateKey(void);
        const char *getPublicKey(void);
        const char *getAddress(void);

};




Wallet::Wallet() {
    RAND_poll();
    this->ctx = secp256k1_context_create(SECP256K1_CONTEXT_SIGN | SECP256K1_CONTEXT_VERIFY);
    if(RAND_bytes(privateKey, 64) != 1){
        std::cout << "Error on generate private key!" << std::endl;
    }    
    int result = secp256k1_ec_pubkey_create(ctx, &pubkey, privateKey);
    size_t length = 65;
    secp256k1_ec_pubkey_serialize(ctx, publicKey, &length, &pubkey, SECP256K1_EC_UNCOMPRESSED);
    printf("public key :  ");
    for (size_t i=0; i<65; ++i){
        printf("%02x", publicKey[i]);
    }
    std::cout << std::endl;
}

Wallet::~Wallet() {
    secp256k1_context_destroy(ctx);
}

void Wallet::createAddress(){
    unsigned char output_sha[SHA256_DIGEST_LENGTH];
    SHA256(publicKey, sizeof(publicKey), output_sha);
    unsigned char output_ripemd160[RIPEMD160_DIGEST_LENGTH];
    RIPEMD160(output_sha, sizeof(output_sha), output_ripemd160);
    unsigned char checksum_ripemd160[1 + RIPEMD160_DIGEST_LENGTH];
    checksum_ripemd160[0] = 0x00; // byte de versão da criptomoeda
    for (uint8_t i = 0; i < RIPEMD160_DIGEST_LENGTH; i++) {
        checksum_ripemd160[i + 1] = output_ripemd160[i];
    }
    unsigned char sha256_step1[SHA256_DIGEST_LENGTH];
    SHA256(checksum_ripemd160, sizeof(checksum_ripemd160), sha256_step1);
    unsigned char sha256_step2[SHA256_DIGEST_LENGTH];
    SHA256(sha256_step1, sizeof(sha256_step1), sha256_step2);
    
}


uint8_t* Wallet::getPrivateKey(){
    int i;
    printf("private key : ");
    for (i=0; i<32; ++i){
        printf("%02x", privateKey[i]);
    }
    std::cout << std::endl;
    return privateKey;
}

int main(void){

    Wallet *wallet = new Wallet;
    wallet->getPrivateKey();
    return 0;
}
