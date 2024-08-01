package functions

import (
	"html/template"
	"mwalimu/blockchain"
	"net/http"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

func LearnerHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "session-name")
	user, ok := session.Values["user"].(blockchain.User)
	if !ok {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
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

	// Get updated user data from blockchain
	updatedUser, err := userStorage.GetUserByEmail(user.Email)
	if err != nil {
		http.Error(w, "Failed to retrieve user data", http.StatusInternalServerError)
		return
	}

	data := struct {
		Name              string
		ID                string
		CompletedSessions int64
		UpcomingSessions  int64
		TotalHours        int64
		AverageRating     float64
		WalletBalance     int64
	}{
		Name:              updatedUser.Name,
		ID:                "STU001", // You might want to generate this differently
		CompletedSessions: updatedUser.CompletedSessions.Int64(),
		UpcomingSessions:  updatedUser.UpcomingSessions.Int64(),
		TotalHours:        updatedUser.TotalHours.Int64(),
		AverageRating:     float64(updatedUser.AverageRating.Int64()) / 10, // Assuming rating is stored as an integer with one decimal place
		WalletBalance:     updatedUser.WalletBalance.Int64(),
	}

	tmpl, err := template.ParseFiles("templates/learner.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	tmpl.Execute(w, data)
}
