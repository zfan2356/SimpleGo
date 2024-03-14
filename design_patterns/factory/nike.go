package factory

type Nike struct {
}

func (n *Nike) MakeShoes(size int) IShoes {
	return &NikeShoes{Shoes{
		logo: "nike",
		size: size,
	}}
}

func (n *Nike) MakeShirt(size string, color string) IShirt {
	return &NikeShirt{Shirt{
		logo:  "nike",
		size:  size,
		color: color,
	}}
}

type NikeShoes struct {
	Shoes
}

type NikeShirt struct {
	Shirt
}
