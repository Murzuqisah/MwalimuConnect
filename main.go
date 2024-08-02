package main

import (
	"crypto/rand"
	"database/sql"
	"encoding/base64"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"mwalimu/blockchain"
	"mwalimu/ethaccounts"
	"mwalimu/functions"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

type Config struct {
	EthLink         string
	EthPrivateKey   string
	ContractAddress string
}

func NewConfig() *Config {
	return &Config{
		EthLink:       "https://mainnet.infura.io/v3/436051433bd54480aca5957b7d0c77bb",
		EthPrivateKey: "436051433bd54480aca5957b7d0c77bb",
		ContractAddr:  "0x...",
	}
}

func main() {
	// Initialize configuration
	cfg := NewConfig()

	blockchain, err := initBlockchain()
	if err != nil {
		log.Fatal(err)
	}

	client, err := connectToEthereumClient(cfg.EthLink)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	contract, err := loadContract(client, cfg.ContractAddr, cfg.EthPrivateKey)
	if err != nil {
		log.Fatal(err)
	}

	// Serve static files (CSS, JS)
	http.HandleFunc("/scripts/", fileServer)
	http.HandleFunc("/static/", fileServer)

	// Serve dynamic HTML templates
	http.HandleFunc("/", servePage)

	// // Handle routes
	http.HandleFunc("/forgot-password", functions.ServeForgotPasswordForm)
	http.HandleFunc("/handle-forgot-password", functions.HandleForgotPassword)
	http.HandleFunc("/reset-password", functions.ServeResetPasswordForm)
	http.HandleFunc("/handle-reset-password", functions.HandleResetPassword)

	// serve signup
	http.HandleFunc("/signup", functions.SignupHandler)

	// Start the server
	log.Printf("Server started at http://localhost:8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}

func initBlockchain() (*functions.Blockchain, error) {
	// Initialize blockchain with genesis block
	blockchain := functions.Blockchain{}
	genesisBlock := functions.Block{
		Timestamp:    time.Now().Unix(),
		Data:         []byte("Genesis Block"),
		PreviousHash: []byte{},
	}
	genesisBlock.Hash = functions.CalculateHash(genesisBlock)
	blockchain.Blocks = append(blockchain.Blocks, genesisBlock)
	return &blockchain, nil
}

func connectToEthereumClient(link string) (ethclient.Client, error) {
	// Connect to the Ethereum client
	client, err := ethclient.Dial(link)
	return *client, err
}

func loadContract(client *ethclient.Client, contractAddr string, privateKey string) (*blockchain.Contract, error) {
	contractAddress := common.HexToAddress(contractAddr)
	privateKeyECDSA, err := crypto.HexToECDSA(privateKey)
	if err != nil {
		return nil, err
	}
	contract, err := blockchain.NewContract(contractAddress, privateKeyECDSA)
	return contract, err
}

// Serve static files (CSS, JS)
func fileServer(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./static/"+r.URL.Path[len("/scripts/"):])
}

// Serve HTML pages based on referral
func servePage(w http.ResponseWriter, r *http.Request) {
	page := r.URL.Path[1:]
	if page == "" || page == "/" {
		page = "index.html"
	}

	// Define template path
	tmplPath := "./templates/" + page

	// Parse and execute the template
	tmpl, err := template.ParseFiles(tmplPath)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		log.Printf("Template parsing error: %v", err)
		return
	}
	err = tmpl.Execute(w, nil)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		log.Printf("Template execution error: %v", err)
		return
	}
}

// The handler function for signup route
func handler(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case "/":
		functions.SignupHandler(w, r)
	default:
		code := http.StatusNotFound
		codefin := strconv.Itoa(code)
		fmt.Printf(codefin, "Page not found")
		os.Exit(1)
	}
}

func SignupHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.ServeFile(w, r, "templates/signup.html")
		return
	}

	// name := r.FormValue("name")
	// email := r.FormValue("email")
	// password := r.FormValue("password")
	// role := r.FormValue("role")

	// Create a new Ethereum account for the user
	account, err := blockchain.createEthereumAccount() // <--- create this function
	if err != nil {
		http.Error(w, "Failed to create Ethereum account", http.StatusInternalServerError)
		return
	}

	// Store user credentials
	err = functions.storeUserCredentials(r, account) // <--- create this function
	if err != nil {
		http.Error(w, "Failed to store user credentials", http.StatusInternalServerError)
		return
	}

	session, err := createSession(account)
	if err != nil {
		http.Error(w, "Failed to create session", http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:  "session",
		Value: session,
	})

	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

func createEthereumAccount() (string, error) {
	// Implement Ethereum account creation logic here
	// For example, using the go-ethereum library:
	account, err := ethaccounts.NewAccount(rand.Reader)
	return account.Address.Hex(), err
}

func storeUserCredentials(r *http.Request, account string) error {
	// 	// Implement user credential storage logic here
	// 	// For example, using a database:
	db, err := sql.Open("sqlite3", "./users.db")
	if err != nil {
		return err
	}
	defer db.Close()
	_, err = db.Exec("INSERT INTO users (address, password) VALUES (?, ?)", account, r.FormValue("password"))
	return err
}

func createSession(account string) (string, error) {
	// 	// Implement session creation logic here
	// 	// For example, using a random session ID:
	session := make([]byte, 32)
	_, err := rand.Read(session)
	return base64.StdEncoding.EncodeToString(session), err
}

// func session() {
// 	userStorage, err := blockchain.NewUserStorageContract(contractAddress, client)
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	// Create a new instance of the contract
// 	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, contract)
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	tx, err := userStorage.CreateUser(auth, "John Doe", "john@example.com", "Learner")
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	// Create an instance of SessionChain
// 	sessionChain, err := blockchain.NewSessionChain(client, contractAddr, EthPrivateKey) // <--- create function
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	// Retrieve session ID
// 	sessionID, err := sessionChain.GetSessionID()
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	fmt.Println("Session ID:", sessionID)

// 	// Store a session ID
// 	tx, err = sessionChain.StoreSessionID("my_session_id")
// 	if err != nil {
// 		log.Fatalf("Failed to store session ID: %v", err)
// 	}
// 	fmt.Printf("Transaction sent: %s\n", tx.Hash().Hex())

// 	sessionID, err = sessionChain.GetSessionID()
// 	if err != nil {
// 		log.Fatalf("Failed to retrieve session ID: %v", err)
// 	}

// 	log.Printf("Retrieved session ID: %s\n", sessionID)

// Connect to Ethereum client and get user storage contract
// client, err := ethclient.Dial(ethLink)
// if err != nil {
// 	http.Error(w, "Failed to connect to Ethereum client", http.StatusInternalServerError)
// 	return
// }

// contractAddress := common.HexToAddress(ethPrivateKey)
// userStorage, err := blockchain.NewUserStorage(contractAddress, client)
// if err != nil {
// 	http.Error(w, "Failed to initialize user storage", http.StatusInternalServerError)
// 	return
// }

// // Create transaction options
// auth, err := bind.NewKeyedTransactorWithChainID(account.PrivateKey, err)
// if err != nil {
// 	log.Fatal(err)
// }

// // Create user on the blockchain
// err = userStorage.CreateUser(auth, name, email, role)
// if err != nil {
// 	http.Error(w, "Failed to create user on blockchain", http.StatusInternalServerError)
// 	return
// }

// }
