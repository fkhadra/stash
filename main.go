package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/AlecAivazis/survey/v2"
	"github.com/AlecAivazis/survey/v2/terminal"
	"github.com/manifoldco/promptui"
)

const (
	Apply = iota
	Save
	Drop
)

type Stash struct {
	Label       string
	Description string
}

func main() {
	cmd := exec.Command("git", "status")
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		os.Exit(1)
	}

	var choice int
	selectAction := &survey.Select{
		Message: "What do you want to do:",
		Options: []string{
			"Apply - apply the changes recorded in the stash",
			"Save - save your local modifications to a new stash",
			"Remove - remove one or more entries from the stash list",
		},
	}

	ask(selectAction, &choice)

	switch choice {
	case Apply:
		applyStash()
	case Drop:
		dropStash()
	case Save:
		saveStash()
	}
}

func dropStash() {
	stashes := getStashList(false)
	var options []string

	for _, stash := range stashes {
		options = append(options, stash.Label)
	}

	var choices []int
	selectStashes := &survey.MultiSelect{
		Message: "Select the entries to remove",
		Options: options,
	}

	ask(selectStashes, &choices)
	for idx, v := range choices {
		output, err := exec.Command("git", "stash", "drop", stashAtIndex(v-idx)).Output()
		if err != nil {
			panic(err.Error())
		}
		println(string(output))
	}
}

func saveStash() {
	msg := ""
	prompt := &survey.Input{
		Message: "Type a description for your stash or press enter if you don't care:",
		Help:    "Useful to remember what the stash is about",
	}

	ask(prompt, &msg)

	var cmd *exec.Cmd

	if len(msg) > 0 {
		cmd = exec.Command("git", "stash", "save", msg)
	} else {
		cmd = exec.Command("git", "stash", "save")
	}
	if stdout, err := cmd.Output(); err != nil {
		panic(err)
	} else {
		fmt.Println(string(stdout))
		os.Exit(0)
	}
}

func getStashList(withDescription bool) []Stash {
	stdout, err := exec.Command("git", "stash", "list").Output()
	if err != nil {
		panic(err.Error())
	} else if len(stdout) == 0 {
		fmt.Println("This repo does not contains any stash in it")
		os.Exit(0)
	}

	var stashes []Stash
	for index, line := range strings.Split(string(stdout), "\n") {
		if len(line) > 0 {
			description := ""

			if withDescription {
				output, err := exec.Command("git", "stash", "show", stashAtIndex(index)).Output()
				if err != nil {
					panic(err.Error())
				}
				description = string(output)
			}

			stashes = append(stashes, Stash{Description: description, Label: line})
		}
	}

	return stashes
}

func stashAtIndex(index int) string {
	return fmt.Sprintf("stash@{%d}", index)
}

func applyStash() {
	stashes := getStashList(true)
	prompt := promptui.Select{
		Label: "Select a stash from the list",
		Items: stashes,
		Templates: &promptui.SelectTemplates{
			Active:   "👉 {{ .Label | cyan }}",
			Inactive: "{{ .Label | cyan }}",
			Selected: "\U0001f7e2 {{ .Label | green }}",
			Details: `
📝 Stash content
-----------------
{{ .Description }}
`,
		},
	}

	choice, _, err := prompt.Run()

	if err == promptui.ErrInterrupt {
		os.Exit(-1)
	}

	output, err := exec.Command("git", "stash", "apply", fmt.Sprint(choice)).Output()

	if err != nil {
		panic(err.Error())
	}

	fmt.Println(string(output))
	os.Exit(0)
}

func ask(prompt survey.Prompt, choice interface{}) {
	if err := survey.AskOne(prompt, choice); err != nil {
		if err == terminal.InterruptErr {
			fmt.Println("👋 Ciao!")
			os.Exit(1)
		}
	}
}
