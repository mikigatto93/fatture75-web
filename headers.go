package main

type FixtureHeaders struct {
	WidthCol        string
	HeightCol       string
	TypeCol         string
	DescriptionCol  string
	QuantityCol     string
	PriceCol        string
	OptionFamilyRow string
	OptionNameRow   string
	OptionMinCol    string
	OptionsMaxCol   string
}

type ExpenseHeaders struct {
	DescriptionCol string
	PriceCol       string
}

var OtherExpenseHeaders ExpenseHeaders = ExpenseHeaders{
	DescriptionCol: "B",
	PriceCol:       "G",
}

var FixtureHeadersMap map[FixtureGroup]FixtureHeaders = map[FixtureGroup]FixtureHeaders{
	GroupA: {
		WidthCol:       "D",
		HeightCol:      "F",
		TypeCol:        "I",
		DescriptionCol: "J",
		QuantityCol:    "C",
		PriceCol:       "AF",
		OptionMinCol:   "K",
		OptionsMaxCol:  "AE",
	},

	GroupB: {
		WidthCol:       "AI",
		HeightCol:      "AK",
		TypeCol:        "AN",
		DescriptionCol: "AO",
		QuantityCol:    "AH",
		PriceCol:       "CP",
		OptionMinCol:   "AP",
		OptionsMaxCol:  "CO",
	},

	GroupC: {
		WidthCol:       "CS",
		HeightCol:      "CU",
		TypeCol:        "CX",
		DescriptionCol: "CY",
		QuantityCol:    "CR",
		PriceCol:       "ED",
		OptionMinCol:   "CZ",
		OptionsMaxCol:  "EC",
	},

	GroupD: {
		WidthCol:       "EG",
		HeightCol:      "EI",
		TypeCol:        "EL",
		DescriptionCol: "EM",
		QuantityCol:    "EF",
		PriceCol:       "EN",
		OptionMinCol:   "",
		OptionsMaxCol:  "",
	},
}

const (
	MinFixtureRow int = 7
	MaxFixtureRow int = 56

	MinComplementaryWorksRow int = 22
	MaxComplementaryWorksRow int = 35

	MinOptionalServicesRow int = 39
	MaxOptionalServicesRow int = 43

	ProfessionalExpensesRow int = 35

	OptionFamilyRow int = 6
	OptionNameRow   int = 5

	RollerShuttersHeader string = "BK"
)
