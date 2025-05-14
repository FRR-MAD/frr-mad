package common

import (
	"encoding/base64"
	"fmt"
	"os"
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

func copyOSC52(data string) error {
	// 1) Base64-encode
	enc := base64.StdEncoding.EncodeToString([]byte(data))

	// 2) Build the OSC 52 sequence
	seq := fmt.Sprintf("\x1b]52;c;%s\a", enc)

	// 3) If running inside tmux, wrap it so tmux passes it through
	if os.Getenv("TMUX") != "" {
		// tmux escape: P + tmux; + ESC + sequence + ESC + \
		seq = fmt.Sprintf("\x1bPtmux;\x1b%s\x1b\\", seq)
	}

	// 4) Print to stdout
	_, err := os.Stdout.WriteString(seq)
	return err
}
