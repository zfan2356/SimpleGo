package factory

import "fmt"

type IFactory interface {
	MakeShoes(size int) IShoes
	MakeShirt(size, color string) IShirt
}

func GetFactory(brand string) (IFactory, error) {
	if brand == "adidas" {
		return &Adidas{}, nil
	}
	if brand == "nike" {
		return &Nike{}, nil
	}
	return nil, fmt.Errorf("Wrong brand type.")
}
