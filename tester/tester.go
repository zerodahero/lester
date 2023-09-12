package tester

import (
	"errors"
	"lester/files"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/pterm/pterm"
)

type ITestRunner interface {
	RunTests([]string) bool
	SetConfigArg([]string)
}

type TestRunner struct {
	PhpunitPath string
	ConfigArg   []string
	WorkingDir  string
}

func (t *TestRunner) RunTests(args []string) bool {
	args = normalizeFileArgs(args)
	parts := append(t.ConfigArg, args...)

	cmd := exec.Command(t.PhpunitPath, parts...)
	cmd.Dir = t.WorkingDir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	pterm.Debug.Printfln("Running %s from %s", cmd, cmd.Dir)
	err := cmd.Run()
	if err != nil {
		if errors.Is(err, &exec.ExitError{}) {
			pterm.Fatal.Println(err)
		}
		return false
	}
	return true
}

func NewTestRunnerFromPath(workingPath files.IProjectPath) *TestRunner {
	phpunitPath, configArg := getPhpunitBinaryFrom(workingPath)
	return &TestRunner{
		PhpunitPath: phpunitPath,
		ConfigArg:   configArg,
		WorkingDir:  ".",
	}
}

func (t *TestRunner) SetConfigArg(configArg []string) {
	t.ConfigArg = configArg
}

func getPhpunitBinaryFrom(workingPath files.IProjectPath) (string, []string) {

	if workingPath.Join("composer.json").FileDoesNotExist() {
		// Merged to root, run from root

		return getRootPhpunit(), getXmlConfigurationArgFor(workingPath)
	}

	maybePhpunit := workingPath.Join("vendor", "bin", "phpunit")
	if maybePhpunit.FileDoesNotExist() {
		log.Fatalln("Looks like we may need to install some composer deps. I wasn't able to find " + maybePhpunit.GetPath())
	}

	return maybePhpunit.MakeAbsolute().GetPath(), []string{}
}

func getXmlConfigurationArgFor(workingPath files.IProjectPath) []string {
	return []string{
		"--configuration",
		workingPath.Join("phpunit.xml").MakeAbsolute().GetPath(),
	}
}

func getRootPhpunit() string {
	maybePhpunit := files.NewFromCwd("vendor", "bin", "phpunit")
	if maybePhpunit.FileDoesNotExist() {
		log.Fatalln("Can't find phpunit! Last place we looked was " + maybePhpunit.GetPath())
	}

	return maybePhpunit.MakeAbsolute().GetPath()
}

func normalizeFileArgs(args []string) []string {
	for i, arg := range args {
		if strings.HasSuffix(arg, "Test.php") {
			args[i] = files.NewFromPath(arg).MakeAbsolute().GetPath()
		}
	}

	return args
}
