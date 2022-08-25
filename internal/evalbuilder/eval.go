package evalbuilder

import (
	"fmt"

	counter "github.com/DARKEMPIRESL/GODARKCHESS/pkg/eval/counter"
	fast "github.com/DARKEMPIRESL/GODARKCHESS/pkg/eval/fast"
	linear "github.com/DARKEMPIRESL/GODARKCHESS/pkg/eval/linear"
	material "github.com/DARKEMPIRESL/GODARKCHESS/pkg/eval/material"
	pesto "github.com/DARKEMPIRESL/GODARKCHESS/pkg/eval/pesto"
	weiss "github.com/DARKEMPIRESL/GODARKCHESS/pkg/eval/weiss"
)

func Build(key string) interface{} {
	switch key {
	case "counter":
		return counter.NewEvaluationService()
	case "weiss":
		return weiss.NewEvaluationService()
	case "linear":
		return linear.NewEvaluationService()
	case "pesto":
		return pesto.NewEvaluationService()
	case "material":
		return material.NewEvaluationService()
	case "fast":
		return fast.NewEvaluationService()
	}
	panic(fmt.Errorf("bad eval %v", key))
}
