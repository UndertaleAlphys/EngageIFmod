package operations

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
)

func CreatePersonStatsTable() *widget.Table {
	table := widget.NewTable(
		func() (int, int) { return 11, 4 }, // 标題行+9个属性+2个合计
		func() fyne.CanvasObject {
			return widget.NewLabel("Template")
		},
		func(id widget.TableCellID, obj fyne.CanvasObject) {
			label := obj.(*widget.Label)
			label.SetText("")
		},
	)

	// 设置列宽以更好地利用空间
	table.SetColumnWidth(0, 100) // 属性列
	table.SetColumnWidth(1, 150) // 基础值列
	table.SetColumnWidth(2, 150) // 上限列
	table.SetColumnWidth(3, 150) // 成长率列

	return table
}
