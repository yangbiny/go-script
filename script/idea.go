package main

import "os/exec"

func main() {
	path := "idea"
	command := exec.Command(path, "/Volumes/workspace/acm-api")

	err := command.Start()
	err = command.Wait()
	if err != nil {
		println(err)
	}

}
