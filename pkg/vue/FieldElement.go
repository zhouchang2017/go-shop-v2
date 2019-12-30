package vue

type FieldElement struct {
	BasicElement
	Panel          string `json:"panel"`
	Readonly       bool   `json:"readonly"`
	showOnIndex    bool   `json:"-"`
	showOnDetail   bool   `json:"-"`
	showOnCreation bool   `json:"-"`
	showOnUpdate   bool   `json:"-"`
}

func (f FieldElement) GetPanel() string {
	return f.Panel
}

func (f *FieldElement) SetPanel(name string) {
	f.Panel = name
}

func (f FieldElement) ShowOnIndex() bool {
	return f.showOnIndex
}

func (f FieldElement) ShowOnDetail() bool {
	return f.showOnDetail
}

func (f FieldElement) ShowOnCreation() bool {
	return f.showOnCreation
}

func (f FieldElement) ShowOnUpdate() bool {
	return f.showOnUpdate
}

func (f *FieldElement) SetShowOnIndex(show bool) {
	f.showOnIndex = show
}

func (f *FieldElement) SetShowOnDetail(show bool) {
	f.showOnDetail = show
}

func (f *FieldElement) SetShowOnUpdate(show bool) {
	f.showOnUpdate = show
}

func (f *FieldElement) SetShowOnCreation(show bool) {
	f.showOnCreation = show
}

func (f *FieldElement) HideFromIndex(cb func() bool) {
	f.showOnIndex = !cb()
}

func (f *FieldElement) HideFromDetail(cb func() bool) {
	f.showOnDetail = !cb()
}

func (f *FieldElement) OnlyOnIndex() {
	f.showOnIndex = true
	f.showOnCreation = false
	f.showOnDetail = false
	f.showOnUpdate = false
}

func (f *FieldElement) OnlyOnDetail() {
	f.showOnIndex = false
	f.showOnCreation = false
	f.showOnDetail = true
	f.showOnUpdate = false
}

func (f *FieldElement) OnlyOnForm() {
	f.showOnIndex = false
	f.showOnCreation = true
	f.showOnDetail = false
	f.showOnUpdate = true
}

func (f *FieldElement) ExceptOnForms() {
	f.showOnIndex = true
	f.showOnCreation = false
	f.showOnDetail = true
	f.showOnUpdate = false
}

func SetShowOnIndex(show bool) FieldOption {
	return func(field interface{}) {
		basicField, err := resolveBasicField(field)
		if err == nil {
			basicField.showOnIndex = show
		}
	}
}

func SetShowOnDetail(show bool) FieldOption {
	return func(field interface{}) {
		basicField, err := resolveBasicField(field)
		if err == nil {
			basicField.showOnDetail = show
		}
	}
}

func SetShowOnUpdate(show bool) FieldOption {
	return func(field interface{}) {
		basicField, err := resolveBasicField(field)
		if err == nil {
			basicField.showOnUpdate = show
		}
	}
}

func SetShowOnCreation(show bool) FieldOption {
	return func(field interface{}) {
		basicField, err := resolveBasicField(field)
		if err == nil {
			basicField.showOnCreation = show
		}
	}
}

func OnlyOnIndex() FieldOption {
	return func(field interface{}) {
		basicField, err := resolveBasicField(field)
		if err == nil {
			basicField.OnlyOnIndex()
		}
	}
}

func OnlyOnDetail() FieldOption {
	return func(field interface{}) {
		basicField, err := resolveBasicField(field)
		if err == nil {
			basicField.OnlyOnDetail()
		}
	}
}

func OnlyOnForm() FieldOption {
	return func(field interface{}) {
		basicField, err := resolveBasicField(field)
		if err == nil {
			basicField.OnlyOnForm()
		}
	}
}

func ExceptOnForms() FieldOption {
	return func(field interface{}) {
		basicField, err := resolveBasicField(field)
		if err == nil {
			basicField.ExceptOnForms()
		}
	}
}
