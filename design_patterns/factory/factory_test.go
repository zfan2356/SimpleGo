package factory

import (
	"fmt"
	"testing"
)

func Test01(t *testing.T) {
	AdidasFactory, _ := GetFactory("adidas")
	NikeFactory, _ := GetFactory("nike")

	as1 := AdidasFactory.MakeShoes(10)
	as2 := AdidasFactory.MakeShirt("XXL", "black")

	ns1 := NikeFactory.MakeShoes(20)
	ns2 := NikeFactory.MakeShirt("XXXL", "green")

	fmt.Println(as1, as2, ns1, ns2)
}
