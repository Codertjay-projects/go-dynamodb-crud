package health

import (
	"errors"
	"go-dynamodb-crud/internal/handlers"
	"go-dynamodb-crud/internal/repository/adapter"
	HttpStatus "go-dynamodb-crud/utils/http"
	"net/http"
)

type Handler struct {
	// uses an interface which is used to access the function containing the
	// request type
	handlers.Interface
	// used to communicate with the database, and it contains all
	// the function need to communicate to the database
	Repository adapter.Interface
}

func NewHandler(repository adapter.Interface) handlers.Interface {
	return &Handler{
		Repository: repository,
	}
}

func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
	if !h.Repository.Health() {

		HttpStatus.StatusInternalServerError(w, r,
			errors.New("Relational database not alive"))
		return
	}
	HttpStatus.StatusOk(w, r, "Service Ok")
}

func (h *Handler) Post(w http.ResponseWriter, r *http.Request) {
	HttpStatus.StatusMethodNotAllowed(w, r)
}
func (h *Handler) Put(w http.ResponseWriter, r *http.Request) {
	HttpStatus.StatusMethodNotAllowed(w, r)
}
func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	HttpStatus.StatusMethodNotAllowed(w, r)
}
func (h *Handler) Options(w http.ResponseWriter, r *http.Request) {
	HttpStatus.StatusMethodNotAllowed(w, r)
}
