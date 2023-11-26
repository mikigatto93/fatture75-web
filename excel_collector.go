package main

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/xuri/excelize/v2"
)

func cell(row int, col string) string {
	return fmt.Sprintf("%s%d", col, row)
}

func getCellInt(file *excelize.File, sheet string, cellCoords string) (int, error) {
	val, err := file.GetCellValue(sheet, cellCoords)
	if err != nil {
		return 0, err
	}

	//remove the , that gets added at the thousands ex: 1,000
	formattedVal := strings.ReplaceAll(val, ",", "")

	intValue, err := strconv.Atoi(formattedVal)
	if err != nil {
		return 0, err
	}

	return intValue, nil
}

func getCellFloat(file *excelize.File, sheet string, cellCoords string) (float32, error) {
	val, err := file.GetCellValue(sheet, cellCoords)
	if err != nil {
		return 0, err
	}
	//fmt.Println(val)
	floatValue, err := strconv.ParseFloat(val, 32)
	if err != nil {
		return 0, err
	}

	return float32(floatValue), nil
}

type ExcelCollector struct {
	filePath             string
	file                 *excelize.File
	Costumer             CostumerData
	Products             []Fixture
	ProfessionalExpenses ServiceExpense
	OtherExpenses        []ServiceExpense
}

func NewExcelCollector(filePath string) *ExcelCollector {
	col := ExcelCollector{
		filePath:      filePath,
		Products:      make([]Fixture, 0),
		OtherExpenses: make([]ServiceExpense, 0),
	}
	return &col
}

func (c *ExcelCollector) LoadData() error {
	file, err := excelize.OpenFile(c.filePath)
	if err != nil {
		return err
	}

	defer func() {
		// Close the spreadsheet.
		if err := file.Close(); err != nil {
			fmt.Println(err)
		}
	}()

	c.file = file

	c.loadCostumer()
	c.loadProducts()
	c.loadOtherExpenses()
	return nil

}

func (c *ExcelCollector) loadCostumer() {
	name, err1 := c.file.GetCellValue(InizioSheet, ConstumerNameCell)
	address, err2 := c.file.GetCellValue(InizioSheet, ConstumerAddressCell)
	municipality, err3 := c.file.GetCellValue(InizioSheet, ConstumerMunicipalityCell)

	errs := [3]error{err1, err2, err3}
	for i := 0; i < 3; i++ {
		if errs[i] != nil {
			fmt.Println(fmt.Errorf("error in loading constumer data %d: %v", i, errs[i]))
			return
		}
	}

	c.Costumer = NewCostumer(name, address, municipality)

}

func (c *ExcelCollector) loadProducts() {

	for i := MinFixtureRow; i <= MaxFixtureRow; i++ {

		fixtureGroup := c.getRowFixtureGroup(i)

		if fixtureGroup != "" {
			prod, err := c.buildFixture(i, fixtureGroup)

			if err != nil {
				fmt.Println(fmt.Errorf("error in loading the product line %d: %v", i, err))
			} else {
				c.Products = append(c.Products, prod)
			}
		}

	}

}

func (c *ExcelCollector) loadOtherExpenses() {
	c.loadComplementaryWorks("Opere complementari")
	c.loadOptionalServices("Servizi eventuali")
}

func (c *ExcelCollector) loadComplementaryWorks(expenseType string) {
	for i := MinComplementaryWorksRow; i <= MaxComplementaryWorksRow; i++ {
		p, err1 := getCellFloat(c.file,
			Check1Sheet, cell(i, OtherExpenseHeaders.PriceCol))
		if err1 != nil {
			fmt.Println(fmt.Errorf("error in loading the expense line %d: %v", i, err1))
		} else {
			if p > 0 {
				desc, err2 := c.file.GetCellValue(
					Check1Sheet, cell(i, OtherExpenseHeaders.DescriptionCol))

				if err2 != nil {
					fmt.Println(fmt.Errorf("error in loading the expense line %d: %v", i, err2))
				} else {
					ex := ServiceExpense{
						Price: p,

						Description: ApplyRules(desc, i, OtherExpenseHeaders.DescriptionCol, Check1Sheet),

						Type: expenseType,
					}

					if i != ProfessionalExpensesRow {
						c.OtherExpenses = append(c.OtherExpenses, ex)
					} else {
						c.ProfessionalExpenses = ex
					}
				}
			}
		}
	}
}

