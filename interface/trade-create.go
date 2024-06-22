package api

import (
	Logus "backend/logging"
	"database/sql"
	"encoding/json"
	"net/http"
)

type TRADE_MODEL struct {
	ID            int64    `json:"id"`
	UID           string   `json:"uid"`
	KID           *int64   `json:"kid"`
	Exchange      int      `json:"exchange"`
	Symbol        string   `json:"symbol"`
	Tags          *string  `json:"tags"`
	Rating        *int     `json:"rating"`
	EntryTime     string   `json:"entry_time"`
	ExitTime      string   `json:"exit_time"`
	Duration      int      `json:"duration"`
	Deposit       int      `json:"deposit"`
	Side          string   `json:"side"`
	Procent       float64  `json:"procent"`
	Funding       *float64 `json:"funding"`
	Income        float64  `json:"income"`
	Profit        float64  `json:"profit"`
	Turnover      float64  `json:"turnover"`
	MaxVolume     float64  `json:"max_volume"`
	Volume        float64  `json:"volume"`
	Commission    float64  `json:"commission"`
	AvgEntryPrice float64  `json:"avg_entry_price"`
	AvgExitPrice  float64  `json:"avg_exit_price"`
	Deals         string   `json:"deals"`
}

func TradesCreate(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	uid := r.Header.Get("X-TRADIFY-UID")
	if uid == "" {
		json.NewEncoder(w).Encode("Нельзя :3")
		return
	}

	var trades []TRADE_MODEL

	if err := json.NewDecoder(r.Body).Decode(&trades); err != nil {
		json.NewEncoder(w).Encode(nil)
		Logus.Logus.Error(err)
		return
	}

	tx, err := db.Begin()
	if err != nil {
		json.NewEncoder(w).Encode(nil)
		Logus.Logus.Error(err)
		return
	}

	stmt, err := tx.Prepare(`INSERT INTO trades 
	(uid, kid, exchange, symbol, tags, rating, entry_time, exit_time, duration, deposit, side, procent, funding, income, profit, turnover, max_volume, volume, commission, avg_entry_price, avg_exit_price, deals) 
	VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`)
	if err != nil {
		json.NewEncoder(w).Encode(nil)
		Logus.Logus.Error(err)
		return
	}
	defer stmt.Close()

	count := 0
	for _, trade := range trades {
		_, err := stmt.Exec(
			uid, trade.KID, trade.Exchange, trade.Symbol, trade.Tags, trade.Rating, trade.EntryTime, trade.ExitTime, trade.Duration, trade.Deposit, trade.Side, trade.Procent, trade.Funding, trade.Income, trade.Profit, trade.Turnover, trade.MaxVolume, trade.Volume, trade.Commission, trade.AvgEntryPrice, trade.AvgExitPrice, trade.Deals,
		)
		if err != nil {
			json.NewEncoder(w).Encode(nil)
			Logus.Logus.Error(err)
			tx.Rollback()
			return
		}
		count++
	}

	if err := tx.Commit(); err != nil {
		json.NewEncoder(w).Encode(nil)
		Logus.Logus.Error(err)
		return
	}

	json.NewEncoder(w).Encode(count)
}
