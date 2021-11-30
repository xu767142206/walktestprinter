package main

import (
	"fmt"
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
	"strconv"
)

type Printer struct {
	Name    string
	Spec    string
	Num     string
	Brand   string
	Pepole  string
	Date    string
	Desc    string
	SendNum int
}

func main() {

	go task()

	print := new(Printer)

	var name, spec, brand, people, date, desc, num *walk.LineEdit

	var sendNum *walk.NumberEdit
	//var wv *walk.WebView
	//var mw *walk.MainWindow

	var lebs *walk.Label

	MainWindow{
		Title: "打印送样单",
		// 指定窗口的大小
		Size: Size{Width: 300, Height: 230},

		Layout: Grid{
			Columns:   4,
			Alignment: AlignHCenterVCenter,
		},
		Children: []Widget{
			Label{
				Text:       "样品名称:",
				ColumnSpan: 1,
				Alignment:  AlignHCenterVCenter,
			},
			LineEdit{AssignTo: &name, Text: print.Name, CueBanner: "请输入样品名称", TextAlignment: AlignCenter, ColumnSpan: 3},

			Label{
				ColumnSpan: 4,
			},

			Label{
				Text:       "规格型号:",
				ColumnSpan: 1,
				Alignment:  AlignHCenterVCenter,
			},
			LineEdit{AssignTo: &spec, Text: print.Spec, CueBanner: "请输入规格型号", TextAlignment: AlignCenter, ColumnSpan: 3},

			Label{
				ColumnSpan: 4,
			},

			Label{
				Text:       "数量:",
				ColumnSpan: 1,
				Alignment:  AlignHCenterVCenter,
			},
			LineEdit{AssignTo: &num, Text: print.Num, CueBanner: "请输入数量", TextAlignment: AlignCenter, ColumnSpan: 3},

			Label{
				ColumnSpan: 4,
			},

			Label{
				Text:       "品牌:",
				ColumnSpan: 1,
				Alignment:  AlignHCenterVCenter,
			},
			LineEdit{AssignTo: &brand, CueBanner: "请输入品牌", TextAlignment: AlignCenter, ColumnSpan: 3},

			Label{
				ColumnSpan: 4,
			},

			Label{
				Text:       "送样人:",
				ColumnSpan: 1,
				Alignment:  AlignHCenterVCenter,
			},
			LineEdit{AssignTo: &people, CueBanner: "请输入送样人", TextAlignment: AlignCenter, ColumnSpan: 3},

			Label{
				ColumnSpan: 4,
			},

			Label{
				Text:       "送样日期:",
				ColumnSpan: 1,
				Alignment:  AlignHCenterVCenter,
			},
			LineEdit{AssignTo: &date, CueBanner: "请输入送样日期", TextAlignment: AlignCenter, ColumnSpan: 3},

			Label{
				ColumnSpan: 4,
			},

			Label{
				Text:       "备注:",
				ColumnSpan: 1,
				Alignment:  AlignHCenterVCenter,
			},
			LineEdit{AssignTo: &desc, CueBanner: "请输入备注", TextAlignment: AlignCenter, ColumnSpan: 3},

			Label{
				ColumnSpan: 4,
			},

			Label{
				Text:       "打印数量:",
				ColumnSpan: 1,
				Alignment:  AlignHCenterVCenter,
			},
			NumberEdit{
				AssignTo:   &sendNum,
				Value:      1,
				Decimals:   1,
				Alignment:  AlignHCenterVCenter,
				ColumnSpan: 3,
			},

			Label{
				ColumnSpan: 4,
			},
			VSeparator{
				ColumnSpan: 4,
			},

			PushButton{
				Text:       "打印",
				ColumnSpan: 4,
				OnClicked: func() {
					print.Name = name.Text()
					print.Spec = spec.Text()
					print.Num = num.Text()
					print.Brand = brand.Text()
					print.Pepole = people.Text()
					print.Date = date.Text()
					print.Desc = desc.Text()
					print.SendNum = int(sendNum.Value())

					job := &Job{
						printerNum: strconv.Itoa(print.SendNum),
						vArr:       [7]string{},
					}

					job.vArr[0] = print.Name
					job.vArr[1] = print.Spec
					job.vArr[2] = print.Num
					job.vArr[3] = print.Brand
					job.vArr[4] = print.Pepole
					job.vArr[5] = print.Date
					job.vArr[6] = print.Desc

					JobChan <- job

					lebs.SetText("开始打印.....")

					fmt.Println(print)

				},
			},

			Label{
				AssignTo:   &lebs,
				ColumnSpan: 4,
			},
		},
	}.Run()
}
