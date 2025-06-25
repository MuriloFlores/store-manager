package user

type PromoteUserRequest struct {
	Role string `json:"role" validate:"required,oneof=admin manager salesperson client stock_person cashier"`
}
