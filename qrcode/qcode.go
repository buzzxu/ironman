package qrcode

import (
	"encoding/base64"
	"text/template"
	"strconv"
	"net/http"
)
var ImageTemplate string = `<!DOCTYPE html>
<html lang="en"><head></head>
<body><img src="data:image/jpg;base64,{{.Image}}"></body>`

//HTML 生成HTML片段
func HTML(w http.ResponseWriter,param *QRParam) error  {
	val,err := String(param)
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
func Image(w http.ResponseWriter,param *QRParam) error {
	bytes,err := QrCode(param)
	if err != nil {
		return err
	}
	w.Header().Set("Content-Type", "image/png")
	w.Header().Set("Content-Length", strconv.Itoa(len(bytes)))
	if _, err := w.Write(bytes); err != nil {
		return err
	}
	return nil

}
//String 生成Base64
func String(param *QRParam)(string,error)  {
	bytes,err:=QrCode(param)
	if err != nil {
		return "",err
	}
	return base64.StdEncoding.EncodeToString(bytes),nil
}



