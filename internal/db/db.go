package db

import (
	"context"
	"os"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/micahasowata/blog/internal/config"
)

func NewProduction(cfg *config.Config) (*pgxpool.Pool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	db, err := pgxpool.New(ctx, cfg.ProdDSN)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func NewTest() (*pgxpool.Pool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	db, err := pgxpool.New(ctx, os.Getenv("TEST_DSN"))
	if err != nil {
		return nil, err
	}

	return db, nil
}

func Clean(db *pgxpool.Pool) error {
	query := `
	DO $$ DECLARE r RECORD;
		BEGIN FOR r IN (
			SELECT 
				tablename 
			FROM 
				pg_tables 
			WHERE 
				schema_name = 'public'
			) LOOP EXECUTE 'drop table if exists ' || quote_ident(r.tablename) || ' cascade';
		END LOOP;
	END $$;
	`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	tx, err := db.BeginTx(ctx, pgx.TxOptions{
		IsoLevel:       pgx.Serializable,
		AccessMode:     pgx.ReadWrite,
		DeferrableMode: pgx.NotDeferrable,
	})

	if err != nil {
		return err
	}

	defer tx.Rollback(ctx)

	_, err = tx.Exec(ctx, query)
	if err != nil {
		return err
	}

	return nil
}
