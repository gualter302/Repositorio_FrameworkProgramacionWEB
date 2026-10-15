package config

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

// ConectarRedis crea el cliente y comprueba que Redis esté disponible antes de iniciar la API.
func ConectarRedis(ctx context.Context) (*redis.Client, error) {
	opciones, err := opcionesRedis()
	if err != nil {
		return nil, err
	}

	cliente := redis.NewClient(opciones)
	ctxPing, cancelar := context.WithTimeout(ctx, 5*time.Second)
	defer cancelar()

	if err := cliente.Ping(ctxPing).Err(); err != nil {
		_ = cliente.Close()
		return nil, fmt.Errorf("Redis no respondió: %w", err)
	}

	return cliente, nil
}

func opcionesRedis() (*redis.Options, error) {
	if url := os.Getenv("REDIS_URL"); url != "" {
		opciones, err := redis.ParseURL(url)
		if err != nil {
			return nil, fmt.Errorf("configuración de Redis inválida: %w", err)
		}
		return opciones, nil
	}

	base, err := strconv.Atoi(valorODefecto("REDIS_DB", "0"))
	if err != nil || base < 0 {
		return nil, fmt.Errorf("REDIS_DB debe ser un número entero igual o mayor que cero")
	}

	return &redis.Options{
		Addr:     fmt.Sprintf("%s:%s", valorODefecto("REDIS_HOST", "localhost"), valorODefecto("REDIS_PORT", "6379")),
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       base,
	}, nil
}
