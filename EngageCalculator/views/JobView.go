package views

import (
	"EngageCalculator/models"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"strconv"
)

type JobView struct {
	Jobs     map[string]models.Job
	JobNames []string

	// UI组件
	JobIDLabel     *widget.Label
	StyleNameLabel *widget.Label
	MaxLevelLabel  *widget.Label
	WeaponTable    *widget.Table
	StatsTable     *widget.Table
	JobSelect      *widget.Select
}

func NewJobView(jobs map[string]models.Job, jobNames []string) *JobView {
	jv := &JobView{
		Jobs:     jobs,
		JobNames: jobNames,

		JobIDLabel:     widget.NewLabel("职业ID: "),
		StyleNameLabel: widget.NewLabel("职业名称: "),
		MaxLevelLabel:  widget.NewLabel("最高等级: "),
		WeaponTable:    createJobWeaponTable(),
		StatsTable:     createJobStatsTable(),
	}

	jv.JobSelect = widget.NewSelect(jobNames, jv.onJobSelected)

	return jv
}

func (jv *JobView) GetContent() *fyne.Container {
	// 创建职业顶部信息区域
	jobTopInfo := container.NewGridWithColumns(3, jv.JobIDLabel, jv.StyleNameLabel, jv.MaxLevelLabel)

	// 创建职业标签页容器
	jobTabs := container.NewAppTabs(
		container.NewTabItem("武器等级", jv.WeaponTable),
		container.NewTabItem("属性详情", jv.StatsTable),
	)

	// 创建职业部分
	jobSection := container.NewBorder(
		container.NewVBox(
			widget.NewLabel("选择职业:"),
			jv.JobSelect,
			jobTopInfo,
		),
		nil, nil, nil,
		jobTabs,
	)

	return jobSection
}

func (jv *JobView) onJobSelected(selected string) {
	// 根据中文名称查找对应的JID
	var selectedJid string
	for jid, name := range models.JobCanRead {
		if name == selected {
			selectedJid = jid
			break
		}
	}

	job, exists := jv.Jobs[selectedJid]
	if !exists {
		return
	}

	// 更新职业基本信息
	jv.JobIDLabel.SetText("职业ID: " + job.Jid)
	jv.StyleNameLabel.SetText("职业名称: " + job.StyleName)
	jv.MaxLevelLabel.SetText("最高等级: " + strconv.Itoa(job.MaxLevel))

	// 更新职业表格
	updateJobWeaponTable(jv.WeaponTable, job)
	updateJobStatsTable(jv.StatsTable, job)
}

func createJobWeaponTable() *widget.Table {
	return widget.NewTable(
		func() (int, int) { return 10, 2 },
		func() fyne.CanvasObject {
			return widget.NewLabel("Template")
		},
		func(id widget.TableCellID, obj fyne.CanvasObject) {
			label := obj.(*widget.Label)
			label.SetText("")
		},
	)
}

func updateJobWeaponTable(table *widget.Table, job models.Job) {
	table.Length = func() (int, int) { return len(job.MaxWeaponLevels) + 1, 2 }
	table.UpdateCell = func(id widget.TableCellID, obj fyne.CanvasObject) {
		label := obj.(*widget.Label)
		if id.Row == 0 {
			if id.Col == 0 {
				label.SetText("武器类型")
			} else {
				label.SetText("最高等级")
			}
		} else {
			row := id.Row - 1
			if row < len(job.MaxWeaponLevels) {
				if id.Col == 0 {
					if name, ok := models.WeaponNames[row]; ok {
						label.SetText(name)
					} else {
						label.SetText("Unknown")
					}
				} else {
					label.SetText(job.MaxWeaponLevels[row])
				}
			}
		}
	}
	table.Refresh()
}

func createJobStatsTable() *widget.Table {
	return widget.NewTable(
		func() (int, int) { return 9, 4 },
		func() fyne.CanvasObject {
			return widget.NewLabel("Template")
		},
		func(id widget.TableCellID, obj fyne.CanvasObject) {
			label := obj.(*widget.Label)
			label.SetText("")
		},
	)
}

func updateJobStatsTable(table *widget.Table, job models.Job) {
	table.Length = func() (int, int) { return len(job.JobBaseStats) + 1, 4 }
	table.UpdateCell = func(id widget.TableCellID, obj fyne.CanvasObject) {
		label := obj.(*widget.Label)
		if id.Row == 0 {
			switch id.Col {
			case 0:
				label.SetText("属性")
			case 1:
				label.SetText("基础值")
			case 2:
				label.SetText("上限值")
			case 3:
				label.SetText("成长值")
			}
		} else {
			row := id.Row - 1
			if row < len(job.JobBaseStats) {
				switch id.Col {
				case 0:
					if name, ok := models.StatNames[row]; ok {
						label.SetText(name)
					} else {
						label.SetText("Unknown")
					}
				case 1:
					label.SetText(strconv.Itoa(job.JobBaseStats[row]))
				case 2:
					label.SetText(strconv.Itoa(job.JobLimitStats[row]))
				case 3:
					label.SetText(strconv.Itoa(job.JobGrowStats[row]))
				}
			}
		}
	}
	table.Refresh()
}
