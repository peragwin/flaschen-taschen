package flaschen

import (
	"fmt"
	"image"
	"image/color"
	"net"
)

// Flaschen represents a client connection to a Flaschen server
type Flaschen struct {
	remote string
	conn   net.Conn

	width       int
	height      int
	layer       int
	transparent bool

	header []byte
	footer []byte

	image *image.RGBA
}

func NewFlaschen(width, height, layer int, remote string) (*Flaschen, error) {
	conn, err := net.Dial("udp", remote)
	if err != nil {
		return nil, err
	}
	header := fmt.Sprintf("P6\n%d %d\n255\n", width, height)
	footer := fmt.Sprintf("0\n0\n%d\n", layer)
	return &Flaschen{
		remote:      remote,
		conn:        conn,
		width:       width,
		height:      height,
		layer:       0,
		transparent: false,
		image:       image.NewRGBA(image.Rect(0, 0, width, height)),
		header:      []byte(header),
		footer:      []byte(footer),
	}, nil
}

func (f *Flaschen) Rect() image.Rectangle {
	return f.image.Rect
}

func (f *Flaschen) Pixel(x, y int, col color.RGBA) {
	if x >= f.width || y >= f.height {
		panic("pixel out of range")
	}
	if col.R == 0 && !f.transparent {
		col.R = 1
	}
	if col.G == 0 && !f.transparent {
		col.G = 1
	}
	if col.B == 0 && !f.transparent {
		col.B = 1
	}
	// fmt.Println("set pixel", col)
	f.image.Set(x, y, col)
}

func (f *Flaschen) Show() error {
	data := make([]byte, len(f.header)+3*f.width*f.height)
	copy(data, f.header)
	data = append(data, f.footer...)

	for x := 0; x < f.width; x++ {
		for y := 0; y < f.height; y++ {
			offset := (x+y*f.width)*3 + len(f.header)
			col := f.image.RGBAAt(x, y)
			data[offset] = col.R
			data[offset+1] = col.G
			data[offset+2] = col.B
			// fmt.Println("set color", col)
		}
	}

	_, err := f.conn.Write(data)
	// fmt.Println(string(data))
	return err
}

func (f *Flaschen) Close() error {
	return f.conn.Close()
}
