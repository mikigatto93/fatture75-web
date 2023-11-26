package main

import (
	"math"

	fattureincloud "github.com/fattureincloud/fattureincloud-go-sdk/v2/model"
)

type vatCode int

const (
	VatCode10     vatCode = 3
	VatCode22     vatCode = 0
	VatCode4for75 vatCode = 7981652 // specific for 75% -> Dpr n. 633/1972 - Tabella A, Parte II, punto 41-ter
)

// not used for now
var vatCodeGroupMap = map[FixtureGroup]vatCode{
	"A": VatCode22,
	"B": VatCode22,
	"C": VatCode10,
	"D": VatCode10,
}

func getVatFromVatCode(code vatCode) float32 {
	if code == VatCode10 {
		return 10
	} else if code == VatCode22 {
		return 22
	} else if code == VatCode4for75 {
		return 4
	}

	return 0

}

type Document struct {
	Client             *fattureincloud.Entity
	IssuedDocument     fattureincloud.IssuedDocument
	Fixtures           map[FixtureGroup][]Fixture
	Expenses           []ServiceExpense
	TotalBeforeTaxes22 float32
	TotalBeforeTaxes10 float32
	TotalBeforeTaxes4  float32
	Total              float32
	IsBonus75          bool
}

func NewDocument(docType fattureincloud.IssuedDocumentType, collector *ExcelCollector, isBonus75 bool) *Document {

	doc := Document{
		IssuedDocument: *fattureincloud.NewIssuedDocument(),
		IsBonus75:      isBonus75,
		Fixtures:       make(map[FixtureGroup][]Fixture),
	}

	client := *fattureincloud.NewEntity().
		SetName(collector.Costumer.Name).
		SetAddressStreet(collector.Costumer.Address).
		SetAddressCity(collector.Costumer.Municipality).
		SetAddressProvince(collector.Costumer.Discrict).
		SetCountry("Italia")

	doc.IssuedDocument.
		SetEntity(client).
		SetType(docType).
		SetCurrency(*fattureincloud.NewCurrency().SetId("EUR")).
		SetLanguage(*fattureincloud.NewLanguage().SetCode("it").SetName("italiano"))

	doc.fillFixtureList(collector)
	doc.fillExpenseList(collector)

	return &doc
}

func (d *Document) fillFixtureList(collector *ExcelCollector) {

	for _, fixture := range collector.Products {
		if !d.IsBonus75 {
			fixture.VatCode = vatCodeGroupMap[fixture.Group]

			if fixture.VatCode == VatCode10 {
				d.TotalBeforeTaxes10 += fixture.Price
			} else if fixture.VatCode == VatCode22 {
				d.TotalBeforeTaxes22 += fixture.Price
			}

		} else {
			fixture.VatCode = VatCode4for75
			d.TotalBeforeTaxes4 += fixture.Price
		}

		//excorporate fixture with type SERRAMENTO_TAPPARELLA
		if fixture.Type == SERRAMENTO_TAPPARELLA {
			processFixtureShutters(fixture)
		}

		d.Fixtures[fixture.Group] = append(d.Fixtures[fixture.Group], fixture)

	}
}

func processFixtureShutters(fix Fixture) {
	return
}

func (d *Document) fillExpenseList(collector *ExcelCollector) {

	for _, expense := range collector.OtherExpenses {
		if !d.IsBonus75 {
			expense.VatCode = VatCode10
			d.TotalBeforeTaxes10 += expense.Price
		} else {
			expense.VatCode = VatCode4for75
			d.TotalBeforeTaxes4 += expense.Price
		}
		d.Expenses = append(d.Expenses, expense)
	}

	// special professional expenses must be taxed with vat at 22
	// it does not concur to the TotalBeforeTaxes22 as it must be taxed at 22
	// without any type of deduction
	profExpenses := collector.ProfessionalExpenses
	profExpenses.VatCode = VatCode22
	d.Expenses = append(d.Expenses, profExpenses)
}

func (d *Document) FillItems() {

	if d.TotalBeforeTaxes10 != 0 || d.TotalBeforeTaxes22 != 0 {
		d.applyVatCheck()
	}

	d.calculateTotal()

	itemsList := []fattureincloud.IssuedDocumentItemsListItem{}

	groups := []FixtureGroup{GroupA, GroupB, GroupC, GroupD}
	//fill fixture
	for _, group := range groups {
		for _, fixture := range d.Fixtures[group] {

			newItem := *fattureincloud.NewIssuedDocumentItemsListItem().
				SetName(fixture.Category).
				SetDescription(fixture.GetExtensiveDescription()).
				SetNetPrice(fixture.GetUnitPrice()).
				SetDiscount(0).
				SetQty(float32(fixture.Quantity)).
				SetVat(*fattureincloud.NewVatType().SetId(int32(fixture.VatCode)))
			itemsList = append(itemsList, newItem)
		}
	}
	//fill expenses
	for _, expense := range d.Expenses {
		newItem := *fattureincloud.NewIssuedDocumentItemsListItem().
			SetName(expense.Type).
			SetDescription(expense.Description).
			SetNetPrice(expense.Price).
			SetDiscount(0).
			SetQty(1).
			SetVat(*fattureincloud.NewVatType().SetId(int32(expense.VatCode)))

		itemsList = append(itemsList, newItem)
	}

	d.IssuedDocument.SetItemsList(itemsList)

	if d.IsBonus75 {
		d.ApplyDiscount(75)
	}
}

func (d *Document) applyVatCheck() {
	// see https://www.commercialistatelematico.com/articoli/2023/03/beni-significativi-problematiche-iva.html

	var amountToApply float32

	if d.TotalBeforeTaxes22-d.TotalBeforeTaxes10 > 0 {
		amountToApply = d.TotalBeforeTaxes10
	} else {
		amountToApply = d.TotalBeforeTaxes22
	}

	deduction := ServiceExpense{
		Description: "Detrazione per diversa imputazione IVA beni significativi",
		Type:        "Detrazione",
		Price:       -amountToApply,
		VatCode:     VatCode22,
	}

	addition := ServiceExpense{
		Description: "Riaddebito per diversa imputazione IVA agevolata beni significativi",
		Type:        "Riaddebito",
		Price:       amountToApply,
		VatCode:     VatCode10,
	}

	d.Expenses = append(d.Expenses, deduction, addition)
}

func (d *Document) calculateTotal() {

	groups := []FixtureGroup{GroupA, GroupB, GroupC, GroupD}
	for _, group := range groups {
		for _, fix := range d.Fixtures[group] {
			d.Total +=
				fix.Price + (fix.Price * getVatFromVatCode(fix.VatCode) / 100)
		}
	}
	for _, ex := range d.Expenses {
		d.Total +=
			ex.Price + (ex.Price * getVatFromVatCode(ex.VatCode) / 100)
	}

}

func (d *Document) ApplyDiscount(amount float32) {
	discount := -amount * d.Total / 100
	roundedDiscount := math.Round(float64(discount*100)) / 100
	d.IssuedDocument.SetAmountDueDiscount(float32(roundedDiscount))
}
