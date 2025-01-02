package db

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/pdupub/go-pdu/internal/core"

	_ "github.com/mattn/go-sqlite3"
)

func insertQuantum(db *sql.DB, sq *core.SignedQuantum) error {

	// 1) 插入 quantum
	_, err := db.Exec(`
        INSERT INTO quantum (signature, last, nonce, type, signer, timestamp)
        VALUES (?, ?, ?, ?, ?, ?)`,
		sq.Signature, sq.Last, sq.Nonce, sq.Type, sq.Signer, time.Now().Unix())
	if err != nil {
		return fmt.Errorf("insert quantum error: %w", err)
	}

	// 2) 插入 contents
	for _, c := range sq.Contents {
		_, err := db.Exec(`
            INSERT INTO content (quantum_signature, data, format)
            VALUES (?, ?, ?)`,
			sq.Signature, c.Data, c.Format)
		if err != nil {
			return fmt.Errorf("insert content error: %w", err)
		}
	}

	// 3) 对 references，每个 ref_text 先检查是否已存在；若不存在就插入
	for _, ref := range sq.References {
		var refID int64
		// 先查是否已存在
		err := db.QueryRow(`SELECT id FROM reference WHERE ref_text = ?`, ref).Scan(&refID)
		if err != nil {
			if err == sql.ErrNoRows {
				// 不存在则插入
				res, err := db.Exec(`INSERT INTO reference (ref_text) VALUES (?)`, ref)
				if err != nil {
					return fmt.Errorf("insert reference error: %w", err)
				}
				refID, _ = res.LastInsertId()
			} else {
				return fmt.Errorf("query reference error: %w", err)
			}
		}

		// 4) 插入 quantum_references
		_, err = db.Exec(`
            INSERT INTO quantum_reference (quantum_signature, reference_id)
            VALUES (?, ?)`,
			sq.Signature, refID)
		if err != nil {
			return fmt.Errorf("insert quantum_reference error: %w", err)
		}
	}

	return nil
}

func queryQuantumsByReference(db *sql.DB, refText string) ([]core.SignedQuantum, error) {
	// 多表join： quantum + quantum_references + references
	rows, err := db.Query(`
        SELECT q.signature, q.last, q.nonce, q.type, q.signer, q.timestamp
        FROM quantum q
        JOIN quantum_reference qr ON q.signature = qr.quantum_signature
        JOIN reference r ON qr.reference_id = r.id
        WHERE r.ref_text = ?
    `, refText)
	if err != nil {
		return nil, fmt.Errorf("query error: %w", err)
	}
	defer rows.Close()

	var results []core.SignedQuantum
	for rows.Next() {
		var sq core.SignedQuantum
		var t int64
		err := rows.Scan(&sq.Signature, &sq.Last, &sq.Nonce, &sq.Type, &sq.Signer, &t)
		if err != nil {
			return nil, fmt.Errorf("scan error: %w", err)
		}
		// contents、references 暂时留空，需要另外查
		results = append(results, sq)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	// 如果还要一次性把 contents、references 补齐，也可以在这里对每个 signature 再做额外查询:
	// fetchContents(db, sq.Signature)
	// fetchReferences(db, sq.Signature)
	// 此处演示就不写这么详细了，概念上是一致的：通过 signature 去 contents, quantum_references, references 表取。

	return results, nil
}
