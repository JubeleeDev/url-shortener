package shortener

import (
	"context"
	"errors"
	"os"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
)

func TestPostgresStore(t *testing.T) {

	dsn := os.Getenv("TEST_DATABASE_URL")

	if dsn == "" {
		t.Skip("TEST_DATABASE_URL not set")
	}

	pool, err := pgxpool.New(t.Context(), dsn)
	t.Cleanup(pool.Close)

	if err != nil {
		t.Skip("pool create error")
	}

	store := NewPostgresStore(pool)

	t.Run("save_and_find", func(t *testing.T) {
		err := TruncateLinks(t.Context(), pool)
		if err != nil {
			t.Fatalf("failed to truncate links: %v", err)
		}
		url := "https://google.com/"
		codeLength := 15
		link, err := NewLink(url, codeLength)

		if err != nil {
			t.Fatalf("error while creating link %v", err)
		}
		err = store.Save(t.Context(), link)

		if err != nil {
			t.Errorf("expect saved link, got err: %v", err)
		}

		foundedLink, err := store.Find(t.Context(), link.Code)
		if err != nil {
			t.Errorf("expect founded link, got err: %v", err)
		}

		if foundedLink.Code != link.Code {
			t.Errorf("expect code %s, got %s", foundedLink.Code, link.Code)
		}

		if foundedLink.OriginalURL != link.OriginalURL {
			t.Errorf("expect original url %s, got %s", foundedLink.OriginalURL, link.OriginalURL)
		}

	})

	t.Run("find_missing", func(t *testing.T) {
		err := TruncateLinks(t.Context(), pool)
		if err != nil {
			t.Fatalf("failed to truncate links: %v", err)
		}

		_, err = store.Find(t.Context(), "missing")
		if !errors.Is(err, ErrNotFound) {
			t.Error("expect not found by missing code, got link")
		}

	})

	t.Run("duplicate", func(t *testing.T) {
		err := TruncateLinks(t.Context(), pool)
		if err != nil {
			t.Fatalf("failed to truncate links: %v", err)
		}
		url := "https://google.com/"
		codeLength := 15
		link, err := NewLink(url, codeLength)

		if err != nil {
			t.Fatalf("error while creating link %v", err)
		}
		err = store.Save(t.Context(), link)

		if err != nil {
			t.Errorf("expect saved link, got err: %v", err)
		}

		err = store.Save(t.Context(), link)

		if err != nil {
			if !errors.Is(err, ErrUniqueConflict) {
				t.Errorf("expect error unique conflict, got %v", err)
			}
		} else {
			t.Error("expect error, got nil")
		}

	})

}

func TruncateLinks(ctx context.Context, pool *pgxpool.Pool) error {
	query := "TRUNCATE links"
	_, err := pool.Exec(ctx, query)
	return err
}
