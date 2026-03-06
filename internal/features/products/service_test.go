package products

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockRepository is a mock implementation of the Repository interface.
type MockRepository struct {
	mock.Mock
}

func (m *MockRepository) ListProducts(ctx context.Context, search string, limit, offset int32) ([]Product, int64, error) {
	args := m.Called(ctx, search, limit, offset)
	return args.Get(0).([]Product), args.Get(1).(int64), args.Error(2)
}

func (m *MockRepository) GetProductByID(ctx context.Context, id string) (Product, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(Product), args.Error(1)
}

func (m *MockRepository) CreateProduct(ctx context.Context, id, name, price string) error {
	args := m.Called(ctx, id, name, price)
	return args.Error(0)
}

func (m *MockRepository) UpdateProduct(ctx context.Context, id, name, price string) error {
	args := m.Called(ctx, id, name, price)
	return args.Error(0)
}

func (m *MockRepository) DeleteProduct(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockRepository) CheckNameExists(ctx context.Context, name string) (bool, error) {
	args := m.Called(ctx, name)
	return args.Bool(0), args.Error(1)
}

func (m *MockRepository) CheckNameExistsForOther(ctx context.Context, name, id string) (bool, error) {
	args := m.Called(ctx, name, id)
	return args.Bool(0), args.Error(1)
}

func TestListProducts(t *testing.T) {
	mockRepo := new(MockRepository)
	service := NewService(mockRepo, nil)

	ctx := context.Background()
	expectedProducts := []Product{
		{ID: "1", Name: "Product A", Price: "10.00"},
		{ID: "2", Name: "Product B", Price: "20.00"},
	}

	// Setup expectation: ListProducts(ctx, search, limit, offset)
	mockRepo.On("ListProducts", ctx, "test", int32(10), int32(0)).Return(expectedProducts, int64(2), nil)

	// Execute
	products, pagination, err := service.ListProducts(ctx, "test", 1, 10)

	// Assertions
	assert.NoError(t, err)
	assert.Len(t, products, 2)
	assert.Equal(t, int64(2), pagination.Total)
	assert.Equal(t, "Product A", products[0].Name)

	mockRepo.AssertExpectations(t)
}

func TestCreateProduct_DuplicateName(t *testing.T) {
	mockRepo := new(MockRepository)
	service := NewService(mockRepo, nil)

	ctx := context.Background()

	// Setup expectation: CheckNameExists returns true
	mockRepo.On("CheckNameExists", ctx, "Existing Product").Return(true, nil)

	// Execute
	err := service.CreateProduct(ctx, "Existing Product", "10.00")

	// Assertions
	assert.Error(t, err)
	assert.Equal(t, ErrProductNameExists, err)
	mockRepo.AssertExpectations(t)
}
