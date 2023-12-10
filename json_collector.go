package main

import (
	"encoding/json"
	"fmt"
	"io"
)

type JsonProductData struct {
	Group              FixtureGroup
	Position           int
	Casing             string
	RollerShutterPrice float32
	Height             int
	Width              int
	Notes              string
	TotalPrice         float32
	ProductId          string
	Quantity           int
	Reference          string
	Depth              int
	Uuid               string
}

type JsonCollector struct {
	Products      []JsonProductData
	ExcelFileName string
}

type rawJsonRequestBody struct {
	Products      map[string]rawJsonProductData `json:"products"`
	ExcelFileName string                        `json:"excel_file_name"`
}

type rawJsonProductData struct {
	Group              FixtureGroup `json:"group"`
	Position           int          `json:"position"`
	Casing             string       `json:"casing"`
	RollerShutterPrice float32      `json:"roller_shutter_price"`

	ProductData struct {
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
		Products:      make([]JsonProductData, 0),
		ExcelFileName: "",
	}
	return &col
}

func (c *JsonCollector) LoadData(dataReader io.Reader) error {
	jsonData := rawJsonRequestBody{
		Products:      make(map[string]rawJsonProductData),
		ExcelFileName: "test",
	}

	decoder := json.NewDecoder(dataReader)
	err := decoder.Decode(&jsonData)
	if err != nil {
		fmt.Println("Error unmarchalling json quote data")
		return err
	}

	c.ExcelFileName = jsonData.ExcelFileName
	c.Products = parseJsonData(jsonData.Products)
	return nil
}

func parseJsonData(data map[string]rawJsonProductData) []JsonProductData {

	prodList := make([]JsonProductData, 0)

	for key, product := range data {
		fmt.Println(key)
		newProdData := JsonProductData{
			Group:              product.Group,
			Position:           product.Position,
			Casing:             product.Casing,
			RollerShutterPrice: product.RollerShutterPrice,
			Height:             product.ProductData.Height,
			Width:              product.ProductData.Width,
			Notes:              product.ProductData.Notes,
			TotalPrice:         product.ProductData.TotalPrice,
			ProductId:          product.ProductData.ProductId,
			Quantity:           product.ProductData.Quantity,
			Reference:          product.ProductData.Reference,
			Depth:              product.ProductData.Depth,
			Uuid:               key,
		}
		prodList = append(prodList, newProdData)
	}

	return prodList

}
