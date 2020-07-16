package main

import (
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"runtime"
	"strconv"
	"strings"

	"github.com/chzyer/readline"
	"github.com/manifoldco/promptui"
	"github.com/urfave/cli/v2"
)

type instance struct {
	Number int
	IP     string
	Name   string
	Zone   string
}

var (
	version               = "v0.0.1"
	appAuthor *cli.Author = &cli.Author{
		Name:  "LuciferInLove",
		Email: "lucifer.in.love@ya.ru",
	}
)

func main() {
	app := &cli.App{
		Name:            "dynamic-sshmenu-aws",
		Usage:           "builds dynamic aws instances addresses list like sshmenu",
		Authors:         []*cli.Author{appAuthor},
		Action:          action,
		UsageText:       "dynamic-sshmenu-aws [-- <args>]",
		Version:         version,
		HideHelpCommand: true,
	}

	app.Commands = []*cli.Command{}

	app.Flags = []cli.Flag{
		&cli.StringFlag{
			Name:     "search-key",
			Aliases:  []string{"k"},
			Usage:    "key of instance tag to search",
			Required: true,
		},
		&cli.StringFlag{
			Name:     "search-value",
			Aliases:  []string{"s"},
			Usage:    "value of instance tag to search",
			Required: true,
		},
		&cli.StringFlag{
			Name:    "display-name",
			Aliases: []string{"d"},
			Usage:   "key of instance tag to display its values in results",
			Value:   "Name",
		},
	}

	app.RunAndExitOnError()
}

func promptSelect(instances []instance) (string, error) {
	var selectionSymbol string

	searcher := func(input string, index int) bool {
		instance := instances[index]
		name := strings.Replace(strings.ToLower(instance.Name), " ", "", -1)
		input = strings.Replace(strings.ToLower(input), " ", "", -1)

		return strings.Contains(name, input)
	}

	promptui.IconInitial = ""

	if runtime.GOOS == "windows" {
		selectionSymbol = "->"
	} else {
		selectionSymbol = "➜"
	}

	templates := promptui.SelectTemplates{
		Label:    "{{ . | cyan }}",
		Active:   fmt.Sprintf("%s %s", selectionSymbol, "{{ (printf \"%v.\t%v\t| %v (%v)\" .Number .IP .Name .Zone) | green }}"),
		Inactive: "  {{ (printf \"%v.\t%v\t| %v (%v)\" .Number .IP .Name .Zone) | white }}",
	}

	prompt := promptui.Select{
		Label:        "Select a target (press \"q\" for exit)",
		Items:        instances,
		Templates:    &templates,
		Size:         20,
		Searcher:     searcher,
		HideSelected: true,
		Keys: &promptui.SelectKeys{
			Next: promptui.Key{
				Code:    readline.CharNext,
				Display: "↓",
			},
			Prev: promptui.Key{
				Code:    readline.CharPrev,
				Display: "↑",
			},
			PageUp: promptui.Key{
				Code:    readline.CharForward,
				Display: "→",
			},
			PageDown: promptui.Key{
				Code:    readline.CharBackward,
				Display: "←",
			},
			Search: promptui.Key{
				Code:    '/',
				Display: "/",
			},
			Exit: promptui.Key{
				Code:    'q',
				Display: "q",
			},
		},
	}

	_, result, err := prompt.Run()

	if err != nil {
		if err == promptui.ErrInterrupt {
			return result, fmt.Errorf("Interrupted by \"%w\"", err)
		} else if err == promptui.ErrEOF {
			return result, fmt.Errorf("Unexpected end of file: \"%w\"", err)
		} else {
			return result, err
		}
	}

	return result, nil
}

func parseResult(result string) (instance, error) {
	var (
		parsedResult       []string
		instanceFromResult instance
	)

	re := regexp.MustCompile(`{(\d+\s+\d+\.\d+\.\d+\.\d+\s+\S+\s+\S+)}`)
	parsedResult = strings.Fields(re.ReplaceAllString(result, "$1"))

	if instanceNumber, err := strconv.Atoi(parsedResult[0]); err == nil {
		instanceFromResult = instance{
			Number: instanceNumber,
			IP:     string(parsedResult[1]),
			Name:   string(parsedResult[2]),
			Zone:   string(parsedResult[3]),
		}
	} else {
		return instanceFromResult, err
	}

	return instanceFromResult, nil
}

func action(c *cli.Context) error {
	instances, err := getSliceOfInstances(c.String("search-key"), c.String("search-value"), c.String("display-name"))

	if err != nil {
		return fmt.Errorf("There was an error listing instances in:\n%w", err)
	}

	result, err := promptSelect(instances)

	if err != nil {
		return err
	}

	instanceFromResult, err := parseResult(result)

	if err != nil {
		return err
	}

	cmd := exec.Command("ssh", instanceFromResult.IP)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()

	if err != nil {
		return fmt.Errorf("Command finished with an error, ssh: %w", err)
	}

	return nil
}
