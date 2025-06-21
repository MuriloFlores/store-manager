package user

type UpdateUserRequest struct {
	Name *string `json:"name" validate:"omitempty,min=3,max=100"`
	Role *string `json:"role" validate:"omitempty,oneof=admin manager salesperson client stock_person cashier"`
}
