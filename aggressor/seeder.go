package aggressor

import (
	"bufio"
	"bytes"
	"fmt"
	"lester/files"
	"os/exec"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/pterm/pterm"
)

var execLookPath = exec.LookPath
var execCommand = exec.Command

func SeedFilesFromGit(a *Aggressor, spinner *pterm.SpinnerPrinter) error {
	if _, err := execLookPath("git"); err != nil {
		return fmt.Errorf("could not find git (needed to seed initial modified files): %w", err)
	}

	var result error
	spinner.UpdateText("Seeding committed difference")
	if err := seedCommittedDifferences(a); err != nil {
		result = multierror.Append(result, err)
	}
	spinner.UpdateText("Seeding untracked files")
	if err := seedUntrackedFiles(a); err != nil {
		result = multierror.Append(result, err)
	}
	spinner.UpdateText("Seeding unstaged files")
	if err := seedTrackedFiles(a, false); err != nil {
		result = multierror.Append(result, err)
	}
	spinner.UpdateText("Seeding staged files")
	if err := seedTrackedFiles(a, true); err != nil {
		result = multierror.Append(result, err)
	}

	return result
}

func seedCommittedDifferences(a *Aggressor) error {
	branch, err := execCommand("git", "rev-parse", "--abbrev-ref", "HEAD").Output()
	if err != nil {
		return fmt.Errorf("error detecting current branch for git diff: %w", err)
	}

	// Committed differences between this branch and main
	out, err := execCommand("git", "diff", "--name-only", strings.TrimSpace(string(branch)), "main").Output()
	if err != nil {
		return fmt.Errorf("error attempting git diff for committed files: %w", err)
	}

	return addOutputToModifiedFiles(a, out)
}

func seedUntrackedFiles(a *Aggressor) error {
	out, err := execCommand("git", "ls-files", "--other", "--exclude-standard").Output()
	if err != nil {
		return fmt.Errorf("error attempting to find untracked files: %w", err)
	}

	return addOutputToModifiedFiles(a, out)
}

func seedTrackedFiles(a *Aggressor, staged bool) error {
	args := []string{"diff", "--name-only"}
	if staged {
		args = append(args, "--staged")
	}

	out, err := execCommand("git", args...).Output()
	if err != nil {
		return fmt.Errorf("error attempting to find tracked files (staged: %t): %w", staged, err)
	}

	return addOutputToModifiedFiles(a, out)
}

func addOutputToModifiedFiles(a *Aggressor, output []byte) error {
	scanner := bufio.NewScanner(bytes.NewReader(output))
	for scanner.Scan() {
		file := files.NewFromPath(scanner.Text())
		if file.IsPhpFile() {
			a.AddModifiedFile(file)
		}
	}

	return nil
}
