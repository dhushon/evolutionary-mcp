package repository

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

func TestPostgresMemoryStore(t *testing.T) {
	ctx := context.Background()

	pgContainer, err := postgres.Run(ctx,
		"postgres:16-alpine",
		postgres.WithDatabase("test-db"),
		postgres.WithUsername("user"),
		postgres.WithPassword("password"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2)),
	)
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		if err := pgContainer.Terminate(ctx); err != nil {
			t.Fatalf("failed to terminate container: %s", err)
		}
	}()

	connStr, err := pgContainer.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		t.Fatal(err)
	}

	pool, err := pgxpool.New(ctx, connStr)
	if err != nil {
		t.Fatal(err)
	}
	defer pool.Close()

	store := NewPostgresMemoryStore(pool)

	_, err = pool.Exec(ctx, `CREATE EXTENSION IF NOT EXISTS "uuid-ossp"; CREATE TABLE memories (
		id UUID PRIMARY KEY,
		content TEXT NOT NULL,
		embedding VECTOR(3),
		confidence FLOAT NOT NULL,
		version INT NOT NULL
	);`)
	if err != nil {
		t.Fatal(err)
	}

	t.Run("Save and Get", func(t *testing.T) {
		id := uuid.New().String()
		memory := &Memory{
			ID:         id,
			Content:    "test content",
			Confidence: 0.9,
			Version:    1,
		}

		err := store.Save(ctx, memory)
		assert.NoError(t, err)

		retrieved, err := store.Get(ctx, id)
		assert.NoError(t, err)
		assert.Equal(t, memory.ID, retrieved.ID)
		assert.Equal(t, memory.Content, retrieved.Content)
		assert.Equal(t, memory.Confidence, retrieved.Confidence)
		assert.Equal(t, memory.Version, retrieved.Version)
	})
}
