package main

import (
	"context"
	"database/sql"
	"testing"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	_ "github.com/jackc/pgx/v5/stdlib"
	_ "github.com/lib/pq"

	pgxdb "github.com/aitva/postgres_bench/db/pgx"
	"github.com/aitva/postgres_bench/db/pq"
)

func BenchmarkGetPage_pq(b *testing.B) {
	db, err := sql.Open("postgres", postgresURI)
	if err != nil {
		b.Fatalf("fail to open db: %v", err)
	}
	defer db.Close()

	queries := pq.New(db)
	ids, err := queries.ListPageIDs(context.Background(), pq.ListPageIDsParams{
		Limit: 10000,
	})
	if err != nil {
		b.Fatalf("fail to list ids: %v", err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		id := ids[i%len(ids)]
		p, err := queries.GetPage(context.Background(), id)
		if err != nil {
			b.Fatalf("fail to get page %v: %v", id, err)
		}
		_ = p
	}
}

func BenchmarkGetPage_pgxstdlib(b *testing.B) {
	db, err := sql.Open("pgx", postgresURI)
	if err != nil {
		b.Fatalf("fail to open db: %v", err)
	}
	defer db.Close()

	queries := pq.New(db)
	ids, err := queries.ListPageIDs(context.Background(), pq.ListPageIDsParams{
		Limit: 10000,
	})
	if err != nil {
		b.Fatalf("fail to list ids: %v", err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		id := ids[i%len(ids)]
		p, err := queries.GetPage(context.Background(), id)
		if err != nil {
			b.Fatalf("fail to get page %v: %v", id, err)
		}
		_ = p
	}
}

func BenchmarkGetPage_pgx(b *testing.B) {
	conn, err := pgx.Connect(context.Background(), postgresURI)
	if err != nil {
		b.Fatalf("fail to connect to db: %v", err)
	}
	defer conn.Close(context.Background())

	queries := pgxdb.New(conn)
	ids, err := queries.ListPageIDs(context.Background(), pgxdb.ListPageIDsParams{
		Limit: 10000,
	})
	if err != nil {
		b.Fatalf("fail to list ids: %v", err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		id := ids[i%len(ids)]
		p, err := queries.GetPage(context.Background(), id)
		if err != nil {
			b.Fatalf("fail to get page %v: %v", id, err)
		}
		_ = p
	}
}

func BenchmarkListPages_pq(b *testing.B) {
	db, err := sql.Open("postgres", postgresURI)
	if err != nil {
		b.Fatalf("fail to open db: %v", err)
	}
	defer db.Close()
	queries := pq.New(db)

	b.ResetTimer()
	var cursor uuid.NullUUID
	for i := 0; i < b.N; i++ {
		pages, err := queries.ListPages(context.Background(), pq.ListPagesParams{
			Cursor: cursor,
			Limit:  1000,
		})
		if err != nil {
			b.Fatalf("fail to list pages: %v", err)
		}

		cursor.UUID = pages[len(pages)-1].ID
		cursor.Valid = true
	}
}

func BenchmarkListPages_pgxstlib(b *testing.B) {
	db, err := sql.Open("pgx", postgresURI)
	if err != nil {
		b.Fatalf("fail to open db: %v", err)
	}
	defer db.Close()
	queries := pq.New(db)

	b.ResetTimer()
	var cursor uuid.NullUUID
	for i := 0; i < b.N; i++ {
		pages, err := queries.ListPages(context.Background(), pq.ListPagesParams{
			Cursor: cursor,
			Limit:  1000,
		})
		if err != nil {
			b.Fatalf("fail to list pages: %v", err)
		}

		cursor.UUID = pages[len(pages)-1].ID
		cursor.Valid = true
	}
}

func BenchmarkListPages_pgx(b *testing.B) {
	conn, err := pgx.Connect(context.Background(), postgresURI)
	if err != nil {
		b.Fatalf("fail to open conn: %v", err)
	}
	defer conn.Close(context.Background())
	queries := pgxdb.New(conn)

	b.ResetTimer()
	var cursor uuid.NullUUID
	for i := 0; i < b.N; i++ {
		pages, err := queries.ListPages(context.Background(), pgxdb.ListPagesParams{
			Cursor: cursor,
			Limit:  1000,
		})
		if err != nil {
			b.Fatalf("fail to list pages: %v", err)
		}

		cursor.UUID = pages[len(pages)-1].ID
		cursor.Valid = true
	}
}
