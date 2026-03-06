package orders

import (
	"nds-go-starter/internal/json"
	"net/http"

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

type createOrderItemRequest struct {
	ProductID string  `json:"product_id" validate:"required"`
	Quantity  int32   `json:"quantity" validate:"required,gt=0"`
	Price     float64 `json:"price" validate:"required,gt=0"`
}

type createOrderRequest struct {
	CustomerName string                   `json:"customer_name" validate:"required,min=3"`
	Items        []createOrderItemRequest `json:"items" validate:"required,min=1,dive"`
}

// GetOrder godoc
// @Summary      Get Order
// @Description  Get order details including items and product names by ID
// @Tags         orders
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "Order ID"
// @Security     BearerAuth
// @Success      200  {object}  json.Response{data=Order}
// @Failure      400  {object}  json.Response
// @Failure      401  {object}  json.Response
// @Failure      404  {object}  json.Response
// @Failure      500  {object}  json.Response
// @Router       /orders/{id} [get]
func (h *handler) GetOrder(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		json.Write(w, r, http.StatusBadRequest, "id is required")
		return
	}

	order, err := h.service.GetOrder(r.Context(), id)
	if err != nil {
		json.WriteError(w, r, err)
		return
	}

	json.Write(w, r, http.StatusOK, order)
}

// CreateOrder godoc
// @Summary      Create Order
// @Description  Create a new order with multiple items (Transactional)
// @Tags         orders
// @Accept       json
// @Produce      json
// @Param        request  body      createOrderRequest  true  "Order details"
// @Security     BearerAuth
// @Success      201      {object}  json.Response{data=map[string]string} "Order created successfully"
// @Failure      400      {object}  json.Response
// @Failure      401      {object}  json.Response
// @Failure      500      {object}  json.Response
// @Router       /orders [post]
func (h *handler) CreateOrder(w http.ResponseWriter, r *http.Request) {
	var req createOrderRequest
	if err := json.DecodeAndValidate(r, &req); err != nil {
		json.WriteError(w, r, err)
		return
	}

	var items []OrderItem
	for _, it := range req.Items {
		items = append(items, OrderItem{
			ProductID: it.ProductID,
			Quantity:  it.Quantity,
			Price:     it.Price,
		})
	}

	orderID, err := h.service.CreateOrder(r.Context(), req.CustomerName, items)
	if err != nil {
		json.WriteError(w, r, err)
		return
	}

	json.Write(w, r, http.StatusCreated, map[string]string{
		"order_id": orderID,
		"message":  "Order created successfully",
	})
}
