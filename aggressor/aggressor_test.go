package aggressor

import (
	"fmt"
	"lester/files"
	"testing"
	"testing/fstest"

	"github.com/stretchr/testify/assert"
)

func TestInitComponentsFromConfig(t *testing.T) {
	config := map[string]string{
		"banana": "path/to/banana",
		"orange": "some/other/path",
	}
	a := NewAggressor()
	a.InitComponentsFromConfig(config)

	for c, p := range config {
		got, ok := a.Components[c]
		if !ok {
			t.Errorf("missing %s", c)
		}
		if got.getRootPath().GetPath() != p {
			t.Errorf("for %s, got %v, want %q", c, got, p)
		}
	}
}

func TestAddModifiedFile_SourceFile(t *testing.T) {
	fs := &fstest.MapFS{
		// path/banana prefix removed for fs scope
		"tests/sourceTest.php": {Data: []byte("@group chips")},
	}
	a := Aggressor{
		Components: map[string]IComponent{
			"banana": &Component{
				rootPath:      files.NewFromPath("path/banana"),
				fs:            fs,
				matchingTests: make(map[string]bool),
			},
			"orange": &Component{
				rootPath: files.NewFromPath("some/path/orange"),
			},
		},
		ModifiedFiles:  make(map[string]bool),
		MatchingGroups: make(map[string]bool),
	}

	path := files.NewMockIProjectPath(t)
	path.On("MakeRelative").Return(path)
	path.On("HasParentPath", a.Components["banana"].getRootPath()).Return(true)
	path.On("HasParentPath", a.Components["orange"].getRootPath()).Maybe().Return(false)
	path.On("GetPath").Return("path/banana/source.php")
	path.On("IsTestFile").Return(false)
	path.On("GetFileStringWithoutExt").Return("source")

	a.AddModifiedFile(path)

	if _, ok := a.ModifiedFiles[path.GetPath()]; !ok {
		t.Errorf("did not add %v to the modified files", path)
	}

	if _, ok := a.Components["banana"].getMatchingTests()["tests/sourceTest.php"]; !ok {
		t.Errorf("did not add matching test")
	}

	if _, ok := a.MatchingGroups["chips"]; !ok {
		t.Errorf("did not add matching group")
	}
}

func TestAddModifiedFile_TestFile(t *testing.T) {
	fs := &fstest.MapFS{
		// path/banana prefix removed for fs scope
		"tests/sourceTest.php": {Data: []byte("@group chips")},
	}
	a := Aggressor{
		Components: map[string]IComponent{
			"banana": &Component{
				rootPath:      files.NewFromPath("path/banana"),
				fs:            fs,
				matchingTests: make(map[string]bool),
			},
			"orange": &Component{
				rootPath: files.NewFromPath("some/path/orange"),
			},
		},
		ModifiedFiles:  make(map[string]bool),
		MatchingGroups: make(map[string]bool),
	}

	path := files.NewMockIProjectPath(t)
	path.On("MakeRelative").Return(path)
	path.On("HasParentPath", a.Components["banana"].getRootPath()).Return(true)
	path.On("HasParentPath", a.Components["orange"].getRootPath()).Maybe().Return(false)
	path.On("GetPath").Return("path/banana/tests/sourceTest.php")
	path.On("IsTestFile").Return(true)

	a.AddModifiedFile(path)

	if _, ok := a.ModifiedFiles[path.GetPath()]; !ok {
		t.Errorf("did not add %v to the modified files", path)
	}

	if _, ok := a.Components["banana"].getMatchingTests()["tests/sourceTest.php"]; !ok {
		t.Errorf("did not add matching test")
	}

	if _, ok := a.MatchingGroups["chips"]; !ok {
		t.Errorf("did not add matching group")
	}
}

func TestMakesMatchingTestName(t *testing.T) {
	cases := []struct {
		code string
		test string
	}{
		{"/path/to/SomeCode.php", "SomeCodeTest.php"},
		{"/path/to/Some/Nested/Code.php", "CodeTest.php"},
	}

	for _, test := range cases {
		t.Run(fmt.Sprintf("%s matches to test name %s", test.code, test.test), func(t *testing.T) {
			got := makeMatchingTestClassName(files.NewFromPath(test.code))
			if got.GetPath() != test.test {
				t.Errorf("got %q, want %q", got.GetPath(), test.test)
			}
		})
	}
}

