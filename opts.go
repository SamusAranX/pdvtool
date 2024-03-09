package main

import "github.com/jessevdk/go-flags"

type DitherType string
type DitherDiffMode string

const (
	DitherNone           DitherType = "none"
	DitherBayer          DitherType = "bayer"
	DitherHeckbert       DitherType = "heckbert"
	DitherFloydSteinberg DitherType = "floyd_steinberg"
	DitherSierra2        DitherType = "sierra2"
	DitherSierra2Lite    DitherType = "sierra2_4a"
	DitherSierra3        DitherType = "sierra3"
	DitherBurkes         DitherType = "burkes"
	DitherAtkinson       DitherType = "atkinson"

	DitherDiffNone DitherDiffMode = "none"
	DitherDiffRect DitherDiffMode = "rectangle"
)

type DitherOptions struct {
	Type       DitherType     `long:"dither-type" short:"t" default:"sierra2" description:"The dithering mode" choice:"none" choice:"bayer" choice:"heckbert" choice:"floyd_steinberg" choice:"sierra2" choice:"sierra2_4a" choice:"sierra3" choice:"burkes" choice:"atkinson"`
	BayerScale int            `long:"dither-scale" short:"s" default:"0" description:"The dither scale (only for type bayer)" choice:"0" choice:"1" choice:"2" choice:"3" choice:"4" choice:"5"`
	DiffMode   DitherDiffMode `long:"dither-mode" short:"m" default:"rectangle" description:"The dithering diff mode" choice:"none" choice:"rectangle"`
}

type DefaultOptions struct {
	Debug bool `long:"debug" description:"Debug mode"`

	InFile flags.Filename `long:"infile" short:"i" required:"true" description:"The input file (.pdv or video file)"`
}

type DecodeOptions struct {
	DefaultOptions

	Pos struct {
		OutDir flags.Filename `positional-arg-name:"outdir" description:"The output directory"`
	} `positional-args:"yes" required:"yes"`
}

type EncodeOptions struct {
	DefaultOptions
	DitherOptions

	MakeMOV bool   `short:"x" description:"Generates a lossless x265 video instead of a pdv file. For testing purposes only."`
	Black   uint32 `short:"b" base:"16" description:"black" default:"000000"`
	White   uint32 `short:"w" base:"16" description:"white" default:"FFFFFF"`

	Pos struct {
		OutFile flags.Filename `positional-arg-name:"outfile" description:"The output file"`
	} `positional-args:"yes" required:"yes"`
}
