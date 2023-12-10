package main

import (
	"bytes"
	"encoding/xml"
	"strconv"
	"strings"
)

type rawXmlDescriptionItem struct {
	Value string `xml:"Value,attr"`
	Name  string `xml:"Name,attr"`
}

type rawXmlProductData struct {
	ProductId           string
	UnitListPrice       float32
	QuantityListPrice   float32
	Quantity            int
	Reference           string
	Position            int
	DescriptionExtended []rawXmlDescriptionItem `xml:"DescriptionExtended>DescriptionItem"`
}

type rawXmlDocumentData struct {
	XMLName xml.Name             `xml:"Order"`
	Header  rawXmlDocumentHeader `xml:"Header"`
	Rows    []rawXmlProductData  `xml:"Rows>Row"`
}

type rawXmlDocumentHeader struct {
	Title      string
	Customer   string
	Commission string
}

type XmlProductData struct {
	ProductId string
	Height    int
	Width     int
	Quantity  int
	Price     float32
	Notes     string
	Reference string
	Position  int
	Depth     int //only for casing
}

type XmlHeaderData struct {
	Title      string
	Customer   string
	Commission string
}

type XmlCollector struct {
	fileBuffer  bytes.Buffer
	ProductData []XmlProductData
	HeaderData  XmlHeaderData
}

func NewXmlCollector(fileBuf bytes.Buffer) *XmlCollector {
	return &XmlCollector{
		fileBuffer:  fileBuf,
		ProductData: make([]XmlProductData, 0),
	}
}

func (c *XmlCollector) LoadData() error {
	xmlData := rawXmlDocumentData{}

	err := xml.Unmarshal(c.fileBuffer.Bytes(), &xmlData)

	if err != nil {
		return err
	}

	c.HeaderData = XmlHeaderData{
		Title:      xmlData.Header.Title,
		Customer:   xmlData.Header.Customer,
		Commission: xmlData.Header.Commission,
	}

	c.ProductData = parseXmlData(xmlData)
	return nil
}

func parseXmlData(data rawXmlDocumentData) []XmlProductData {

	prodList := make([]XmlProductData, 0)

	for _, p := range data.Rows {

		//set custom default values
		width := -1
		height := -1
		depth := -1 //only for casing
		var notes string

		for _, val := range p.DescriptionExtended {

			if val.Name == "Altezza" && height < 0 {
				//get only the first occurrence
				parsedVal := strings.Replace(val.Value, "mm", "", 1)
				height, _ = strconv.Atoi(parsedVal)
			} else if val.Name == "Larghezza" && width < 0 {
				//get only the first occurrence
				parsedVal := strings.Replace(val.Value, "mm", "", 1)
				width, _ = strconv.Atoi(parsedVal)
			} else if val.Name == "Note" {
				notes = val.Value
			} else if val.Name == "Profondità" {
				parsedVal := strings.Replace(val.Value, "mm", "", 1)
				depth, _ = strconv.Atoi(parsedVal)
			}
		}

		if strings.ToUpper(notes) == "NO 75" ||
			p.ProductId == "WND.CEL.10" {
			continue
		} else {

			prodData := XmlProductData{
				ProductId: p.ProductId,
				Height:    height,
				Width:     width,
				Price:     p.QuantityListPrice, // total price
				Quantity:  p.Quantity,
				Notes:     notes,
				Reference: p.Reference,
				Depth:     depth,
				Position:  p.Position,
			}

			prodList = append(prodList, prodData)
		}
	}

	return prodList
}
