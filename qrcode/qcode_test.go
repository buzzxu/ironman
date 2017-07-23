package qrcode

import (
	"testing"
	"github.com/boombuler/barcode/qr"
	"github.com/boombuler/barcode"
	"encoding/base64"
	"image/png"
	"bytes"
	"fmt"
)

func TestBase64(t *testing.T)  {


	qrCode, _ := qr.Encode("Hello World", qr.H, qr.Auto)

	qrCode, _ = barcode.Scale(qrCode, 200, 200)

	emptyBuff := bytes.NewBuffer(nil)
	png.Encode(emptyBuff,qrCode)
	fmt.Println(base64.StdEncoding.EncodeToString(emptyBuff.Bytes()))
}
