package main

/*
  In this file we store the miner details. And also we have functions to generate keys, convert keys to hex, verify keys, convert hex to keys.
  1-		Miner is a struct to store the miner details
  2-		selfMiner is a struct to store the self miner details
  3-		MemPool is a struct to store the memory pool details
  4-		hexToPublicKey is a function to convert hex string to public key
  5-		newMiner is a function to create a new miner object
  6-		generateKeys is a function to generate public and private keys
  7-		keyToHex is a function to convert public and private key to hex string
  8-		keyVerification is a function to verify the public and private key
  9-		hexToPrivateKey is a function to convert hex string to private key
*/

import (
	"bufio"
	"context"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"math/big"
	"net"
	"sync"
)

/*
Miner is a struct to store the miner details
 1. conn: connection object
 2. write: writer object
 3. read: reader object
 4. pubKey: public key of the miner
*/
type Miner struct {
	conn   net.Conn
	write  *bufio.Writer
	read   *bufio.Reader
	pubKey *ecdsa.PublicKey
}

/*
selfMiner is a struct to store the self miner details
*/
type selfMiner struct {
	pubKey             *ecdsa.PublicKey
	prvKey             *ecdsa.PrivateKey
	pubKeyStr          string
	prvKeyStr          string
	nonce              int
	context            context.Context
	cancel             context.CancelFunc
	CurrentlyMineBlock Block
	interuptStatus     bool
	role               string
	connectionAlive    bool
	serviceMachineAddr string
	transactionFile    string
	connListen         net.Listener
	mu                 sync.Mutex
	blockLength        int
	powLenght          int
	LedgerFile         string
	readLedger         bool
}

/*
MemPool is a struct to store the memory pool details
transactions: list of transactions
*/
type MemPool struct {
	transactions []Transaction
}

/*
hexToPublicKey is a function to convert hex string to public key
hexStr: hex string
returns public key and error
*/
func hexToPublicKey(hexStr string) (*ecdsa.PublicKey, error) {

	pubBytes, err := hex.DecodeString(hexStr)
	if err != nil {
		return nil, fmt.Errorf("invalid hex string: %v", err)
	}

	if len(pubBytes) != 65 {
		return nil, fmt.Errorf("invalid public key length: expected 65 bytes, got %d", len(pubBytes))
	}

	if pubBytes[0] != 0x04 {
		return nil, fmt.Errorf("invalid public key format: expected uncompressed key")
	}

	x := new(big.Int).SetBytes(pubBytes[1:33])
	y := new(big.Int).SetBytes(pubBytes[33:65])

	pubKey := &ecdsa.PublicKey{
		Curve: elliptic.P256(),
		X:     x,
		Y:     y,
	}

	return pubKey, nil
}

/*
newMiner is a function to create a new miner object
 1. conn: connection object
 2. returns miner object
*/
func newMiner(conn net.Conn) *Miner {
	return &Miner{
		conn:  conn,
		write: bufio.NewWriter(conn),
		read:  bufio.NewReader(conn),
	}
}

/*
generateKeys is a function to generate public and private keys
returns public key, private key and error
*/
func generateKeys() (*ecdsa.PublicKey, *ecdsa.PrivateKey, error) {
	priv, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return nil, nil, fmt.Errorf("error generating key: %v", err)
	}
	pub := &priv.PublicKey

	return pub, priv, nil
}

/*
keyToHex is a function to convert public and private key to hex string
 1. pub: public key
 2. priv: private key
 3. returns public key and private key in hex format
*/
func keyToHex(pub *ecdsa.PublicKey, priv *ecdsa.PrivateKey) (string, string) {

	pubBytes := append([]byte{0x04}, pub.X.Bytes()...)
	pubBytes = append(pubBytes, pub.Y.Bytes()...)
	pubHex := hex.EncodeToString(pubBytes)
	privHex := hex.EncodeToString(priv.D.Bytes())

	return pubHex, privHex
}

/*
keyVerification is a function to verify the public and private key
 1. pub: public key
 2. priv: private key
 3. returns true if the keys are verified, false otherwise
*/
func keyVerification(pub *ecdsa.PublicKey, priv *ecdsa.PrivateKey) bool {
	if pub == nil || priv == nil {
		return false
	}

	if pub.Curve != priv.PublicKey.Curve {
		return false
	}

	if pub.X.Cmp(priv.PublicKey.X) == 0 && pub.Y.Cmp(priv.PublicKey.Y) == 0 {
		if pub.Curve.IsOnCurve(pub.X, pub.Y) {
			return true
		}
	}
	return false
}

/*
hexToPrivateKey is a function to convert hex string to private key
 1. hexStr: hex string
 2. returns private key and error
*/
func hexToPrivateKey(hexStr string) (*ecdsa.PrivateKey, error) {
	privateKeyBytes, err := hex.DecodeString(hexStr)
	if err != nil {
		return nil, fmt.Errorf("Invalid hex string: %v", err)
	}

	priv := new(ecdsa.PrivateKey)
	priv.PublicKey.Curve = elliptic.P256() // Define the curve
	priv.D = new(big.Int).SetBytes(privateKeyBytes)
	priv.PublicKey.X, priv.PublicKey.Y = priv.PublicKey.Curve.ScalarBaseMult(privateKeyBytes)

	return priv, nil
}
