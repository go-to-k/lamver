package io

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/AlecAivazis/survey/v2"
)

const checkboxesPageSize = 30

func GetCheckboxes(label string, opts []string) ([]string, bool) {
	var checkboxes []string

	for {
		prompt := &survey.MultiSelect{
			Message:  label,
			Options:  opts,
			PageSize: checkboxesPageSize,
		}
		survey.AskOne(prompt, &checkboxes)

		if len(checkboxes) == 0 {
			Logger.Warn().Msg("Select values!")
			ok := GetYesNo("Do you want to finish?")
			if ok {
				Logger.Info().Msg("Finished...")
				return checkboxes, false
			}
			continue
		}

		ok := GetYesNo("OK?")
		if ok {
			return checkboxes, true
		}
	}
}

func InputKeywordForFilter(label string) string {
	reader := bufio.NewReader(os.Stdin)

	fmt.Fprintf(os.Stderr, "%s", label)
	s, _ := reader.ReadString('\n')
	fmt.Fprintln(os.Stderr)

	s = strings.TrimSpace(s)

	return s
}

func GetYesNo(label string) bool {
	choices := "Y/n"
	r := bufio.NewReader(os.Stdin)
	var s string

	for {
		fmt.Fprintf(os.Stderr, "%s (%s) ", label, choices)
		s, _ = r.ReadString('\n')
		fmt.Fprintln(os.Stderr)

		s = strings.TrimSpace(s)
		if s == "" {
			return true
		}
		s = strings.ToLower(s)
		if s == "y" || s == "yes" {
			return true
		}
		if s == "n" || s == "no" {
			return false
		}
	}
}
