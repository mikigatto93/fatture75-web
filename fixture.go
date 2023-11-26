package main

import (
	"fmt"
	"math"
)

type FixtureGroup string
type FixtureType int

const (

	//groups
	GroupA FixtureGroup = "A"
	GroupB FixtureGroup = "B"
	GroupC FixtureGroup = "C"
	GroupD FixtureGroup = "D"

	//types
	CASSONETTO            FixtureType = iota
	SERRAMENTO            FixtureType = iota
	SERRAMENTO_TAPPARELLA FixtureType = iota
	BLINDATO              FixtureType = iota
	PORTONCINO            FixtureType = iota
	TAPPARELLA            FixtureType = iota
	UNKNOWN               FixtureType = iota
)

func getFixtureTypeFromCategory(category string) FixtureType {
	if category == "Cassonetti" {
		return CASSONETTO
	} else if category == "Serramenti" {
		return SERRAMENTO
	} else if category == "SerrTapp" {
		return SERRAMENTO_TAPPARELLA
	} else if category == "PorteBlindate" {
		return BLINDATO
	} else if category == "Portoncino" {
		return PORTONCINO
	} else {
		return UNKNOWN
	}
}

type option struct {
	//Family string
	Value string
	Name  string
}

type Fixture struct {
	Height      int
	Width       int
	Quantity    int
	Description string
	Category    string
	Price       float32 // total price
	VatCode     vatCode
	Options     []option
	Group       FixtureGroup
	Type        FixtureType
}

func (f Fixture) GetExtensiveDescription() string {
	var desc string
	if f.Description != "" {
		desc = f.Description + ", "
	} else {
		desc = ""
	}

	return fmt.Sprintf("%s dimensioni: %dx%d", desc, f.Width, f.Height)
}

func (f Fixture) GetUnitPrice() float32 {
	//approx at 2 decimals
	unitPrice := f.Price / float32(f.Quantity)
	return float32(math.Round(float64(unitPrice)*100) / 100)
}
