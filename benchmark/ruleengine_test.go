package benchmark

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"testing"

	ruleenginecore "github.com/niharrathod/ruleengine-core"
)

var singleRuleEngine ruleenginecore.RuleEngine

func init() {
	wd, _ := os.Getwd()
	valBytes, _ := os.ReadFile(wd + "/SingleRule.json")
	var engineConfig ruleenginecore.RuleEngineConfig
	err := json.Unmarshal(valBytes, &engineConfig)

	if err != nil {
		fmt.Println("json unmarshal failed : ", err.Error())
		panic("No")
	}

	engine, reErr := ruleenginecore.New(&engineConfig)

	if reErr != nil {
		fmt.Println("json unmarshal failed : ", reErr.Error())
		panic("No")
	}

	singleRuleEngine = engine
}

var singleRuleInputTable = []ruleenginecore.Input{
	{
		"TotalAmount":    "25000",
		"IsHotelBooking": "true",
		"PaxCount":       "10",
	},
}
var result any

func BenchmarkRuleEngineComplete(b *testing.B) {
	var r any
	for i, input := range singleRuleInputTable {
		b.Run(fmt.Sprintf("SingleRule_input_%d", i), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				r, _ = singleRuleEngine.Evaluate(context.TODO(), input, ruleenginecore.EvaluateOptions().Complete())
			}
		})
		result = r
	}
}
