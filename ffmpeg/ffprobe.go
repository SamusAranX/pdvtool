package ffmpeg

import (
	"encoding/json"
	"errors"
	"fmt"
)

func Probe(inputFile string) (*ProbeInfo, error) {
	ffprobeArgs := []string{
		"ffprobe",
		"-v", "quiet",
		"-print_format", "json",
		"-show_streams",
		"-select_streams", "v",
		// "-count_frames",
		inputFile,
	}

	ffprobeStruct := &ProbeInfo{}

	// run ffprobe and get stdout
	stdout, err := execCmd(ffprobeArgs)
	if err != nil {
		return nil, fmt.Errorf("can't decode file: %s", err.Error())
	}

	// parse stdout into struct
	err = json.Unmarshal(stdout, ffprobeStruct)
	if err != nil {
		return nil, errors.New("can't make sense of ffprobe output")
	}

	// error out if the number of streams isn't exactly one
	numStreams := len(ffprobeStruct.Streams)
	if numStreams == 0 {
		return nil, errors.New("no video streams found")
	} else if numStreams > 1 {
		return nil, errors.New("more than one video stream found")
	}

	return ffprobeStruct, nil
}
