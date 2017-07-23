package qrcode

import (
	"testing"
	"github.com/boombuler/barcode/qr"
	"github.com/boombuler/barcode"
	"encoding/base64"
	"image/png"
	"bytes"
	"fmt"
	"os"
)

func TestBase64(t *testing.T)  {


	qrCode, _ := qr.Encode("Hello World", qr.Q, qr.Auto)

	qrCode, _ = barcode.Scale(qrCode, 200, 200)

	emptyBuff := bytes.NewBuffer(nil)
	png.Encode(emptyBuff,qrCode)
	fmt.Println(base64.StdEncoding.EncodeToString(emptyBuff.Bytes()))
}

func TestImg(t *testing.T)  {


	qrCode, _ := qr.Encode("http://www.baidu.com/com/fdf/fdh.html?a=1", qr.M, qr.Auto)

	// Scale the barcode to 200x200 pixels
	qrCode, _ = barcode.Scale(qrCode, 200, 200)

	// create the output file
	file, _ := os.Create("qrcode.png")
	defer file.Close()

	// encode the barcode as png
	png.Encode(file, qrCode)

}