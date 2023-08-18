package controllers

import "net/http"

// This interface decouples the controllers and views modules.
type Template interface {
	Execute(w http.ResponseWriter, r *http.Request, data interface{}, errs ...error)
}
