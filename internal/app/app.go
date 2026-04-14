package app

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"superindo-test/internal/product"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

type Config struct {
	AppPort string

	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	DBSSLMode  string

	RedisAddr       string
	RedisPassword   string
	RedisDB         int
	RedisTTLSeconds int
}

type App struct {
	Config Config
	Router http.Handler
}

func NewApp() (*App, func(), error) {
	cfg := loadConfig()

	dbpool, err := newPGXPool(cfg)
	if err != nil {
		return nil, nil, err
	}

	rdb := newRedis(cfg)

	repo := product.NewRepository(dbpool)
	cache := product.NewCache(rdb, time.Duration(cfg.RedisTTLSeconds)*time.Second)
	svc := product.NewService(repo, cache)
	h := product.NewHandler(svc)

	router := newRouter(h)

	cleanup := func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_ = rdb.Close()
		dbpool.Close()
		_ = ctx
	}

	return &App{Config: cfg, Router: router}, cleanup, nil
}

func loadConfig() Config {
	return Config{
		AppPort: getEnv("APP_PORT", "8080"),

		DBHost:     getEnv("DB_HOST", "postgres"),
		DBPort:     getEnv("DB_PORT", "5432"),
		DBUser:     getEnv("DB_USER", "superindo"),
		DBPassword: getEnv("DB_PASSWORD", "superindo"),
		DBName:     getEnv("DB_NAME", "superindo"),
		DBSSLMode:  getEnv("DB_SSLMODE", "disable"),

		RedisAddr:       getEnv("REDIS_ADDR", "redis:6379"),
		RedisPassword:   getEnv("REDIS_PASSWORD", ""),
		RedisDB:         getEnvInt("REDIS_DB", 0),
		RedisTTLSeconds: getEnvInt("REDIS_TTL_SECONDS", 30),
	}
}

func newRouter(h *product.Handler) http.Handler {
	r := chi.NewRouter()

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	})

	r.Route("/product", func(r chi.Router) {
		r.Post("/", h.CreateProduct)
		r.Get("/", h.ListProducts)
	})

	return r
}

func newPGXPool(cfg Config) (*pgxpool.Pool, error) {
	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=%s",
		urlQueryEscape(cfg.DBUser),
		urlQueryEscape(cfg.DBPassword),
		cfg.DBHost,
		cfg.DBPort,
		cfg.DBName,
		cfg.DBSSLMode,
	)
	poolCfg, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, err
	}
	poolCfg.MaxConns = 10
	poolCfg.MinConns = 1
	poolCfg.MaxConnIdleTime = 5 * time.Minute

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	pool, err := pgxpool.NewWithConfig(ctx, poolCfg)
	if err != nil {
		return nil, err
	}
	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, err
	}
	log.Printf("postgres connected")
	return pool, nil
}

func newRedis(cfg Config) *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:     cfg.RedisAddr,
		Password: cfg.RedisPassword,
		DB:       cfg.RedisDB,
	})
	return rdb
}

func getEnv(key, def string) string {
	v := strings.TrimSpace(os.Getenv(key))
	if v == "" {
		return def
	}
	return v
}

func getEnvInt(key string, def int) int {
	v := strings.TrimSpace(os.Getenv(key))
	if v == "" {
		return def
	}
	n, err := strconv.Atoi(v)
	if err != nil {
		return def
	}
	return n
}

func urlQueryEscape(s string) string {
	// simple escape for user/pass in DSN without pulling net/url everywhere
	return strings.ReplaceAll(s, "@", "%40")
}

