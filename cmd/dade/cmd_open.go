package main

import (
	"os/exec"
)

var openBrowserFunc = openBrowserDefault

func openBrowserDefault(url string) error {
	return exec.Command("open", url).Start()
}
