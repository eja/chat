// Copyright (C) 2023-2024 by Ubaldo Porcheddu <ubaldo@eja.it>

package ff

import (
	"encoding/json"
	"os/exec"

	"github.com/eja/tibula/log"
)

func ffmpeg(args []string) error {
	baseArgs := []string{"-y", "-nostdin", "-hide_banner"}
	cmd := exec.Command("ffmpeg", append(baseArgs, args...)...)
	log.Trace("[FF]", "ffmpeg", args)
	return cmd.Run()
}

func ffprobe(args []string) ([]byte, error) {
	baseArgs := []string{"-y", "-nostdin", "-hide_banner", "-v", "error"}
	cmd := exec.Command("ffprobe", append(baseArgs, args...)...)
	log.Trace("[FF]", "ffprobe", args)
	return cmd.Output()
}

func MpegAudioGoogle(fileIn string, fileOut string) error {
	return ffmpeg([]string{
		"-i", fileIn,
		"-vn", "-ar", "48000", "-ac", "1", "-c:a", "libopus", "-f", "ogg", fileOut,
	})
}

func MpegAudioMeta(fileIn string, fileOut string) error {
	return ffmpeg([]string{
		"-i", fileIn,
		"-vn", "-ar", "48000", "-b:a", "12k", "-ac", "1", "-c:a", "libopus", "-f", "ogg", fileOut,
	})
}

func MpegAudioWhisper(fileIn string, fileOut string) error {
	return ffmpeg([]string{
		"-i", fileIn,
		"-ar", "16000", "-ac", "1", "-c:a", "pcm_s16le", "-f", "wav", fileOut,
	})
}

func ProbeAudio(file string) (map[string]interface{}, error) {
	output, err := ffprobe([]string{
		"-print_format", "json", "-show_format", "-show_streams", file,
	})
	if err != nil {
		return nil, err
	}

	var data map[string]interface{}
	err = json.Unmarshal(output, &data)
	if err != nil {
		return nil, err
	}

	var audioStream map[string]interface{}
	if streams, ok := data["streams"].([]interface{}); ok {
		for _, stream := range streams {
			if streamMap, ok := stream.(map[string]interface{}); ok {
				if codecType, ok := streamMap["codec_type"].(string); ok && codecType == "audio" {
					audioStream = streamMap
					break
				}
			}
		}
	}

	return audioStream, nil
}
