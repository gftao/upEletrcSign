package textCvtImge

import (
	"image"
	"os"
	"bufio"
	"image/png"
	"fmt"
	"container/list"
	"io/ioutil"
	"strings"
	"github.com/golang/freetype/truetype"
	"image/jpeg"
	"log"
	"image/draw"
	"github.com/golang/freetype"
	"io"
)

// TextConvert 文本转化
type TextConvert struct {
	FontPath string
	Text     string
	rgba     *image.RGBA
	png      string
	Font     string
	FontSize map[int]float64
}

// textLine 单行数据
type textLine struct {
	Text  string
	Index int
}

func (tc *TextConvert) NewRGBA(image2 image.Image) {

}

func (tc *TextConvert) ResizeRGBA(image2 image.Image) *image.RGBA {
	i := image.NewRGBA(image.Rect(0, 0, tc.rgba.Bounds().Dx(), tc.rgba.Bounds().Dy()+image2.Bounds().Dx()))
	return i
}

// NewTextConvert 新建转换主体
func NewTextConvert(ft string) (*TextConvert, error) {
	_, err := os.Stat(ft)
	if err != nil {
		// 获取默认字体,打印输出错误
		return &TextConvert{}, err
	}
	return &TextConvert{Font: ft}, nil
}

// sliptString 字符串拆解
func (tc *TextConvert) sliptString(lines *list.List, lineSize int) (countAll int) {
	//log.Println(tc.Text)
	countAll = 0
	text := strings.Replace(tc.Text, "\t", "    ", 0)
	texts := strings.Split(text, "\n")

	for _, v := range texts {

		txt := []rune(v)
		// 处理回车和换行
		c := len(txt)
		var count int

		if n := c / lineSize; n <= 0 {
			count = 1
		} else {
			count = n + 1
		}

		for j := 0; j < count; j++ {
			start := j * lineSize
			end := (j + 1) * lineSize
			if j == count-1 {
				lines.PushBack(textLine{string(txt[start:]), countAll + j})
			} else {
				lines.PushBack(textLine{string(txt[start:end]), countAll + j})
			}
		}

		countAll += count
	}

	return
}

// wrap 字符换行处理
func (tc *TextConvert) wrap(lines *list.List, lineSize int) (countAll int) {
	text := WrapString(tc.Text, uint(lineSize))
	texts := strings.Split(text, "\n")
	countAll = len(texts)

	for k, v := range texts {
		fmt.Println(v)
		lines.PushBack(textLine{v, k})
	}
	return
}

// Write 实现 io.Writer 接口 写入文字实体
func (tc *TextConvert) Write(p []byte) (n int, err error) {
	tc.Text += string(p)
	return 0, nil
}

// EncodeImg 处理图片生成
func (tc *TextConvert) EncodeImg() *TextConvert {
	var size = 20.0
	var dx = 500
	var dy = 2000

	lineSize := int(dx/int(size)) * 4
	fmt.Println(lineSize)

	// 计算行数以及拆解字符串
	lines := list.New()
	//tl := tc.sliptString(lines, lineSize)
	tl := tc.wrap(lines, lineSize)
	dy = tl*int(size) + 20

	fontb, err := ioutil.ReadFile(tc.Font)
	if err != nil {
		log.Println(err)
	}
	fontf, err := truetype.Parse(fontb)
	if err != nil {
		log.Println(err)
	}

	fg, bg := image.Black, image.White
	tc.rgba = image.NewRGBA(image.Rect(0, 0, dx, dy))
	draw.Draw(tc.rgba, tc.rgba.Bounds(), bg, image.ZP, draw.Src)
	c := freetype.NewContext()
	c.SetDPI(72)
	c.SetFont(fontf)
	c.SetFontSize(size)
	c.SetClip(tc.rgba.Bounds())
	c.SetDst(tc.rgba)
	c.SetSrc(fg)
	//c.SetHinting(font.HintingVertical)

	opts := truetype.Options{}
	opts.Size = size
	//face := truetype.NewFace(f, &opts)
	//faceWidth,ok := face.GlyphAdvance(rune("中"))
	intSize := int(size)

	var txtNext *list.Element

	for txt := lines.Front(); txt != nil; txt = txtNext {
		txtNext = txt.Next()
		t := txt.Value.(textLine)
		pt := freetype.Pt(16, (t.Index+2)*intSize)
		c.SetFontSize(tc.FontSize[t.Index])
		//c.SetDPI(tc.DPI[t.Index])
		c.DrawString(t.Text, pt)
	}

	return tc
}

func (tc *TextConvert) WriteTo(w io.Writer) {
	//bf := bufio.NewWriter(w)
	//err := png.Encode(w, tc.rgba)
	err := jpeg.Encode(w, tc.rgba, &jpeg.Options{Quality: 80})
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
}

func (tc *TextConvert) SaveImg() {
	outFile, _ := os.Create("demo2.png")
	defer outFile.Close()
	bf := bufio.NewWriter(outFile)
	err := png.Encode(bf, tc.rgba)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	err = bf.Flush()
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
}
