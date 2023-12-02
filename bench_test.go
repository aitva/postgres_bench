package main

import (
	"context"
	"database/sql"
	"testing"

	"github.com/google/uuid"
	_ "github.com/jackc/pgx/v5/stdlib"
	_ "github.com/lib/pq"

	"github.com/aitva/postgres_bench/db/pq"
)

func BenchmarkGetPage_PQ(b *testing.B) {
	db, err := sql.Open("postgres", postgresURI)
	if err != nil {
		b.Fatalf("fail to open db: %v", err)
	}
	defer db.Close()

	queries := pq.New(db)
	ids, err := queries.ListIDs(context.Background(), uuid.NullUUID{})
	if err != nil {
		b.Fatalf("fail to list ids: %v", err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		id := ids[i%len(ids)]
		_, err = queries.GetPage(context.Background(), id)
		if err != nil {
			b.Fatalf("fail to get page %v: %v", id, err)
		}
	}
}

func BenchmarkGetPage_PGX(b *testing.B) {
	db, err := sql.Open("pgx", postgresURI)
	if err != nil {
		b.Fatalf("fail to open db: %v", err)
	}
	defer db.Close()

	queries := pq.New(db)
	ids, err := queries.ListIDs(context.Background(), uuid.NullUUID{})
	if err != nil {
		b.Fatalf("fail to list ids: %v", err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		id := ids[i%len(ids)]
		_, err = queries.GetPage(context.Background(), id)
		if err != nil {
			b.Fatalf("fail to get page %v: %v", id, err)
		}
	}
}
