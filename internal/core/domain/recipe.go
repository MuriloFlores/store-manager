package domain

type Recipe struct {
	id          string
	productID   string
	ingredients []RecipeLine
}

func NewRecipe(id, productID string) (*Recipe, error) {
	if id == "" {
		return nil, &ErrInvalidInput{FieldName: "id", Reason: "id is required"}
	}
	if productID == "" {
		return nil, &ErrInvalidInput{FieldName: "productID", Reason: "productID is required"}
	}

	return &Recipe{
		id:          id,
		productID:   productID,
		ingredients: make([]RecipeLine, 0),
	}, nil
}

func HydrateRecipe(id, productID string, ingredients []RecipeLine) *Recipe {
	return &Recipe{
		id:          id,
		productID:   productID,
		ingredients: ingredients,
	}
}

func (r *Recipe) ID() string {
	return r.id
}

func (r *Recipe) ProductID() string {
	return r.productID
}

func (r *Recipe) Ingredients() []RecipeLine {
	return r.ingredients
}

func (r *Recipe) AddIngredient(newLine RecipeLine) error {
	for _, existingLine := range r.ingredients {
		if existingLine.MaterialID() == newLine.MaterialID() {
			return &ErrInvalidInput{FieldName: "ingredient", Reason: "ingredient already exists in the recipe"}
		}
	}

	r.ingredients = append(r.ingredients, newLine)
	return nil
}

func (r *Recipe) RemoveIngredient(materialID string) error {
	foundIndex := -1
	for i, line := range r.ingredients {
		if line.MaterialID() == materialID {
			foundIndex = i
			break
		}
	}

	if foundIndex == -1 {
		return &ErrNotFound{ResourceName: "ingredient", ResourceID: materialID}
	}

	r.ingredients = append(r.ingredients[:foundIndex], r.ingredients[foundIndex+1:]...)
	return nil
}

func (r *Recipe) UpdateIngredientQuantity(materialID string, newQuantity float64) error {
	if newQuantity < 0 {
		return &ErrInvalidInput{FieldName: "quantity", Reason: "quantity must be non-negative"}
	}

	for i := range r.ingredients {
		if r.ingredients[i].materialID == materialID {
			r.ingredients[i].quantity = newQuantity
			return nil
		}
	}

	return &ErrNotFound{ResourceName: "ingredient", ResourceID: materialID}
}

func (r *Recipe) CalculateTotalCost(materialCosts map[string]int64) (int64, error) {
	var totalCost int64

	for _, ingredient := range r.ingredients {
		cost, ok := materialCosts[ingredient.MaterialID()]
		if !ok {
			return 0, &ErrInvalidInput{FieldName: "cost", Reason: "no cost found for material"}
		}

		lineCost := int64(float64(cost) * ingredient.Quantity())
		totalCost += lineCost
	}

	return totalCost, nil
}
