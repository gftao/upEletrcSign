package trans

import (
	"crypto/rsa"
	"mygolib/modules/config"
	"errors"
	"github.com/golang/freetype/truetype"
	"io/ioutil"
	"fmt"
	"github.com/golang/freetype"
	"bytes"
	"image"
	"golang.org/x/image/bmp"
	"github.com/nfnt/resize"
	"github.com/gwpp/tinify-go/tinify"
)

const TermKeyOutTime = 100

type globArgv struct {
	PubKey   *rsa.PublicKey
	PriKey   *rsa.PrivateKey
	FontPath string
	PicPath  string

	TFont *truetype.Font
	Logo  image.Image
}

var GlobA *globArgv

func InitArgv() error {
	if !config.HasConfigInit() {
		return errors.New("配置文件未初始化，请先初始化")
	}

	GlobA = new(globArgv)
	config.SetSection("glob")
	fp, ok := config.String("fontPath")
	if !ok {
		return errors.New("FontPath 未配置")
	}

	pp, ok := config.String("picPath")
	if !ok {
		return errors.New("PicPath 未配置")
	}
	GlobA.FontPath = fp
	GlobA.PicPath = pp
	fontBytes, err := ioutil.ReadFile(fp)
	if err != nil {
		return fmt.Errorf("%s", err)
	}
	font, err := freetype.ParseFont(fontBytes)
	if err != nil {
		return fmt.Errorf("%s", err)
	}
	GlobA.TFont = font

	lp, ok := config.String("logo")
	if !ok {
		return errors.New("FontPath 未配置")
	}
	logoBytes, err := ioutil.ReadFile(lp)
	if err != nil {
		return fmt.Errorf("%s", err)
	}
	//img, s, err := image.Decode(bytes.NewBuffer(logoBytes))
	//fmt.Println(img, s, err)
	img, err := bmp.Decode(bytes.NewBuffer(logoBytes))
	if err != nil {
		return fmt.Errorf("%s", err)
	}
	m := resize.Resize(220, 50, img, resize.Lanczos3)
	GlobA.Logo = m
	Tinify.SetKey("xz9CbfmSEgaVIkvj8NDTvbJ1SqHjK5Uv")
	return nil
}