func (c *ExcelCollector) loadOptionalServices(expenseType string) {
	for i := MinOptionalServicesRow; i <= MaxOptionalServicesRow; i++ {
		p, err1 := getCellFloat(c.file,
			Check1Sheet, cell(i, OtherExpenseHeaders.PriceCol))
		if err1 != nil {
			fmt.Println(fmt.Errorf("error in loading the expense line %d: %v", i, err1))
		} else {
			if p > 0 {
				desc, err2 := c.file.GetCellValue(
					Check1Sheet, cell(i, OtherExpenseHeaders.DescriptionCol))

				if err2 != nil {
					fmt.Println(fmt.Errorf("error in loading the expense line %d: %v", i, err2))
				} else {
					ex := ServiceExpense{
						Price: p,

						Description: ApplyRules(desc, i, OtherExpenseHeaders.DescriptionCol, Check1Sheet),

						Type: expenseType,
					}

					c.OtherExpenses = append(c.OtherExpenses, ex)
				}
			}
		}
	}
}

func (c *ExcelCollector) getRowFixtureGroup(rowIndex int) FixtureGroup {

	for group, headers := range FixtureHeadersMap {

		value, err := c.file.GetCellValue(
			SerramentiSheet, cell(rowIndex, headers.QuantityCol))

		if err == nil && value != "" {
			return group
		}
	}
	return ""
}

func (c *ExcelCollector) buildFixture(rowIndex int, fixtureGroup FixtureGroup) (Fixture, error) {

	headers := FixtureHeadersMap[fixtureGroup]

	h, err1 := getCellInt(c.file,
		SerramentiSheet, cell(rowIndex, headers.HeightCol))

	w, err2 := getCellInt(c.file,
		SerramentiSheet, cell(rowIndex, headers.WidthCol))

	q, err3 := getCellInt(c.file,
		SerramentiSheet, cell(rowIndex, headers.QuantityCol))

	d, err4 := c.file.GetCellValue(
		SerramentiSheet, cell(rowIndex, headers.DescriptionCol))

	cat, err5 := c.file.GetCellValue(
		SerramentiSheet, cell(rowIndex, headers.TypeCol))

	p, err6 := getCellFloat(c.file,
		SerramentiSheet, cell(rowIndex, headers.PriceCol))

	ops, err7 := c.getOptions(headers, rowIndex)

	errs := [7]error{err1, err2, err3, err4, err5, err6, err7}
	for i := 0; i < 7; i++ {
		if errs[i] != nil {
			return Fixture{}, errs[i]
		}
	}

	prod := Fixture{
		Height:   h,
		Width:    w,
		Quantity: q,
		Price:    p,

		Description: ApplyRules(d, rowIndex, headers.DescriptionCol, SerramentiSheet),

		Category: ApplyRules(cat, rowIndex, headers.TypeCol, SerramentiSheet),

		Options: ops,
		Group:   fixtureGroup,
		Type:    getFixtureTypeFromCategory(cat),
	}

	return prod, nil

}

func (c *ExcelCollector) getOptions(headers FixtureHeaders, currentRow int) ([]option, error) {

	options := make([]option, 0)

	if headers.OptionMinCol == "" && headers.OptionsMaxCol == "" {
		return options, nil
	}

	//ignore errors in conversion because the cells are fixed
	colStart, _, _ := excelize.CellNameToCoordinates(cell(1, headers.OptionMinCol))

	colEnd, _, _ := excelize.CellNameToCoordinates(cell(1, headers.OptionsMaxCol))

	for i := colStart; i < colEnd+1; i++ {
		cellName, err := excelize.CoordinatesToCellName(i, 1)

		if err != nil {
			return nil, err
		}

		currentCol := strings.Replace(cellName, "1", "", 1) //remove the row number to get curretn column in string format

		optionVal, err := c.file.GetCellValue(SerramentiSheet, cell(currentRow, currentCol))

		if err != nil {
			return nil, err
		}

		if optionVal != "" || optionVal == "No" {
			optionName, err := c.file.GetCellValue(SerramentiSheet, cell(OptionNameRow, currentCol))

			if err != nil {
				return nil, err
			}

			if optionVal == "Sì" {
				optionVal = ""
			}

			newOption := option{
				Value: optionVal,
				Name:  optionName,
			}

			options = append(options, newOption)
		}
	}

	return options, nil

}
