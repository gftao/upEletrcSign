package iServer

import (
	"upEletrcSign/trans"
	"mygolib/gerror"
	"encoding/json"
	"mygolib/modules/myLogger"
	"github.com/vdobler/chart/imgg"
	"image/color"
	"github.com/vdobler/chart"
	"strconv"
	"image/png"
	"os"
	"github.com/boombuler/barcode/code128"
	"github.com/boombuler/barcode"
	"image"
	"image/draw"
	"bytes"
	"github.com/boombuler/barcode/qr"
)

type T8262 struct {
	reqMsg  *trans.TransMessage
	Poses   *trans.Spos
	Font    string
	picPath string
	logo    image.Image
}

func (t *T8262) Init() gerror.IError {
	t.Font = trans.GlobA.FontPath
	t.picPath = trans.GlobA.PicPath
	t.logo = trans.GlobA.Logo
	myLogger.Info("FontPath:", t.Font)
	myLogger.Info("PicPath: ", t.picPath)

	return nil
}

func (t *T8262) DoTrans(msg *trans.TransMessage) (gerror.IError) {
	t.reqMsg = msg
	//t.Poses = make(trans.Spos, 10)
	t.Poses = &trans.Spos{}
	err := json.Unmarshal([]byte(msg.MsgBody.UPos_sign), &t.Poses)
	if err != nil {
		return gerror.NewR(2001, err, "解析Poses失败")
	}
	myLogger.Debugf("Poses: %+v", len(t.Poses.Pos))
	t.picPath = t.picPath + msg.MsgBody.Orig_sys_order_id[:8] + "/"
	exist := t.FileIfExit()
	if !exist {
		err := os.MkdirAll(t.picPath, os.ModePerm)
		if err != nil {
			return gerror.NewR(2009, err, "MkdirAll失败：%s", t.picPath)
		}
	}
	t.picPath = t.picPath + msg.MsgBody.Orig_sys_order_id + ".png"
	myLogger.Infoln("文件是路径: ", t.picPath)

	exist = t.FileIfExit()
	myLogger.Infoln("文件是否已存在: ", exist)
	if exist {
		return nil
	}
	//white := color.RGBA{220, 220, 220, 255}
	white := color.RGBA{255, 255, 255, 255}
	height := 0
	for _, v := range t.Poses.Pos {
		if v != nil {
			switch v.Content_type {
			case "txt":
				height += 18
			case "one-dimension":
				height += 28
			case "two-dimension":
				height += 50
			case "jpg":
				if v.Content == "<%pic_type:logo%>" {
					height += 45
				} else if v.Content == "<%pic_type:sign%>" {
					height += 70
				} else {
					height += 28
				}
			default:
			}
		}
	}
	myLogger.Debugf("height : %d", height)
	g := imgg.New(280, height, white, trans.GlobA.TFont, imgg.ConstructFontSizes(12))
	x, y := 10, 2
	for _, v := range t.Poses.Pos {
		if v != nil {
			myLogger.Debugf(">>%+v", *v)

			myLogger.Debugln(x, y)
			switch v.Content_type {
			case "txt":
				s, _ := strconv.Atoi(v.Size)
				x, al := checkPosition(v.Position, x)
				s -= 2
				g.Text(x, y, v.Content, al, 0, chart.Font{Size: chart.FontSize(s), Color: color.Black})
			case "one-dimension":
				qrCode, _ := code128.Encode(v.Content)
				qr, _ := barcode.Scale(qrCode, 220, 30)
				//qr := resize.Resize(100, 30, qrCode, resize.Lanczos3)
				Y := qr.Bounds().Dy()
				X := qr.Bounds().Dx()
				//fmt.Println("one:", X, Y)
				draw.Draw(g.Image, image.Rect(30, y-2, X, y+Y), qr, image.ZP, draw.Src)
				y += 12
			case "two-dimension":
				qrCode, _ := qr.Encode(v.Content, qr.M, qr.Auto)
				br, _ := barcode.Scale(qrCode, 100, 100)
				Y := br.Bounds().Dy()
				X := br.Bounds().Dx()
				//fmt.Println("one:", Y)
				draw.Draw(g.Image, image.Rect(10, y-2, X, y+Y), br, image.ZP, draw.Src)
				y += 20
			case "jpg":
				if v.Content == "<%pic_type:logo%>" {
					si := t.logo.Bounds()
					sp := si.Size()
					//fmt.Println("logo:", sp.Y)
					draw.Draw(g.Image, image.Rect(28, y, sp.X+40, y+sp.Y), t.logo, image.ZP, draw.Src)
					y = y + sp.Y - 26
				}else if v.Content == "<%pic_type:sign%>" {
					m, err := t.decode()
					if err != nil {
						return err
					}
					si := m.Bounds()
					sp := si.Size()
					//fmt.Println("sign:", sp.Y)
					draw.Draw(g.Image, image.Rect(20, y-18, sp.X, y+sp.Y), m, image.ZP, draw.Src)
					y = y + sp.Y - 26
				}
			default:
				continue

			}
			y = y + 18
		}
	}
	file, err := os.Create(t.picPath)
	if err != nil {
		return gerror.NewR(2008, err, "Create file failed.")
	}
	defer file.Close()

	png.Encode(file, g.Image)

	return nil
}

func checkPosition(pos string, x int) (int, string) {
	al := ""
	if pos == "left" {
		x = 20
		al = "tl"
	} else if pos == "center" {
		x = 140
		al = "tc"
	} else if pos == "right" {
		x = 120
		al = "tl"
	}
	return x, al
}

func (t *T8262) decode() (image.Image, gerror.IError) {

	s, err := trans.DecodeBase64([]byte(t.reqMsg.MsgBody.Sign_img))
	if err != nil {
		return nil, gerror.NewR(2005, err, "解析Sign_img失败")
	}
	im, err := png.Decode(bytes.NewBufferString(s))
	if err != nil {
		return nil, gerror.NewR(2006, err, "jpeg Decode失败")
	}
	t.reqMsg.MsgBody.Pos_sign = ""
	t.reqMsg.MsgBody.Sign_img = ""
	return im, nil
}

func (t *T8262) FileIfExit() bool {
	_, err := os.Lstat(t.picPath)
	if err == nil {
		return true
	}
	return false
}
