package contracts

type Metric interface {
	Card
	Name() string
}
