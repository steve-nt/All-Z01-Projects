package utils

import (
	"log"
	"os/exec"
)

func RemoveImage(path string) error {

	path = "../../static" + path
	command := exec.Command("rm", path)
	err := command.Run()
	if err != nil {
		log.Println(err)
	}
	return nil
}
