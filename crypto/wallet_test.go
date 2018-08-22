package crypto

import (
	"testing"
)

func TestWallet(t *testing.T) {
	wallet := NewWallet()
	wallet.SetPassword("pass")
	if len(wallet.GetKeyPairs()) > 0 {
		t.Fatal("0 key")
	}
	wallet.NewKeyPair()
	if len(wallet.GetKeyPairs()) != 1 {
		t.Fatal("1 key")
	}
	priv1, _ := NewRandomPrivateKey()
	pub := priv1.PublicKey()
	wif := priv1.String()
	wallet.ImportPrivateKey(wif)
	if len(wallet.GetKeyPairs()) != 2 {
		t.Fatal("2 keys")
	}
	priv2, _ := wallet.GetPrivateKey(*pub)
	if priv2.String() != wif {
		t.Fatal("private key should be the same")
	}
	wallet2 := NewWallet()
	wallet2.SetPassword("pass")
	wallet2.LoadWalletFile()
	if len(wallet2.GetKeyPairs()) != 2 {
		t.Fatal("2 keys")
	}
	priv3, _ := wallet2.GetPrivateKey(*pub)
	if priv3.String() != wif {
		t.Fatal("private key should be the same")
	}
}
