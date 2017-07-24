package qrcode

import (
	"testing"
	"fmt"
	"os"
	"io/ioutil"
	"github.com/labstack/gommon/color"
)

func TestBase64(t *testing.T)  {

	param := &QRParam{
		Content:"http://www.google.com",
		Size:256,
		BgColor:color.Blu,
		Logo:"https://static.oschina.net/uploads/user/410/820033_50.jpg?t=1350367163000",
		BgMaxSize:1,
	}
	qrCode, _ := String(param)

	fmt.Println(qrCode)
}

func TestImg(t *testing.T)  {


	param := &QRParam{
		Content:"http://www.google.com",
		Size:256,
		BgColor:color.Blu,
		Logo:"https://static.oschina.net/uploads/user/410/820033_50.jpg?t=1350367163000",
		BgMaxSize:1,
	}
	qrCode,_ := QrCode(param)

	ioutil.WriteFile("qrcode.png",qrCode,os.FileMode(0644))


}