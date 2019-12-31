package contracts

type Router interface {
	WithMeta(key string, value interface{})
	AddChild(r Router)
	Component() string
	RouterName() string
	Path() string
}
