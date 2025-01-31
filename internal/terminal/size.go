package terminal

import (
	"fmt"
	"os"

	"github.com/matchaprof/fogborne/internal/core/logging"
	"golang.org/x/term"
)

func GetTerminalSize() (width, height int, err error) {
	if !term.IsTerminal(int(os.Stdin.Fd())) {
		logging.Logger.Error("You are not using a terminal!")
		return 0, 0, fmt.Errorf("not a terminal")
	}

	width, height, err = term.GetSize(int(os.Stdin.Fd()))
	logging.Logger.Infof("Player's Terminal Size is %dx%d", width, height)
	return width, height, err
}
