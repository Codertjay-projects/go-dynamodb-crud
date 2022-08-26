package product

import (
	"errors"
	"github.com/go-chi/chi"
	"github.com/google/uuid"
	productController "go-dynamodb-crud/internal/controllers"
	product "go-dynamodb-crud/internal/entities/product"
	"go-dynamodb-crud/internal/handlers"
	"go-dynamodb-crud/internal/repository/adapter"
	RulesProduct "go-dynamodb-crud/internal/rules/product"
	// Rules "go-dynamodb-crud/internal/rules"
	HttpStatus "go-dynamodb-crud/utils/http"
	"net/http"
	"time"
)

// the handler
type Handler struct {
	handlers.Interface
	Controller productController.Interface
	// i changed this to RulesProduct from product
	Rules RulesProduct.Interface
}

func NewHandler(repository adapter.Interface) handlers.Interface {
	return &Handler{
		Controller: productController.NewController(repository),
		Rules:      RulesProduct.NewRules(),
	}
}

func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
	// check the url params before moving to th url
	if chi.URLParam(r, "ID") != "" {
		h.GetOne(w, r)
	} else {
		h.GetAll(w, r)
	}
}

func (h *Handler) GetOne(w http.ResponseWriter, r *http.Request) {
	ID, err := uuid.Parse(chi.URLParam(r, "ID"))
	if err != nil {
		HttpStatus.StatusBadRequest(w, r,
			errors.New("ID is not uuid valid"))
		return
	}
	response, err := h.Controller.ListOne(ID)
	if err != nil {
		HttpStatus.StatusInternalServerError(w, r, err)
		return
	}
	HttpStatus.StatusOk(w, r, response)
}
func (h *Handler) GetAll(w http.ResponseWriter, r *http.Request) {
	response, err := h.Controller.ListAll()
	if err != nil {
		HttpStatus.StatusInternalServerError(w, r, err)
		return
	}
	HttpStatus.StatusOk(w, r, response)
}
func (h *Handler) Post(w http.ResponseWriter, r *http.Request) {
	productBody, err := h.GetBodyAndValidate(r, uuid.New())
	if err != nil {
		HttpStatus.StatusBadRequest(w, r, err)
		return
	}
	ID, err := h.Controller.Create(productBody)
	if err != nil {
		HttpStatus.StatusInternalServerError(w, r, err)
		return
	}
	HttpStatus.StatusOk(w, r, map[string]interface{}{"id": ID.String()})
}
func (h *Handler) Put(w http.ResponseWriter, r *http.Request) {
	ID, err := uuid.Parse(chi.URLParam(r, "ID"))
	if err != nil {
		HttpStatus.StatusBadRequest(w, r, errors.New("ID is not UUID valid"))
		return
	}
	productBody, err := h.GetBodyAndValidate(r, ID)
	if err != nil {
		HttpStatus.StatusBadRequest(w, r, err)
		return
	}
	if err := h.Controller.Update(productBody); err != nil {
		HttpStatus.StatusInternalServerError(w, r, err)
		return
	}
	HttpStatus.StatusNoContent(w, r)
}
func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	ID, err := uuid.Parse(chi.URLParam(r, "ID"))
	if err != nil {
		HttpStatus.StatusBadRequest(w, r, errors.New("ID is not uuid valid"))
		return
	}
	if err = h.Controller.Remove(ID); err != nil {
		HttpStatus.StatusInternalServerError(w, r, err)
		return
	}
	HttpStatus.StatusNoContent(w, r)
}
func (h *Handler) Options(w http.ResponseWriter, r *http.Request) {
	HttpStatus.StatusNoContent(w, r)
}

func (h *Handler) GetBodyAndValidate(r *http.Request, ID uuid.UUID) (response *product.Product, err error) {
	// this valid date and unmarshal
	productBody := product.Product{}
	body, err := h.Rules.ConvertIOReaderToStruct(r.Body, productBody)
	if err != nil {
		return &product.Product{}, errors.New("body is required")
	}
	productParsed, err := product.InterfaceToModel(body)
	if err != nil {
		return &product.Product{}, errors.New("error on converting body to model")
	}
	SetDefaultValues(productParsed, ID)
	return productParsed, h.Rules.Validate(productParsed)
}

func SetDefaultValues(product *product.Product, ID uuid.UUID) {
	currentTime, _ := time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	product.UpdatedAt = currentTime
	if ID == uuid.Nil {
		product.ID = uuid.New()
		product.CreatedAt = currentTime
	} else {
		product.ID = ID
	}
}
