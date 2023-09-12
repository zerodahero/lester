package aggressor

import (
	"io/fs"
	"lester/files"
	"os"
	"regexp"
	"strings"

	"golang.org/x/exp/maps"
)

type Aggressor struct {
	Components     map[string]IComponent
	ModifiedFiles  map[string]bool
	MatchingGroups map[string]bool
}

var groupRegexp = regexp.MustCompile(`(?m)^[^@]*?@group ([a-zA-Z0-9_\-]+?)\s*?$`)

func NewAggressor() *Aggressor {
	return &Aggressor{
		Components:     make(map[string]IComponent),
		ModifiedFiles:  make(map[string]bool),
		MatchingGroups: make(map[string]bool),
	}
}

func (a *Aggressor) InitComponentsFromConfig(cfg map[string]string) {
	for name, path := range cfg {
		rootPath := files.NewFromPath(path)
		a.Components[name] = &Component{
			rootPath:      rootPath,
			matchingTests: make(map[string]bool),
			testRunner:    nil,
		}
	}
}

func (a *Aggressor) AddModifiedFile(path files.IProjectPath) {
	component := a.getComponentFromPath(path)
	a.ModifiedFiles[path.GetPath()] = true

	if path.IsTestFile() {
		a.addModifiedTest(component, path)
		return
	}

	a.addMatchingTestsAndGroups(component, path)
}

func (a *Aggressor) addModifiedTest(component IComponent, path files.IProjectPath) {
	testFile := strings.TrimPrefix(path.GetPath(), component.getRootPath().GetPath())
	testFile = strings.TrimPrefix(testFile, string(os.PathSeparator))
	component.addMatchingTest(testFile)

	a.addMatchingGroupsFromTest(component.getFs(), files.NewFromPath(testFile))
}

func (a *Aggressor) addMatchingTestsAndGroups(component IComponent, path files.IProjectPath) {
	fsys := component.getFs()
	testToFind := makeMatchingTestClassName(path)
	fs.WalkDir(fsys, ".", func(dirPath string, d fs.DirEntry, err error) error {
		if d.Name() == testToFind.GetPath() {
			component.addMatchingTest(dirPath)
			a.addMatchingGroupsFromTest(fsys, files.NewFromPath(dirPath))
		}
		return nil
	})
}

func (a *Aggressor) addMatchingGroupsFromTest(fs fs.FS, path files.IProjectPath) {
	// NOTE: temporarily skipping since this is reductive rather than expansive
	// contents, err := path.ReadFile(fs)
	// if err != nil {
	// 	pterm.Error.Printfln("Could not read file for matching groups at %s", path.GetPath())
	// 	pterm.Error.Println(err)
	// 	return
	// }

	// groups := groupRegexp.FindAllStringSubmatch(string(contents), -1)
	// for _, group := range groups {
	// 	if group[1] == "" {
	// 		continue
	// 	}
	// 	a.MatchingGroups[group[1]] = true
	// }
}

func makeMatchingTestClassName(path files.IProjectPath) files.IProjectPath {
	return files.NewFromPath(path.GetFileStringWithoutExt() + "Test.php")
}

func (a *Aggressor) getComponentFromPath(path files.IProjectPath) IComponent {
	path = path.MakeRelative()
	for _, component := range a.Components {
		if path.HasParentPath(component.getRootPath()) {
			return component
		}
	}
	// Default to root
	root := a.Components["root"]
	return root
}

func (a *Aggressor) RunTests() bool {
	// printModified(a.ModifiedFiles)

	groups := maps.Keys(a.MatchingGroups)

	allPassed := true
	for _, component := range a.Components {
		passed := component.runTests(groups)
		allPassed = allPassed && passed
	}

	return allPassed
}

func (a *Aggressor) HasTestsToRun() bool {
	if len(a.ModifiedFiles) == 0 {
		return false
	}

	for _, component := range a.Components {
		if component.hasTestsToRun() {
			return true
		}
	}

	return false
}

func (a *Aggressor) GetModified() []string {
	return maps.Keys(a.ModifiedFiles)
}

// func printModified(modified map[string]bool) {
// 	// FUTURE: tree?

// 	pterm.Info.Println("Detected modifications to these files:")
// 	pterm.DefaultBulletList.WithItems(items).Render()
// }
