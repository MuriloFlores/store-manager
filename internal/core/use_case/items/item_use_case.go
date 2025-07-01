package items

type ItemsUseCases struct {
	Create
	Find
	Update
	Delete
	List
}

func NewItemUseCases() *ItemsUseCases {
	return &ItemsUseCases{}
}
