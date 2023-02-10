package FFmpeg

import (
	"bytes"
	"fmt"
	ffmpeg "github.com/u2takey/ffmpeg-go"
	"log"
	"os"
)

func GetCover(buf *bytes.Buffer, videoURL string) error {

	err := ffmpeg.Input(videoURL).
		Filter("select", ffmpeg.Args{fmt.Sprintf("gte(n,%d)", 5)}).
		Output("pipe:", ffmpeg.KwArgs{"vframes": 1, "format": "image2", "vcodec": "mjpeg"}).
		WithOutput(buf, os.Stdout).
		Run()
	if err != nil {
		log.Printf("cover maker failed : %v", err)
	}
	log.Printf("cover maker succeed")

	return err

}
