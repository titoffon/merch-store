package db

import (
	"context"
	"fmt"
)

type PurchaseCount struct {
    MerchItem string
    Quantity  int64
}

type ReceivedTransaction struct {
    FromUser string
    Amount   int64
}

type SentTransaction struct {
    ToUser string
    Amount int64
}

func (r *DB) GetUserPurchases(ctx context.Context, username string) ([]PurchaseCount, error) {
	q := `
			SELECT merch_item, COUNT(*) as quantity
			FROM purchases
			WHERE username = $1
			GROUP BY merch_item
		`
	rows, err := r.DBPool.Query(ctx, q, username)
    if err != nil {
        return nil, fmt.Errorf("failed to query purchases: %w", err)
    }
    defer rows.Close()

    var results []PurchaseCount
    for rows.Next() {
        var pc PurchaseCount
        if err := rows.Scan(&pc.MerchItem, &pc.Quantity); err != nil {
            return nil, fmt.Errorf("failed to scan user purchases: %w", err)
        }
        results = append(results, pc)
    }
    if err := rows.Err(); err != nil {
        return nil, err
    }
    return results, nil
}

func (r *DB) GetTransactionsReceived(ctx context.Context, username string) ([]ReceivedTransaction, error) {
    q := `
        SELECT sender, amount
        FROM transaction_log
        WHERE recipient = $1
        ORDER BY created_at DESC
    `
    rows, err := r.DBPool.Query(ctx, q, username)
    if err != nil {
        return nil, fmt.Errorf("failed to query received transactions: %w", err)
    }
    defer rows.Close()

    var results []ReceivedTransaction
    for rows.Next() {
        var rt ReceivedTransaction
        if err := rows.Scan(&rt.FromUser, &rt.Amount); err != nil {
            return nil, fmt.Errorf("failed to scan received transaction: %w", err)
        }
        results = append(results, rt)
    }
    if err := rows.Err(); err != nil {
        return nil, err
    }
    return results, nil
}

func (r *DB) GetTransactionsSent(ctx context.Context, username string) ([]SentTransaction, error) {
    q := `
        SELECT recipient, amount
        FROM transaction_log
        WHERE sender = $1
        ORDER BY created_at DESC
    `
    rows, err := r.DBPool.Query(ctx, q, username)
    if err != nil {
        return nil, fmt.Errorf("failed to query sent transactions: %w", err)
    }
    defer rows.Close()

    var results []SentTransaction
    for rows.Next() {
        var st SentTransaction
        if err := rows.Scan(&st.ToUser, &st.Amount); err != nil {
            return nil, fmt.Errorf("failed to scan sent transaction: %w", err)
        }
        results = append(results, st)
    }
    if err := rows.Err(); err != nil {
        return nil, err
    }
    return results, nil
}