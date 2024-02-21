package models

type InventoryItem struct{
	//omitempty is used to exlcuded those values while serializing the data to json or bson.
	ID string `json:"id,omitempty" bson:"_id,omitempty"`
	UserID string `json:"userID,omitempty" bson:"userID,omitempty"`
	ProductName string `json:"productName" bson:"productName"`
	Units int `json:"units" bson:"units"`
	Price float64 `json:"price" bson:"price"`
}