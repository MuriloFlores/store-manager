package value_objects

type Role string

const (
	Admin       = "admin"
	Manager     = "manager"
	Salesperson = "salesperson"
	Client      = "client"
	StockPerson = "stock_person"
	Cashier     = "cashier"
)

func (r Role) IsValid() bool {
	switch r {
	case Admin, Manager, Salesperson, Client, StockPerson, Cashier:
		return true
	}

	return false
}

func (r Role) ToString() string {
	return string(r)
}

func (r Role) ExpectedValue() string {
	return "admin, manager, salesperson, client, stock_person, cashier"
}

func (r Role) IsStockEmployee() bool {
	switch r {
	case Admin, Manager, StockPerson:
		return true
	default:
		return false
	}
}
