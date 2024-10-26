package wallet

import (
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/sha256"
	"encoding/binary"
	"github.com/ethereum/go-ethereum/crypto"
	"log"
	"math"
	"math/big"
	"sync"
)

// Signature represents the R, S, V components of an ECDSA signature
type Signature struct {
	R *big.Int
	S *big.Int
	V uint8
}

// MultiSigWallet represents a multi-signature wallet
type MultiSigWallet struct {
	PublicKeys   []*ecdsa.PublicKey // List of Public Keys of participants
	Threshold    int                // Minimum number of signatures required
	mutex        sync.Mutex
	Transactions []*Transaction // List of pending Transactions
}

// Transaction represents a wallet transaction
type Transaction struct {
	To         string      // Recipient address
	Amount     float64     // Amount to send
	Signatures []Signature // Collected signatures
}

// NewWallet initializes a new multi-sig wallet
func NewWallet(keys []*ecdsa.PublicKey, threshold int) *MultiSigWallet {
	return &MultiSigWallet{
		PublicKeys: keys,
		Threshold:  threshold,
	}
}

// AddTransaction creates a new pending transaction
func (w *MultiSigWallet) AddTransaction(to string, amount float64) *Transaction {
	tx := &Transaction{To: to, Amount: amount}
	w.mutex.Lock()
	w.Transactions = append(w.Transactions, tx)
	w.mutex.Unlock()
	return tx
}

// Hash returns a SHA-256 hash of the transaction
func (tx *Transaction) Hash() []byte {
	// Create a buffer to store transaction data
	data := make([]byte, 8) // 8 bytes to encode float64 amount
	binary.LittleEndian.PutUint64(data, math.Float64bits(tx.Amount))

	// Concatenate transaction details
	payload := append([]byte(tx.To), data...)

	// Calculate and return the SHA-256 hash
	hash := sha256.Sum256(payload)
	return hash[:]
}

// SignTransaction allows a user to sign the transaction with their private key
func (w *MultiSigWallet) SignTransaction(tx *Transaction, privKey *ecdsa.PrivateKey) error {
	hash := tx.Hash()                                   // Generate the hash of the transaction
	r, s, err := ecdsa.Sign(rand.Reader, privKey, hash) // Sign the hash
	if err != nil {
		return err
	}

	// Calculate the recovery ID (V value)
	v := uint8(0)
	if s.Cmp(new(big.Int).Div(crypto.S256().Params().N, big.NewInt(2))) > 0 {
		s.Sub(crypto.S256().Params().N, s)
		v = 1
	}

	// Store the signature as a Signature struct
	sig := Signature{R: r, S: s, V: v}
	tx.Signatures = append(tx.Signatures, sig)

	if len(tx.Signatures) >= w.Threshold {
		log.Println("Transaction fully signed and ready to be executed")
	}
	return nil
}

// VerifySignatures verifies if the collected signatures meet the threshold
func (w *MultiSigWallet) VerifySignatures(tx *Transaction) bool {
	if len(tx.Signatures) < w.Threshold {
		return false
	}
	// Add actual verification logic if needed
	return true
}
