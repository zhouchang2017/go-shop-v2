package policies

func init() {
	register(NewInventoryPolicy)
}

type InventoryPolicy struct {

}

func NewInventoryPolicy() *InventoryPolicy {
	return &InventoryPolicy{}
}


