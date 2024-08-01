package functions

import (
	"math/big"
	"net/http"
	"strconv"

	"mwalimu/blockchain"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

func TopUpWalletHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	session, _ := store.Get(r, "session-name")
	user, ok := session.Values["user"].(blockchain.User)
	if !ok {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	amountStr := r.FormValue("amount")
	amount, err := strconv.ParseInt(amountStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid amount", http.StatusBadRequest)
		return
	}

	// Connect to Ethereum client and get user storage contract
	client, err := ethclient.Dial("https://mainnet.infura.io/v3/YOUR-PROJECT-ID")
	if err != nil {
		http.Error(w, "Failed to connect to Ethereum client", http.StatusInternalServerError)
		return
	}

	contractAddress := common.HexToAddress("YOUR_CONTRACT_ADDRESS")
	userStorage, err := blockchain.NewUserStorage(contractAddress, client)
	if err != nil {
		http.Error(w, "Failed to initialize user storage", http.StatusInternalServerError)
		return
	}

	// Get the user's Ethereum account
	account, err := getEthereumAccount(user.Email)
	if err != nil {
		http.Error(w, "Failed to retrieve Ethereum account", http.StatusInternalServerError)
		return
	}

	// Create transaction options
	auth := bind.NewKeyedTransactor(account.PrivateKey)

	// Update wallet balance on the blockchain
	newBalance := new(big.Int).Add(user.WalletBalance, big.NewInt(amount))
	err = userStorage.UpdateWalletBalance(auth, newBalance)
	if err != nil {
		http.Error(w, "Failed to update wallet balance", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/learner", http.StatusSeeOther)
}
