package main

import (
	"flag"
	"log"
	"os"
	"time"

	"github.com/zhoukk/joy5/av"
	"github.com/zhoukk/joy5/format"
	"github.com/zhoukk/joy5/format/mp4"
	"github.com/zhoukk/joy5/format/rtsp"
)

func main() {
	var src string
	var dst string

	flag.StringVar(&src, "s", "rtsp://admin:xxx@192.168.1.64//Streaming/Channels/1", "rtsp src")
	flag.StringVar(&dst, "d", "out.mp4", "mp4 output")
	flag.Parse()

	var rtspc *rtsp.Client
	rtspc, err := rtsp.Dial(src)
	if err != nil {
		log.Fatal(err)
	}
	defer rtspc.Close()

	var f *os.File
	if f, err = os.Create(dst); err != nil {
		log.Fatal(err)
	}

	var m *mp4.Muxer
	m, err = mp4.NewMuxer(format.NewStreamsWriteSeeker(f, rtspc))
	if err != nil {
		log.Fatal(err)
	}
	err = m.WriteFileHeader()
	if err != nil {
		log.Fatal(err)
	}

	defer func() {
		if err := m.WriteTrailer(); err != nil {
			log.Fatal(err)
		}
		if err := f.Close(); err != nil {
			log.Fatal(err)
		}
	}()

	for {
		var pkt av.Packet
		if pkt, err = rtspc.ReadPacket(); err != nil {
			log.Fatal(err)
		}
		log.Println(pkt.CTime, pkt.Time)
		if err = m.WritePacket(pkt); err != nil {
			log.Fatal(err)
		}
		if pkt.Time > 10*time.Second {
			break
		}
	}
}
