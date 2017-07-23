package qrcode

import (
	"image"
	"image/png"
	"github.com/buzzxu/go-qrcode"
	"bytes"
)

type QRDiy struct {
	Arg QRArg
}

func (q *QRDiy) Encode() ([]byte, error) {
	var code *qrcode.QRCode
	//	code, err := qrcode.New(q.Arg.content, q.Arg.level)
	code, err := qrcode.NewWithColor(q.Arg.Content, q.Arg.level, q.Arg.bgcolor, q.Arg.forecolor)
	if err != nil {
		return nil, err
	}
	var img image.Image
	if q.Arg.bdmaxsize <= 0 {
		img = code.Image(q.Arg.size)
	} else {
		img = code.ImageWithBorderMaxSize(q.Arg.size, q.Arg.bdmaxsize)
	}

	if q.Arg.bgimg != nil {
		q.embgimg(img, q.Arg.bgimg)
	}
	if q.Arg.logo != nil {
		q.emlogo(img, q.Arg.logo)
	}
	var b bytes.Buffer
	err = png.Encode(&b, img)
	if err != nil {
		return nil, err
	}
	buf := b.Bytes()
	return buf, nil
}

func (q *QRDiy) emlogo(rst, logo image.Image) {
	offset := rst.Bounds().Max.X/2 - logo.Bounds().Max.X/2
	for x := 0; x < logo.Bounds().Max.X; x++ {
		for y := 0; y < logo.Bounds().Max.Y; y++ {
			rst.(*image.RGBA).Set(x+offset, y+offset, logo.At(x, y))
		}
	}
	return
}
func (q *QRDiy) embgimg(rst, bgimg image.Image) {
	if rst.Bounds().Max.X > q.Arg.size {
		return
	}
	qsx, qsy := 0, 0
	br, bg, bb, _ := q.Arg.bgcolor.RGBA()
	fr, fg, fb, _ := q.Arg.forecolor.RGBA()
	qex, qey := 0, 0
	oks, oke := false, false
	for z := 0; z < rst.Bounds().Max.X; z++ {
		cs := rst.(*image.RGBA).At(z, z)
		ce := rst.(*image.RGBA).At(rst.Bounds().Max.X-1-z, z)
		r, g, b, _ := cs.RGBA()
		if r == fr && g == fg && b == fb && !oks {
			qsx, qsy = z, z
			oks = true
		}
		r, g, b, _ = ce.RGBA()
		if r == fr && g == fg && b == fb && !oke {
			qex, qey = rst.Bounds().Max.X-1-z, rst.Bounds().Max.Y-1-z
			oke = true
		}
		if oks && oke {
			break
		}
	}

	for x := 0; x < rst.Bounds().Max.X; x++ {
		for y := 0; y < rst.Bounds().Max.Y; y++ {
			if x < qsx || y < qsy || x > qex || y > qey {
				rst.(*image.RGBA).Set(x, y, bgimg.At(x, y))
			} else {
				//				r, g, b, _ := rst.(*image.RGBA).At(x, y).RGBA()
				//				if r == fr && g == fg && b == fb {
				//					rst.(*image.RGBA).Set(x, y, bgimg.At(x, y))
				//				}
				r, g, b, _ := rst.(*image.RGBA).At(x, y).RGBA()
				if r == br && g == bg && b == bb {
					rst.(*image.RGBA).Set(x, y, bgimg.At(x, y))
				}
			}
		}
		if x >= qsx-2 && x <= qex+2 {
			rst.(*image.RGBA).Set(x, qsy-1, q.Arg.bgcolor)
			rst.(*image.RGBA).Set(x, qey+1, q.Arg.bgcolor)
			rst.(*image.RGBA).Set(x, qsy-2, q.Arg.bgcolor)
			rst.(*image.RGBA).Set(x, qey+2, q.Arg.bgcolor)
		}
	}

	for y := qsy - 2; y <= qey+2; y++ {
		rst.(*image.RGBA).Set(qsx-1, y, q.Arg.bgcolor)
		rst.(*image.RGBA).Set(qex+1, y, q.Arg.bgcolor)
		rst.(*image.RGBA).Set(qsx-2, y, q.Arg.bgcolor)
		rst.(*image.RGBA).Set(qex+2, y, q.Arg.bgcolor)
	}
	return
}
