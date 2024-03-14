package factory

type IShoes interface {
	SetLogo(logo string)
	SetSize(size int)
	GetLogo() string
	GetSize() int
}

type Shoes struct {
	logo string
	size int
}

func (s *Shoes) SetLogo(logo string) {
	s.logo = logo
}
func (s *Shoes) GetLogo() string {
	return s.logo
}
func (s *Shoes) SetSize(size int) {
	s.size = size
}
func (s *Shoes) GetSize() int {
	return s.size
}

type IShirt interface {
	SetLogo(logo string)
	SetSize(size string)
	SetColor(color string)
	GetLogo() string
	GetSize() string
	GetColor() string
}
type Shirt struct {
	logo  string
	size  string
	color string
}

func (s *Shirt) SetLogo(logo string) {
	s.logo = logo
}
func (s *Shirt) GetLogo() string {
	return s.logo
}
func (s *Shirt) SetSize(size string) {
	s.size = size
}
func (s *Shirt) GetSize() string {
	return s.size
}
func (s *Shirt) SetColor(color string) {
	s.color = color
}
func (s *Shirt) GetColor() string {
	return s.color
}
