package domain

type RecipeLine struct {
	id         string
	materialID string
	quantity   float64
}

func NewRecipeLine(id, materialID string, quantity float64) (*RecipeLine, error) {
	if id == "" {
		return nil, &ErrInvalidInput{FieldName: "id", Reason: "id is required"}
	}
	if materialID == "" {
		return nil, &ErrInvalidInput{FieldName: "materialID", Reason: "materialID is required"}
	}
	if quantity <= 0 {
		return nil, &ErrInvalidInput{FieldName: "quantity", Reason: "quantity must be greater than zero"}
	}

	return &RecipeLine{
		id:         id,
		materialID: materialID,
		quantity:   quantity,
	}, nil
}

func (rl *RecipeLine) ID() string { return rl.id }

func (rl *RecipeLine) MaterialID() string { return rl.materialID }

func (rl *RecipeLine) Quantity() float64 { return rl.quantity }
