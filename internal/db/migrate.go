package db

import (
	"database/sql"
	_ "embed"
	"log"
)

//go:embed schema.sql
var schema string

// Migrate executa o schema SQL no banco. Usa IF NOT EXISTS em todas as
// operações, então é seguro rodar múltiplas vezes.
func Migrate(db *sql.DB) error {
	log.Println("→ Rodando migrations...")
	if _, err := db.Exec(schema); err != nil {
		return err
	}
	log.Println("✓ Migrations concluídas")
	return nil
}
