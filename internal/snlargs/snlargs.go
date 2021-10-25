package snlargs

import (
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/kassybas/shannel/internal/snlapi"
	"github.com/sirupsen/logrus"
)

var argWsSeparatedRegex, argEqSeparatedRegex *regexp.Regexp

func init() {
	argWsSeparatedRegex = regexp.MustCompile(`^-{1,2}\w*$`)
	argEqSeparatedRegex = regexp.MustCompile(`^-{1,2}\w*=.*$`)
}

// LoadArgs evaluates the arguments to variables based on the passed CLI flags
func LoadArgs(argConfs []snlapi.Arg, cliArgs []string) (map[string]*string, error) {
	argRes := map[string]*string{}

	am := map[string]snlapi.Arg{}
	posSl := []snlapi.Arg{}
	for _, ac := range argConfs {
		if ac.Type == "" || ac.Type == "bool" || ac.Type == "named" {
			am[ac.Name] = ac
			continue
		}
		if ac.Type == "pos" {
			posSl = append(posSl, ac)
			continue
		}
		//TODO: Validation
		return nil, fmt.Errorf("unknown argument type: name: '%s', type: '%s' should be one of: pos, bool, named", ac.Name, ac.Type)
	}

	prevHandled := false
	passedPosArg := 0
	for i, cliArg := range cliArgs {
		if prevHandled {
			prevHandled = false
			continue
		}
		// Whitespace separated named or bool flag
		if match := argWsSeparatedRegex.MatchString(cliArg); match {
			argName := strings.TrimLeft(cliArg, "-")
			c, exists := am[argName]
			if !exists {
				return nil, fmt.Errorf("named argument does not exist: name '%s'", argName)
			}
			// Bool flag
			if c.Type != "" && c.Type == "bool" {
				trueVal := "1"
				argRes[argName] = &trueVal
				continue
			}
			if c.Type == "" || c.Type == "named" {
				// Named flag whitespace separated
				if i+1 >= len(cliArgs) {
					return nil, fmt.Errorf("missing value after last named argument: name '%s'", cliArg)
				}
				argRes[argName] = &cliArgs[i+1]
				prevHandled = true
				continue
			}
		}
		// Equal sign separated named argument
		if match := argEqSeparatedRegex.MatchString(cliArg); match {
			tmpS := strings.TrimLeft(cliArg, "-")
			spl := strings.SplitN(tmpS, "=", 2)
			argName, argValue := spl[0], spl[1]
			c, exists := am[argName]
			if !exists {
				return nil, fmt.Errorf("named argument does not exist: '%s'", argName)
			}
			if !(c.Type == "" || c.Type == "named") {
				return nil, fmt.Errorf("value provided for non-named argument %s: '%s'", argName, cliArg)
			}
			argRes[argName] = &argValue
			continue
		}

		// Positional arguments
		if len(posSl) > passedPosArg {
			a := posSl[passedPosArg]
			argRes[a.Name] = &cliArgs[i]
			passedPosArg++
		} else {
			return nil, fmt.Errorf("too many positional arguments given: '%s'", cliArg)
		}
	}

	for i := range argConfs {
		argName := argConfs[i].Name
		if _, exists := argRes[argName]; !exists {
			if argConfs[i].FromEnvVar != nil {
				value, isSet := os.LookupEnv(argName)
				if isSet {
					argRes[argName] = &value
					continue
				}
			}
			if argConfs[i].Default != nil {
				argRes[argName] = argConfs[i].Default
				continue
			}
			if argConfs[i].Type == "bool" {
				falseVal := "0"
				argRes[argName] = &falseVal
				continue
			}
			if argConfs[i].Type == "pos" {
				return nil, fmt.Errorf("value for positional argument missing: '%s'", argName)
			}
			return nil, fmt.Errorf("value for argument missing: not given via parameter, environment variable or default value: '%s'", argName)
		}
	}
	logrus.Trace("parsed args", argRes)

	return argRes, nil
}
