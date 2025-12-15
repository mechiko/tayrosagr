package main

import "github.com/lxn/walk"

func logMessage(s string, t *RichEdit, _ *walk.Dialog) {
	txt := color(s, Blue)
	msg, styles := parseString(txt)

	// dlg.Synchronize(func() {
	// 	t.AppendText(msg, styles...)
	// })
	t.AppendText(msg+"\n", styles...)
}

func logErrMessage(s string, t *RichEdit, _ *walk.Dialog) {
	txt := color(s, Red)
	msg, styles := parseString(txt)
	// dlg.Synchronize(func() {
	// 	t.AppendText(msg, styles...)
	// })
	t.AppendText(msg+"\n", styles...)
}
