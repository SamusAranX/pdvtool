package ffmpeg

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

type ProbeInfo struct {
	Streams []struct {
		Index              int    `json:"index"`
		CodecName          string `json:"codec_name"`
		CodecLongName      string `json:"codec_long_name"`
		Profile            string `json:"profile"`
		CodecType          string `json:"codec_type"`
		CodecTagString     string `json:"codec_tag_string"`
		CodecTag           string `json:"codec_tag"`
		Width              int    `json:"width"`
		Height             int    `json:"height"`
		CodedWidth         int    `json:"coded_width"`
		CodedHeight        int    `json:"coded_height"`
		ClosedCaptions     int    `json:"closed_captions"`
		HasBFrames         int    `json:"has_b_frames"`
		SampleAspectRatio  string `json:"sample_aspect_ratio"`
		DisplayAspectRatio string `json:"display_aspect_ratio"`
		PixFmt             string `json:"pix_fmt"`
		Level              int    `json:"level"`
		ChromaLocation     string `json:"chroma_location"`
		Refs               int    `json:"refs"`
		IsAvc              string `json:"is_avc"`
		NalLengthSize      string `json:"nal_length_size"`
		RFrameRate         string `json:"r_frame_rate"`
		AvgFrameRate       string `json:"avg_frame_rate"`
		TimeBase           string `json:"time_base"`
		StartPts           int    `json:"start_pts"`
		StartTime          string `json:"start_time"`
		DurationTs         int    `json:"duration_ts"`
		Duration           string `json:"duration"`
		BitRate            string `json:"bit_rate"`
		BitsPerRawSample   string `json:"bits_per_raw_sample"`
		NbFrames           string `json:"nb_frames"`
		NbReadFrames       string `json:"nb_read_frames"`
		Disposition        struct {
			Default         int `json:"default"`
			Dub             int `json:"dub"`
			Original        int `json:"original"`
			Comment         int `json:"comment"`
			Lyrics          int `json:"lyrics"`
			Karaoke         int `json:"karaoke"`
			Forced          int `json:"forced"`
			HearingImpaired int `json:"hearing_impaired"`
			VisualImpaired  int `json:"visual_impaired"`
			CleanEffects    int `json:"clean_effects"`
			AttachedPic     int `json:"attached_pic"`
			TimedThumbnails int `json:"timed_thumbnails"`
		} `json:"disposition"`
		Tags struct {
			Language    string `json:"language"`
			HandlerName string `json:"handler_name"`
			VendorID    string `json:"vendor_id"`
		} `json:"tags"`
	} `json:"streams"`
}

// FrameRate returns the first stream's frame rate (in frames per second)
func (fi ProbeInfo) FrameRate() (frameRate float64, err error) {
	var fps float64

	fpsString := fi.Streams[0].RFrameRate
	fps, err = strconv.ParseFloat(fpsString, 64)

	if err != nil {
		// frame rate value is probably a fraction

		if !strings.Contains(fpsString, "/") {
			// frame rate value is not a fraction
			return 0, errors.New("can't parse frame rate")
		}

		fpsStringParts := strings.Split(fpsString, "/")
		numerator, _ := strconv.ParseFloat(fpsStringParts[0], 64)
		determinator, _ := strconv.ParseFloat(fpsStringParts[1], 64)
		fps = numerator / determinator
	}

	return fps, nil
}

// Width returns the first stream's width in pixels
func (fi ProbeInfo) Width() int {
	return fi.Streams[0].Width
}

// Height returns the first stream's height in pixels
func (fi ProbeInfo) Height() int {
	return fi.Streams[0].Height
}

// Duration returns the first stream's duration
func (fi ProbeInfo) Duration() (time.Duration, error) {
	duration, err := strconv.ParseFloat(fi.Streams[0].Duration, 64)
	if err != nil {
		return 0, errors.New("can't parse video duration")
	}

	return time.Duration(duration * float64(time.Second)), nil
}

// NumFrames returns the first stream's total number of frames
func (fi ProbeInfo) NumFrames() int64 {
	numReadFrames, err := strconv.ParseInt(fi.Streams[0].NbReadFrames, 10, 64)
	if err != nil {
		return -1
	}

	return numReadFrames
}

// FrameNameFormat returns a format string for frame extraction. If Num
func (fi ProbeInfo) FrameNameFormat(ext string) string {
	nf := fi.NumFrames()
	if nf <= 0 {
		return fmt.Sprintf("frame%%06d.%s", ext)
	}

	return fmt.Sprintf("frame%%0%dd.%s", len(fmt.Sprintf("%d", nf)), ext)
}
