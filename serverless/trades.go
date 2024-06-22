package resourses

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	Logus "backend/logging"
)

type TRADE_MODEL struct {
	ID            uint     `db:"id" json:"id"`
	UID           string   `db:"uid" json:"uid"`
	KID           *uint    `db:"kid" json:"kid"`
	Exchange      uint8    `db:"exchange" json:"exchange"`
	Symbol        string   `db:"symbol" json:"symbol"`
	Tags          *string  `db:"tags" json:"tags"`
	Rating        *uint8   `db:"rating" json:"rating"`
	EntryTime     uint     `db:"entry_time" json:"entry_time"`
	ExitTime      uint     `db:"exit_time" json:"exit_time"`
	Duration      uint     `db:"duration" json:"duration"`
	Deposit       uint16   `db:"deposit" json:"deposit"`
	Side          string   `db:"side" json:"side"`
	Procent       float32  `db:"procent" json:"procent"`
	Funding       *float32 `db:"funding" json:"funding"`
	Income        float32  `db:"income" json:"income"`
	Profit        float32  `db:"profit" json:"profit"`
	Turnover      float32  `db:"turnover" json:"turnover"`
	MaxVolume     float32  `db:"max_volume" json:"max_volume"`
	Volume        float32  `db:"volume" json:"volume"`
	Commission    float32  `db:"commission" json:"commission"`
	AvgEntryPrice float32  `db:"avg_entry_price" json:"avg_entry_price"`
	AvgExitPrice  float32  `db:"avg_exit_price" json:"avg_exit_price"`
	Deals         string   `db:"deals" json:"deals"`
}

