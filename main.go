package main

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"github.com/chzyer/readline"
	"github.com/manifoldco/promptui"
	"github.com/urfave/cli/v2"
)

var (
	connectTo string
	version               = "v0.1.2"
	appAuthor *cli.Author = &cli.Author{
		Name:  "LuciferInLove",
		Email: "lucifer.in.love@protonmail.com",
	}
)

const (
	sshExecutable = "ssh"
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
			Name:    "tags",
			Aliases: []string{"t"},
			Usage:   "instance tags in \"key1:value1,value2;key2:value1\" format. If undefined, full list will be shown",
			Value:   "",
		},
		&cli.StringFlag{
			Name:    "display-name",
			Aliases: []string{"d"},
			Usage:   "key of instance tag to display its values in results",
			Value:   "Name",
		},
		&cli.BoolFlag{
			Name:    "public-ip",
			Aliases: []string{"p"},
			Usage:   "use public ip instead of private",
			Value:   false,
		},
		&cli.StringFlag{
			Name:    "ssh-username",
			Aliases: []string{"u"},
			Usage:   "ssh username. If undefined, the current user will be used",
			Value:   "",
		},
	}

	app.RunAndExitOnError()
}

func promptSelect(instances []string) (string, error) {
	var selectionSymbol string

	searcher := func(input string, index int) bool {
		instance, err := parseInstance(instances[index])
		if err != nil {
			return false
		}

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

	var funcMap = promptui.FuncMap
	funcMap["parse"] = parseInstance

	templates := promptui.SelectTemplates{
		Label: "{{ . | cyan }}",
		Active: fmt.Sprintf("%s %s", selectionSymbol,
			"{{ (printf \"%v.\t%v\t| %v (%v)\" (. | parse).Number (. | parse).IP (. | parse).Name (. | parse).Zone) | green }}",
		),
		Inactive: "  {{ (printf \"%v.\t%v\t| %v (%v)\" (. | parse).Number (. | parse).IP (. | parse).Name (. | parse).Zone) | white }}",
		FuncMap:  funcMap,
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
		return result, err
	}

	return result, nil
}

func action(c *cli.Context) error {
	username := c.String("ssh-username")

	instances, err := getSliceOfInstances(c.String("tags"), c.String("display-name"), c.Bool("public-ip"))
	if err != nil {
		if err.Error() == "WrongTagDefinition" {
			cli.ShowAppHelp(c)
			return fmt.Errorf("\nIncorrect Usage. Wrong tag definition in flag -t")
		}
		//return fmt.Errorf("Listing AWS instances:\n%w", err)
		instances = []string{`{"Number":1,"IP":"172.16.0.11","Name":"test-instance","Zone":"us-east-1a"}`}
	}

	result, err := promptSelect(instances)
	if err != nil {
		switch err {
		case promptui.ErrInterrupt:
			return nil
		case promptui.ErrEOF:
			return fmt.Errorf("Unexpected end of file: \"%w\"", err)
		default:
			return err
		}
	}

	instanceFromResult, err := parseInstance(result)
	if err != nil {
		return err
	}

	sshPath, err := exec.LookPath(sshExecutable)
	if err != nil {
		return fmt.Errorf("%w", err)
	}

	if username == "" {
		connectTo = instanceFromResult.IP
	} else {
		connectTo = username + "@" + instanceFromResult.IP
	}

	cmd := exec.Command(sshPath, connectTo)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err = cmd.Run()
	if err != nil {
		return fmt.Errorf("Command finished with an error, ssh: %w", err)
	}

	return nil
}
