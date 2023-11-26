package main

import "strings"

const (
	ConstumerNameCell         string = "E7"
	ConstumerAddressCell      string = "E9"
	ConstumerMunicipalityCell string = "E11"
)

type CostumerData struct {
	Name         string
	Address      string
	Municipality string
	Discrict     string
}

func NewCostumer(name string, address string, municipality string) CostumerData {

	parsedMunicipality, district := parseMunicipality(municipality)

	c := CostumerData{
		Name:         name,
		Address:      address,
		Municipality: parsedMunicipality,
		Discrict:     district,
	}

	return c

}

func parseMunicipality(municipality string) (string, string) {
	if strings.Contains(municipality, "(") {
		res := strings.Split(municipality, "(")
		return res[0], strings.Replace(res[1], ")", "", 1)
	} else {
		return municipality, ""
	}
}
