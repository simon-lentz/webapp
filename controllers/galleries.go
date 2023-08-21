package controllers

import (
	"fmt"
	"net/http"
	"net/url"
	"path/filepath"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/simon-lentz/webapp/context"
	"github.com/simon-lentz/webapp/errors"
	"github.com/simon-lentz/webapp/models"
)

type Galleries struct {
	Templates struct {
		New   Template
		Edit  Template
		Index Template
		Show  Template
	}
	GalleryService *models.GalleryService
}

func (g Galleries) New(w http.ResponseWriter, r *http.Request) {
	var data struct {
		Title string
	}
	data.Title = r.FormValue("title")
	g.Templates.New.Execute(w, r, data)
}

func (g Galleries) Create(w http.ResponseWriter, r *http.Request) {
	var data struct {
		UserID uint
		Title  string
	}
	data.UserID = context.User(r.Context()).ID
	data.Title = r.FormValue("title")
	gallery, err := g.GalleryService.Create(data.Title, data.UserID)
	if err != nil {
		g.Templates.New.Execute(w, r, data, err)
		return
	}
	editPath := fmt.Sprintf("/galleries/%d/edit", gallery.ID)
	http.Redirect(w, r, editPath, http.StatusFound)
}

func (g Galleries) Edit(w http.ResponseWriter, r *http.Request) {
	gallery, err := g.galleryByID(w, r, userMustOwnGalleryOpt)
	if err != nil {
		return
	}
	type Image struct {
		GalleryID       uint
		Filename        string
		FilenameEscaped string
	}
	var data struct {
		ID     uint
		Title  string
		Images []Image
	}
	data.ID = gallery.ID
	data.Title = gallery.Title
	images, err := g.GalleryService.Images(gallery.ID)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Something Went Wrong", http.StatusInternalServerError)
		return
	}
	for _, image := range images {
		data.Images = append(data.Images, Image{
			GalleryID:       image.GalleryID,
			Filename:        image.Filename,
			FilenameEscaped: url.PathEscape(image.Filename),
		})
	}
	g.Templates.Edit.Execute(w, r, data)
}

func (g Galleries) Update(w http.ResponseWriter, r *http.Request) {
	gallery, err := g.galleryByID(w, r, userMustOwnGalleryOpt)
	if err != nil {
		return
	}
	gallery.Title = r.FormValue("title")
	if err = g.GalleryService.Update(gallery); err != nil {
		http.Error(w, "Something went wrong.", http.StatusInternalServerError)
		return
	}
	editPath := fmt.Sprintf("/galleries/%d/edit", gallery.ID)
	http.Redirect(w, r, editPath, http.StatusFound)
}

func (g Galleries) DeleteImage(w http.ResponseWriter, r *http.Request) {
	filename := chi.URLParam(r, "filename")
	gallery, err := g.galleryByID(w, r, userMustOwnGalleryOpt)
	if err != nil {
		return
	}
	if err = g.GalleryService.DeleteImage(gallery.ID, filename); err != nil {
		http.Error(w, "Something Went Wrong", http.StatusInternalServerError)
		return
	}
	editPath := fmt.Sprintf("/galleries/%d/edit", gallery.ID)
	http.Redirect(w, r, editPath, http.StatusFound)
}
func (g Galleries) Index(w http.ResponseWriter, r *http.Request) {
	// The new type is helpful in case client side rendering diverges from server side representation.
	type Gallery struct {
		ID    uint
		Title string
	}
	var data struct {
		Galleries []Gallery
	}
	user := context.User(r.Context())
	galleries, err := g.GalleryService.ByUserID(user.ID)
	if err != nil {
		http.Error(w, "Something went wrong.", http.StatusInternalServerError)
		return
	}
	for _, gallery := range galleries {
		data.Galleries = append(data.Galleries, Gallery{
			ID:    gallery.ID,
			Title: gallery.Title,
		})
	}
	g.Templates.Index.Execute(w, r, data)
}

func (g Galleries) Show(w http.ResponseWriter, r *http.Request) {
	gallery, err := g.galleryByID(w, r)
	if err != nil {
		return
	}
	type Image struct {
		GalleryID uint
		Filename  string
		// URL safe filename
		FilenameEscaped string
	}
	var data struct {
		ID     uint
		Title  string
		Images []Image
	}
	data.ID = gallery.ID
	data.Title = gallery.Title
	images, err := g.GalleryService.Images(gallery.ID)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Something Went Wrong", http.StatusInternalServerError)
		return
	}
	for _, image := range images {
		data.Images = append(data.Images, Image{
			GalleryID:       image.GalleryID,
			Filename:        image.Filename,
			FilenameEscaped: url.PathEscape(image.Filename),
		})
	}
	g.Templates.Show.Execute(w, r, data)
}

func (g Galleries) Delete(w http.ResponseWriter, r *http.Request) {
	gallery, err := g.galleryByID(w, r, userMustOwnGalleryOpt)
	if err != nil {
		return
	}
	if err = g.GalleryService.Delete(gallery.ID); err != nil {
		http.Error(w, "Something Went Wrong", http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/galleries", http.StatusFound)
}

func (g Galleries) Image(w http.ResponseWriter, r *http.Request) {
	filename := g.filename(w, r)
	galleryID, err := strconv.Atoi(chi.URLParam(r, "id")) //ascii to int
	if err != nil {
		if errors.Is(err, models.ErrNotFound) {
			http.Error(w, "Image Not Found", http.StatusNotFound)
			return
		}
		http.Error(w, "Invalid ID", http.StatusNotFound)
		return
	}
	image, err := g.GalleryService.Image(uint(galleryID), filename)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Something Went Wrong", http.StatusInternalServerError)
		return
	}
	http.ServeFile(w, r, image.Path)
}

type galleryOpt func(http.ResponseWriter, *http.Request, *models.Gallery) error

// Combine with functional options pattern.
func (g Galleries) galleryByID(w http.ResponseWriter, r *http.Request, opts ...galleryOpt) (*models.Gallery, error) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusNotFound)
		return nil, err
	}
	gallery, err := g.GalleryService.ByID(uint(id))
	if err != nil {
		if errors.Is(err, models.ErrNotFound) {
			http.Error(w, "Gallery Not Found", http.StatusNotFound)
			return nil, err
		}
		http.Error(w, "Something Went Wrong", http.StatusInternalServerError)
		return nil, err
	}
	for _, opt := range opts {
		if err = opt(w, r, gallery); err != nil {
			return nil, err
		}
	}

	return gallery, nil
}

func userMustOwnGalleryOpt(w http.ResponseWriter, r *http.Request, gallery *models.Gallery) error {
	user := context.User(r.Context())
	if gallery.UserID != user.ID {
		http.Error(w, "Resource Not Found", http.StatusNotFound)
		return fmt.Errorf("Resource Not Found")
	}
	return nil
}

// Prevent url injection by trimming to base path.
func (g Galleries) filename(w http.ResponseWriter, r *http.Request) string {
	filePath := chi.URLParam(r, "filename")
	fileName := filepath.Base(filePath)
	return fileName
}
