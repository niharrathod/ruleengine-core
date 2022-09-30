package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"time"

	ruleenginecore "github.com/niharrathod/ruleengine-core"
)

var intStep = int64(1000)
var floatStep = float64(1000.1358)
var fieldIntValueMap = map[string]int64{}
var fieldFloatValueMap = map[string]float64{}

func main() {

	intCount, floatCount, stringCount, boolCount := 5, 5, 5, 5
	ruleCount := 1000
	outFile := "ruleEngineConfig_" + strconv.Itoa(ruleCount) + ".json"

	fields := fieldGenerate(intCount, floatCount, stringCount, boolCount)
	rules, condTypes := generateRules(ruleCount, fields)
	result := &ruleenginecore.RuleEngineConfig{
		Fields:         fields,
		ConditionTypes: condTypes,
		Rules:          rules,
	}

	writeJsonToFile(result, outFile)
	writeJsonToFile(generateValidInput(fields, fieldIntValueMap, fieldFloatValueMap), "ValidInput.json")
	writeJsonToFile(generateInvalidInput(fields, fieldIntValueMap, fieldFloatValueMap), "InvalidInput.json")
}

func generateRules(ruleCount int, fields ruleenginecore.Fields) (map[string]*ruleenginecore.Rule, map[string]*ruleenginecore.ConditionType) {
	ruleMap := map[string]*ruleenginecore.Rule{}
	conditionTypeMap := map[string]*ruleenginecore.ConditionType{}

	for i := 0; i < ruleCount; i++ {
		subConditions := []*ruleenginecore.Condition{}

		for fieldname, fieldtype := range fields {

			switch fieldtype {
			case ruleenginecore.IntType:
				min := fieldIntValueMap[fieldname]
				max := min + intStep

				// prepare two conditions :  GreaterThanEqual to for min and LessThanEqual to for max value
				condMap := generateIntConditions(fieldname, min, max)

				for condName, cond := range condMap {
					// adding to customConditionType
					conditionTypeMap[condName] = cond
					// adding to currentRuleSubConditions
					ruleSubCond := ruleenginecore.Condition{
						ConditionType: condName,
					}
					subConditions = append(subConditions, &ruleSubCond)

				}
				fieldIntValueMap[fieldname] = max

			case ruleenginecore.FloatType:
				min := fieldFloatValueMap[fieldname]
				max := min + floatStep

				// prepare two conditions :  GreaterThanEqual to for min and LessThanEqual to for max value
				condMap := generateFloatConditions(fieldname, min, max)

				for condName, cond := range condMap {
					// adding to customConditionType
					conditionTypeMap[condName] = cond
					// adding to currentRuleSubConditions
					ruleSubCond := ruleenginecore.Condition{
						ConditionType: condName,
					}
					subConditions = append(subConditions, &ruleSubCond)

				}
				fieldFloatValueMap[fieldname] = max
			case ruleenginecore.StringType:
				condMap := generateStringConditions(fieldname, conditionTypeMap)

				for condName, cond := range condMap {
					// adding to customConditionType
					conditionTypeMap[condName] = cond
					// adding to currentRuleSubConditions
					ruleSubCond := ruleenginecore.Condition{
						ConditionType: condName,
					}
					subConditions = append(subConditions, &ruleSubCond)

				}

			case ruleenginecore.BoolType:
				condMap := generateBoolConditions(fieldname, conditionTypeMap)

				for condName, cond := range condMap {
					// adding to customConditionType
					conditionTypeMap[condName] = cond
					// adding to currentRuleSubConditions
					ruleSubCond := ruleenginecore.Condition{
						ConditionType: condName,
					}
					subConditions = append(subConditions, &ruleSubCond)

				}
			}
		}

		rulename := "rule_" + strconv.Itoa(i)
		rootCondition := ruleenginecore.Condition{
			ConditionType: ruleenginecore.AndOperator,
			SubConditions: subConditions,
		}

		rule := ruleenginecore.Rule{
			Priority:      i,
			RootCondition: &rootCondition,
			Result: map[string]any{
				rulename: "yes",
			},
		}

		ruleMap[rulename] = &rule
	}

	return ruleMap, conditionTypeMap
}

