package main

import (
	"errors"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	d "github.com/sqweek/dialog"
	"golang.design/x/clipboard"
	"golang.org/x/image/colornames"
)

var os_select int

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

func showStandaloneOverview(w fyne.Window) {
	type StandaloneOptions struct {
		os_select     int
		self_destruct int
	}
	var Options StandaloneOptions
	os_selector := widget.NewSelect([]string{"Linux", "Mac", "Windows"}, func(s string) {
		switch s {
		case "Linux":
			Options.os_select = 1
		case "Mac":
			Options.os_select = 2
		case "Windows":
			Options.os_select = 2
		default:
			dialog.ShowError(errors.New("invalid option: "+s), w)
		}
	})
	self_destruct := widget.NewEntry()
	self_destruct.SetPlaceHolder("No. Of Failed Attempts (0 To Disable)")
	formItems := []*widget.FormItem{
		widget.NewFormItem("Operating System", os_selector),
		widget.NewFormItem("Self Destruct", self_destruct),
	}

	formDialog := dialog.NewForm("Standalone Options", "Done", "", formItems, func(b bool) {
	}, w)
	formDialog.Resize(fyne.NewSize(400, 300))
	fyne.DoAndWait(func() { formDialog.Show() })
	file, err := d.Directory().Title("Choose File").Browse()
	if errors.Is(err, d.ErrCancelled) {
		return
	}

	entry := widget.NewMultiLineEntry()
	checkbox_deleteWorkspace := widget.NewCheck("Delete Workspace", func(b bool) {})
	dialog.ShowForm("Enter Text", "Encrypt", "Cancel", []*widget.FormItem{widget.NewFormItem("Text", entry), widget.NewFormItem("", checkbox_deleteWorkspace)}, func(ok bool) {
		if !ok {
			return
		}
		createPasswordBox(w, func(password string, ok bool) {
			if !ok {
				return
			}

			if err != nil {
				dialog.ShowError(err, w)
				return
			}
			makeNewStandalone(entry.Text, false, password, checkbox_deleteWorkspace.Checked, file)
		}, true)
	}, w)

}
