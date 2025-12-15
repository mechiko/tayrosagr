package main

import (
	"fmt"
	"tayrosagr/process"
	"tayrosagr/utility"

	"github.com/lxn/walk"
	dcl "github.com/lxn/walk/declarative"
	"github.com/mechiko/dbscan"
)

func startDialog(dbs *dbscan.Dbs) (out string, err error) {
	var dlg *walk.Dialog
	var page *walk.Composite
	var acceptPB, cancelPB *walk.PushButton
	var textRich *RichEdit
	var fileName *walk.Label

	icon, err := walk.Resources.Icon("3")
	if err != nil {
		return "", fmt.Errorf("%w", err)
	}
	if err := (dcl.Dialog{
		AssignTo:      &dlg,
		Title:         "Агрегация для Тайрос",
		Size:          dcl.Size{Width: 700, Height: 400},
		Icon:          icon,
		Layout:        dcl.VBox{Spacing: 10, Margins: dcl.Margins{Left: 5, Top: 5, Right: 5, Bottom: 5}},
		DefaultButton: &acceptPB,
		CancelButton:  &cancelPB,
		Children: []dcl.Widget{
			dcl.Composite{
				Layout:   dcl.VBox{MarginsZero: true, SpacingZero: false, Margins: dcl.Margins{Left: 0, Top: 0, Right: 0, Bottom: 0}},
				Border:   true,
				AssignTo: &page,
				MinSize:  dcl.Size{Height: 200},
				Children: []dcl.Widget{},
			},
			dcl.Composite{
				Layout:  dcl.HBox{MarginsZero: true, SpacingZero: false, Margins: dcl.Margins{Left: 0, Top: 0, Right: 0, Bottom: 0}},
				Border:  false,
				MinSize: dcl.Size{Width: 500},
				Children: []dcl.Widget{
					dcl.Label{
						Text: "Файл для обработки:",
					},
					dcl.Label{
						AssignTo: &fileName,
						Text:     "",
					},
					dcl.PushButton{
						Text: "Выбрать",
						OnClicked: func() {
							fileName, err := utility.DialogOpenFile([]utility.FileType{utility.Csv, utility.Txt}, "", ".")
							if err != nil {
								logErrMessage(err.Error(), textRich, dlg)
								// utility.MessageBox32("ошибка", err.Error())
								return
							}
							process, err := process.New(dbs)
							if err != nil {
								logErrMessage(err.Error(), textRich, dlg)
								// utility.MessageBox32("ошибка", err.Error())
								return
							}
							err = process.ReadFile(fileName)
							if err != nil {
								logErrMessage(err.Error(), textRich, dlg)
								// utility.MessageBox32("ошибка", err.Error())
								return
							}
							msg := fmt.Sprintf("обработано %d марок и %d коробок", len(process.Cis), len(process.Pallet))
							logMessage(msg, textRich, dlg)
							err = process.WritePaletsForce()
							if err != nil {
								logErrMessage(err.Error(), textRich, dlg)
								// utility.MessageBox32("ошибка", err.Error())
								return
							}
							logMessage("агрегация записана в АлкоХелп", textRich, dlg)
						},
					},
					dcl.HSpacer{},
				},
			},
			dcl.Composite{
				Border: false,
				Layout: dcl.HBox{MarginsZero: true, SpacingZero: true, Margins: dcl.Margins{Left: 5, Top: 5, Right: 5, Bottom: 0}},
				Children: []dcl.Widget{
					// dcl.PushButton{
					// 	AssignTo: &acceptPB,
					// 	Text:     "OK",
					// 	OnClicked: func() {
					// 		dlg.Accept()
					// 	},
					// },
					dcl.PushButton{
						AssignTo:  &cancelPB,
						Text:      "Выход",
						OnClicked: func() { dlg.Cancel() },
					},
					dcl.HSpacer{},
				},
			},
			dcl.VSpacer{},
		},
	}).Create(nil); err != nil {
		return "", fmt.Errorf("%w", err)
	}
	textRich, err = NewRichEdit(page)
	if err != nil {
		return "", fmt.Errorf("%w", err)
	}
	if textRich == nil {
		return "", fmt.Errorf("nil richtext")
	}
	dlg.SetBounds(walk.Rectangle{
		X:     300,
		Y:     300,
		Width: 400,
		// Height: 200,
	})
	logMessage("выберите файл", textRich, dlg)
	if ret := dlg.Run(); ret != 1 {
		return "", fmt.Errorf("dialog return %d", ret)
	}
	return out, nil
}