func TestGetComponentFromPath(t *testing.T) {
	a := Aggressor{
		Components: map[string]IComponent{
			"banana": &Component{
				rootPath: files.NewFromPath("path/banana"),
			},
			"orange": &Component{
				rootPath: files.NewFromPath("some/path/orange"),
			},
		},
	}
	path := files.NewMockIProjectPath(t)
	path.Mock.On("MakeRelative").Return(path)
	path.Mock.On("HasParentPath", files.NewFromPath("path/banana")).Return(true)
	path.Mock.On("HasParentPath", files.NewFromPath("some/path/orange")).Maybe().Return(true)

	got := a.getComponentFromPath(path)

	want := a.Components["banana"]
	assert.Equal(t, got.getRootPath(), want.getRootPath())
}

func TestGetComponentFromPath_FallBackToRoot(t *testing.T) {
	a := Aggressor{
		Components: map[string]IComponent{
			"banana": &Component{
				rootPath: files.NewFromPath("path/banana"),
			},
			"root": &Component{
				rootPath: files.NewFromPath("root/path"),
			},
		},
	}
	path := files.NewMockIProjectPath(t)
	path.Mock.On("MakeRelative").Return(path)
	path.Mock.On("HasParentPath", files.NewFromPath("path/banana")).Return(false)
	path.Mock.On("HasParentPath", files.NewFromPath("root/path")).Return(false)

	got := a.getComponentFromPath(path)

	want := a.Components["root"]
	assert.Equal(t, got.getRootPath(), want.getRootPath())
}

func TestRunTests_Empty_Passes(t *testing.T) {
	a := &Aggressor{}

	result := a.RunTests()

	assert.Equal(t, true, result)
}

func TestRunTests_WithTests_Passes(t *testing.T) {
	bananaComponent := NewMockIComponent(t)
	a := &Aggressor{
		Components: map[string]IComponent{
			"banana": bananaComponent,
		},
		ModifiedFiles:  map[string]bool{"file1": true},
		MatchingGroups: map[string]bool{"group1": true},
	}
	bananaComponent.Mock.On("runTests", []string{"group1"}).Return(true)

	result := a.RunTests()

	assert.Equal(t, true, result)
}

func TestHasTestsToRun(t *testing.T) {
	cases := []struct {
		name          string
		modifiedFiles []string
		matchingTests map[string][]string
		expected      bool
	}{
		{
			name:          "Has files and tests",
			modifiedFiles: []string{"file1", "file2"},
			matchingTests: map[string][]string{
				"banana": {"file1", "file2"},
			},
			expected: true,
		},
		{
			name:          "Has only files",
			modifiedFiles: []string{"file1", "file2"},
			matchingTests: map[string][]string{},
			expected:      false,
		},
		// Has only tests does't happen since tests
		// are always part of "modifiedFiles"
	}

	for _, tt := range cases {
		a := NewAggressor()
		a.Components = make(map[string]IComponent)
		a.ModifiedFiles = make(map[string]bool)

		for _, mf := range tt.modifiedFiles {
			a.ModifiedFiles[mf] = true
		}
		for c, mtg := range tt.matchingTests {
			a.Components[c] = &Component{
				matchingTests: make(map[string]bool),
			}
			for _, mt := range mtg {
				a.Components[c].getMatchingTests()[mt] = true
			}
		}

		t.Run(tt.name, func(t *testing.T) {
			got := a.HasTestsToRun()
			if got != tt.expected {
				t.Errorf("got %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestGetModified(t *testing.T) {
	a := NewAggressor()
	a.ModifiedFiles = map[string]bool{
		"file1": true,
		"file2": true,
		"file3": true,
	}

	got := a.GetModified()
	assert.ElementsMatch(t, []string{"file1", "file2", "file3"}, got)
}
