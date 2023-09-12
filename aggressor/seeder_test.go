package aggressor

import (
	"fmt"
	"lester/files"
	"os"
	"os/exec"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func fakeExecCommand(command string, args ...string) *exec.Cmd {
	cs := []string{"-test.run=TestHelperProcess", "--", command}
	cs = append(cs, args...)
	cmd := exec.Command(os.Args[0], cs...)
	cmd.Env = []string{
		"GO_WANT_HELPER_PROCESS=1",
		"TEST_COMMAND=" + command,
		"TEST_ARGS=" + strings.Join(args, ";;"),
	}
	return cmd
}

const fileResult = "foo.php\nfooTest.php"

func TestSeedsFromGit(t *testing.T) {
	cases := []struct {
		name   string
		seeder func(*Aggressor) error
	}{
		{name: "committed differences", seeder: seedCommittedDifferences},
		{name: "untracked files", seeder: seedUntrackedFiles},
		{
			name: "staged files",
			seeder: func(a *Aggressor) error {
				return seedTrackedFiles(a, true)
			},
		},
		{
			name: "unstaged files",
			seeder: func(a *Aggressor) error {
				return seedTrackedFiles(a, true)
			},
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			a := Aggressor{
				Components: map[string]IComponent{
					"root": &Component{
						rootPath:      files.NewFromPath("."),
						matchingTests: make(map[string]bool),
					},
				},
				ModifiedFiles: make(map[string]bool),
			}
			execCommand = fakeExecCommand
			defer func() { execCommand = exec.Command }()

			tt.seeder(&a)

			assert.ElementsMatch(t, []string{"foo.php", "fooTest.php"}, a.GetModified())
		})
	}
}

func TestHelperProcess(t *testing.T) {
	if os.Getenv("GO_WANT_HELPER_PROCESS") != "1" {
		return
	}
	// some code here to check arguments perhaps?
	fmt.Fprintf(os.Stdout, fileResult)
	os.Exit(0)
}
