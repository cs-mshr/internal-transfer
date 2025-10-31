package database

import (
	"context"
	"fmt"
	"net"
	"net/url"
	"strconv"
	"time"

	"github.com/chandra-shekhar/internal-transfers/internal/config"
	loggerConfig "github.com/chandra-shekhar/internal-transfers/internal/logger"
	pgxzero "github.com/jackc/pgx-zerolog"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/tracelog"
	"github.com/rs/zerolog"
)

type postgresDB struct {
	pool *pgxpool.Pool
	log  *zerolog.Logger
}

const DatabasePingTimeout = 10

func New(cfg *config.Config, logger *zerolog.Logger) (DB, error) {
	hostPort := net.JoinHostPort(cfg.Database.Host, strconv.Itoa(cfg.Database.Port))

	// URL-encode the password
	encodedPassword := url.QueryEscape(cfg.Database.Password)
	dsn := fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=%s",
		cfg.Database.User,
		encodedPassword,
		hostPort,
		cfg.Database.Name,
		cfg.Database.SSLMode,
	)

	pgxPoolConfig, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to parse pgx pool config: %w", err)
	}

	// Set connection pool configuration
	pgxPoolConfig.MaxConns = int32(cfg.Database.MaxOpenConns)
	pgxPoolConfig.MinConns = int32(cfg.Database.MaxIdleConns)
	pgxPoolConfig.MaxConnLifetime = time.Duration(cfg.Database.ConnMaxLifetime) * time.Second
	pgxPoolConfig.MaxConnIdleTime = time.Duration(cfg.Database.ConnMaxIdleTime) * time.Second

	// Enable query logging in local environment
	if cfg.Primary.Env == "local" {
		globalLevel := logger.GetLevel()
		pgxLogger := loggerConfig.NewPgxLogger(globalLevel)
		pgxPoolConfig.ConnConfig.Tracer = &tracelog.TraceLog{
			Logger:   pgxzero.NewLogger(pgxLogger),
			LogLevel: tracelog.LogLevel(loggerConfig.GetPgxTraceLogLevel(globalLevel)),
		}
	}

	pool, err := pgxpool.NewWithConfig(context.Background(), pgxPoolConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create pgx pool: %w", err)
	}

	database := &postgresDB{
		pool: pool,
		log:  logger,
	}

	ctx, cancel := context.WithTimeout(context.Background(), DatabasePingTimeout*time.Second)
	defer cancel()
	if err = pool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	logger.Info().Msg("connected to the database")

	return database, nil
}

func (db *postgresDB) Exec(ctx context.Context, sql string, arguments ...interface{}) (pgconn.CommandTag, error) {
	return db.pool.Exec(ctx, sql, arguments...)
}

func (db *postgresDB) Query(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error) {
	return db.pool.Query(ctx, sql, args...)
}

func (db *postgresDB) QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row {
	return db.pool.QueryRow(ctx, sql, args...)
}

func (db *postgresDB) Begin(ctx context.Context) (pgx.Tx, error) {
	return db.pool.Begin(ctx)
}

func (db *postgresDB) BeginTx(ctx context.Context, txOptions pgx.TxOptions) (pgx.Tx, error) {
	return db.pool.BeginTx(ctx, txOptions)
}

func (db *postgresDB) Ping(ctx context.Context) error {
	return db.pool.Ping(ctx)
}

func (db *postgresDB) Close() error {
	db.log.Info().Msg("closing database connection pool")
	db.pool.Close()
	return nil
}
