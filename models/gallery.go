package models

import (
	"database/sql"
	"fmt"
	"io"
	"io/fs"
	"os"
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
	if err := os.RemoveAll(service.galleryDir(id)); err != nil {
		return fmt.Errorf("delete gallery images: %w", err)
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

func (service *GalleryService) Image(galleryID uint, filename string) (Image, error) {
	galleryDir := service.galleryDir(galleryID)
	imagePath := filepath.Join(galleryDir, filename)
	_, err := os.Stat(imagePath)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return Image{}, ErrNotFound
		}
		return Image{}, fmt.Errorf("querying for image: %w", err)
	}
	return Image{
		Filename:  filename,
		GalleryID: galleryID,
		Path:      imagePath,
	}, nil
}

func (service *GalleryService) CreateImage(galleryID uint, filename string, contents io.ReadSeeker) error {
	if err := checkContentType(contents, service.imageContentTypes()); err != nil {
		return fmt.Errorf("creating image %v: %w", filename, err)
	}
	if err := checkExtension(filename, service.extensions()); err != nil {
		return fmt.Errorf("creating image %v: %w", filename, err)
	}
	galleryDir := service.galleryDir(galleryID)
	if err := os.MkdirAll(galleryDir, 0755); err != nil { // make sure directory exists and is accessible
		return fmt.Errorf("create gallery-%d image directory: %w", galleryID, err)
	}
	imagePath := filepath.Join(galleryDir, filename)
	dst, err := os.Create(imagePath)
	if err != nil {
		return fmt.Errorf("creating image file: %w", err)
	}
	defer dst.Close()
	if _, err := io.Copy(dst, contents); err != nil {
		return fmt.Errorf("copying to image file: %w", err)
	}
	return nil
}

func (service *GalleryService) DeleteImage(galleryID uint, filename string) error {
	image, err := service.Image(galleryID, filename)
	if err != nil {
		return fmt.Errorf("deleting image: %w", err)
	}
	if err = os.Remove(image.Path); err != nil {
		return fmt.Errorf("deleting image: %w", err)
	}
	return nil
}

func (service *GalleryService) extensions() []string {
	return []string{".jpg", ".jpeg", ".png", ".gif"}
}

func (service *GalleryService) imageContentTypes() []string {
	return []string{"image/png", "image/jpg", "image/gif", "image/jpeg"}
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
