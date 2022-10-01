package benchmark

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"testing"

	ruleenginecore "github.com/niharrathod/ruleengine-core"
)

var ruleEngine_10_Rule ruleenginecore.RuleEngine
var ruleEngine_100_Rule ruleenginecore.RuleEngine
var ruleEngine_1000_Rule ruleenginecore.RuleEngine

var validInput = ruleenginecore.Input{}
var invalidInput = ruleenginecore.Input{}

var resultTemp any

func init() {
	ruleEngine_10_Rule = prepareRuleEngine("generator/ruleEngineConfig_10.json")
	ruleEngine_100_Rule = prepareRuleEngine("generator/ruleEngineConfig_100.json")
	ruleEngine_1000_Rule = prepareRuleEngine("generator/ruleEngineConfig_1000.json")
	validInput = prepareInput("generator/ValidInput.json")
	invalidInput = prepareInput("generator/InvalidInput.json")
}

func Benchmark_RuleEngine_10_Rule_EvaluateComplete(b *testing.B) {
	var r any

	b.Run("validInput", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			r, _ = ruleEngine_10_Rule.Evaluate(context.TODO(), validInput, ruleenginecore.EvaluateOptions().Complete())
		}
	})
	resultTemp = r

	b.Run("invalidInput", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			r, _ = ruleEngine_10_Rule.Evaluate(context.TODO(), invalidInput, ruleenginecore.EvaluateOptions().Complete())
		}
	})
	resultTemp = r

}

func Benchmark_RuleEngine_100_Rule_EvaluateComplete(b *testing.B) {
	var r any

	b.Run("validInput", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			r, _ = ruleEngine_100_Rule.Evaluate(context.TODO(), validInput, ruleenginecore.EvaluateOptions().Complete())
		}
	})
	resultTemp = r

	b.Run("invalidInput", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			r, _ = ruleEngine_100_Rule.Evaluate(context.TODO(), invalidInput, ruleenginecore.EvaluateOptions().Complete())
		}
	})
	resultTemp = r
}

func Benchmark_RuleEngine_1000_Rule_EvaluateComplete(b *testing.B) {
	var r any

	b.Run("validInput", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			r, _ = ruleEngine_1000_Rule.Evaluate(context.TODO(), validInput, ruleenginecore.EvaluateOptions().Complete())
		}
	})
	resultTemp = r

	b.Run("invalidInput", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			r, _ = ruleEngine_1000_Rule.Evaluate(context.TODO(), invalidInput, ruleenginecore.EvaluateOptions().Complete())
		}
	})
	resultTemp = r

}

func prepareRuleEngine(configJsonFile string) ruleenginecore.RuleEngine {
	wd, _ := os.Getwd()
	valBytes, _ := os.ReadFile(wd + "/" + configJsonFile)
	var engineConfig ruleenginecore.RuleEngineConfig
	err := json.Unmarshal(valBytes, &engineConfig)

	if err != nil {
		fmt.Println("json unmarshal failed : ", err.Error())
		panic("No")
	}

	engine, reErr := ruleenginecore.New(&engineConfig)

	if reErr != nil {
		fmt.Println("ruleengine new failed : ", reErr.Error())
		panic("No")
	}

	return engine
}

func prepareInput(filename string) ruleenginecore.Input {
	input := ruleenginecore.Input{}

	err := json.Unmarshal(readJsonFile(filename), &input)
	if err != nil {
		fmt.Println(filename+" json unmarshal failed : ", err.Error())
		panic("No")
	}

	return input
}

func readJsonFile(filename string) []byte {
	wd, _ := os.Getwd()
	valBytes, _ := os.ReadFile(wd + "/" + filename)
	return valBytes
}
