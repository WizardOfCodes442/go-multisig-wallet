package main

import (
	"crypto/ecdsa"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"io/ioutil"
	"log"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

// Load compiled contract ABI and Bytecode
func loadContractFiles(abiPath, binPath string) (abi.ABI, []byte, error) {
	// Read the ABI file
	abiFile, err := ioutil.ReadFile(abiPath)
	if err != nil {
		return abi.ABI{}, nil, fmt.Errorf("failed to read ABI: %v", err)
	}

	// Parse the ABI from file content using strings.NewReader
	parsedABI, err := abi.JSON(strings.NewReader(string(abiFile)))
	if err != nil {
		return abi.ABI{}, nil, fmt.Errorf("failed to parse ABI: %v", err)
	}

	// Read the bytecode file
	binFile, err := ioutil.ReadFile(binPath)
	if err != nil {
		return abi.ABI{}, nil, fmt.Errorf("failed to read ByteCode: %v", err)
	}

	return parsedABI, binFile, nil
}

// DeployContract deploys the MultiSig wallet contract
func DeployContract(client *ethclient.Client, privateKey *ecdsa.PrivateKey, abiPath, binPath string, owners []common.Address, requiredSignatures *big.Int) (string, error) {
	contractABI, contractBin, err := loadContractFiles(abiPath, binPath)
	if err != nil {
		return "", err
	}

	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, big.NewInt(11155111)) // Mainnet
	if err != nil {
		return "", fmt.Errorf("failed to create transactor: %v", err)
	}

	auth.Value = big.NewInt(0)      // No ether sent
	auth.GasLimit = uint64(3000000) // Gas limit

	address, tx, _, err := bind.DeployContract(auth, contractABI, contractBin, client, owners, requiredSignatures)
	if err != nil {
		return "", fmt.Errorf("failed to deploy contract: %v", err)
	}

	log.Printf("Contract deployed! TxHash: %s\n", tx.Hash().Hex())
	return address.Hex(), nil
}

func main() {
	// Connect to Ethereum client (e.g., Infura or local node)
	client, err := ethclient.Dial("https://sepolia.infura.io/v3/6f9fd646619f45379636ad9df58f3d70")
	if err != nil {
		log.Fatalf("Failed to connect to Ethereum Client: %v", err)
	}
	owners := []common.Address{
		common.HexToAddress("0x123..."),
		common.HexToAddress("0x456..."),
	}
	requiredSignatures := big.NewInt(2)

	// Load private key
	privateKey, err := crypto.HexToECDSA("f418c51dc455d3f1a6ddcb67e42892b6efbbf6171eef5ab1b415cdae61186694")
	if err != nil {
		log.Fatalf("Failed to load private key: %v", err)
	}

	// Deploy the contract
	contractAddress, err := DeployContract(
		client,
		privateKey,
		"./eth-contracts/compiled/MultisigWallet.abi",
		"./eth-contracts/compiled/MultisigWallet.bin",
		owners,
		requiredSignatures)
	if err != nil {
		log.Fatalf("Failed to deploy contract: %v", err)
	}

	log.Printf("Contract deployed at address: %s\n", contractAddress)
}
