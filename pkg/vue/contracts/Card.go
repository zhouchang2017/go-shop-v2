package contracts

type Card interface {
	Element
	Width() string
	ShowOnIndex() bool
	ShowOnDetail() bool
	Grid() bool
}

type MoreLink interface {
	Link() VueRouterOption
}