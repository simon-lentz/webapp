package models

import (
	"database/sql"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/simon-lentz/webapp/errors"
)

type Image struct {
	GalleryID uint
	Path      string
	Filename  string
}

type Gallery struct {
	ID     uint
	UserID uint
	Title  string
}

type GalleryService struct {
	DB        *sql.DB
	ImagesDir string
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

// Glob images in ImagesDir
func (service *GalleryService) Images(galleryID uint) ([]Image, error) {
	globPattern := filepath.Join(service.galleryDir(galleryID), "*")
	allFiles, err := filepath.Glob(globPattern)
	if err != nil {
		return nil, fmt.Errorf("retrieving gallery images: %w", err)
	}
	var images []Image
	for _, file := range allFiles {
		if hasExtension(file, service.extensions()) {
			images = append(images,
				Image{
					GalleryID: galleryID,
					Path:      file,
					Filename:  filepath.Base(file),
				},
			)
		}
	}
	return images, nil

}

func (service *GalleryService) extensions() []string {
	return []string{".jpg", ".jpeg", ".png", ".gif"}
}

func (service *GalleryService) galleryDir(id uint) string {
	imagesDir := service.ImagesDir
	if imagesDir == "" {
		imagesDir = "images"
	}
	return filepath.Join(imagesDir, fmt.Sprintf("gallery-%d", id))
}

func hasExtension(file string, extensions []string) bool {
	for _, ext := range extensions {
		file = strings.ToLower(file)
		ext = strings.ToLower(ext)
		if filepath.Ext(file) == ext {
			return true
		}
	}
	return false
}
