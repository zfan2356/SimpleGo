package main

import (
	"design_patterns/factory"
	"fmt"
)

func main() {
	AdidasFactory, _ := factory.GetFactory("adidas")
	NikeFactory, _ := factory.GetFactory("nike")

	as1 := AdidasFactory.MakeShoes(10)
	as2 := AdidasFactory.MakeShirt("XXL", "black")

	ns1 := NikeFactory.MakeShoes(20)
	ns2 := NikeFactory.MakeShirt("XXXL", "green")

	fmt.Println(as1, as2, ns1, ns2)

}
