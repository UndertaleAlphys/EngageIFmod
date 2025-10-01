package operations

import (
	"EngageCalculator/models"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
	"math"
	"strconv"
)

// 辅助函数，确保数值不小于0
func nonNegative(value int) int {
	if value < 0 {
		return 0
	}
	return value
}

func UpdatePersonStatsTable(table *widget.Table, person models.Person, job models.Job, simState *models.PersonSimState) {
	table.Length = func() (int, int) {
		return 11, 4
	}

	table.CreateCell = func() fyne.CanvasObject {
		label := &widget.Label{}
		label.Alignment = fyne.TextAlignCenter
		label.Wrapping = fyne.TextWrapWord
		return label
	}

	table.UpdateCell = func(id widget.TableCellID, template fyne.CanvasObject) {
		label := template.(*widget.Label)
		label.Alignment = fyne.TextAlignCenter
		label.Wrapping = fyne.TextWrapWord

		if id.Row == 0 {
			label.Importance = widget.HighImportance
			switch id.Col {
			case 0:
				label.SetText(models.LabelAttr)
			case 1:
				label.SetText(models.LabelBaseStats)
			case 2:
				label.SetText(models.LabelLimitStats)
			case 3:
				label.SetText(models.LabelGrowStats)
			}
		} else if id.Row >= 1 && id.Row <= 9 {
			row := id.Row - 1
			if row < len(person.PersonBaseStats) {
				switch id.Col {
				case 0:
					if name, ok := models.StatNames[row]; ok {
						label.SetText(name)
					} else {
						label.SetText("Unknown")
					}
				case 1:
					personGrow := person.PersonGrowStats[row]
					levelBonus := (person.Level - 1) * personGrow
					roundedBonus := int(math.Round(float64(levelBonus) / 100.0))

					baseValue := person.PersonBaseStats[row] + job.JobBaseStats[row] + roundedBonus + simState.GrowthCount[row]

					totalPoints := simState.TotalPoints[row] + float64(person.PersonBaseStats[row]) + float64(job.JobBaseStats[row]) + float64(roundedBonus)

					if baseValue >= job.JobLimitStats[row]+person.PersonLimitStats[row] {
						//绿色字体
						baseValue = job.JobLimitStats[row] + person.PersonLimitStats[row]
						totalPoints = float64(job.JobLimitStats[row] + person.PersonLimitStats[row])
					}

					// 确保基础值不小于0
					baseValue = nonNegative(baseValue)

					if totalPoints > 0 {
						label.Importance = widget.MediumImportance
						if baseValue >= job.JobLimitStats[row]+person.PersonLimitStats[row] {
							baseValue = job.JobLimitStats[row] + person.PersonLimitStats[row]
							totalPoints = float64(job.JobLimitStats[row] + person.PersonLimitStats[row])
							label.Importance = widget.SuccessImportance

						}
						label.SetText(fmt.Sprintf("%d(%.2f)", baseValue, totalPoints))
						label.Refresh()
					} else {
						label.SetText(fmt.Sprintf("%d", baseValue))
					}

				case 2:
					personLimit := person.PersonLimitStats[row]
					jobLimit := job.JobLimitStats[row]

					// 确保最终显示的数值不小于0
					finalLimit := nonNegative(jobLimit + personLimit)

					if personLimit > 0 {
						label.SetText(strconv.Itoa(finalLimit) + "(+" + strconv.Itoa(personLimit) + ")")
					} else if personLimit < 0 {
						label.SetText(strconv.Itoa(finalLimit) + "(" + strconv.Itoa(personLimit) + ")")
					} else {
						label.SetText(strconv.Itoa(finalLimit))
					}
				case 3:
					personGrow := person.PersonGrowStats[row]
					jobGrow := job.JobGrowStats[row]

					displayGrow := personGrow
					if simState.Spellbook && row != 8 {
						displayGrow += 10
					}

					displayGrow = nonNegative(displayGrow)

					if jobGrow > 0 {
						if simState.EffortTalent {
							label.SetText(fmt.Sprintf("%d(+%d)", displayGrow, jobGrow*2))
						} else {
							label.SetText(fmt.Sprintf("%d(+%d)", displayGrow, jobGrow))
						}
					} else {
						label.SetText(fmt.Sprintf("%d", displayGrow))
					}
				}
			}
		} else if id.Row == 10 {
			// 总合计行
			switch id.Col {
			case 0:
				label.Importance = widget.HighImportance
				label.SetText(models.LabelSumStats)
			case 1:
				// 计算所有属性的基础值总和
				totalBase := 0
				for i := 0; i < len(person.PersonBaseStats); i++ {
					personGrow := person.PersonGrowStats[i]
					levelBonus := (person.Level - 1) * personGrow
					roundedBonus := int(math.Round(float64(levelBonus) / 100.0))
					baseValue := person.PersonBaseStats[i] + job.JobBaseStats[i] + roundedBonus + simState.GrowthCount[i]
					totalBase += baseValue
				}
				// 确保总计值不小于0
				totalBase = nonNegative(totalBase)
				label.SetText(strconv.Itoa(totalBase))
			case 2:
				totalLimit := 0
				for i := 0; i < len(person.PersonLimitStats); i++ {
					personLimit := person.PersonLimitStats[i]
					jobLimit := job.JobLimitStats[i]
					limitValue := jobLimit
					if personLimit > 0 {
						limitValue += personLimit
					} else if personLimit < 0 {
						limitValue += personLimit
					}
					totalLimit += limitValue
				}
				totalLimit = nonNegative(totalLimit)
				label.SetText(strconv.Itoa(totalLimit))
			case 3:
				totalGrow := 0
				for i := 0; i < len(person.PersonGrowStats); i++ {
					personGrow := person.PersonGrowStats[i]
					jobGrow := job.JobGrowStats[i]

					if simState.Spellbook && i != 8 {
						personGrow += 10
					}

					if simState.EffortTalent {
						jobGrow *= 2
					}
					growValue := personGrow + jobGrow

					totalGrow += nonNegative(growValue)
				}
				label.SetText(strconv.Itoa(totalGrow))
			}
		}
	}
	table.Refresh()
}