func generateIntConditions(fieldname string, min, max int64) map[string]*ruleenginecore.ConditionType {
	result := map[string]*ruleenginecore.ConditionType{}
	minStr := strconv.FormatInt(min, 10)
	maxStr := strconv.FormatInt(max, 10)
	gteConditionName := fieldname + "_" + ruleenginecore.GreaterEqualOperator + "_" + minStr
	result[gteConditionName] = &ruleenginecore.ConditionType{
		Operator:    ruleenginecore.GreaterEqualOperator,
		OperandType: ruleenginecore.IntType,
		Operands: []*ruleenginecore.Operand{
			{
				OperandAs: ruleenginecore.OperandAsField,
				Val:       fieldname,
			},
			{
				OperandAs: ruleenginecore.OperandAsConstant,
				Val:       minStr,
			},
		},
	}

	lteConditionName := fieldname + "_" + ruleenginecore.LessEqualOperator + "_" + maxStr
	result[lteConditionName] = &ruleenginecore.ConditionType{
		Operator:    ruleenginecore.LessEqualOperator,
		OperandType: ruleenginecore.IntType,
		Operands: []*ruleenginecore.Operand{
			{
				OperandAs: ruleenginecore.OperandAsField,
				Val:       fieldname,
			},
			{
				OperandAs: ruleenginecore.OperandAsConstant,
				Val:       maxStr,
			},
		},
	}
	return result
}

func generateFloatConditions(fieldname string, min, max float64) map[string]*ruleenginecore.ConditionType {
	result := map[string]*ruleenginecore.ConditionType{}
	minStr := strconv.FormatFloat(min, 'E', -1, 64)
	maxStr := strconv.FormatFloat(max, 'E', -1, 64)
	gteConditionName := fieldname + "_" + ruleenginecore.GreaterEqualOperator + "_" + minStr
	result[gteConditionName] = &ruleenginecore.ConditionType{
		Operator:    ruleenginecore.GreaterEqualOperator,
		OperandType: ruleenginecore.FloatType,
		Operands: []*ruleenginecore.Operand{
			{
				OperandAs: ruleenginecore.OperandAsField,
				Val:       fieldname,
			},
			{
				OperandAs: ruleenginecore.OperandAsConstant,
				Val:       minStr,
			},
		},
	}

	lteConditionName := fieldname + "_" + ruleenginecore.LessEqualOperator + "_" + maxStr
	result[lteConditionName] = &ruleenginecore.ConditionType{
		Operator:    ruleenginecore.LessEqualOperator,
		OperandType: ruleenginecore.FloatType,
		Operands: []*ruleenginecore.Operand{
			{
				OperandAs: ruleenginecore.OperandAsField,
				Val:       fieldname,
			},
			{
				OperandAs: ruleenginecore.OperandAsConstant,
				Val:       maxStr,
			},
		},
	}
	return result
}

func generateStringConditions(fieldname string, existingCondition map[string]*ruleenginecore.ConditionType) map[string]*ruleenginecore.ConditionType {
	result := map[string]*ruleenginecore.ConditionType{}
	if RandomBool() {
		containConditionName := fieldname + "_" + ruleenginecore.ContainOperator + "_" + "apple"

		if val, ok := existingCondition[containConditionName]; ok {
			result[containConditionName] = val
			return result
		}

		result[containConditionName] = &ruleenginecore.ConditionType{
			Operator:    ruleenginecore.ContainOperator,
			OperandType: ruleenginecore.StringType,
			Operands: []*ruleenginecore.Operand{
				{
					OperandAs: ruleenginecore.OperandAsConstant,
					Val:       "This should be big story or comma separated values or any kind of string which ends with apple",
				},
				{
					OperandAs: ruleenginecore.OperandAsField,
					Val:       fieldname,
				},
			},
		}
	} else {
		containConditionName := fieldname + "_" + ruleenginecore.EqualOperator + "_" + "apple"

		if val, ok := existingCondition[containConditionName]; ok {
			result[containConditionName] = val
			return result
		}

		result[containConditionName] = &ruleenginecore.ConditionType{
			Operator:    ruleenginecore.EqualOperator,
			OperandType: ruleenginecore.StringType,
			Operands: []*ruleenginecore.Operand{
				{
					OperandAs: ruleenginecore.OperandAsField,
					Val:       fieldname,
				},
				{
					OperandAs: ruleenginecore.OperandAsConstant,
					Val:       "apple",
				},
			},
		}
	}

	return result
}

