package camforchat

import (
	"context"
	"fmt"
	"github.com/pion/rtp"
	"github.com/pion/webrtc/v2/pkg/media/ivfwriter"
	"log"
	"time"
)

// BroadcastScreenshot keeps video screen
type BroadcastScreenshot struct {
	BroadcastID string
	// Receives stream
	StreamBuf chan []byte

	ivffile *ivfwriter.IVFWriter
}

// NewBroadcastScreenshot creates new BroadcastScreenshot
func NewBroadcastScreenshot(broadcastID string) (*BroadcastScreenshot, error) {
	bs := &BroadcastScreenshot{BroadcastID: broadcastID, StreamBuf: make(chan []byte, 2)}
	if err := bs.rotateIVFFile(); err != nil {
		return nil, err
	}
	return bs, nil
}

// TakeScreenshots takes screenshot periodically from stream
func (bs *BroadcastScreenshot) TakeScreenshots(ctx context.Context) {
	go func() {
		defer bs.Close()

		tick := time.Tick(time.Second * 60)

		for {
			select {
			case <-tick:
				if err := bs.rotateIVFFile(); err != nil {
					log.Println("Rotation file: ", err)
				}
			case rtpBuf := <-bs.StreamBuf:
				if bs.ivffile == nil {
					continue
				}

				if err := bs.writeIVFFile(rtpBuf); err != nil {
					log.Println("Write file: ", err)
				}
			case <-ctx.Done():
				break
			}
		}
	}()
}

// Close closes current screenshot
func (bs *BroadcastScreenshot) Close() {
	bs.closeIVFFile()
}

func (bs *BroadcastScreenshot) closeIVFFile() {
	if bs.ivffile != nil {
		bs.ivffile.Close()
	}
}

func (bs *BroadcastScreenshot) rotateIVFFile() error {
	var err error

	bs.closeIVFFile()

	filename := fmt.Sprintf("/app/videos/vid-%s-%d.ivf", bs.BroadcastID, time.Now().Unix())
	bs.ivffile, err = ivfwriter.New(filename)
	if err != nil {
		bs.ivffile = nil
		return err
	}

	return nil
}

func (bs *BroadcastScreenshot) writeIVFFile(rtpBuf []byte) error {
	r := &rtp.Packet{}
	if err := r.Unmarshal(rtpBuf); err != nil {
		return err
	}
	if err := bs.ivffile.WriteRTP(r); err != nil {
		return err
	}

	return nil
}
