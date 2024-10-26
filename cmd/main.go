package main

import (
	"context"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"fmt"
	"github.com/WizardOfCodes442/go-multisig-wallet/wallet"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"log"
	"math/big"
	"os"
	"strings"
	"time"
)

//TIP See GoLand help at <a href="https://www.jetbrains.com/help/go/">jetbrains.com/help/go/</a>.
// Also, you can try interactive lessons for GoLand by selecting 'Help | Learn IDE Features' from the main menu.

// Constants fo deployed contract
const (
	ContractAddress = "0x8928e59eFC1F62a230331a566D59F8457CE80bb0" // Your deployed address
	InfuraURL       = "https://sepolia.infura.io/v3/YOUR_INFURA_PROJECT_ID"
	PrivateKey      = "f418c51dc455d3f1a6ddcb67e42892b6efbbf6171eef5ab1b415cdae61186694" // Replace with your private key
)

// LoadFile reads the content of a given
func LoadFile(filePath string) (string, error) {
	data, err := os.ReadFile(filePath) //
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// Load the ABI to interact with the deployed contract
func loadABI(abiPath string) abi.ABI {
	abiData, err := LoadFile(abiPath)
	if err != nil {
		log.Fatalf("Failed to load ABI: %v", err)
	}

	parsedABI, err := abi.JSON(strings.NewReader(abiData))
	if err != nil {
		log.Fatalf("Failed to parse ABI: %v", err)
	}
	return parsedABI

}

// Interact with Ethereum contract: submit a transaction
func submitTransaction(client *ethclient.Client, parsedABI abi.ABI, privateKey *ecdsa.PrivateKey, to string, value *big.Int) {
	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, big.NewInt(11155111)) // Sepolia TestNetwork
	if err != nil {
		log.Fatalf("Failed to create transactor: %v", err)
	}

	auth.GasLimit = uint64(3000000)
	printTransactionDetails(auth, to, value)

	contractAddress := common.HexToAddress(ContractAddress)
	instance := bind.NewBoundContract(contractAddress, parsedABI, client, client, client)

	//Call submitTransaction method on the deployed contract
	tx, err := instance.Transact(auth, "submitTransaction", common.HexToAddress(to), value)
	if err != nil {
		log.Fatalf("Failed to submit transaction:L %v", err)
	}
	fmt.Printf("Transaction submitted! TX %s\n", tx.Hash().Hex())
}

// Function to print the details of a transaction
func printTransactionDetails(tx *bind.TransactOpts, recipient string, value *big.Int) {
	fmt.Println("\n=== Transaction Details ===")
	fmt.Printf("From: %s\n", tx.From.Hex())
	fmt.Printf("To: %s\n", recipient)
	fmt.Printf("Value: %s Wei\n", value.String())
	fmt.Printf("Nonce: %d\n", tx.Nonce.Uint64())
	fmt.Printf("Gas Price: %s\n", tx.GasPrice.String())
	fmt.Printf("Gas Limit: %d\n", tx.GasLimit)
	//fmt.Printf("Raw Data: %x\n", tx.Data)
	fmt.Println("===========================\n")
}

// Function to print the details of a signature
func printSignatureDetails(sig wallet.Signature) {
	r, s, v := sig.R, sig.S, sig.V

	fmt.Println("\n=== Signature Details ===")
	fmt.Printf("R: %x\n", r)
	fmt.Printf("S: %x\n", s)
	fmt.Printf("V: %d\n", v)
	fmt.Println("==========================\n")
}

//Main function: Test both Go and Solidity muiltisig wallet

func main() {
	//==== Pat 1 Test GO-Based MultiSig Wallet

	//1. Generate keys for three owners
	privKeys := make([]*ecdsa.PrivateKey, 3)
	pubKeys := make([]*ecdsa.PublicKey, 3)
	for i := 0; i < 3; i++ {
		priv, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
		if err != nil {
			log.Fatalf("Failed to generate key: %v", err)
		}
		privKeys[i] = priv
		pubKeys[i] = &priv.PublicKey

		// Print Private and Public Keys
		fmt.Printf("\nOwner %d:\n", i+1)
		fmt.Printf("Private Key: %x\n", priv.D) // Private key in hex
		fmt.Printf("Public Key (X, Y): (%x, %x)\n", priv.PublicKey.X, priv.PublicKey.Y)

	}

	//2. Create a new Go-based multisig wallet with 2 out of 3 signatures required
	mswallet := wallet.NewWallet(pubKeys, 2)

	//3. Add a new transaction to the Go wallet

	tx := mswallet.AddTransaction("", 10.5)
	fmt.Println("\nTransaction added to Go-based multisig wallet:")
	fmt.Printf("Recipient: 0xRecipientAddress\nAmount: 10.5 ETH\n")
	//4. Sign the transaction using two private keys
	for i := 0; i < 2; i++ {
		err := mswallet.SignTransaction(tx, privKeys[i])
		if err != nil {
			log.Fatalf("failed to sign transaction: %v", err)
		}
		// Print the signature details for each key
		printSignatureDetails(tx.Signatures[i])

	}

	//5. verify if the transaction is fully signed and ready to execute
	if mswallet.VerifySignatures(tx) {
		fmt.Println("Go-based wallet: Transaction verified successfully !")
		fmt.Println("Transaction tx ")
	} else {
		fmt.Println("Go-based wallet: Transaction verification failed ")
	}

	//====== PART 2: INTERACT WITH DEPLOYED ETHEREUM CONTRACT  ====

	//1. Connect to sepolia testnet
	client, err := ethclient.Dial(InfuraURL)
	if err != nil {
		log.Fatalf("Failed to connect to Ethereum client: %v", err)
	}

	//2. Load private key to interact with the contract
	privateKey, err := crypto.HexToECDSA(PrivateKey)
	if err != nil {
		log.Fatalf("Failed to load private key: %v", err)
	}
	fmt.Println("\nLoaded Ethereum private key for contract interaction.")
	fmt.Printf("Private Key: %x\n", privateKey.D)

	//Load ABI to interact with contract
	parsedABI := loadABI("./eth-contracts/compiled/Multisigwallet.abi")

	//4. submit a new transaction to the ethereum contract
	fmt.Println("Submitting a transactions to the Ethereum contract ")
	submitTransaction(client, parsedABI, privateKey, "0xRecipentaddress", big.NewInt(10000))

	//5. verify the submitted transaction(optional)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	receipt, err := client.TransactionReceipt(ctx, common.HexToHash("0xf57dcf225496ea48972a92e5924dd3c7a7e6ebc03a313dae29aec585cd7e7aea"))
	if err != nil {
		log.Fatalf("Failed to get transaction receipt: %v", err)
	}

	if receipt.Status == 1 {
		fmt.Println("Ethereum contract: Transaction executed successfully!")
	} else {
		fmt.Println("Ethereum cxontract: Transaction execution failed .")
	}
}
