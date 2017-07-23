package qrcode

import (
	"github.com/boombuler/barcode"
	"github.com/boombuler/barcode/qr"
	"image/png"
	"bytes"
	"encoding/base64"
	"text/template"
	"github.com/labstack/echo"
	"strconv"
)
var ImageTemplate string = `<!DOCTYPE html>
<html lang="en"><head></head>
<body><img src="data:image/jpg;base64,{{.Image}}"></body>`

//HTML 生成HTML片段
func HTML(c echo.Context,content string,width, height int)  {
	val,err := String(content,width,height)
	if err != nil {
		c.Logger().Error(err)
	}

	if tmpl, err := template.New("image").Parse(ImageTemplate); err != nil {
		c.Logger().Error("unable to parse image template.")
	} else {
		data := map[string]interface{}{"Image": val}
		if err = tmpl.Execute(c.Response().Writer, data); err != nil {
			c.Logger().Error("unable to execute template.")
		}
	}
}
//Image 直接生成图片
func Image(c echo.Context,content string,width, height int)  {
	qrcode,err := QrCode(content,width,height)
	if err != nil {
		c.Logger().Error(err)
	}
	buffer := new(bytes.Buffer)
	if err := png.Encode(buffer, qrcode); err != nil {
		c.Logger().Error("unable to encode image.")
	}
	c.Response().Header().Set("Content-Length", strconv.Itoa(len(buffer.Bytes())))
	c.Blob(200,"image/png",buffer.Bytes())

}
//String 生成Base64
func String(content string,width, height int)(string,error)  {
	qrCode,err:=QrCode(content,width,height)
	if err != nil {
		return "",err
	}
	buffer := bytes.NewBuffer(nil)
	png.Encode(buffer,qrCode)
	return base64.StdEncoding.EncodeToString(buffer.Bytes()),nil
}

//qrcode 生成二维码
func QrCode(content string,width, height int)(barcode.Barcode,error)  {
	qrCode, err := qr.Encode(content, qr.H, qr.Auto)
	if err!= nil {
		return nil,err
	}
	qrCode, err = barcode.Scale(qrCode, width, height)
	if err!= nil {
		return nil,err
	}
	return qrCode,err
}

