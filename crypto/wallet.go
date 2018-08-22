package crypto

import (
	"os"
	"fmt"
	"io/ioutil"
	"encoding/json"
	"crypto/aes"
	"io"
	"crypto/rand"
	"crypto/cipher"
	"crypto/sha512"
	"strings"
	"blockchain/crypto"
)

type PlainKeys struct {
	Keys map[string]string
	Checksum [64]byte
}

type keyPairs map[string]string //[public_key]private_key

type Wallet struct {
	cipher 		[]byte
	name 		string
	keyPairs
	password 	[64]byte
}

func NewWallet() *Wallet {
	return &Wallet{
		cipher: 	make([]byte, 0),
		name:		"wallet.json",
		keyPairs: 	make(map[string]string, 0),
		password: 	[64]byte{},
	}
}

func (wallet *Wallet) SetPassword(password string) {
	wallet.password = sha512.Sum512([]byte(password))
}

func (wallet *Wallet) GetPrivateKey(publicKey PublicKey) (*PrivateKey, error) {
	if wif, ok := wallet.keyPairs[publicKey.String()]; ok {
		return NewPrivateKey(wif)
	}
	return nil, fmt.Errorf("private key doesn't exist")
}

func (wallet *Wallet) DecryptKeyPairs() error {
	block, err := aes.NewCipher(wallet.password[0:32])
	if err != nil {
		return err
	}
	decryptedData := make([]byte, len(wallet.cipher[aes.BlockSize:]))
	decryptStream := cipher.NewCTR(block, wallet.cipher[:aes.BlockSize])
	decryptStream.XORKeyStream(decryptedData, wallet.cipher[aes.BlockSize:])
	err = json.Unmarshal(decryptedData, &wallet.keyPairs)
	if err != nil {
		return err
	}
	return nil
}

func (wallet *Wallet) GetKeyPairs() map[string]string{
	return wallet.keyPairs
}

func (wallet *Wallet) PublicKeys() []*PublicKey {
	pubKeys := make([]*PublicKey, 0)
	for pubKeyString := range wallet.keyPairs {
		pubKey, _ := NewPublicKey(pubKeyString)
		pubKeys = append(pubKeys, pubKey)
	}
	return pubKeys
}

func (wallet *Wallet) NewKeyPair() string {
	privateKey,_ := crypto.NewRandomPrivateKey()
	publicKey := privateKey.PublicKey().String()
	wallet.keyPairs[publicKey] = privateKey.String()
	return publicKey
}

func (wallet *Wallet) ImportPrivateKey(wif string) error {
	privateKey, err := NewPrivateKey(wif)
	if err != nil {
		return err
	}
	publicKey := privateKey.PublicKey().String()
	if _, ok := wallet.keyPairs[publicKey]; ok {
		return fmt.Errorf("duplicated key")
	}
	wallet.keyPairs[publicKey] = privateKey.String()
	return nil
}

func (wallet *Wallet) RemoveKeyPair(publicKey string) bool {
	if _, ok := wallet.keyPairs[publicKey]; ok {
		delete(wallet.keyPairs, publicKey)
		wallet.SaveToFile()
		return true
	}
	return false
}

func (wallet *Wallet) SaveToFile() error {
	wallet.EncryptKeyPairs()
	keyData, err := json.Marshal(wallet.cipher)
	if err != nil {
		return err
	}
	f, err := os.Create(wallet.name)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = io.Copy(f, strings.NewReader(string(keyData)))
	if err != nil {
		return err
	}
	return nil
}

func (wallet *Wallet) LoadWalletFile() error {
	f, err := os.Open(wallet.name)
	if err != nil {
		return err
	}
	defer f.Close()
	data, _ := ioutil.ReadAll(f)
	json.Unmarshal(data, &wallet.cipher)
	wallet.DecryptKeyPairs()
	return nil
}

func (wallet *Wallet) EncryptKeyPairs() error {
	keyData, _ := json.Marshal(wallet.keyPairs)
	block, err := aes.NewCipher(wallet.password[0:32])
	if err != nil {
		return err
	}
	wallet.cipher = make([]byte, aes.BlockSize+len(keyData))
	iv := wallet.cipher[:aes.BlockSize]
	if _, err = io.ReadFull(rand.Reader, iv); err != nil {
		return err
	}
	encryptStream := cipher.NewCTR(block, iv)
	encryptStream.XORKeyStream(wallet.cipher[aes.BlockSize:], keyData)
	return nil
}
