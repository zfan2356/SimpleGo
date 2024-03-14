package factory

type Adidas struct {
}

func (a *Adidas) MakeShoes(size int) IShoes {
	return &AdidasShoes{
		Shoes{
			logo: "adidas",
			size: size,
		},
	}
}

func (a *Adidas) MakeShirt(size string, color string) IShirt {
	return &AdidasShirt{
		Shirt{
			logo:  "adidas",
			size:  size,
			color: color,
		},
	}
}

type AdidasShoes struct {
	Shoes
}

type AdidasShirt struct {
	Shirt
}
