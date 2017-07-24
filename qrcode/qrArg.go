package qrcode

import (
	"image"
	_ "image/jpeg"
	_ "image/png"
	"net/http"
	"github.com/buzzxu/go-qrcode"
	"strconv"
	"image/color"
)

type QRParam struct {
	Content string
	Size int
	BgMaxSize int
	BgColor string
	ForeColor string
	Logo string
	BgImg string
}
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


//parse 解析参数
func  parse(param *QRParam) *QRArg {
	arg := &QRArg{}
	arg.Content = param.Content
	arg.parseSize(param.Size)
	arg.parseBdmaxsize(param.BgMaxSize)
	arg.parseBGColor(param.BgColor)
	arg.parseForeColor(param.ForeColor)
	arg.parseLogo(param.Logo)
	arg.parseBGImg(param.BgImg)
	if arg.logo == nil {
		arg.level = qrcode.Highest
	}else{
		arg.level = qrcode.Medium
	}

	if arg.bgimg != nil {
		if arg.bgimg.Bounds().Max.X > arg.bgimg.Bounds().Max.Y {
			arg.size = arg.bgimg.Bounds().Max.Y
		} else {
			arg.size = arg.bgimg.Bounds().Max.X
		}
		//		q.level = qrcode.Highest
	}
	return arg
}

func (q *QRArg) parseSize(size int)  {

	if size != 0 {
		q.size = size
	}else{
		q.size = 256
	}
}

func (q *QRArg) parseBdmaxsize(size int)  {

	if size != 0 {

		q.bdmaxsize = size
	}else{
		q.bdmaxsize = -1
	}
}

func (q *QRArg) parseBGColor(str string)  {
	s, err := strconv.ParseUint(str, 16, 32)
	if err != nil {
		q.bgcolor = color.White
	}else{
		q.bgcolor = color.RGBA{R: uint8(s & 0xff0000 >> 16),
			G: uint8(s & 0xff00 >> 8),
			B: uint8(s & 0xff),
			A: uint8(0xff)}
	}

}

func (q *QRArg) parseForeColor(str string)  {
	s, err := strconv.ParseUint(str, 16, 32)
	if err != nil {
		q.forecolor =  color.Black
	}else{
		q.forecolor = color.RGBA{R: uint8(s & 0xff0000 >> 16),
			G: uint8(s & 0xff00 >> 8),
			B: uint8(s & 0xff),
			A: uint8(uint8(0xff))}
	}

}

func (q *QRArg) parseLogo(str string) {
	if len(str) == 0 {
		q.logo = nil
	}else{
		q.logo = downImg(str)
	}
}

func (q *QRArg) parseBGImg(str string)  {
	if len(str) == 0 {
		q.bgimg = nil
	}else{
		q.bgimg = downImg(str)
	}
}

func  downImg(str string) image.Image {
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
