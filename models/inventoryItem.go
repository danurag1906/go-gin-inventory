package models

type InventoryItem struct{
	ID string `json:"id,omitempty" bson:"_id,omitempty"`
	ProductName string `json:"productName" bson:"productName"`
	Units int `json:"units" bson:"units"`
	Price float64 `json:"price" bson:"price"`
}