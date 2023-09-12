package aggressor

import (
	"lester/files"
	"lester/tester"
	"os"
	"testing"
	"testing/fstest"

	"github.com/stretchr/testify/assert"
	"golang.org/x/exp/maps"
)

func TestGetRootPath(t *testing.T) {
	path := files.NewMockIProjectPath(t)
	c := Component{
		rootPath: path,
	}

	got := c.getRootPath()
	assert.EqualValues(t, path, got)
}

func TestGetFs(t *testing.T) {
	fs := fstest.MapFS{}
	c := Component{
		fs: fs,
	}

	got := c.getFs()
	assert.EqualValues(t, fs, got)
}

func TestAddMatchingTest_GetTests(t *testing.T) {
	c := Component{
		matchingTests: make(map[string]bool),
	}

	c.addMatchingTest("test1.php")

	got := c.getMatchingTests()
	assert.Equal(t, []string{"test1.php"}, maps.Keys(got))
}

func TestRunComponentTests_NoTests_Passes(t *testing.T) {
	c := Component{
		matchingTests: make(map[string]bool),
	}

	got := c.runTests([]string{})
	assert.Equal(t, true, got)
}

func TestRunComponentTests_WithTests_Passes(t *testing.T) {
	tr := tester.NewMockITestRunner(t)
	fs := fstest.MapFS{
		// prefix is not necessary because of fs root
		"phpunit.xml": {Data: []byte("<phpunit></phpunit>")},
	}
	called := 0
	osWriteFile = func(name string, data []byte, perm os.FileMode) error {
		called++
		return nil
	}
	c := Component{
		rootPath:      files.NewFromPath("banana/sundae"),
		matchingTests: map[string]bool{"test1.php": true},
		testRunner:    tr,
		fs:            fs,
	}
	tr.Mock.On("RunTests", []string{}).Return(true)
	tr.Mock.On("SetConfigArg", []string{"--configuration", "banana/sundae/.lester-phpunit.xml"})

	got := c.runTests([]string{})
	assert.Equal(t, true, got)
	assert.Equal(t, 1, called)
}

func TestHasTestsToRun_NoTest_Fails(t *testing.T) {
	c := Component{
		matchingTests: make(map[string]bool),
	}
	got := c.hasTestsToRun()
	assert.Equal(t, false, got)
}

func TestHasTestsToRun_WithTest_Passes(t *testing.T) {
	c := Component{
		matchingTests: map[string]bool{
			"test1.php": true,
		},
	}
	got := c.hasTestsToRun()
	assert.Equal(t, true, got)
}
