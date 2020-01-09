package contracts

type Card interface {
	Element
	Width() string
	ShowOnIndex() bool
	ShowOnDetail() bool
}
