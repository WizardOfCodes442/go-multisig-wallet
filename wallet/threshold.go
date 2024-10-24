package wallet

import (
	"github.com/hashicorp/vault/shamir"
)

// SplitKey splits a private key
// into n p arts with a threshold t.
func SplitKey(secret []byte, n, t int) ([][]byte, error) {
	parts, err := shamir.Split(secret, n, t)
	if err != nil {
		return nil, err
	}
	return parts, nil
}

// CombineKey combines parts of a private key to reconstruct it.
func CombineKey(parts [][]byte) ([]byte, error) {
	secret, err := shamir.Combine(parts)
	if err != nil {
		return nil, err
	}
	return secret, nil
}
