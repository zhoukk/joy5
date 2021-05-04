package main

import (
	"log"
	"os"

	"github.com/zhoukk/joy5/av"
	"github.com/zhoukk/joy5/format"
	"github.com/zhoukk/joy5/format/mp4"
	"github.com/zhoukk/joy5/format/rtsp"
)

func main() {
	dst := "out.mp4"
	src := "rtsp://admin:xxx@192.168.1.64//Streaming/Channels/1"

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

	for i := 0; i < 200; i++ {
		var pkt av.Packet
		if pkt, err = rtspc.ReadPacket(); err != nil {
			log.Fatal(err)
		}
		if err = m.WritePacket(pkt); err != nil {
			log.Fatal(err)
		}

	}
}
