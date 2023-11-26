package main

import (
	"encoding/json"
	"fmt"
)

type JsonProductData struct {
	Group            FixtureGroup
	Position         int
	Casing           string
	HasRollerShutter bool
	Height           int
	Width            int
	Notes            string
	TotalPrice       float32
	ProductId        string
	Quantity         int
	Reference        string
	Depth            int
	Uuid             string
}

type JsonCollector struct {
	Products []JsonProductData
}

type rawJsonQuoteData map[string]rawJsonProductData

type rawJsonProductData struct {
	Group            FixtureGroup `json:"group"`
	Position         int          `json:"position"`
	Casing           string       `json:"casing"`
	HasRollerShutter bool         `json:"has_roller_shutter"`
	ProductData      struct {
		Height     int     `json:"height"`
		Width      int     `json:"width"`
		Notes      string  `json:"notes"`
		TotalPrice float32 `json:"tot_price"`
		ProductId  string  `json:"product_id"`
		Quantity   int     `json:"quantity"`
		Reference  string  `json:"reference"`
		Depth      int     `json:"depth"`
	} `json:"product_data"`
}

func NewJsonCollector() *JsonCollector {
	col := JsonCollector{
		Products: make([]JsonProductData, 0),
	}
	return &col
}

func (c *JsonCollector) LoadData(data []byte) error {
	jsonData := make(rawJsonQuoteData)

	err := json.Unmarshal(data, &jsonData)
	if err != nil {
		fmt.Println("Error unmarchalling json quote data")
		return err
	}

	c.Products = parseJsonData(jsonData)
	return nil
}

func parseJsonData(data rawJsonQuoteData) []JsonProductData {

	prodList := make([]JsonProductData, 0)

	for key, product := range data {
		newProdData := JsonProductData{
			Group:            product.Group,
			Position:         product.Position,
			Casing:           product.Casing,
			HasRollerShutter: product.HasRollerShutter,
			Height:           product.ProductData.Height,
			Width:            product.ProductData.Width,
			Notes:            product.ProductData.Notes,
			TotalPrice:       product.ProductData.TotalPrice,
			ProductId:        product.ProductData.ProductId,
			Quantity:         product.ProductData.Quantity,
			Reference:        product.ProductData.Reference,
			Depth:            product.ProductData.Depth,
			Uuid:             key,
		}
		prodList = append(prodList, newProdData)
	}

	return prodList

}
