package postgres

import (
	"GoNews/pkg/storage"
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

// Хранилище данных.
type Store struct {
	db *pgxpool.Pool
}

// Конструктор объекта хранилища.
func New(constr string) (*Store, error) {
	db, err := pgxpool.New(context.Background(), constr)
	if err != nil {
		return nil, err
	}
	s := Store{
		db: db,
	}
	return &s, nil
}

// curl http://localhost:8080/posts
func (s *Store) Posts() ([]storage.Post, error) {
	rows, err := s.db.Query(context.Background(), `
	SELECT
		id,
		author_id,
		title,
		content, 
		created_at
	FROM posts
	ORDER BY id;
`)
	if err != nil {
		return nil, err
	}
	var tasks []storage.Post
	for rows.Next() {
		var t storage.Post
		err = rows.Scan(
			&t.ID,
			&t.AuthorID,
			&t.Title,
			&t.Content,
			&t.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		// добавление переменной в массив результатов
		tasks = append(tasks, t)

	}
	// ВАЖНО не забыть проверить rows.Err()
	return tasks, rows.Err()
}

// json example
// {title:"json title",content:"json content",author_id:1,"created_at":55}

// пример команды для теста
/*
curl -X POST http://localhost:8080/posts
-H "Content-Type: application/json"
-d '{"title":"json title","content":"json content","created_at":55,"author_id":1}'
*/

// ! надо доработать. добавить транзакцию, проверку автора и добавления автора если такого нет в таблице.
func (s *Store) AddPost(p storage.Post) error {
	query := `INSERT INTO posts (title, content, author_id, created_at) VALUES ($1, $2, $3, $4);`

	err := s.db.QueryRow(context.Background(), query, p.Title, p.Content, p.AuthorID, p.CreatedAt).Scan()
	return err

}

// curl request example
/*
curl -X PUT http://localhost:8080/posts -H "Content-Type: application/json" -d '{"id":15,"title":"updated title","content":"new_content","created_at":7766,"author_id":1}'
*/
// должны совпадать author_id и id, меняет title и content
func (pg *Store) UpdatePost(p storage.Post) error {
	query := `UPDATE posts SET title = $3, content = $4 WHERE author_id = $2 AND id=$1;`
	_, err := pg.db.Exec(context.Background(), query, p.ID, p.AuthorID, p.Title, p.Content)
	return err
}

// request example
/*
	curl -X DELETE http://localhost:8080/posts -H "Content-Type: application/json" -d '{"id":9}'
*/
func (s *Store) DeletePost(p storage.Post) error {
	query := `DELETE FROM posts WHERE id = $1;	`
	err := s.db.QueryRow(context.Background(), query, p.ID).Scan()
	return err
}
