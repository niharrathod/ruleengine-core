package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	ruleenginecore "github.com/niharrathod/ruleengine-core"
)

func main() {
	/*
		For example, We need rule engine for coupon and discount service, where customer would apply coupon and we need to validate coupon and determine the discount.
		lets say we have below two rules:
			rule 1:
				RuleName : RuleDiscount20Per
				Rule : (totalAmount > 20000) and (only for flightBooking) and (fightDestination are Bangalore,Delhi,Mumbai,Chennai)
				RuleResult:
					discountPer : 20

			rule 2:
				RuleName : RuleDiscount10Per
				Rule : (totalAmount <= 20000) and (only for flightBooking) and (fightDestination are Bangalore,Delhi,Mumbai,Chennai)
				RuleResult:
					discountPer : 10

		lets prepare RuleEngineConfig for about rule engine
			Step.1 : determine list of fields and appropriate fieldtypes
				as per above rule below list would cover all
				1. totalAmount  		-> int type
				2. isFlightBooking 		-> bool type
				3. flightDestination 	-> string type

			Step.2: define custom conditions
				1. (totalAmount > 20000)
						greaterCondition with Int operand type

				2. (isFlightBooking = true)
						equalCondition with bool operand type

				3. (flightDestination are Bangalore,Delhi,Mumbai,Chennai)
						containCondition with string operand type

			Step.3: define rule combining custom conditions
				rule condition combination supports 'and','or','not' operations
				take a look at exampleConfig.json as an example

		exampleConfig.json represent RuleEngineConfig as a json form
	*/

	ruleEngineConfig := &ruleenginecore.RuleEngineConfig{}

	// exampleConfig.json file has RuleEngineConfig for this example
	readJsonFileAndUnmarshal("exampleConfig.json", &ruleEngineConfig)

	ruleengine, err := ruleenginecore.New(ruleEngineConfig)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	input1 := ruleenginecore.Input{
		"totalAmount":       "25000",
		"isFlightBooking":   "true",
		"flightDestination": "Bangalore",
	}

	outputs, err := ruleengine.Evaluate(context.TODO(), input1, ruleenginecore.EvaluateOptions().Complete())
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	for _, output := range outputs {
		fmt.Printf("%+v\n", output)
	}

	input2 := ruleenginecore.Input{
		"totalAmount":       "15000",
		"isFlightBooking":   "true",
		"flightDestination": "Bangalore",
	}

	outputs, err = ruleengine.Evaluate(context.TODO(), input2, ruleenginecore.EvaluateOptions().Complete())
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	for _, output := range outputs {
		fmt.Printf("%+v\n", output)
	}

}

func readJsonFileAndUnmarshal(configJsonFile string, val any) {
	wd, _ := os.Getwd()
	valBytes, _ := os.ReadFile(wd + "/" + configJsonFile)
	err := json.Unmarshal(valBytes, val)

	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}
