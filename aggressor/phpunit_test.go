package aggressor

import (
	"encoding/xml"
	"lester/files"
	"os"
	"testing"
	"testing/fstest"

	"github.com/stretchr/testify/assert"
)

func TestReadPhpunitXml(t *testing.T) {
	fs := fstest.MapFS{}
	path := files.NewMockIProjectPath(t)
	path.On("ReadFile", fs).Return([]byte("<phpunit></phpunit>"), nil)

	got := readPhpunitXml(fs, path)
	assert.EqualValues(t, Phpunit{XMLName: xml.Name{Local: "phpunit"}}, got)
}

func TestReset(t *testing.T) {
	p := Phpunit{
		Groups: &Groups{},
	}
	p.Testsuites.Testsuite = []Testsuite{
		{Name: "SomeName", Directory: []string{"Dir1", "Dir2"}, File: []string{"File1", "File2"}},
	}
	p.reset()
	assert.Equal(t, []Testsuite{}, p.Testsuites.Testsuite)
	assert.Nil(t, p.Groups)
	assert.Equal(t, Coverage{}, p.Coverage)
}

func TestAddTests(t *testing.T) {
	p := Phpunit{}
	p.addTests([]string{"test1.php", "test2.php"})

	assert.Equal(t, "lester-auto", p.Testsuites.Testsuite[0].Name)
	assert.ElementsMatch(t, []string{"test1.php", "test2.php"}, p.Testsuites.Testsuite[0].File)
}

func TestAddGroups(t *testing.T) {
	p := Phpunit{}
	p.addGroups([]string{"group1", "group2"})
	assert.ElementsMatch(t, []string{"group1", "group2"}, p.Groups.Include.Group)
}

func TestWrite(t *testing.T) {
	p := Phpunit{}
	called := 0
	osWriteFile = func(name string, data []byte, perm os.FileMode) error {
		assert.Equal(t, "banana.xml", name)
		out, _ := xml.MarshalIndent(p, " ", "  ")
		assert.Equal(t, out, data)
		called++
		return nil
	}
	p.write("banana.xml")
	assert.Equal(t, 1, called)
}
