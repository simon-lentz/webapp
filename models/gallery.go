package models

import (
	"database/sql"
	"fmt"

	"github.com/simon-lentz/webapp/errors"
)

type Gallery struct {
	ID     uint
	UserID uint
	Title  string
}

type GalleryService struct {
	DB *sql.DB
}

func (service *GalleryService) Create(title string, userID uint) (*Gallery, error) {
	gallery := Gallery{
		Title:  title,
		UserID: userID,
	}

	row := service.DB.QueryRow(`
	INSERT INTO galleries (title, user_id)
	VALUES ($1, $2) RETURNING id;`, gallery.Title, gallery.UserID)
	if err := row.Scan(&gallery.ID); err != nil {
		return nil, fmt.Errorf("create: %w", err)
	}

	return &gallery, nil
}

func (service *GalleryService) ByID(id uint) (*Gallery, error) {
	// TODO: add id validation
	gallery := Gallery{
		ID: id,
	}

	row := service.DB.QueryRow(`
	SELECT title, user_id
	FROM galleries
	WHERE id = $1;`, gallery.ID)
	if err := row.Scan(&gallery.Title, &gallery.UserID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("byID: %w", err)
	}

	return &gallery, nil
}

func (service *GalleryService) ByUserID(userID uint) ([]Gallery, error) {
	rows, err := service.DB.Query(`
	SELECT id, title
	FROM galleries
	WHERE user_id = $1;`, userID)
	if err != nil {
		return nil, fmt.Errorf("query gallery by userid: %w", err)
	}
	var galleries []Gallery
	for rows.Next() {
		gallery := Gallery{
			UserID: userID,
		}
		if err := rows.Scan(&gallery.ID, &gallery.Title); err != nil {
			return nil, fmt.Errorf("query gallery by userid: %w", err)
		}
		galleries = append(galleries, gallery)
	}
	if rows.Err() != nil {
		return nil, fmt.Errorf("query gallery by userid: %w", err)
	}

	return galleries, nil
}

func (service *GalleryService) Update(gallery *Gallery) error {
	if _, err := service.DB.Exec(`
	UPDATE galleries
	SET title = $2
	WHERE id = $1;`, gallery.ID, gallery.Title); err != nil {
		return fmt.Errorf("update: %w", err)
	}

	return nil
}

func (service *GalleryService) Delete(id uint) error {
	if _, err := service.DB.Exec(`
	DELETE FROM galleries
	WHERE id = $1;`, id); err != nil {
		return fmt.Errorf("delete: %w", err)
	}

	return nil
}