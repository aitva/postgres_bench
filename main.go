//go:generate sqlc generate --file sqlc.yaml
//go:generate sqlc generate --file sqlc.pgx.yaml

package main

import (
	"cmp"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"io"
	"os"
	"slices"

	"github.com/cheggaaa/pb/v3"
	"github.com/google/uuid"
	_ "github.com/lib/pq"

	"github.com/aitva/postgres_bench/dataset"
	"github.com/aitva/postgres_bench/db/pq"
	"github.com/aitva/postgres_bench/decoder"
)

const postgresURI = `postgres://postgres_bench:postgres_bench@localhost/postgres_bench?sslmode=disable`
const postgresSchema = `schema.sql`

var spinner pb.ProgressBarTemplate = `{{with string . "prefix"}}{{.}} {{end}} {{ cycle . "⠋" "⠙" "⠹" "⠸" "⠼" "⠴" "⠦" "⠧" "⠇" "⠏" }} {{counters .}} {{speed . "%s p/s"}} {{with string . "suffix"}} {{.}}{{end}}`
var spinnerETA pb.ProgressBarTemplate = `{{with string . "prefix"}}{{.}} {{end}} {{ cycle . "⠋" "⠙" "⠹" "⠸" "⠼" "⠴" "⠦" "⠧" "⠇" "⠏" }} {{counters .}} {{speed . "%s p/s"}} {{rtime .}}{{with string . "suffix"}} {{.}}{{end}}`

func main() {
	fmt.Printf("Loading Wikipedia dataset...\n")
	datasets, err := loadDatasets()
	if err != nil {
		fatalf("fail to load datasetis: %v", err)
	}
	defer datasets.Close()

	fmt.Printf("Setting up database...\n")
	db, err := setupDatabase()
	if err != nil {
		fatalf("fail to setup database: %v", err)
	}
	defer db.Close()

	fmt.Printf("Loading dataset into SQLite...\n")
	count := 0
	for _, d := range datasets {
		bar := spinner.Start(0).Set("prefix", "  "+d.Name())
		c, err := insertDataset(bar, db, d)
		if err != nil {
			fatalf("fail to insert dataset: %v", err)
		}
		count += c
		bar.Finish()
	}

	fmt.Printf("Completed, %d pages created.\n", count)
}

func loadDatasets() (dataset.Datasets, error) {
	type result struct {
		Dataset *dataset.Dataset
		Err     error
	}

	pool := pb.NewPool()
	results := make(chan result)
	for _, name := range dataset.Names {
		bar := spinnerETA.New(0).
			Set(pb.Bytes, true).
			Set("prefix", "  "+name)
		pool.Add(bar)

		go func(name string, bar *pb.ProgressBar) {
			defer bar.Finish()

			d, err := dataset.Load(name)
			if errors.Is(err, dataset.ErrNotFound) {
				d, err = dataset.Download(name, func(n int64, r io.Reader) io.Reader {
					return bar.AddTotal(n).NewProxyReader(r)
				})
				results <- result{d, err}
				return
			}
			if err != nil {
				results <- result{nil, err}
				return
			}

			bar.AddTotal(d.Size())
			bar.Add64(d.Size())
			results <- result{d, nil}
		}(name, bar)
	}

	err := pool.Start()
	if err != nil {
		return nil, fmt.Errorf("start pool: %v", err)
	}
	defer pool.Stop()

	datasets := make(dataset.Datasets, len(dataset.Names))
	for i := range datasets {
		result := <-results
		if result.Err != nil {
			return nil, result.Err
		}
		datasets[i] = result.Dataset
	}

	slices.SortFunc(datasets, func(a, b *dataset.Dataset) int {
		return cmp.Compare(a.Name(), b.Name())
	})
	return datasets, nil
}

func insertDataset(bar *pb.ProgressBar, db *sql.DB, dataset *dataset.Dataset) (int, error) {
	d, err := decoder.New(dataset)
	if err != nil {
		return 0, fmt.Errorf("create decoder: %v\n", err)
	}

	tx, err := db.Begin()
	if err != nil {
		return 0, fmt.Errorf("begin transaction: %v\n", err)
	}
	defer tx.Rollback()

	var count int
	queries := pq.New(tx)
	for d.Next() {
		var p decoder.Page

		err := d.Scan(&p)
		if err != nil {
			return 0, fmt.Errorf("scan page: %v\n", err)
		}

		err = queries.CreatePage(context.TODO(), pq.CreatePageParams{
			ID:        uuid.Must(uuid.NewRandom()),
			UpdatedAt: p.UpdatedAt,
			Title:     p.Title,
			Text:      p.Text,
		})
		if err != nil {
			return 0, fmt.Errorf("create page: %v\n", err)
		}

		count++
		bar.Increment()
	}
	if err := d.Err(); err != nil && !errors.Is(err, io.EOF) {
		return 0, fmt.Errorf("scan pages: %v\n", err)
	}
	if err := tx.Commit(); err != nil {
		return 0, fmt.Errorf("commit transaction: %v\n", err)
	}
	return count, nil
}

func setupDatabase() (*sql.DB, error) {
	db, err := sql.Open("postgres", postgresURI)
	if err != nil {
		return nil, fmt.Errorf("open db: %v", err)
	}

	// Drop existing tables.
	username := ""
	err = db.QueryRow("SELECT USER").Scan(&username)
	if err != nil {
		return nil, err
	}

	_, err = db.Exec(fmt.Sprintf("DROP OWNED BY %q", username))
	if err != nil {
		return nil, err
	}

	// Create schema.
	schema, err := os.ReadFile(postgresSchema)
	if err != nil {
		return nil, fmt.Errorf("read file %v: %v", postgresSchema, err)
	}

	_, err = db.Exec(string(schema))
	if err != nil {
		return nil, err
	}

	return db, nil
}

func fatalf(format string, v ...any) {
	fmt.Printf(format, v...)
	os.Exit(1)
}
