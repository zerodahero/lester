package aggressor

import (
	"io/fs"
	"lester/files"
	"lester/tester"
	"os"

	"golang.org/x/exp/maps"
)

type IComponent interface {
	getFs() fs.FS
	runTests([]string) bool
	getTestRunner() tester.ITestRunner
	setTestConfigArg([]string)
	hasTestsToRun() bool
	generatePhpunitXmlFile([]string)
	getRootPath() files.IProjectPath
	addMatchingTest(string)
	getMatchingTests() map[string]bool
}

type Component struct {
	rootPath      files.IProjectPath
	fs            fs.FS
	matchingTests map[string]bool
	testRunner    tester.ITestRunner
}

func (c *Component) getRootPath() files.IProjectPath {
	return c.rootPath
}

func (c *Component) getFs() fs.FS {
	if c.fs != nil {
		return c.fs
	}

	c.fs = os.DirFS(c.rootPath.GetPath())
	return c.fs
}

func (c *Component) addMatchingTest(testFile string) {
	c.matchingTests[testFile] = true
}

func (c *Component) getMatchingTests() map[string]bool {
	return c.matchingTests
}

func (c *Component) runTests(groups []string) bool {
	if !c.hasTestsToRun() {
		return true
	}

	c.generatePhpunitXmlFile(groups)
	tr := c.getTestRunner()
	return tr.RunTests([]string{})
}

func (c *Component) getTestRunner() tester.ITestRunner {
	if c.testRunner != nil {
		return c.testRunner
	}

	c.testRunner = tester.NewTestRunnerFromPath(c.rootPath)
	return c.testRunner
}

func (c *Component) setTestConfigArg(configArgs []string) {
	tr := c.getTestRunner()
	tr.SetConfigArg(configArgs)
}

func (c *Component) hasTestsToRun() bool {
	return len(c.matchingTests) > 0
}

func (c *Component) generatePhpunitXmlFile(groups []string) {
	phpunit := readPhpunitXml(c.getFs(), files.NewFromPath("phpunit.xml"))
	phpunit.reset()
	phpunit.addTests(maps.Keys(c.matchingTests))
	phpunit.addGroups(groups)

	xmlFileName := c.rootPath.Join(".lester-phpunit.xml").GetPath()
	phpunit.write(xmlFileName)

	c.setTestConfigArg([]string{"--configuration", xmlFileName})
}
