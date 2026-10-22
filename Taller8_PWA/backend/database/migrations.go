package database

import (
	"context"
	"embed"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

//go:embed *.sql
var archivos embed.FS

var migraciones = []string{
	"init.sql",
	"004_registro_referidos.sql",
	"005_productos_activos.sql",
	"006_comisiones_liquidaciones.sql",
	"007_pwa.sql",
}

const bloqueoMigraciones int64 = 807302

// AplicarMigraciones prepara una base nueva y conserva las bases existentes.
// El bloqueo evita que dos instancias intenten cambiar el esquema al mismo tiempo.
func AplicarMigraciones(ctx context.Context, pool *pgxpool.Pool) error {
	conexion, err := pool.Acquire(ctx)
	if err != nil {
		return fmt.Errorf("no se pudo reservar la conexión para migraciones: %w", err)
	}
	defer conexion.Release()

	if _, err := conexion.Exec(ctx, "SELECT pg_advisory_lock($1)", bloqueoMigraciones); err != nil {
		return fmt.Errorf("no se pudo bloquear la ejecución de migraciones: %w", err)
	}
	defer func() {
		_, _ = conexion.Exec(context.Background(), "SELECT pg_advisory_unlock($1)", bloqueoMigraciones)
	}()

	if _, err := conexion.Exec(ctx, `CREATE TABLE IF NOT EXISTS schema_migrations (
		version TEXT PRIMARY KEY,
		aplicada_en TIMESTAMPTZ NOT NULL DEFAULT NOW()
	)`); err != nil {
		return fmt.Errorf("no se pudo preparar el historial de migraciones: %w", err)
	}

	for _, nombre := range migraciones {
		var aplicada bool
		if err := conexion.QueryRow(ctx, "SELECT EXISTS (SELECT 1 FROM schema_migrations WHERE version = $1)", nombre).Scan(&aplicada); err != nil {
			return fmt.Errorf("no se pudo comprobar la migración %s: %w", nombre, err)
		}
		if aplicada {
			continue
		}

		contenido, err := archivos.ReadFile(nombre)
		if err != nil {
			return fmt.Errorf("no se pudo leer la migración %s: %w", nombre, err)
		}

		transaccion, err := conexion.Begin(ctx)
		if err != nil {
			return fmt.Errorf("no se pudo iniciar la migración %s: %w", nombre, err)
		}
		if _, err = transaccion.Exec(ctx, string(contenido)); err == nil {
			_, err = transaccion.Exec(ctx, "INSERT INTO schema_migrations (version) VALUES ($1)", nombre)
		}
		if err != nil {
			_ = transaccion.Rollback(ctx)
			return fmt.Errorf("no se pudo aplicar la migración %s: %w", nombre, err)
		}
		if err := transaccion.Commit(ctx); err != nil {
			return fmt.Errorf("no se pudo confirmar la migración %s: %w", nombre, err)
		}
	}

	return nil
}
