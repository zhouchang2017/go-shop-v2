package vue

type CustomCreateTitle interface {
	CreateTitle() string
}

type CustomUpdateTitle interface {
	UpdateTitle() string
}

type CustomIndexTitle interface {
	IndexTitle() string
}

type CustomDetailTitle interface {
	DetailTitle() string
}

type CustomCreateButtonName interface {
	CreateButtonName() string
}

type Router struct {
	Path      string    `json:"path"`
	Component string    `json:"component"`
	Name      string    `json:"name,omitempty"` // 命名路由
	Children  []*Router `json:"children,omitempty"`
	Meta      MetaItems `json:"meta,omitempty"`
	Hidden    bool      `json:"hidden"`
}

func (m *Router) WithMeta(key string, value interface{}) {
	m.Meta = append(m.Meta, &metaItem{key, value})
}

func (m *Router) AddChild(r *Router) {
	m.Children = append(m.Children, r)
}
