package common

import (
	"fmt"
	"os/exec"
	"runtime"
)

// CopyToClipboard copies the given text to the clipboard.
func CopyToClipboard(text string) error {
	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "darwin":
		cmd = exec.Command("pbcopy")
	case "linux":
		cmd = exec.Command("xclip", "-selection", "clipboard")
	case "windows":
		cmd = exec.Command("cmd", "/c", "clip")
	default:
		return fmt.Errorf("unsupported platform")
	}

	in, err := cmd.StdinPipe()
	if err != nil {
		return err
	}

	if err := cmd.Start(); err != nil {
		return err
	}

	_, err = in.Write([]byte(text))
	if err != nil {
		return err
	}

	in.Close()
	return cmd.Wait()
}
