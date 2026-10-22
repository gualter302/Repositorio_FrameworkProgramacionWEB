package config

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

// CargarEnv lee un archivo .env (si existe) y coloca cada variable en el entorno
// solo cuando todavía no está definida. Evita depender de una librería externa.
func CargarEnv() {
	archivo, err := os.Open(".env")
	if err != nil {
		return // Sin archivo .env usamos las variables del entorno o los valores por defecto.
	}
	defer archivo.Close()

	lector := bufio.NewScanner(archivo)
	for lector.Scan() {
		linea := strings.TrimSpace(lector.Text())
		if linea == "" || strings.HasPrefix(linea, "#") {
			continue
		}

		clave, valor, ok := strings.Cut(linea, "=")
		if !ok {
			continue
		}
		clave = strings.TrimSpace(clave)
		valor = strings.Trim(strings.TrimSpace(valor), `"'`)

		if _, existe := os.LookupEnv(clave); !existe {
			_ = os.Setenv(clave, valor)
		}
	}
}

// construirURL arma la cadena de conexión a PostgreSQL. DATABASE_URL tiene prioridad;
// en caso contrario se construye a partir de las variables individuales con valores por defecto.
func construirURL() string {
	if url := os.Getenv("DATABASE_URL"); url != "" {
		return url
	}

	host := valorODefecto("DB_HOST", "localhost")
	puerto := valorODefecto("DB_PORT", "5432")
	usuario := valorODefecto("DB_USER", "multicatalogo")
	password := valorODefecto("DB_PASSWORD", "multicatalogo")
	base := valorODefecto("DB_NAME", "multicatalogo")

	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		usuario, password, host, puerto, base)
}

func valorODefecto(clave, porDefecto string) string {
	if valor := os.Getenv(clave); valor != "" {
		return valor
	}
	return porDefecto
}

// ConectarBaseDeDatos crea el pool de conexiones y verifica que PostgreSQL responda.
// El pool se devuelve para que main pueda cerrarlo al finalizar la aplicación.
func ConectarBaseDeDatos(ctx context.Context) (*pgxpool.Pool, error) {
	cfg, err := pgxpool.ParseConfig(construirURL())
	if err != nil {
		return nil, fmt.Errorf("configuración de conexión inválida: %w", err)
	}
	cfg.MaxConns = 10
	cfg.MaxConnLifetime = time.Hour

	pool, err := pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		return nil, fmt.Errorf("no se pudo crear el pool de conexiones: %w", err)
	}

	ctxPing, cancelar := context.WithTimeout(ctx, 5*time.Second)
	defer cancelar()
	if err := pool.Ping(ctxPing); err != nil {
		pool.Close()
		return nil, fmt.Errorf("PostgreSQL no respondió: %w", err)
	}

	return pool, nil
}
