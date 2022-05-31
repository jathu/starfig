package command

import (
	"fmt"
	"os"
	"strings"

	"github.com/jathu/starfig/internal/evaluator"
	"github.com/jathu/starfig/internal/starverse"
	"github.com/jathu/starfig/internal/target"
	"github.com/jathu/starfig/internal/util"
	"github.com/logrusorgru/aurora"
	"go.starlark.net/starlark"
)

func Build(args []string, keepGoing bool) error {
	starverseDir, err := starverse.FindStarverseDirectory()
	if err != nil {
		return err
	}

	evaluatedOutput := new(starlark.Dict)
	summary := buildResult{results: map[string]*[]error{}}

	for _, arg := range args {
		buildTargets, err := target.ParseBuildTarget(starverseDir, arg)
		if err != nil {
			if keepGoing {
				summary.note(arg, err)
			} else {
				return err
			}
		}
		for _, buildTarget := range buildTargets {
			evaluateResults, err := evaluator.EvaluateBuildTarget(starverseDir, buildTarget)
			if err != nil {
				if keepGoing {
					summary.note(buildTarget.Target(), err)
				} else {
					return err
				}
			}

			for _, evaluateResult := range evaluateResults {
				summary.note(evaluateResult.Target.Target(), nil)
				evaluatedOutput.SetKey(
					starlark.String(evaluateResult.Target.Target()),
					evaluateResult.Result.Evaluated,
				)
			}
		}
	}

	fmt.Println(value2json(evaluatedOutput))
	if keepGoing {
		printSummary(summary)
	}
	return nil
}

func printSummary(summary buildResult) {
	count := summary.count()
	countComponents := []string{"\nSummary:"}
	if count.failed > 0 {
		countComponents = append(countComponents, aurora.Red(fmt.Sprintf("%d FAIL", count.failed)).String())
	}
	if count.ok > 0 {
		countComponents = append(countComponents, aurora.Green(fmt.Sprintf("%d OK", count.ok)).String())
	}
	countComponents = append(countComponents, fmt.Sprintf("%d TOTAL", count.total))
	fmt.Fprintln(os.Stderr, strings.Join(countComponents, " "))
	for name, errs := range summary.results {
		if len(*errs) == 0 {
			fmt.Fprintln(os.Stderr, fmt.Sprintf("  %s %s", aurora.Green("  OK"), name))
		} else {
			errorMessages := util.JoinErrors(*errs, fmt.Sprintf("%s", aurora.Red("     * ")), "\n")
			fmt.Fprintln(os.Stderr, fmt.Sprintf("  %s %s\n%s", aurora.Red("FAIL"), name, errorMessages))
		}
	}
}

func value2json(rawValue starlark.Value) string {
	switch value := rawValue.(type) {
	case starlark.Bool:
		if value == starlark.True {
			return "true"
		} else {
			return "false"
		}
	case starlark.Tuple:
		return fmt.Sprintf("%s: %s", value2json(value.Index(0)), value2json(value.Index(1)))
	case *starlark.List:
		items := []string{}
		for i := 0; i < value.Len(); i++ {
			items = append(items, value2json(value.Index(i)))
		}
		return fmt.Sprintf("[%s]", strings.Join(items, ", "))
	case *starlark.Dict:
		items := []string{}
		for _, tuple := range value.Items() {
			items = append(items, value2json(tuple))
		}
		return fmt.Sprintf("{%s}", strings.Join(items, ", "))
	default:
		return value.String()
	}
}

type buildResult struct {
	results map[string]*[]error
}

func (br buildResult) note(name string, err error) {
	_, found := br.results[name]
	if !found {
		br.results[name] = &[]error{}
	}
	if err != nil {
		*br.results[name] = append(*br.results[name], err)
	}
}

type buildResultCountSummary struct {
	ok     int
	failed int
	total  int
}

func (br buildResult) count() buildResultCountSummary {
	summary := buildResultCountSummary{ok: 0, failed: 0, total: 0}
	for _, errs := range br.results {
		if len(*errs) == 0 {
			summary.ok += 1
		} else {
			summary.failed += 1
		}
		summary.total += 1
	}
	return summary
}
