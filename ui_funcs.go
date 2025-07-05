package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"golang.design/x/clipboard"
	"golang.org/x/image/colornames"
)

func createPasswordBox(w fyne.Window, callback func(password string, ok bool), show bool) {
	passwordEntry := widget.NewPasswordEntry()

	strength := widget.NewLabelWithStyle(analyzePasswordStrength(""), fyne.TextAlignCenter, fyne.TextStyle{Bold: true})
	if show {
		dialog.ShowForm("Enter Password", "OK", "Cancel", []*widget.FormItem{widget.NewFormItem("Password", passwordEntry), widget.NewFormItem("Strength: ", strength)}, func(ok bool) {
			if !ok {
				callback("", false)
				return
			}
			callback(passwordEntry.Text, true)
		}, w)
		passwordEntry.OnChanged = func(s string) {
			strength.SetText(canvas.Text{
				Alignment: fyne.TextAlignCenter,
				Color:     colornames.Red,
				Text:      analyzePasswordStrength(s),
				TextSize:  15,
				TextStyle: fyne.TextStyle{},
			}.Text)
		}
	} else {
		dialog.ShowForm("Enter Password", "OK", "Cancel", []*widget.FormItem{widget.NewFormItem("Password", passwordEntry)}, func(ok bool) {
			if !ok {
				callback("", false)
				return
			}
			callback(passwordEntry.Text, true)
		}, w)
	}

}

func showDecryptedContent(w fyne.Window, decryptedText string) {
	text := widget.NewMultiLineEntry()
	text.SetText(decryptedText)
	text.Disable()

	copyButton := widget.NewButton("Copy", func() {
		clipboard.Write(clipboard.FmtText, []byte(decryptedText))
	})

	formItems := []*widget.FormItem{
		widget.NewFormItem("Decrypted Text", text),
		widget.NewFormItem("", copyButton),
	}

	d := dialog.NewForm("Decrypted Content", "Done", "", formItems, func(b bool) {}, w)
	d.Resize(fyne.NewSize(400, 300))
	d.Show()
}
