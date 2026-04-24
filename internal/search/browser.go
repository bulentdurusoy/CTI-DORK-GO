package search

import (
	"fmt"
	"net/url"
	"os/exec"
	"runtime"
)

// OpenInBrowser opens a Google search query in the user's default web browser.
// This is the safe execution mode that avoids scraping Google directly.
func OpenInBrowser(query string) error {
	searchURL := "https://www.google.com/search?q=" + url.QueryEscape(query)

	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "windows":
		cmd = exec.Command("rundll32", "url.dll,FileProtocolHandler", searchURL)
	case "darwin":
		cmd = exec.Command("open", searchURL)
	case "linux":
		cmd = exec.Command("xdg-open", searchURL)
	default:
		return fmt.Errorf("unsupported operating system: %s", runtime.GOOS)
	}

	return cmd.Start()
}
