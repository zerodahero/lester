package files

import (
	"errors"
	"io"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"
)

type IProjectPath interface {
	MakeRelative() IProjectPath
	MakeAbsolute() IProjectPath
	IsTestFile() bool
	IsPhpFile() bool
	FileDoesNotExist() bool
	Join(...string) IProjectPath
	GetPath() string
	ReadFile(fs.FS) ([]byte, error)
	HasParentPath(IProjectPath) bool
	GetFileStringWithoutExt() string
}

type ProjectPath struct {
	Path string
}

var cwd IProjectPath

func NewFromCwd(paths ...string) IProjectPath {
	return GetCwd().Join(paths...)
}

func NewFromPath(path string) IProjectPath {
	return ProjectPath{Path: path}
}

func GetCwd() IProjectPath {
	if cwd != nil {
		return cwd
	}

	wd, err := os.Getwd()
	if err != nil {
		log.Fatalln(err)
	}

	cwd = ProjectPath{Path: wd}
	return cwd
}

func (f ProjectPath) FileDoesNotExist() bool {
	_, err := os.Stat(f.Path)
	return errors.Is(err, os.ErrNotExist)
}

func (f ProjectPath) MakeAbsolute() IProjectPath {
	absPath, err := filepath.Abs(f.Path)
	if err != nil {
		log.Fatalln("Unable to get absolute path for " + f.Path)
	}

	return ProjectPath{Path: absPath}
}

func (f ProjectPath) MakeRelative() IProjectPath {
	cwd, err := os.Getwd()
	if err != nil {
		log.Fatalln("Could not get cwd to make relative file path", err)
	}
	return ProjectPath{strings.TrimPrefix(f.Path, cwd)}
}

func (f ProjectPath) IsTestFile() bool {
	return len(f.Path) > 8 && f.Path[len(f.Path)-8:] == "Test.php"
}

func (f ProjectPath) IsPhpFile() bool {
	return filepath.Ext(f.GetPath()) == ".php"
}

func (f ProjectPath) Join(paths ...string) IProjectPath {
	parts := append([]string{f.Path}, paths...)
	newFile := ProjectPath{Path: filepath.Join(parts...)}
	return newFile
}

func (f ProjectPath) GetPath() string {
	return f.Path
}

func (f ProjectPath) ReadFile(fs fs.FS) ([]byte, error) {
	file, err := fs.Open(f.GetPath())
	if err != nil {
		return []byte{}, err
	}
	defer file.Close()

	contents, err := io.ReadAll(file)
	if err != nil {
		return []byte{}, err
	}

	return contents, nil
}

func (f ProjectPath) HasParentPath(parentPath IProjectPath) bool {
	return strings.HasPrefix(f.GetPath(), parentPath.GetPath())
}

func (f ProjectPath) GetFileStringWithoutExt() string {
	return strings.TrimSuffix(filepath.Base(f.GetPath()), filepath.Ext(f.GetPath()))
}
