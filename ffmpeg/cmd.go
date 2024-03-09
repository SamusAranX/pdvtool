package ffmpeg

import (
	"errors"
	"fmt"
	"os/exec"
)

func execCmd(args []string) ([]byte, error) {
	// fmt.Println(strings.Join(args, " "))

	stdout, err := exec.Command(args[0], args[1:]...).Output()
	if err != nil {
		// fmt.Println(stdout)
		return stdout, errors.New(fmt.Sprintf("error running %s: %s", args[0], err.Error()))
	}

	return stdout, nil
}
