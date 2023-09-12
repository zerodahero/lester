package printer

import (
	"io"

	"github.com/pterm/pterm"
	"github.com/pterm/pterm/putils"
)

type Printer struct {
	out io.Writer
}

func NewPrinter(w io.Writer) *Printer {
	return &Printer{
		out: w,
	}
}

func (p *Printer) PrintLesterIntro() {
	pterm.DefaultCenter.Println("Hello, my name is")
	s, _ := pterm.DefaultBigText.WithLetters(putils.LettersFromString("Lester")).Srender()
	pterm.DefaultCenter.Println(s)
	pterm.DefaultCenter.Println("The Legacy Tester!")
	pterm.DefaultCenter.Println("Now, let's get to testing...")
}
