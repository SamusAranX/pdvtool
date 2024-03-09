package main

import (
	"bytes"
	"compress/zlib"
	"encoding/binary"
	"fmt"
	"image"
	"image/color"
	"image/gif"
	"image/png"
	"io"
	"log"
	"os"
	"path/filepath"
	"strconv"

	"pdvtool/pdv"
)

var (
	colorBlack = color.Gray{}
	colorWhite = color.Gray{Y: 255}
	bwPalette  = color.Palette{colorBlack, colorWhite}

	gifOpt     = gif.Options{NumColors: 2}
	pngEncoder = png.Encoder{CompressionLevel: png.BestCompression}
)

func readStruct[T interface{}](r io.Reader, t *T) error {
	err := binary.Read(r, binary.LittleEndian, t)
	if err != nil {
		return err
	}

	return nil
}

func (opt *DecodeOptions) Execute(args []string) error {
	log.Println("decode!")
	log.Printf("options: %+v", *opt)
	log.Printf("args: %+v", args)

	f, err := os.Open(string(opt.InFile))
	if err != nil {
		return fmt.Errorf("can't open file: %v", err)
	}

	defer f.Close()

	var header pdv.Header
	err = readStruct(f, &header)
	if err != nil {
		return fmt.Errorf("can't read header: %v", err)
	}

	log.Printf("header: %+v", header)
	log.Printf("valid magic: %t", header.ValidateMagic())

	if !header.ValidateMagic() {
		return fmt.Errorf("invalid header")
	}

	outFilePattern := fmt.Sprintf("%%0%dd.png", len(strconv.Itoa(int(header.NumFrames))))
	outFilePattern = filepath.Join(string(opt.Pos.OutDir), outFilePattern)

	frameTable := make([]pdv.FrameTableEntry, header.NumFrames+1)
	for i := 0; i < len(frameTable); i++ {
		var entry pdv.FrameTableEntry
		err = readStruct(f, &entry)
		if err != nil {
			return fmt.Errorf("can't read frame table entry: %v", err)
		}

		frameTable[i] = entry
	}

	if frameTable[len(frameTable)-1].Type() != pdv.FrameTypeEmpty {
		return fmt.Errorf("invalid frame table")
	}
	log.Printf("frame table entries: %d", len(frameTable))

	frameDataOffset, _ := f.Seek(0, io.SeekCurrent)

	log.Printf("frame data offset: %d (0x%[1]X)", frameDataOffset)

	intWidth := int(header.Width)
	intHeight := int(header.Height)
	minP := image.Point{}
	maxP := image.Point{X: intWidth, Y: intHeight}
	imgRect := image.Rectangle{
		Min: minP,
		Max: maxP,
	}

	log.Printf("writing %d frames…", header.NumFrames)

	for i := 0; i < len(frameTable); i++ {
		entry := frameTable[i]
		nextEntry := frameTable[i+1]

		offset := frameDataOffset + int64(entry.Offset())
		length := frameDataOffset + int64(nextEntry.Offset()) - offset

		_, err = f.Seek(offset, io.SeekStart)
		if err != nil {
			return fmt.Errorf("couldn't seek to position %d: %v", offset, err)
		}

		b := make([]byte, length)
		_, err = f.Read(b)

		zf, err := zlib.NewReader(bytes.NewBuffer(b))
		if err != nil {
			return fmt.Errorf("couldn't init zlib reader: %v", err)
		}

		frameData, err := io.ReadAll(zf)
		if err != nil {
			return fmt.Errorf("couldn't decompress frame: %v", err)
		}

		frameImage := image.NewPaletted(imgRect, bwPalette)
		// log.Printf("frameImage %d, pix: %d", i+1, len(frameImage.Pix))

		for j, bits := range frameData {
			pixels := pdv.ExpandBits(bits)

			p := j * 8
			frameImage.Pix[p+0] = pixels[0]
			frameImage.Pix[p+1] = pixels[1]
			frameImage.Pix[p+2] = pixels[2]
			frameImage.Pix[p+3] = pixels[3]
			frameImage.Pix[p+4] = pixels[4]
			frameImage.Pix[p+5] = pixels[5]
			frameImage.Pix[p+6] = pixels[6]
			frameImage.Pix[p+7] = pixels[7]
		}

		ofName := fmt.Sprintf(outFilePattern, i+1)
		of, err := os.Create(ofName)
		if err != nil {
			return fmt.Errorf("couldn't create file %s: %v", ofName, err)
		}

		// err = gif.Encode(of, frameImage, &gifOpt)
		// if err != nil {
		// 	return fmt.Errorf("couldn't encode gif to file %s: %v", ofName, err)
		// }
		err = pngEncoder.Encode(of, frameImage)
		if err != nil {
			return fmt.Errorf("couldn't encode png to file %s: %v", ofName, err)
		}

		err = of.Close()
		if err != nil {
			return fmt.Errorf("couldn't close file %s: %v", ofName, err)
		}

		if nextEntry.Type() == pdv.FrameTypeEmpty {
			log.Printf("frame %d processed, end of table reached", i+1)
			break
		}
	}

	return nil
}
