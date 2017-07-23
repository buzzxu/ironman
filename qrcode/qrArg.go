package qrcode

import (
	"image"
	"net/http"
	"github.com/buzzxu/go-qrcode"
	"strconv"
	"image/color"
)

type QRArg struct {
	Content   string
	size      int
	bgcolor   color.Color
	forecolor color.Color
	logo      image.Image
	level     qrcode.RecoveryLevel
	bgimg     image.Image
	bdmaxsize int
}

type Values interface {
	Param(string) string
}

func (q *QRArg) Parse(query Values) {
	q.Content = query.Param("content")
	q.size = q.parseSize(query.Param("size"))
	q.bgcolor = q.parseBGColor(query.Param("bgcolor"))
	q.forecolor = q.parseForeColor(query.Param("forecolor"))
	q.logo = q.parseLogo(query.Param("logo"))
	q.bgimg = q.parseBGImg(query.Param("bgimg"))
	q.bdmaxsize = q.parseBdmaxsize(query.Param("bdmaxsize"))
	if q.logo == nil {
		q.level = qrcode.Highest
	}else{
		q.level = qrcode.Medium
	}

	if q.bgimg != nil {
		if q.bgimg.Bounds().Max.X > q.bgimg.Bounds().Max.Y {
			q.size = q.bgimg.Bounds().Max.Y
		} else {
			q.size = q.bgimg.Bounds().Max.X
		}
		//		q.level = qrcode.Highest
	}
}

func (q *QRArg) parseSize(str string) int {
	var size int
	if str != "" {
		s, err := strconv.Atoi(str)
		if err != nil {
			size = 256
		}else {
			size = s
		}
	}else{
		size = 256
	}
	return size
}

func (q *QRArg) parseBdmaxsize(str string) int {
	size := -1
	if str != "" {
		s, err := strconv.Atoi(str)
		if err != nil {
			size = -1
		}
		size = s
	}
	return size
}

func (q *QRArg) parseBGColor(str string) color.Color {
	s, err := strconv.ParseUint(str, 16, 32)
	if err != nil {
		return color.White
	}
	return color.RGBA{R: uint8(s & 0xff0000 >> 16),
		G: uint8(s & 0xff00 >> 8),
		B: uint8(s & 0xff),
		A: uint8(0xff)}
}

func (q *QRArg) parseForeColor(str string) color.Color {
	s, err := strconv.ParseUint(str, 16, 32)
	if err != nil {
		return color.Black
	}
	return color.RGBA{R: uint8(s & 0xff0000 >> 16),
		G: uint8(s & 0xff00 >> 8),
		B: uint8(s & 0xff),
		A: uint8(uint8(0xff))}
}

func (q *QRArg) parseLogo(str string) image.Image {
	if len(str) == 0 {
		return nil
	}
	return q.downImg(str)
}

func (q *QRArg) parseBGImg(str string) image.Image {
	if len(str) == 0 {
		return nil
	}
	return q.downImg(str)
}

func (q *QRArg) downImg(str string) image.Image {
	resp, err := http.Get(str)
	if err != nil {
		return nil
	}
	defer resp.Body.Close()
	logo, _, err := image.Decode(resp.Body)
	if err != nil {
		return nil
	}
	return logo
}