func generateBoolConditions(fieldname string, existingCondition map[string]*ruleenginecore.ConditionType) map[string]*ruleenginecore.ConditionType {
	result := map[string]*ruleenginecore.ConditionType{}

	containConditionName := fieldname + "_" + ruleenginecore.EqualOperator + "_" + "true"

	if val, ok := existingCondition[containConditionName]; ok {
		result[containConditionName] = val
		return result
	}

	result[containConditionName] = &ruleenginecore.ConditionType{
		Operator:    ruleenginecore.EqualOperator,
		OperandType: ruleenginecore.BoolType,
		Operands: []*ruleenginecore.Operand{
			{
				OperandAs: ruleenginecore.OperandAsField,
				Val:       fieldname,
			},
			{
				OperandAs: ruleenginecore.OperandAsConstant,
				Val:       "true",
			},
		},
	}

	return result
}

func generateValidInput(fields ruleenginecore.Fields, fieldIntValMap map[string]int64, fieldFloatValMap map[string]float64) *ruleenginecore.Input {
	result := ruleenginecore.Input{}
	for name, typed := range fields {
		switch typed {
		case ruleenginecore.IntType:
			result[name] = strconv.FormatInt(fieldIntValMap[name]-intStep, 10)
		case ruleenginecore.FloatType:
			result[name] = strconv.FormatFloat(fieldFloatValMap[name]-floatStep, 'E', -1, 64)
		case ruleenginecore.StringType:
			result[name] = "apple"
		case ruleenginecore.BoolType:
			result[name] = "true"
		}
	}
	return &result
}

func generateInvalidInput(fields ruleenginecore.Fields, fieldIntValMap map[string]int64, fieldFloatValMap map[string]float64) *ruleenginecore.Input {
	result := ruleenginecore.Input{}
	for name, typed := range fields {
		switch typed {
		case ruleenginecore.IntType:
			result[name] = strconv.FormatInt(fieldIntValMap[name]-intStep, 10)
		case ruleenginecore.FloatType:
			result[name] = strconv.FormatFloat(fieldFloatValMap[name]-floatStep, 'E', -1, 64)
		case ruleenginecore.StringType:
			result[name] = "apple1"
		case ruleenginecore.BoolType:
			result[name] = "false"
		}
	}
	return &result
}

func writeJsonToFile(val any, filename string) {
	jsonBytes, err := json.MarshalIndent(val, "", "    ")

	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	err = os.WriteFile("./"+filename, jsonBytes, os.ModePerm)

	if err != nil {
		panic(err.Error())
	}
}

func fieldGenerate(intCount, floatCount, stringCount, boolCount int) ruleenginecore.Fields {
	fields := ruleenginecore.Fields{}

	// generate int field
	prefix := 'a'

	for i := 0; i < intCount; i++ {
		name := string(prefix) + "_" + ruleenginecore.IntType
		fields[name] = ruleenginecore.IntType
		prefix++
	}

	// generate float field
	prefix = 'a'
	for i := 0; i < floatCount; i++ {
		name := string(prefix) + "_" + ruleenginecore.FloatType
		fields[name] = ruleenginecore.FloatType
		prefix++
	}

	// generate string field
	prefix = 'a'
	for i := 0; i < stringCount; i++ {
		name := string(prefix) + "_" + ruleenginecore.StringType
		fields[name] = ruleenginecore.StringType
		prefix++
	}

	// generate bool field
	prefix = 'a'
	for i := 0; i < boolCount; i++ {
		name := string(prefix) + "_" + ruleenginecore.BoolType
		fields[name] = ruleenginecore.BoolType
		prefix++
	}

	return fields
}

func RandomBool() bool {
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(2) == 1
}
