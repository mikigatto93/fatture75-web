package main

import (
	"encoding/json"
	"fmt"
	"math"
	"os"
	"strconv"

	"github.com/xuri/excelize/v2"
)

func approximatePrice(price float64, factor float64) float64 {
	return math.Ceil(float64(price/factor)) * factor
}

type ConversionMapItem struct {
	Description []string
	Type        string
	Group       FixtureGroup
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

func (w *ExcelWriter) InsertProducts(collector *XmlCollector, newFileName string) error {

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

	rowCounter := 0
	for _, p := range collector.ProductData {

		prodData, ok := w.conversionMap[p.ProductId]
		if ok {
			headers := FixtureHeadersMap[prodData.Group]

			w, err := strconv.Atoi(p.Width)
			if err != nil {
				fmt.Println(err)
				continue
			}

			h, err := strconv.Atoi(p.Height)
			if err != nil {
				fmt.Println(err)
				continue
			}

			// quantity
			file.SetCellInt(SerramentiSheet, cell(MinFixtureRow+rowCounter, headers.QuantityCol), p.Quantity)

			// width
			file.SetCellInt(SerramentiSheet, cell(MinFixtureRow+rowCounter, headers.WidthCol), w)

			//height
			file.SetCellInt(SerramentiSheet, cell(MinFixtureRow+rowCounter, headers.HeightCol), h)

			//type
			file.SetCellStr(SerramentiSheet, cell(MinFixtureRow+rowCounter, headers.TypeCol), prodData.Type)

			//desc
			file.SetCellStr(SerramentiSheet, cell(MinFixtureRow+rowCounter, headers.DescriptionCol), prodData.Description[0])

			//price
			approxPrice := approximatePrice(float64(p.Price), 50)
			file.SetCellFloat(SerramentiSheet, cell(MinFixtureRow+rowCounter, headers.PriceCol), approxPrice, 2, 32)

			rowCounter++

		} else {
			fmt.Printf("Unknown ProductID: %s skipped line %d\n", p.ProductId, rowCounter)

			rowCounter++
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
