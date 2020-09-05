package captcha

import (
	"github.com/golang/freetype"
	"github.com/golang/freetype/truetype"

	"image"
	"image/color"
	"image/draw"
	"io/ioutil"
	"math"
	"math/rand"
	"time"
)

type Captcha struct {
	frontColors []color.Color
	bkgColors   []color.Color
	disturlvl   DisturLevel
	fonts       []*truetype.Font
	size        image.Point
}

type StrType int

const (
	NUM   StrType = 0 // 数字
	LOWER         = 1 // 小写字母
	UPPER         = 2 // 大写字母
	ALL           = 3 // 全部
)

type DisturLevel int

const (
	NORMAL DisturLevel = 4
	MEDIUM             = 8
	HIGH               = 16
)

func New() *Captcha {
	c := &Captcha{
		disturlvl: NORMAL,
		size:      image.Point{82, 32},
	}
	c.frontColors = []color.Color{color.Black}
	c.bkgColors = []color.Color{color.White}
	return c
}

// AddFont 添加一个字体
func (c *Captcha) AddFont(path string) error {
	fontdata, erro := ioutil.ReadFile(path)
	if erro != nil {
		return erro
	}
	font, erro := freetype.ParseFont(fontdata)
	if erro != nil {
		return erro
	}
	if c.fonts == nil {
		c.fonts = []*truetype.Font{}
	}
	c.fonts = append(c.fonts, font)
	return nil
}

//AddFontFromBytes allows to load font from slice of bytes, for example, load the font packed by https://github.com/jteeuwen/go-bindata
func (c *Captcha) AddFontFromBytes(contents []byte) error {
	font, err := freetype.ParseFont(contents)
	if err != nil {
		return err
	}
	if c.fonts == nil {
		c.fonts = []*truetype.Font{}
	}
	c.fonts = append(c.fonts, font)
	return nil
}

// SetFont 设置字体 可以设置多个
func (c *Captcha) SetFont(paths ...string) error {
	for _, v := range paths {
		if erro := c.AddFont(v); erro != nil {
			return erro
		}
	}
	return nil
}

func (c *Captcha) SetDisturbance(d DisturLevel) {
	if d > 0 {
		c.disturlvl = d
	}
}

func (c *Captcha) SetFrontColor(colors ...color.Color) {
	if len(colors) > 0 {
		c.frontColors = c.frontColors[:0]
		for _, v := range colors {
			c.frontColors = append(c.frontColors, v)
		}
	}
}

func (c *Captcha) SetBkgColor(colors ...color.Color) {
	if len(colors) > 0 {
		c.bkgColors = c.bkgColors[:0]
		for _, v := range colors {
			c.bkgColors = append(c.bkgColors, v)
		}
	}
}

func (c *Captcha) SetSize(w, h int) {
	if w < 48 {
		w = 48
	}
	if h < 20 {
		h = 20
	}
	c.size = image.Point{w, h}
}

func (c *Captcha) randFont() *truetype.Font {
	return c.fonts[rand.Intn(len(c.fonts))]
}

// 绘制背景
func (c *Captcha) drawBkg(img *Image) {
	ra := rand.New(rand.NewSource(time.Now().UnixNano()))
	//填充主背景色
	bgcolorindex := ra.Intn(len(c.bkgColors))
	bkg := image.NewUniform(c.bkgColors[bgcolorindex])
	img.FillBkg(bkg)
}

// 绘制噪点
func (c *Captcha) drawNoises(img *Image) {
	ra := rand.New(rand.NewSource(time.Now().UnixNano()))

	// 待绘制图片的尺寸
	size := img.Bounds().Size()
	dlen := int(c.disturlvl)
	// 绘制干扰斑点
	for i := 0; i < dlen; i++ {
		x := ra.Intn(size.X)
		y := ra.Intn(size.Y)
		r := ra.Intn(size.Y/20) + 1
		colorindex := ra.Intn(len(c.frontColors))
		img.DrawCircle(x, y, r, i%4 != 0, c.frontColors[colorindex])
	}

	// 绘制干扰线
	for i := 0; i < dlen; i++ {
		x := ra.Intn(size.X)
		y := ra.Intn(size.Y)
		o := int(math.Pow(-1, float64(i)))
		w := ra.Intn(size.Y) * o
		h := ra.Intn(size.Y/10) * o
		colorindex := ra.Intn(len(c.frontColors))
		img.DrawLine(x, y, x+w, y+h, c.frontColors[colorindex])
		colorindex++
	}

}

// 绘制文字
func (c *Captcha) drawString(img *Image, str string) {

	if c.fonts == nil {
		panic("没有设置任何字体")
	}
	tmp := NewImage(c.size.X, c.size.Y)

	// 文字大小为图片高度的 0.6
	fsize := int(float64(c.size.Y) * 0.6)
	// 用于生成随机角度
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	// 文字之间的距离
	// 左右各留文字的1/4大小为内部边距
	padding := fsize / 4
	gap := (c.size.X - padding*2) / (len(str))

	// 逐个绘制文字到图片上
	for i, char := range str {
		// 创建单个文字图片
		// 以文字为尺寸创建正方形的图形
		str := NewImage(fsize, fsize)
		// str.FillBkg(image.NewUniform(color.Black))
		// 随机取一个前景色
		colorindex := r.Intn(len(c.frontColors))

		//随机取一个字体
		font := c.randFont()
		str.DrawString(font, c.frontColors[colorindex], string(char), float64(fsize))

		// 转换角度后的文字图形
		rs := str.Rotate(float64(r.Intn(40) - 20))
		// 计算文字位置
		s := rs.Bounds().Size()
		left := i*gap + padding
		top := (c.size.Y - s.Y) / 2
		// 绘制到图片上
		draw.Draw(tmp, image.Rect(left, top, left+s.X, top+s.Y), rs, image.ZP, draw.Over)
	}
	if c.size.Y >= 48 {
		// 高度大于48添加波纹 小于48波纹影响用户识别
		tmp.distortTo(float64(fsize)/10, 200.0)
	}

	draw.Draw(img, tmp.Bounds(), tmp, image.ZP, draw.Over)
}

// Create 生成一个验证码图片
func (c *Captcha) Create(num int, t StrType) (*Image, string) {
	if num <= 0 {
		num = 4
	}
	dst := NewImage(c.size.X, c.size.Y)
	//tmp := NewImage(c.size.X, c.size.Y)
	c.drawBkg(dst)
	c.drawNoises(dst)

	str := string(c.randStr(num, int(t)))
	c.drawString(dst, str)
	//c.drawString(tmp, str)

	return dst, str
}

func (c *Captcha) CreateCustom(str string) *Image {
	if len(str) == 0 {
		str = "unkown"
	}
	dst := NewImage(c.size.X, c.size.Y)
	c.drawBkg(dst)
	c.drawNoises(dst)
	c.drawString(dst, str)
	return dst
}

var fontKinds = [][]int{[]int{10, 48}, []int{26, 97}, []int{26, 65}}

// 生成随机字符串
// size 个数 kind 模式
func (c *Captcha) randStr(size int, kind int) []byte {
	ikind, result := kind, make([]byte, size)
	isAll := kind > 2 || kind < 0
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < size; i++ {
		if isAll {
			ikind = rand.Intn(3)
		}
		scope, base := fontKinds[ikind][0], fontKinds[ikind][1]
		result[i] = uint8(base + rand.Intn(scope))
	}
	return result
}
