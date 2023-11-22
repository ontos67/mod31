package pstg

import (
	"GoNews/pkg/storage"
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/jackc/pgx/v4/pgxpool"
)

// База данных.
type PSTG struct {
	pool *pgxpool.Pool
}

func New() (*PSTG, error) {
	connstr := os.Getenv("agrigatordb")
	if connstr == "" {
		return nil, errors.New("не указано подключение к БД")
	}
	pool, err := pgxpool.Connect(context.Background(), connstr)

	if err != nil {
		return nil, err
	}
	db := PSTG{
		pool: pool,
	}
	return &db, nil
}
func (db *PSTG) Posts() ([]storage.Post, error) {

	rows, err := db.pool.Query(context.Background(), `
	SELECT posts.id, posts.author_id,posts.title, posts.content, posts.created_at, authors.name FROM posts
	JOIN authors ON posts.author_id = authors.id
	`,
	)
	if err != nil {
		return nil, err
	}
	var plist []storage.Post
	for rows.Next() {
		var p storage.Post
		err = rows.Scan(
			p.ID,
			p.Title,
			p.Content,
			p.AuthorID,
			p.AuthorName,
			p.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		plist = append(plist, p)
	}
	return plist, rows.Err()
}

func (db *PSTG) AddPost(p storage.Post) error {
	_, err := db.pool.Exec(context.Background(), `
		INSERT INTO posts (author_id, title, content, created_at)
		VALUES ($1, $2, $3, $4)
		`,
		p.AuthorID,
		p.Title,
		p.Content,
		p.CreatedAt,
	)

	if err != nil {
		return err
	}

	return nil
}

func (db *PSTG) UpdatePost(p storage.Post) error {
	r, err := db.pool.Exec(context.Background(), `
		UPDATE posts
		SET author_id =$1, title = $2, content=$3, created_at=$4
		WHERE id =$5;
		`,
		p.AuthorID,
		p.Title,
		p.Content,
		p.CreatedAt,
		p.ID,
	)
	if r.RowsAffected() != 1 {
		return fmt.Errorf("запись с таким id не найдена")
	}
	if err != nil {
		return err
	}

	return nil
}

func (db *PSTG) DeletePost(p storage.Post) error {
	r, err := db.pool.Exec(context.Background(), `
		DELETE FROM posts
		WHERE id =$1;
		`,
		p.ID,
	)
	if r.RowsAffected() != 1 {
		return fmt.Errorf("запись с таким id не найдена")
	}
	if err != nil {
		return err
	}

	return nil
}
