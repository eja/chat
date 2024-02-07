// Copyright (C) 2023-2024 by Ubaldo Porcheddu <ubaldo@eja.it>

package main

import (
	"encoding/json"
	"os/exec"
)

// Function to execute FFmpeg commands
func runFFmpeg(args []string) error {
	baseArgs := []string{"-y", "-nostdin", "-hide_banner"}
	cmd := exec.Command("ffmpeg", append(baseArgs, args...)...)
	Trace("ffmpeg", args)
	return cmd.Run()
}

// Function to execute ffprobe commands
func runFFprobe(args []string) ([]byte, error) {
	baseArgs := []string{"-y", "-nostdin", "-hide_banner", "-v", "error"}
	cmd := exec.Command("ffprobe", append(baseArgs, args...)...)
	Trace("ffprobe", args)
	return cmd.Output()
}

// Function to convert audio using FFmpeg for Google compatibility
func FFAudioGoogle(fileIn string, fileOut string) error {
	return runFFmpeg([]string{
		"-i", fileIn,
		"-vn", "-ar", "48000", "-ac", "1", "-c:a", "libopus", "-f", "ogg", fileOut,
	})
}

// Function to convert audio for metadata extraction
func FFAudioMeta(fileIn string, fileOut string) error {
	return runFFmpeg([]string{
		"-i", fileIn,
		"-vn", "-ar", "48000", "-b:a", "12k", "-ac", "1", "-c:a", "libopus", "-f", "ogg", fileOut,
	})
}

// Function to convert audio for whisper compatibility
func FFAudioWhisper(fileIn string, fileOut string) error {
	return runFFmpeg([]string{
		"-i", fileIn,
		"-ar", "16000", "-ac", "1", "-c:a", "pcm_s16le", "-f", "wav", fileOut,
	})
}

// Function to retrieve audio information using ffprobe
func FFProbeAudio(file string) (map[string]interface{}, error) {
	output, err := runFFprobe([]string{
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
