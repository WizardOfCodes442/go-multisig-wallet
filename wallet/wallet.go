package wallet

import (
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/sha256"
	"encoding/binary"
	"log"
	"math"
	"sync"
)

// MultiSigWallet represents a multi-signaturewallet
type MultiSigWallet struct {
	PublicKeys   []*ecdsa.PublicKey //List of Public Keys of participants
	Threshold    int                // Minimum number of signatures required.
	mutex        sync.Mutex
	Transactions []*Transaction // List of pending Transactions
}

// Transaction represents a wallet transaction
type Transaction struct {
	To         string   //recipient
	Amount     float64  //Amount to send
	Signatures [][]byte // collected signatures
}

// NewWallet initializes a new multi-sig wallet
func NewWallet(keys []*ecdsa.PublicKey, threshold int) *MultiSigWallet {
	return &MultiSigWallet{
		PublicKeys: keys,
		Threshold:  threshold,
	}
}

// AddTransactions creates a new pending transaction
func (w *MultiSigWallet) AddTransaction(to string, amount float64) *Transaction {
	tx := &Transaction{To: to, Amount: amount}
	w.mutex.Lock()
	w.Transactions = append(w.Transactions, tx)
	w.mutex.Unlock()
	return tx
}

// Hash returns a SHA-256 hash of the transaction
func (tx *Transaction) Hash() []byte {
	//Create a buffer to store transaction data
	data := make([]byte, 8) // 8 bytes to encode float64 amount
	binary.LittleEndian.PutUint64(data, math.Float64bits(tx.Amount))

	//Concatenate transaction details
	payload := append([]byte(tx.To), data...)

	//calculate and return the SHA-256 hash.
	hash := sha256.Sum256(payload)
	return hash[:]
}

// SignTransaction allows a user to sign the transaction with their private key
func (w *MultiSigWallet) SignTransaction(tx *Transaction, privKey *ecdsa.PrivateKey) error {
	hash := tx.Hash()                                   //generate the hash of the transaction
	r, s, err := ecdsa.Sign(rand.Reader, privKey, hash) //Sign the hash
	if err != nil {
		return err
	}

	//Combine r and s values into a single byte slice.
	signature := append(r.Bytes(), s.Bytes()...)
	tx.Signatures = append(tx.Signatures, signature)

	if len(tx.Signatures) >= w.Threshold {
		log.Println("Transaction fully signed and ready to be executed ")
	}
	return nil
}

// VerifySignatures verifies that
// the collected signature meet the threshold
func (w *MultiSigWallet) VerifySignatures(tx *Transaction) bool {
	if len(tx.Signatures) < w.Threshold {
		return false
	}
	return true // Add actual verification logic if needed.
}
