package main

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/golang/freetype"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
	"image"
	"image/color"
	"image/jpeg"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"
	"time"
)

const (
	dx     = 473
	dy     = 315
	heigth = 35
	width  = 160
)

type Job struct {
	vArr       [7]string
	printerNum string
}

var JobChan = make(chan *Job, 100)

func task() {

	fileName := "printer.jpeg"
	for {
		select {
		case job := <-JobChan:
			saveImg(fileName, job.vArr)
			StartJob(fileName, job.printerNum)
			os.Remove(fileName)
		default:
			time.Sleep(time.Millisecond * 1000)
		}

	}

}

func saveImg(fileName string, vArr [7]string) {
	file, err := os.Create(fileName)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	img := image.NewRGBA(image.Rect(0, 0, dx, dy))
	i := 0
	flag := false
	//设置每个点的 RGBA (Red,Green,Blue,Alpha(设置透明度))
	for y := 0; y < dy; y++ {
		flag = false
		for x := 0; x < dx; x++ {
			if x <= 2 || y <= 2 || x >= dx-3 || y >= dy-3 {
				//黑色边框
				img.Set(x, y, color.RGBA{0, 0, 0, 255})
				continue
			}

			//竖线
			if y >= 40 && x == width && i < 7 {
				//黑色边框
				img.Set(x, y, color.RGBA{0, 0, 0, 255})
				continue
			}

			//横线
			if y == 40 {
				flag = true
				//黑色边框
				img.Set(x, y, color.RGBA{0, 0, 0, 255})
				continue
			}

			//横线
			if y > 40 && y%heigth == 0 && i < 7 {
				flag = true
				//黑色边框
				img.Set(x, y, color.RGBA{0, 0, 0, 255})
				continue
			}

			//设置一块 白色(255,255,255)不透明的背景
			img.Set(x, y, color.RGBA{255, 255, 255, 255})

		}

		if flag {
			i++
		}
	}

	//读取字体数据
	fontBytes, err := ioutil.ReadFile("SIMLI.TTF")
	if err != nil {
		log.Println(err)
	}
	//载入字体数据
	font, err := freetype.ParseFont(fontBytes)
	if err != nil {
		log.Println("load front fail", err)
	}
	f := freetype.NewContext()
	//设置分辨率
	f.SetDPI(72)
	//设置字体
	f.SetFont(font)
	//设置尺寸
	f.SetFontSize(34)
	f.SetClip(img.Bounds())
	//设置输出的图片
	f.SetDst(img)
	//设置字体颜色(红色)
	f.SetSrc(image.NewUniform(color.RGBA{0, 0, 0, 255}))

	//设置字体的位置
	pt := freetype.Pt(165, 25+int(f.PointToFixed(34))>>8)

	f.DrawString("送 样 单", pt)

	f.SetFontSize(26)

	//设置字体的位置
	pt = freetype.Pt(6, 60+int(f.PointToFixed(26))>>8)
	f.DrawString("样品名称", pt)

	titleArr := [5]string{"规格型号", "数量", "品牌", "送样人", "送样日期"}

	for k, v := range titleArr {
		//设置字体的位置
		pt = freetype.Pt(6, 35*(k+1)+60+int(f.PointToFixed(26))>>8)
		f.DrawString(v, pt)
	}

	name, _ := Utf8ToGbk([]byte(vArr[0]))

	//设置字体的位置
	pt = freetype.Pt(160+(157-strings.Count(string(name), "")*13/2), 60+int(f.PointToFixed(26))>>8)
	f.DrawString(vArr[0], pt)

	for k, v := range vArr[1:6] {

		bt, _ := Utf8ToGbk([]byte(v))

		num := strings.Count(string(bt), "")
		pt = freetype.Pt(160+(157-num*13/2), 35*(k+1)+60+int(f.PointToFixed(26))>>8)
		f.DrawString(v, pt)
	}

	//设置字体的位置
	pt = freetype.Pt(6, 35*6+60+int(f.PointToFixed(26))>>8)
	f.DrawString("备注: "+vArr[6], pt)
	jpeg.Encode(file, img, nil) //将image信息写入文件中
}

func GbkToUtf8(s []byte) ([]byte, error) {
	reader := transform.NewReader(bytes.NewReader(s), simplifiedchinese.GBK.NewDecoder())
	d, e := ioutil.ReadAll(reader)
	if e != nil {
		return nil, e
	}
	return d, nil
}

func Utf8ToGbk(s []byte) ([]byte, error) {
	reader := transform.NewReader(bytes.NewReader(s), simplifiedchinese.GBK.NewEncoder())
	d, e := ioutil.ReadAll(reader)
	if e != nil {
		return nil, e
	}
	return d, nil
}

func StartJob(path, num string) {

	for {
		resultStr := run(path, num)

		if resultStr == "" {
			fmt.Println("打印失败!")
			break
		}

		arr := strings.Split(resultStr, "\n")

		if strings.Contains(arr[len(arr)-2], "please try again slightly") {
			fmt.Println("打印机有任务未完成，重试中!")
			time.Sleep(time.Second * 1)
			continue
		}
		break
	}

}

//开始执行命令
func run(path, num string) string {

	cmd0 := exec.Command("cmd", "/C", "printer.exe", path, strings.Replace(num, "\n", "", -1))
	stdout0, err := cmd0.StdoutPipe() // 获取命令输出内容
	if err != nil {
		fmt.Println(err)
		return ""
	}
	if err := cmd0.Start(); err != nil { //开始执行命令
		fmt.Println(err)
		return ""
	}

	useBufferIO := false
	if !useBufferIO {
		var outputBuf0 bytes.Buffer
		for {
			tempoutput := make([]byte, 256)
			n, err := stdout0.Read(tempoutput)
			if err != nil {
				if err == io.EOF { //读取到内容的最后位置
					break
				} else {
					fmt.Println(err)
					return ""
				}
			}

			if n > 0 {
				outputBuf0.Write(tempoutput[:n])
			}

		}

		return outputBuf0.String()
	} else {
		outputbuf0 := bufio.NewReader(stdout0)
		touput0, _, err := outputbuf0.ReadLine()
		if err != nil {
			return ""
		}

		return string(touput0)
	}

}
