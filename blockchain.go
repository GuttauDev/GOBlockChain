package block

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"
)

const MINING_DIFFICULTY = 3
const MINING_SENDER = "THE BLOCKCHAIN"
const MINING_REWARD = 1.0

// The Stuct of a Block is:
// Nonce a INT
// PreviousHash is a 32 Byte Sha256 SUM
// TimeStamp is a int64 Unix.Nano
// Transactions  is a Array ([]) of Transactions Struct Type

type Block struct {
	timestamp    int64
	nonce        int
	previousHash [32]byte
	transactions []*Transaction
}

// The Struct of Transaction is:
//sender, recipient are string for the walletAddress
// and value is a .1f Float of 32 bytes

type Transaction struct {
	senderBlockchainAddress    string
	recipientBlockchainAddress string
	value                      float32
}

// The Struct of a Blockchain is 2 arrays:
// transactionPool is a Array of Transaction Struct Type
// chain os a Array of Block Struct Type

type Blockchain struct {
	transactionPool   []*Transaction
	chain             []*Block
	blockchainAddress string
}

// the function below creates a newblock

func NewBlock(nonce int, previousHash [32]byte, transactions []*Transaction) *Block {
	b := new(Block)
	b.timestamp = time.Now().UnixNano()
	b.nonce = nonce
	b.previousHash = previousHash
	b.transactions = transactions
	return b
}

// The Function Bellow creates a NewBlockchain returning the Blochain Type

func NewBlockchain(blockchainAddress string) *Blockchain {
	b := &Block{}
	bc := new(Blockchain)
	bc.blockchainAddress = blockchainAddress
	bc.CreateBlock(0, b.Hash())
	return bc

}

// Custom Print Method for Blocks
func (inputB *Block) Print() {

	fmt.Printf("Nonce:              %d\n", inputB.nonce)
	fmt.Printf("previousHash:       %x\n", inputB.previousHash)
	fmt.Printf("timeStamp:          %d\n", inputB.timestamp)
	for _, t := range inputB.transactions {
		t.Print()
	}
}

// This is the Blockchain Type Print Method
func (bc *Blockchain) Print() {
	for i, block := range bc.chain {
		fmt.Printf("%s Chain %d %s\n", strings.Repeat("=", 25), i, strings.Repeat("=", 25))
		block.Print()
	}
	fmt.Printf("%s\n", strings.Repeat("*", 60))
}

// This is the Transaction Type Print Method
func (t *Transaction) Print() {
	fmt.Printf("%s\n", strings.Repeat("-", 40))
	fmt.Printf(" sender_blochain_address         %s\n", t.senderBlockchainAddress)
	fmt.Printf(" recipient_blochain_address      %s\n", t.recipientBlockchainAddress)
	fmt.Printf(" value                           %.1f\n", t.value)
}

// Create a sha256 Hashing Fuction for previousHash variables
// the hash is 32 bytes long and we json.Marshal this
func (b *Block) Hash() [32]byte {
	m, _ := json.Marshal(b) // variave m recebe os bytes e o blank serve para o erro]
	return sha256.Sum256([]byte(m))

}

// Marshaling is used to READ and DECODE Strcuts or Arrays to Byte
// bellow we are modifing the MarshalJson Function because of the LOWERCASE on variables
func (b *Block) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Timestamp    int64          `json:"timestamp"`
		Nonce        int            `json:"nonce"`
		PreviousHash [32]byte       `json:"previous_hash"`
		Transaction  []*Transaction `json:"transactions"`
	}{
		Timestamp:    b.timestamp,
		Nonce:        b.nonce,
		PreviousHash: b.previousHash,
		Transaction:  b.transactions,
	})
}

// the json: part is to how to show the value after mashaling

func (t *Transaction) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Sender    string  `json:"sender_blockhain_address"`
		Recipient string  `json:"recipient_blockchain_address"`
		Value     float32 `json:"value"`
	}{
		Sender:    t.senderBlockchainAddress,
		Recipient: t.recipientBlockchainAddress,
		Value:     t.value,
	})
}

