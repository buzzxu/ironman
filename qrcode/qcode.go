package qrcode

import (
	"github.com/boombuler/barcode"
	"github.com/boombuler/barcode/qr"
	"image/png"
	"bytes"
	"encoding/base64"
	"text/template"
	"strconv"
	"net/http"
)
var ImageTemplate string = `<!DOCTYPE html>
<html lang="en"><head></head>
<body><img src="data:image/jpg;base64,{{.Image}}"></body>`

//HTML 生成HTML片段
func HTML(w http.ResponseWriter,content string,width, height int) error  {
	val,err := String(content,width,height)
	if err != nil {
		return err
	}
	if tmpl, err := template.New("image").Parse(ImageTemplate); err != nil {
		return err
	} else {
		data := map[string]interface{}{"Image": val}
		if err = tmpl.Execute(w, data); err != nil {
			return err
		}
		return nil
	}
}
//Image 直接生成图片
func Image(w http.ResponseWriter,content string,width, height int) error {
	qrcode,err := QrCode(content,width,height)
	if err != nil {
		return err
	}
	buffer := new(bytes.Buffer)
	if err := png.Encode(buffer, qrcode); err != nil {
		return err
	}
	w.Header().Set("Content-Type", "image/png")
	w.Header().Set("Content-Length", strconv.Itoa(len(buffer.Bytes())))
	if _, err := w.Write(buffer.Bytes()); err != nil {
		return err
	}
	return nil

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

