//go:build !windows

package picker

import (
	"os"

	"golang.org/x/term"
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

func termHeight() int {
	_, h, err := term.GetSize(int(os.Stdin.Fd()))
	if err != nil || h <= 0 {
		return 24
	}
	return h
}
