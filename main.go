package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"mwalimu/blockchain"
	"mwalimu/functions"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

const (
	ethLink       = "https://mainnet.infura.io/v3/436051433bd54480aca5957b7d0c77bb"
	ethPrivateKey = "436051433bd54480aca5957b7d0c77bb"
)

func main() {
	// Initialize blockchain with genesis block
	blockchain := functions.Blockchain{}
	genesisBlock := functions.Block{
		Timestamp:    time.Now().Unix(),
		Data:         []byte("Genesis Block"),
		PreviousHash: []byte{},
	}
	genesisBlock.Hash = functions.CalculateHash(genesisBlock)
	blockchain.Blocks = append(blockchain.Blocks, genesisBlock)

	// Connect to the Ethereum client
	client, err := ethclient.Dial(ethLink)
	if err != nil {
		log.Fatalf("Failed to connect to the Ethereum client: %v", err)
	}
	defer client.Close()

	// Load the contract's address and private key
	contractAddress := common.HexToAddress("0x...")     // deployed contract address
	privateKey, err := crypto.HexToECDSA(ethPrivateKey) // API key
	if err != nil {
		log.Fatalf("Failed to parse private key: %v", err)
	}
	userStorage, err := blockchain.NewUserStorageContract(contractAddress, client)
	if err != nil {
		log.Fatal(err)
	}
	contract, err := blockchain.NewContract(contractAddress, privateKey)
	if err != nil {
		log.Fatalf("Failed to deploy contract: %v", err)
	}
	defer contract.Stop()
	// Create a new instance of the contract
	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, contract)
	if err != nil {
		log.Fatal(err)
	}

	tx, err := userStorage.CreateUser(auth, "John Doe", "john@example.com", "Learner")
	if err != nil {
		log.Fatal(err)
	}

	// Create an instance of SessionChain
	sessionChain, err := blockchain.NewSessionChain(client, contractAddress, privateKey) // <--- create function
	if err != nil {
		log.Fatal(err)
	}

	// Retrieve session ID
	sessionID, err := sessionChain.GetSessionID()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Session ID:", sessionID)

	// Store a session ID
	tx, err = sessionChain.StoreSessionID("my_session_id")
	if err != nil {
		log.Fatalf("Failed to store session ID: %v", err)
	}
	fmt.Printf("Transaction sent: %s\n", tx.Hash().Hex())

	sessionID, err = sessionChain.GetSessionID()
	if err != nil {
		log.Fatalf("Failed to retrieve session ID: %v", err)
	}

	log.Printf("Retrieved session ID: %s\n", sessionID)

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

	name := r.FormValue("name")
	email := r.FormValue("email")
	password := r.FormValue("password")
	role := r.FormValue("role")

	// Create a new Ethereum account for the user
	account, err := blockchain.createEthereumAccount() // <--- create this function
	if err != nil {
		http.Error(w, "Failed to create Ethereum account", http.StatusInternalServerError)
		return
	}

	// Store user credentials
	err = functions.storeUserCredentials(email, password, account.Address.Hex()) // <--- create this function
	if err != nil {
		http.Error(w, "Failed to store user credentials", http.StatusInternalServerError)
		return
	}

	// Connect to Ethereum client and get user storage contract
	client, err := ethclient.Dial(ethLink)
	if err != nil {
		http.Error(w, "Failed to connect to Ethereum client", http.StatusInternalServerError)
		return
	}

	contractAddress := common.HexToAddress(ethPrivateKey)
	userStorage, err := blockchain.NewUserStorage(contractAddress, client)
	if err != nil {
		http.Error(w, "Failed to initialize user storage", http.StatusInternalServerError)
		return
	}

	// Create transaction options
	auth, err := bind.NewKeyedTransactorWithChainID(account.PrivateKey, err)
	if err != nil {
		log.Fatal(err)
	}

	// Create user on the blockchain
	err = userStorage.CreateUser(auth, name, email, role)
	if err != nil {
		http.Error(w, "Failed to create user on blockchain", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/login", http.StatusSeeOther)
}
