package contracts

type Panel interface {
	PrepareFields(fields ...Field)
	Element
	GetFields() []Field
}

// 实现该接口的字段，会用panel包裹
type AsPanel interface {
	WarpPanel() Panel
}
