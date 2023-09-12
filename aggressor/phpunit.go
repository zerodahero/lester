package aggressor

import (
	"encoding/xml"
	"io/fs"
	"lester/files"
	"log"
	"os"

	"github.com/pterm/pterm"
)

type Testsuite struct {
	Name      string   `xml:"name,attr"`
	Directory []string `xml:"directory,omitempty"`
	File      []string `xml:"file,omitempty"`
}

type Include struct {
	Group []string `xml:"group,omitempty"`
}
type Exclude struct {
	Group []string `xml:"group,omitempty"`
}

type Groups struct {
	Include *Include `xml:"include,omitempty"`
	Exclude *Exclude `xml:"exclude,omitempty"`
}

type Coverage struct {
	// IncludeUncoveredFiles string `xml:"includeUncoveredFiles,attr,omitempty"`
	// CacheDirectory        string `xml:"cacheDirectory,attr,omitempty"`
	Include struct {
		Directory []struct {
			Text   string `xml:",chardata"`
			Suffix string `xml:"suffix,attr"`
		} `xml:"directory,omitempty"`
	} `xml:"include,omitempty"`
}
type Ini struct {
	Name  string `xml:"name,attr"`
	Value string `xml:"value,attr"`
}

type Phpunit struct {
	XMLName          xml.Name `xml:"phpunit"`
	Bootstrap        string   `xml:"bootstrap,attr,omitempty"`
	BackupGlobals    string   `xml:"backupGlobals,attr,omitempty"`
	Colors           string   `xml:"colors,attr,omitempty"`
	Verbose          string   `xml:"verbose,attr,omitempty"`
	ProcessIsolation string   `xml:"processIsolation,attr,omitempty"`
	FailOnRisky      string   `xml:"failOnRisky,attr,omitempty"`
	StopOnFailure    string   `xml:"stopOnFailure,attr,omitempty"`
	Php              struct {
		Ini []Ini `xml:"ini,omitempty"`
	} `xml:"php,omitempty"`
	Testsuites struct {
		Testsuite []Testsuite `xml:"testsuite"`
	} `xml:"testsuites"`
	Groups   *Groups  `xml:"groups,omitempty"`
	Coverage Coverage `xml:"coverage,omitempty"`
}

var osWriteFile = os.WriteFile

func readPhpunitXml(fs fs.FS, path files.IProjectPath) Phpunit {
	file, err := path.ReadFile(fs)
	if err != nil {
		log.Fatalln(err)
	}

	var p Phpunit
	if err := xml.Unmarshal(file, &p); err != nil {
		log.Fatalln(err)
	}

	return p
}

func (p *Phpunit) reset() {
	p.Testsuites.Testsuite = []Testsuite{}
	p.Groups = nil
	p.Coverage = Coverage{}
}

func (p *Phpunit) addTests(tests []string) {
	p.Testsuites.Testsuite = []Testsuite{
		{File: tests, Name: "lester-auto"},
	}
}

func (p *Phpunit) addGroups(groups []string) {
	if len(groups) == 0 {
		return
	}

	if p.Groups == nil {
		p.Groups = &Groups{}
	}
	if p.Groups.Include == nil {
		p.Groups.Include = &Include{}
	}
	p.Groups.Include.Group = groups
}

func (p *Phpunit) write(xmlFileName string) {
	xmlOut, err := xml.MarshalIndent(p, " ", "  ")
	if err != nil {
		pterm.Fatal.Println(err)
	}

	err = osWriteFile(xmlFileName, xmlOut, 0644)
	if err != nil {
		pterm.Fatal.Println(err)
	}
}
