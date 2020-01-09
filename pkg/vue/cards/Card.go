package cards

import (
	"go-shop-v2/pkg/vue/element"
)

type Card struct {
	*element.Element
	width string
}

func NewCard() *Card {
	return &Card{
		Element: element.NewElement(),
		width:   "1/3",
	}
}

func (this *Card) ShowOnIndex() bool {
	return true
}

func (this *Card) ShowOnDetail() bool {
	return true
}

func (this *Card) Component() string {
	return "card"
}

func (this *Card) PrefixComponent() bool {
	return false
}

func (this *Card) Width() string {
	return this.width
}

func (this *Card) SetWidth50Percent() {
	this.width = "1/2"
}

func (this *Card) SetWidth33Percent() {
	this.width = "1/3"
}

func (this *Card) SetWidth66Percent() {
	this.width = "2/3"
}

func (this *Card) SetWidth25Percent() {
	this.width = "1/4"
}

func (this *Card) SetWidth75Percent() {
	this.width = "3/4"
}

func (this *Card) SetWidth20Percent() {
	this.width = "1/5"
}

func (this *Card) SetWidth40Percent() {
	this.width = "2/5"
}

func (this *Card) SetWidth60Percent() {
	this.width = "3/5"
}

func (this *Card) SetWidth80Percent() {
	this.width = "4/5"
}

func (this *Card) SetWidth16Percent() {
	this.width = "1/6"
}

func (this *Card) SetWidth83Percent() {
	this.width = "5/6"
}

func (this *Card) SetWidthFull() {
	this.width = "full"
}
