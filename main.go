package main

import (
	"fmt"
	"log"
	"os/exec"
	"regexp"
	"strings"
)

type ProgramInfo struct {
	name    string
	variant string
	version string
}

func handle(err error) {
	if err != nil {
		panic(err)
	}
}

func run(args []string) ([]string, error) {
	cmd := exec.Command(args[0], args[1:]...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return []string{}, err
	}

	return strings.Split(string(output), "\n"), nil
}

type Matcher struct {
	line  int
	regex string
}

func match(cmdArgs []string, matchers []Matcher) (string, error) {
	output, err := run(cmdArgs)
	if err != nil {
		return "", err
	}

	for _, matcher := range matchers {
		line := output[matcher.line]

		r, err := regexp.Compile(matcher.regex)
		if err != nil {
			return "", err
		}

		regexMatch := r.FindStringSubmatch(line)
		return regexMatch[r.SubexpIndex("version")], nil
	}

	return "", fmt.Errorf("No match found")
}

func getProgramInfo(programName string) ProgramInfo {
	switch programName {
	case "grep":
		version, err := match([]string{"grep", "--version"}, []Matcher{
			{
				line:  0,
				regex: "^grep \\(GNU grep\\) (?P<version>.*?)$",
			},
		})
		handle(err)

		return ProgramInfo{
			variant: "GNU",
			version: version,
		}
	case "awk":
		version, err := match([]string{"awk", "--version"}, []Matcher{
			{
				line:  0,
				regex: "^GNU Awk (?P<version>.*?),",
			},
		})
		handle(err)

		return ProgramInfo{
			variant: "GNU",
			version: version,
		}
	case "bash":
		version, err := match([]string{"bash", "--version"}, []Matcher{
			{
				line:  0,
				regex: "GNU bash, version (?P<version>.*?) ",
			},
		})
		handle(err)

		return ProgramInfo{
			variant: "GNU",
			version: version,
		}
	case "tar":
		version, err := match([]string{"tar", "--version"}, []Matcher{
			{
				line:  0,
				regex: "tar \\(GNU tar\\) (?<version>.*)$",
			},
		})
		handle(err)

		return ProgramInfo{
			variant: "GNU",
			version: version,
		}
	default:
		log.Fatalln("Failed to match against a supported program")
	}

	return ProgramInfo{}
}

func main() {
	for _, program := range []string{"grep", "awk", "bash", "tar"} {
		info := getProgramInfo(program)
		fmt.Printf("Grep Variant: %s\n", info.variant)
		fmt.Printf("Grep Version: %s\n", info.version)
		fmt.Println()
	}
}
