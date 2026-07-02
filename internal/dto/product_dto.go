package dto

// CreateProductRequest is the payload for creating a product.
type CreateProductRequest struct {
	Name string `json:"name" validate:"required,min=1,max=255"`
}

// UpdateProductRequest is the payload for updating a product.
type UpdateProductRequest struct {
	Name string `json:"name" validate:"required,min=1,max=255"`
}
