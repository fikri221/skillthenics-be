package products

import (
	"fmt"
	"log"
	"nds-go-starter/internal/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type handler struct {
	service Service
}

func NewHandler(service Service) *handler {
	return &handler{
		service: service,
	}
}

type createProductRequest struct {
	Name  string  `json:"name" validate:"required,min=3"`
	Price float64 `json:"price" validate:"required,gt=0"`
}

type updateProductRequest struct {
	Name  string  `json:"name" validate:"required,min=3"`
	Price float64 `json:"price" validate:"required,gt=0"`
}

// ListProducts godoc
// @Summary      List Products
// @Description  Get a paginated list of products with optional search query
// @Tags         products
// @Accept       json
// @Produce      json
// @Param        query  query     string  false  "Search query"
// @Param        page   query     int     false  "Page number (default 1)"
// @Param        size   query     int     false  "Page size (default 10)"
// @Security     BearerAuth
// @Success      200    {object}  json.Response{data=[]Product}
// @Failure      500    {object}  json.Response
// @Router       /products [get]
func (h *handler) ListProducts(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("query")
	sizeStr := r.URL.Query().Get("size")
	pageStr := r.URL.Query().Get("page")

	size, _ := strconv.Atoi(sizeStr)
	if size <= 0 {
		size = 10
	}

	page, _ := strconv.Atoi(pageStr)
	if page <= 0 {
		page = 1
	}

	products, pag, err := h.service.ListProducts(r.Context(), query, page, size)
	if err != nil {
		log.Println(err)
		json.WriteError(w, r, err)
		return
	}

	json.WriteWithPagination(w, r, http.StatusOK, products, pag)
}

// CreateProduct godoc
// @Summary      Create Product
// @Description  Create a new product with name and price
// @Tags         products
// @Accept       json
// @Produce      json
// @Param        request  body      createProductRequest  true  "Product details"
// @Security     BearerAuth
// @Success      201      {object}  json.Response{data=string} "Product created successfully"
// @Failure      400      {object}  json.Response
// @Failure      401      {object}  json.Response
// @Failure      409      {object}  json.Response "Product name exists"
// @Failure      500      {object}  json.Response
// @Router       /products [post]
func (h *handler) CreateProduct(w http.ResponseWriter, r *http.Request) {
	var req createProductRequest
	if err := json.DecodeAndValidate(r, &req); err != nil {
		json.WriteError(w, r, err)
		return
	}

	priceStr := fmt.Sprintf("%.2f", req.Price)
	err := h.service.CreateProduct(r.Context(), req.Name, priceStr)
	if err != nil {
		json.WriteError(w, r, err)
		return
	}

	json.Write(w, r, http.StatusCreated, "Product created successfully")
}

// UpdateProduct godoc
// @Summary      Update Product
// @Description  Update an existing product by ID
// @Tags         products
// @Accept       json
// @Produce      json
// @Param        id       path      string                true  "Product ID"
// @Param        request  body      updateProductRequest  true  "Updated product details"
// @Security     BearerAuth
// @Success      200      {object}  json.Response{data=string} "Product updated successfully"
// @Failure      400      {object}  json.Response
// @Failure      401      {object}  json.Response
// @Failure      404      {object}  json.Response
// @Failure      409      {object}  json.Response
// @Failure      500      {object}  json.Response
// @Router       /products/{id} [put]
func (h *handler) UpdateProduct(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		json.Write(w, r, http.StatusBadRequest, "id is required")
		return
	}

	var req updateProductRequest
	if err := json.DecodeAndValidate(r, &req); err != nil {
		json.WriteError(w, r, err)
		return
	}

	priceStr := fmt.Sprintf("%.2f", req.Price)
	err := h.service.UpdateProduct(r.Context(), id, req.Name, priceStr)
	if err != nil {
		json.WriteError(w, r, err)
		return
	}

	json.Write(w, r, http.StatusOK, "Product updated successfully")
}

// DeleteProduct godoc
// @Summary      Delete Product
// @Description  Delete a product by ID
// @Tags         products
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "Product ID"
// @Security     BearerAuth
// @Success      200  {object}  json.Response{data=string} "Product deleted successfully"
// @Failure      400  {object}  json.Response
// @Failure      401  {object}  json.Response
// @Failure      404  {object}  json.Response
// @Failure      500  {object}  json.Response
// @Router       /products/{id} [delete]
func (h *handler) DeleteProduct(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		json.Write(w, r, http.StatusBadRequest, "id is required")
		return
	}

	if err := h.service.DeleteProduct(r.Context(), id); err != nil {
		json.WriteError(w, r, err)
		return
	}

	json.Write(w, r, http.StatusOK, "Product deleted successfully")
}
