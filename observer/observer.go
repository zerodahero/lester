package observer

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/fsnotify/fsnotify"
	"github.com/pterm/pterm"
)

func WatchForChanges(rootPath string, modified chan<- string) {
	spinner, _ := pterm.DefaultSpinner.Start("Looking for directories to watch...")

	dirs, err := getSubDirectories(rootPath)
	if err != nil {
		pterm.Fatal.Println(err)
	}

	// Create new watcher.
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		pterm.Fatal.Println(err)
	}
	defer watcher.Close()

	// Start listening for events.
	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				if !event.Has(fsnotify.Chmod) && isPhpFile(event.Name) {
					modified <- event.Name
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				pterm.Error.Println("error:", err)
			}
		}
	}()

	if len(dirs) == 0 {
		pterm.Fatal.Println("No directories?!? That's not going to work. Where are my PHP files at?")
	}

	// Add a path.
	for _, dir := range dirs {
		err = watcher.Add(dir)
		if err != nil {
			pterm.Fatal.Println(err)
		}
	}
	spinner.Info(fmt.Sprintf("Watching %d directories\n", len(dirs)))

	// Block main goroutine forever.
	<-make(chan struct{})
}

func getSubDirectories(path string) ([]string, error) {
	files, err := os.ReadDir(path)
	if err != nil {
		return []string{}, err
	}

	paths := []string{}
	if shouldAddPath(path) {
		paths = append(paths, path)
	}
	for _, file := range files {
		if !file.IsDir() {
			continue
		}

		subPath := filepath.Clean(filepath.Join(path, file.Name()))
		if !shouldTraverse(subPath) {
			continue
		}
		subPaths, err := getSubDirectories(subPath)
		if err != nil {
			return []string{}, err
		}
		paths = append(paths, subPaths...)
	}

	return paths, nil
}

func shouldAddPath(path string) bool {
	return containsPhpFiles(path) &&
		!isVendorDir(path)
}

func shouldTraverse(path string) bool {
	return !isVendorDir(path)
}

func isVendorDir(path string) bool {
	base := filepath.Base(path)
	return base == "vendor" || base == "Vendor"
}

func containsPhpFiles(path string) bool {
	files, err := os.ReadDir(path)
	if err != nil {
		return false
	}
	for _, file := range files {
		if isPhpFile(file.Name()) {
			return true
		}
	}

	return false
}

func isPhpFile(path string) bool {
	return filepath.Ext(path) == ".php"
}
