package database

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

// ConfigurarClaveAdministrador sustituye la clave de demostración cuando
// ADMIN_PASSWORD está definido en el alojamiento público.
func ConfigurarClaveAdministrador(ctx context.Context, pool *pgxpool.Pool, password string) error {
	if password == "" {
		return nil
	}
	if len(password) < 12 {
		return fmt.Errorf("ADMIN_PASSWORD debe tener al menos 12 caracteres")
	}

	var hashActual string
	err := pool.QueryRow(ctx, "SELECT password FROM usuarios WHERE email = $1", "admin@upse.edu.ec").Scan(&hashActual)
	if err != nil {
		if err == pgx.ErrNoRows {
			return fmt.Errorf("no se encontró la cuenta administradora")
		}
		return fmt.Errorf("no se pudo consultar la cuenta administradora: %w", err)
	}
	if bcrypt.CompareHashAndPassword([]byte(hashActual), []byte(password)) == nil {
		return nil
	}

	hashNuevo, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("no se pudo proteger la clave administradora: %w", err)
	}
	if _, err := pool.Exec(ctx, "UPDATE usuarios SET password = $1 WHERE email = $2", string(hashNuevo), "admin@upse.edu.ec"); err != nil {
		return fmt.Errorf("no se pudo actualizar la clave administradora: %w", err)
	}
	return nil
}
