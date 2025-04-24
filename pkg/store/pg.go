package store

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/grassrootseconomics/kv-vise-ru/pkg/data"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/tern/v2/migrate"
	"github.com/knadh/goyesql/v2"
)

type (
	PgOpts struct {
		Logg                 *slog.Logger
		DSN                  string
		MigrationsFolderPath string
		QueriesFolderPath    string
	}

	Pg struct {
		logg    *slog.Logger
		db      *pgxpool.Pool
		queries *queries
	}

	queries struct {
		GetSessionData string `query:"get-session-data"`
		Get            string `query:"get"`
	}
)

func NewPgStore(o PgOpts) (Store, error) {
	parsedConfig, err := pgxpool.ParseConfig(o.DSN)
	if err != nil {
		return nil, err
	}

	dbPool, err := pgxpool.NewWithConfig(context.Background(), parsedConfig)
	if err != nil {
		return nil, err
	}

	queries, err := loadQueries(o.QueriesFolderPath)
	if err != nil {
		return nil, err
	}

	// if err := runMigrations(context.Background(), dbPool, o.MigrationsFolderPath); err != nil {
	// 	return nil, err
	// }
	// o.Logg.Info("migrations ran successfully")

	return &Pg{
		logg:    o.Logg,
		db:      dbPool,
		queries: queries,
	}, nil
}

func (pg *Pg) Close() {
	pg.db.Close()
}

func (pg *Pg) GetSessionData(ctx context.Context, prefix []byte) (map[uint16][]string, error) {
	rows, err := pg.db.Query(ctx, pg.queries.GetSessionData, prefix)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make(map[uint16][]string)
	for rows.Next() {
		var (
			key   []byte
			value string
		)
		if err := rows.Scan(&key, &value); err != nil {
			return nil, err
		}
		dataType := uint16(key[len(key)-1])
		result[dataType] = append(result[dataType], value)
	}

	return result, rows.Err()
}

func (pg *Pg) GetAddress(ctx context.Context, sessionID string) (string, error) {
	var value string
	err := pg.db.QueryRow(ctx, pg.queries.Get, data.EncodeKey(sessionID, data.DATA_PUBLIC_KEY)).Scan(&value)
	if err != nil {
		return "", err
	}
	return value, nil
}

func loadQueries(queriesPath string) (*queries, error) {
	parsedQueries, err := goyesql.ParseFile(queriesPath)
	if err != nil {
		return nil, err
	}

	loadedQueries := &queries{}
	if err := goyesql.ScanToStruct(loadedQueries, parsedQueries, nil); err != nil {
		return nil, fmt.Errorf("failed to scan queries %v", err)
	}

	return loadedQueries, nil
}

func runMigrations(ctx context.Context, dbPool *pgxpool.Pool, migrationsPath string) error {
	const migratorTimeout = 15 * time.Second
	ctx, cancel := context.WithTimeout(ctx, migratorTimeout)
	defer cancel()

	conn, err := dbPool.Acquire(ctx)
	if err != nil {
		return err
	}
	defer conn.Release()

	migrator, err := migrate.NewMigrator(ctx, conn.Conn(), "schema_version")
	if err != nil {
		return err
	}

	if err := migrator.LoadMigrations(os.DirFS(migrationsPath)); err != nil {
		return err
	}

	if err := migrator.Migrate(ctx); err != nil {
		return err
	}

	return nil
}
