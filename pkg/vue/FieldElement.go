package vue

type FieldElement struct {
	Element
	Panel          string `json:"panel"`
	ShowOnIndex    bool   `json:"-"`
	ShowOnDetail   bool   `json:"-"`
	ShowOnCreation bool   `json:"-"`
	ShowOnUpdate   bool   `json:"-"`
}

func (f *FieldElement) HideFromIndex(cb func() bool) {
	f.ShowOnIndex = !cb()
}

func (f *FieldElement) HideFromDetail(cb func() bool) {
	f.ShowOnDetail = !cb()
}

func (f *FieldElement) OnlyOnIndex() {
	f.ShowOnIndex = true
	f.ShowOnCreation = false
	f.ShowOnDetail = false
	f.ShowOnUpdate = false
}

func (f *FieldElement) OnlyOnDetail() {
	f.ShowOnIndex = false
	f.ShowOnCreation = false
	f.ShowOnDetail = true
	f.ShowOnUpdate = false
}

func (f *FieldElement) OnlyOnForm() {
	f.ShowOnIndex = false
	f.ShowOnCreation = true
	f.ShowOnDetail = false
	f.ShowOnUpdate = true
}

func (f *FieldElement) ExceptOnForms() {
	f.ShowOnIndex = true
	f.ShowOnCreation = false
	f.ShowOnDetail = true
	f.ShowOnUpdate = false
}
