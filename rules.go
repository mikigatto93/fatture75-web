package main

import (
	"encoding/json"
	"os"
	"strconv"
	"strings"
)

// types for handlig file data
type rawCondSubRule struct {
	Sheet string   `json:"sheet"`
	Cell  string   `json:"cell"`
	Rules []string `json:"rule"`
}

type rawSubRule struct {
	Sheet string `json:"sheet"`
	Cell  string `json:"cell"`
	Rule  string `json:"rule"`
}

type rulesData struct {
	CondSubRules []rawCondSubRule `json:"conditional_substitution_rules"`
	SubRules     []rawSubRule     `json:"substitute_rule"`
}

// types for the actual rules
type cellData struct {
	row string
	col string
}

type conditionalSubstitutionRule struct {
	sheet                       string
	conditionalSubstitutionsMap map[string]string
	cellTarget                  cellData
}

type substitutionRule struct {
	sheet        string
	substitution string
	cellTarget   cellData
}

var (
	condSubRulesList      []conditionalSubstitutionRule
	substitutionRulesList []substitutionRule
)

func LoadRules(filePath string) error {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}

	rawRules := rulesData{}

	err = json.Unmarshal(data, &rawRules)
	if err != nil {
		return err
	}

	loadCondSubRules(rawRules.CondSubRules)
	loadSubRules(rawRules.SubRules)
	return nil

	//fmt.Println(condSubRulesList)
	//fmt.Println(substitutionRulesList)

}

func searchSubRule(list []substitutionRule, cellTarget cellData) (substitutionRule, bool) {
	var rule substitutionRule
	ok := false

	for _, r := range list {

		if r.cellTarget.col == cellTarget.col &&
			r.cellTarget.row == "" {
			rule = r
			ok = true
		} else if r.cellTarget.col == cellTarget.col &&
			r.cellTarget.row == cellTarget.row {
			rule = r
			ok = true
		}

	}

	return rule, ok

}

func searchCondSubRule(list []conditionalSubstitutionRule, cellTarget cellData) (conditionalSubstitutionRule, bool) {
	var rule conditionalSubstitutionRule
	ok := false

	for _, r := range list {

		if r.cellTarget.col == cellTarget.col &&
			r.cellTarget.row == "" {
			rule = r
			ok = true
		}

		if r.cellTarget.col == cellTarget.col &&
			r.cellTarget.row == cellTarget.row {
			rule = r
			ok = true
		}

	}

	return rule, ok
}

func ApplyRules(entry string, row int, col string, sheet string) string {
	returnVal := entry

	cellTarget := cellData{col: col, row: strconv.Itoa(row)}

	subRule, ok1 := searchSubRule(substitutionRulesList, cellTarget)

	if ok1 && subRule.sheet == sheet {
		returnVal = subRule.substitution
	}

	//apply condtional substitution rule
	condSubRule, ok2 := searchCondSubRule(condSubRulesList, cellTarget)

	if ok2 && condSubRule.sheet == sheet {
		for cond, sub := range condSubRule.conditionalSubstitutionsMap {
			if entry == cond {
				returnVal = sub
			}
		}
	}

	return returnVal
}

func loadCondSubRules(rules []rawCondSubRule) {
	condSubRulesList = make([]conditionalSubstitutionRule, 0)
	for _, rule := range rules {

		substitutionsMap := make(map[string]string)

		for _, ruleData := range rule.Rules {
			//condition->substitution
			splittedRule := strings.Split(ruleData, "->")
			substitutionsMap[splittedRule[0]] = splittedRule[1]
		}

		//parse rule.Cell data
		res := strings.Split(rule.Cell, "$")

		cellTarget := cellData{
			col: strings.TrimSpace(res[1]),
			row: strings.TrimSpace(res[2]),
		}

		formattedRule := conditionalSubstitutionRule{
			sheet:                       rule.Sheet,
			conditionalSubstitutionsMap: substitutionsMap,
			cellTarget:                  cellTarget,
		}

		condSubRulesList = append(condSubRulesList, formattedRule)

	}
}

func loadSubRules(rules []rawSubRule) {
	substitutionRulesList = make([]substitutionRule, 0)
	for _, rule := range rules {

		//parse rule.Cell data
		res := strings.Split(rule.Cell, "$")

		cellTarget := cellData{
			col: strings.TrimSpace(res[1]),
			row: strings.TrimSpace(res[2]),
		}

		formattedRule := substitutionRule{
			sheet:        rule.Sheet,
			substitution: rule.Rule,
			cellTarget:   cellTarget,
		}

		substitutionRulesList = append(substitutionRulesList, formattedRule)

	}
}
