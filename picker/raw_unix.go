//go:build !windows

package picker

import (
	"golang.org/x/term"
	"os"
)

var oldState *term.State

func setRaw(enable bool) error {
	if enable {
		state, err := term.MakeRaw(int(os.Stdin.Fd()))
		if err != nil {
			return err
		}
		oldState = state
		return nil
	}
	if oldState != nil {
		return term.Restore(int(os.Stdin.Fd()), oldState)
	}
	return nil
}