// Functions of the Blockchain

// 1. Creates a Block
func (bc *Blockchain) CreateBlock(nonce int, previousHash [32]byte) *Block {
	b := NewBlock(nonce, previousHash, bc.transactionPool)
	bc.chain = append(bc.chain, b)
	bc.transactionPool = []*Transaction{}
	return b
}

// 2. Get the Last Block of the Blockchain
func (bc *Blockchain) LastBlock() *Block {
	return bc.chain[len(bc.chain)-1]
}

// 3. New Transaction function
func NewTransaction(sender string, recipient string, value float32) *Transaction {
	return &Transaction{sender, recipient, value}
}

// 4. Add Transactions to the transaction Pool
func (bc *Blockchain) AddTransaction(sender string, recipient string, value float32) {
	t := NewTransaction(sender, recipient, value)
	bc.transactionPool = append(bc.transactionPool, t)
}

// copy transactions on pool

func (bc *Blockchain) CopyTransactionPool() []*Transaction {
	transactions := make([]*Transaction, 0)
	for _, t := range bc.transactionPool {
		transactions = append(transactions,
			NewTransaction(t.senderBlockchainAddress, t.recipientBlockchainAddress, t.value))
	}
	return transactions
}

func (bc *Blockchain) ValidProof(nonce int, previousHash [32]byte, transactions []*Transaction, difficulty int) bool {
	zeros := strings.Repeat("0", difficulty)
	guessBlock := Block{0, nonce, previousHash, transactions}
	guessHashStr := fmt.Sprintf("%x", guessBlock.Hash())
	return guessHashStr[:difficulty] == zeros
}

func (bc *Blockchain) ProofofWork() int {
	transaction := bc.CopyTransactionPool()
	previousHash := bc.LastBlock().Hash()
	nonce := 0
	for !bc.ValidProof(nonce, previousHash, transaction, MINING_DIFFICULTY) {
		nonce += 1
	}
	return nonce
}

func (bc *Blockchain) Mining() bool {
	bc.AddTransaction(MINING_SENDER, bc.blockchainAddress, MINING_REWARD)
	nonce := bc.ProofofWork()
	previousHash := bc.LastBlock().Hash()
	bc.CreateBlock(nonce, previousHash)
	log.Println("action=mining, status=success")
	return true
}

// This method calculates the Total Amount
func (bc *Blockchain) CalculateTotalAmount(blockchainAddress string) float32 {

	var totalAmount float32 = 0.0

	for _, b := range bc.chain {

		for _, t := range b.transactions {

			value := t.value

			if blockchainAddress == t.recipientBlockchainAddress {
				totalAmount += value
			}

			if blockchainAddress == t.senderBlockchainAddress {
				totalAmount -= value
			}
		}

	}

	return totalAmount
}

func init() {
	log.SetPrefix("Blockchain: ")
}

func main() {
	myBlockchainAddress := "my_blockchain_address"
	// 1. Creating a New Blockchain register the Blocks and the Transactions
	blockchain := NewBlockchain(myBlockchainAddress) // the blockchain is associated to the blockchain var
	blockchain.Print()

	blockchain.AddTransaction("A", "B", 5.0)
	blockchain.Mining() // Adding Transaction to the blockchain since it stars with one block
	blockchain.Print()  // printing the blockchain, it as 3 blocks (block 0, 1 and 2) and 3 transactions

	blockchain.AddTransaction("B", "F", 3.5)
	blockchain.Print()

	fmt.Printf("A has %.1f\n", blockchain.CalculateTotalAmount("A"))
	fmt.Printf("B has %.1f\n", blockchain.CalculateTotalAmount("B"))
	fmt.Printf("F has %.1f\n", blockchain.CalculateTotalAmount("F"))
	fmt.Printf("My Wallet has %.1f\n", blockchain.CalculateTotalAmount("my_blockchain_address"))
}
