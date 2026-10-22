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
	var ultimoError error
	for intento := 1; intento <= 5; intento++ {
		ctxPing, cancelar := context.WithTimeout(ctx, 5*time.Second)
		ultimoError = cliente.Ping(ctxPing).Err()
		cancelar()
		if ultimoError == nil {
			return cliente, nil
		}
		if intento < 5 {
			select {
			case <-ctx.Done():
				_ = cliente.Close()
				return nil, ctx.Err()
			case <-time.After(2 * time.Second):
			}
		}
	}

	_ = cliente.Close()
	return nil, fmt.Errorf("Redis no respondió después de varios intentos: %w", ultimoError)
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
