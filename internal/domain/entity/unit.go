package entity

import "errors"

type Unit string

const (
	Kilograms Unit = "kilo"
	Liters    Unit = "litro"
	Each      Unit = "unidade"
	Boxes     Unit = "caixa"
	Grams     Unit = "gramas"
)

func ValidateUnit(u Unit) error {
	switch u {
	case Kilograms, Liters, Each, Boxes, Grams:
		return nil
	default:
		return errors.New("invalid unit: must be one of 'kilo', 'litro', 'unidade', 'caixas'")
	}

}
