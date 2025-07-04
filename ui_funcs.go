package main

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

func set_selected(index *widget.ListItemID, id int) {
	*index = id
}

func create_password_box(w fyne.Window, callback func(password string, entered bool)) {
	password := widget.NewPasswordEntry()
	items := []*widget.FormItem{
		widget.NewFormItem("Password", password),
	}

	d := dialog.NewForm("Enter Password", "Done", "N", items, func(b bool) {
		fmt.Println("triggered")
		if b {
			callback(password.Text, true)
			return

		} else {
			callback(password.Text, false)
		}
	}, w)
	d.Resize(fyne.NewSize(400, 100))
	d.Show()
}
