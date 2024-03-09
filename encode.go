package main

import (
	"log"

	"pdvtool/ffmpeg"
)

func (opt *EncodeOptions) Execute(args []string) error {
	log.Println("encode!")
	log.Printf("options: %+v", *opt)
	log.Printf("args: %+v", args)

	probeInfo, err := ffmpeg.Probe(string(opt.InFile))
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("%+v", probeInfo)

	return nil
}
