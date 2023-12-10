package main

import (
	"encoding/json"
	"fmt"
	"math"
	"os"

	"github.com/xuri/excelize/v2"
)

func approximatePrice(price float64, factor float64) float64 {
	return math.Ceil(float64(price/factor)) * factor
}

type ConversionMapItem struct {
	Description         []string
	TypeWithShutters    string
	TypeWithoutShutters string
	Group               FixtureGroup
}

type ExcelWriter struct {
	modelFilePath string
	conversionMap map[string]ConversionMapItem
}

func NewExcelWriter(modelFilePath string, conversionMapFilePath string) (*ExcelWriter, error) {

	ew := ExcelWriter{
		modelFilePath: modelFilePath,
	}

	err := ew.fillConversionMap(conversionMapFilePath)

	if err != nil {
		return &ew, err
	}

	return &ew, nil
}

func (w *ExcelWriter) fillConversionMap(conversionMapFilePath string) error {

	data, err := os.ReadFile(conversionMapFilePath)

	if err != nil {
		fmt.Println(err)
		return err
	}

	convMap := make(map[string]ConversionMapItem)

	err = json.Unmarshal(data, &convMap)
	if err != nil {
		fmt.Println(err)
		return err
	}

	//fmt.Println(convMap)

	w.conversionMap = convMap
	return nil
}

func (w *ExcelWriter) InsertProducts(collector *JsonCollector, newFileName string) error {

	file, err := excelize.OpenFile(w.modelFilePath)
	if err != nil {
		return err
	}

	defer func() {
		err := file.Close()
		if err != nil {
			fmt.Println(err)
		}
	}()

	for _, p := range collector.Products {

		prodData, ok := w.conversionMap[p.ProductId]

		if ok {

			headers := FixtureHeadersMap[p.Group]

			// quantity
			file.SetCellInt(SerramentiSheet, cell(MinFixtureRow+p.Position-1, headers.QuantityCol), p.Quantity)

			// width
			file.SetCellInt(SerramentiSheet, cell(MinFixtureRow+p.Position-1, headers.WidthCol), p.Width)

			//height
			file.SetCellInt(SerramentiSheet, cell(MinFixtureRow+p.Position-1, headers.HeightCol), p.Height)

			//type
			var prodType string
			if p.RollerShutterPrice > 0 {
				prodType = prodData.TypeWithShutters

				//shutters type
				file.SetCellStr(SerramentiSheet, cell(MinFixtureRow+p.Position-1, RollerShuttersHeader), "Alluminio coibentato - stecche da 55 mm")

			} else {
				prodType = prodData.TypeWithoutShutters
			}

			file.SetCellStr(SerramentiSheet, cell(MinFixtureRow+p.Position-1, headers.TypeCol), prodType)

			//desc
			if p.Depth > 0 { // it's a casing

				if p.Depth > 110 {
					file.SetCellStr(SerramentiSheet, cell(MinFixtureRow+p.Position-1, headers.DescriptionCol), prodData.Description[0])
				} else {
					file.SetCellStr(SerramentiSheet, cell(MinFixtureRow+p.Position-1, headers.DescriptionCol), prodData.Description[1])
				}

			} else {
				file.SetCellStr(SerramentiSheet, cell(MinFixtureRow+p.Position-1, headers.DescriptionCol), prodData.Description[0])
			}

			//price
			finalPrice := p.TotalPrice + p.RollerShutterPrice
			approxPrice := approximatePrice(float64(finalPrice), 50)
			file.SetCellFloat(SerramentiSheet, cell(MinFixtureRow+p.Position-1, headers.PriceCol), approxPrice, 2, 32)

			fmt.Println(prodData)

		} else {
			fmt.Printf("Unknown ProductID: %s skipped line %d\n", p.ProductId, p.Position)
		}

	}

	err = file.UpdateLinkedValue()
	if err != nil {
		return err
	}

	err = file.SaveAs(newFileName + ".xlsm")
	return err

}

//C:/Users/User/Antenore/AntenoreWebPortal/export/gattogroupsas_TETS_18-ago-2023.xml
