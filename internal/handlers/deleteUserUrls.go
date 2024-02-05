package handlers

import (
	"encoding/json"
	"github.com/grafchitaru/shortener/internal/auth"
	"github.com/grafchitaru/shortener/internal/config"
	"net/http"
)

func DeleteUserUrls(ctx config.HandlerContext, res http.ResponseWriter, req *http.Request) {
	userID, err := auth.GetUserID(req, ctx.Config.SecretKey)
	if err != nil {
		http.Error(res, err.Error(), http.StatusUnauthorized)
		return
	}

	var deleteIDs []string
	if err := json.NewDecoder(req.Body).Decode(&deleteIDs); err != nil {
		http.Error(res, "Invalid JSON array", http.StatusBadRequest)
		return
	}

	message, err := ctx.Repos.DeleteUserURLs(userID, deleteIDs)
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}

	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(http.StatusAccepted)
	json.NewEncoder(res).Encode(map[string]string{"message": message})
}
