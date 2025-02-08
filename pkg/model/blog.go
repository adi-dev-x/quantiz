package model

import (
	"database/sql"
	"time"
)

type Blog struct {
	ID           int64      `json:"id"`
	Title        string     `json:"title"`
	Descriptions *string    `json:"descriptions"`
	Author       string     `json:"author"`
	Content      string     `json:"content"`
	Categories   []Category `json:"categories"`
	UserID       int        `json:"user_id"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
}

func (b *Blog) Scan(rows *sql.Rows) error {
	var descriptions sql.NullString
	err := rows.Scan(
		&b.ID,
		&b.Title,
		&descriptions,
		&b.Content,
		&b.UserID,
		&b.CreatedAt,
		&b.UpdatedAt,

		&b.Author,
	)
	if err != nil {
		return err
	}

	if descriptions.Valid {
		b.Descriptions = &descriptions.String
	} else {
		b.Descriptions = nil
	}

	return nil
}

type Category struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}
type BlogCategory struct {
	BlogId     int64 `json:"blogId"`
	CategoryId int64 `json:"categoryId"`
}

func (c *Category) Scan(rows *sql.Rows) error {
	err := rows.Scan(&c.ID, &c.Name)
	if err != nil {
		return err
	}
	return nil
}