func Trades(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	uid := r.Header.Get("X-TRADIFY-UID")
	if uid == "" {
		json.NewEncoder(w).Encode("Нельзя :3")
		return
	}

	section := r.URL.Query().Get("section")

	switch section {
	case "overview":
		currentTime := time.Now().Unix() * 1000
		startTime := currentTime - 864000000

		rows, err := db.Query(`
				SELECT id, uid, kid, exchange, symbol, tags, rating, entry_time, exit_time, duration, deposit, side, procent, funding, income, profit, turnover, max_volume, volume, commission, avg_entry_price, avg_exit_price, deals
				FROM trades
				WHERE uid = ? AND entry_time >= ? AND entry_time <= ?
				ORDER BY entry_time DESC`, uid, startTime, currentTime)
		if err != nil {
			json.NewEncoder(w).Encode(nil)
			return
		}

		defer rows.Close()

		trades := []TRADE_MODEL{}

		for rows.Next() {
			var trade TRADE_MODEL
			if err := rows.Scan(&trade.ID, &trade.UID, &trade.KID, &trade.Exchange, &trade.Symbol, &trade.Tags, &trade.Rating, &trade.EntryTime, &trade.ExitTime, &trade.Duration, &trade.Deposit, &trade.Side, &trade.Procent, &trade.Funding, &trade.Income, &trade.Profit, &trade.Turnover, &trade.MaxVolume, &trade.Volume, &trade.Commission, &trade.AvgEntryPrice, &trade.AvgExitPrice, &trade.Deals); err != nil {
				json.NewEncoder(w).Encode(nil)
				return
			}
			trades = append(trades, trade)
		}

		if err := rows.Err(); err != nil {
			json.NewEncoder(w).Encode(nil)
			return
		}

		if err := json.NewEncoder(w).Encode(trades); err != nil {
			json.NewEncoder(w).Encode(nil)
		}

	case "table":
		pageStr := r.URL.Query().Get("page")
		pageSizeStr := r.URL.Query().Get("pageSize")

		page, _ := strconv.Atoi(pageStr)
		pageSize, _ := strconv.Atoi(pageSizeStr)

		offset := (page - 1) * pageSize

		var totalRecords int
		_ = db.QueryRow(`SELECT COUNT(*) FROM trades WHERE uid = ?`, uid).Scan(&totalRecords)

		rows, err := db.Query(`
		SELECT id, uid, kid, exchange, symbol, tags, rating, entry_time, exit_time, duration, deposit, side, procent, funding, income, profit, turnover, max_volume, volume, commission, avg_entry_price, avg_exit_price, deals
		FROM trades
		WHERE uid = ?
		ORDER BY entry_time DESC
		LIMIT ? OFFSET ?`, uid, pageSize, offset)
		if err != nil {
			json.NewEncoder(w).Encode(nil)
			return
		}

		defer rows.Close()

		trades := []TRADE_MODEL{}

		for rows.Next() {
			var trade TRADE_MODEL
			if err := rows.Scan(&trade.ID, &trade.UID, &trade.KID, &trade.Exchange, &trade.Symbol, &trade.Tags, &trade.Rating, &trade.EntryTime, &trade.ExitTime, &trade.Duration, &trade.Deposit, &trade.Side, &trade.Procent, &trade.Funding, &trade.Income, &trade.Profit, &trade.Turnover, &trade.MaxVolume, &trade.Volume, &trade.Commission, &trade.AvgEntryPrice, &trade.AvgExitPrice, &trade.Deals); err != nil {
				json.NewEncoder(w).Encode(nil)
				return
			}
			trades = append(trades, trade)
		}

		if err := rows.Err(); err != nil {
			json.NewEncoder(w).Encode(nil)
			return
		}

		response := struct {
			TotalRecords int           `json:"total"`
			TotalPages   int           `json:"last_page"`
			Page         int           `json:"page"`
			PageSize     int           `json:"pageSize"`
			Trades       []TRADE_MODEL `json:"data"`
		}{
			TotalRecords: totalRecords,
			TotalPages:   (totalRecords + pageSize - 1) / pageSize,
			Page:         page,
			PageSize:     pageSize,
			Trades:       trades,
		}

		if err := json.NewEncoder(w).Encode(response); err != nil {
			json.NewEncoder(w).Encode(nil)
		}

	case "times":
		startTimeStr := r.URL.Query().Get("startTime")
		endTimeStr := r.URL.Query().Get("endTime")

		startTime, _ := strconv.ParseInt(startTimeStr, 10, 64)
		endTime, _ := strconv.ParseInt(endTimeStr, 10, 64)

		rows, err := db.Query(`
				SELECT id, uid, kid, exchange, symbol, tags, rating, entry_time, exit_time, duration, deposit, side, procent, funding, income, profit, turnover, max_volume, volume, commission, avg_entry_price, avg_exit_price, deals
				FROM trades
				WHERE uid = ? AND entry_time >= ? AND entry_time <= ?
				ORDER BY entry_time ASC`, uid, startTime, endTime)

		if err != nil {
			json.NewEncoder(w).Encode(nil)
			return
		}

		defer rows.Close()

		trades := []TRADE_MODEL{}

		for rows.Next() {
			var trade TRADE_MODEL
			if err := rows.Scan(&trade.ID, &trade.UID, &trade.KID, &trade.Exchange, &trade.Symbol, &trade.Tags, &trade.Rating, &trade.EntryTime, &trade.ExitTime, &trade.Duration, &trade.Deposit, &trade.Side, &trade.Procent, &trade.Funding, &trade.Income, &trade.Profit, &trade.Turnover, &trade.MaxVolume, &trade.Volume, &trade.Commission, &trade.AvgEntryPrice, &trade.AvgExitPrice, &trade.Deals); err != nil {
				json.NewEncoder(w).Encode(nil)
				return
			}
			trades = append(trades, trade)
		}

		if err := rows.Err(); err != nil {
			json.NewEncoder(w).Encode(nil)
			return
		}

		if err := json.NewEncoder(w).Encode(trades); err != nil {
			json.NewEncoder(w).Encode(500)
			Logus.Logus.Error(err)
		}

	default:
		json.NewEncoder(w).Encode(nil)
	}
}
