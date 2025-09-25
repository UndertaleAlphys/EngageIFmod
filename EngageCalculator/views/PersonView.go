package views

import (
	"EngageCalculator/models"
	"EngageCalculator/operations"
	"EngageCalculator/utils"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"sort"
	"strconv"
)

type PersonView struct {
	Persons       map[string]models.Person
	PersonNames   []string
	Jobs          map[string]models.Job
	CurrentPerson models.Person

	// 人物状态历史
	PersonStateHistory map[string]*models.PersonStateHistory
	PersonSimStates    map[string]*models.PersonSimState

	// UI组件
	PersonIDLabel    *widget.Label
	PersonJidLabel   *widget.Label
	LevelLabel       *widget.Label
	GenderLabel      *widget.Label
	PersonSelect     *widget.Select
	PersonStatsTable *widget.Table
	HistoryStatsList *widget.List

	// 模拟升级相关组件
	LevelUpButton   *widget.Button
	LevelDownButton *widget.Button
	ResetButton     *widget.Button

	// 上位转职相关组件
	JobPromotionContainer *fyne.Container
	JobPromotionButtons   []*widget.Button
	JobPromotionLabel     *widget.Label // 上位转职标签

	// 横转相关组件
	LateralJobSelect *widget.Select
	LateralJobButton *widget.Button

	// 努力才能和术书相关组件
	EffortTalentCheck *widget.Check
	SpellbookCheck    *widget.Check

	LateralJobMap map[string]string

	LevelSlider *widget.Slider // 等级滑动条

}

// 定义一个按PID顺序排序的函数
func (uv *PersonView) getSortedPersonNames() []string {
	// 定义PID有序列表，按照PersonCanRead中定义的顺序
	orderedPIDs := utils.GetOrderedPersonKeys(models.PersonCanRead)

	// 根据orderedPIDs创建排序后的名称列表
	var sortedNames []string
	for _, pid := range orderedPIDs {
		if name, exists := models.PersonCanRead[pid]; exists {
			// 检查该人物是否在当前的Persons中存在
			if _, personExists := uv.Persons[pid]; personExists {
				sortedNames = append(sortedNames, name)
			}
		}
	}

	return sortedNames
}

func NewPersonView(persons map[string]models.Person, personNames []string, jobs map[string]models.Job) *PersonView {
	uv := &PersonView{
		Persons:               persons,
		PersonNames:           personNames,
		Jobs:                  jobs,
		PersonIDLabel:         widget.NewLabel(models.LabelPersonID),
		PersonJidLabel:        widget.NewLabel(models.LabelPersonJid),
		GenderLabel:           widget.NewLabel(models.LabelPersonGender),
		LevelLabel:            widget.NewLabel(models.LabelPersonLevel),
		PersonStateHistory:    make(map[string]*models.PersonStateHistory),
		PersonSimStates:       make(map[string]*models.PersonSimState),
		JobPromotionContainer: container.NewHBox(),
		JobPromotionButtons:   make([]*widget.Button, 0, 2),
		LateralJobMap:         make(map[string]string),
	}

	// 获取按正确顺序排序的人物名称
	sortedPersonNames := uv.getSortedPersonNames()

	// 使用排序后的人物名称创建选择框
	uv.PersonSelect = widget.NewSelect(sortedPersonNames, uv.onPersonSelected)
	uv.PersonSelect.PlaceHolder = "选择人物"

	// 创建按钮
	uv.LevelUpButton = widget.NewButton(models.LabelLevelUp, uv.onLevelUpClicked)
	uv.LevelDownButton = widget.NewButton(models.LabelLevelDown, uv.onLevelDownClicked)
	uv.ResetButton = widget.NewButton(models.LabelReset, uv.onResetClicked)

	uv.JobPromotionContainer = container.NewHBox()

	uv.LevelSlider = widget.NewSlider(1, 40) // 默认范围1-40
	uv.LevelSlider.Step = 1
	uv.LevelSlider.OnChanged = func(value float64) { // 修复引用问题
		uv.onLevelSliderChanged(value)
	}
	// 通过设置滑动条尺寸
	uv.LevelSlider.Resize(fyne.NewSize(300, uv.LevelSlider.MinSize().Height))

	// 创建上位转职标签
	uv.JobPromotionLabel = widget.NewLabel(models.LabelPromotion)

	// 创建横转组件
	uv.LateralJobSelect = widget.NewSelect([]string{}, func(selected string) {})
	uv.LateralJobSelect.PlaceHolder = "选择横转职业"
	uv.LateralJobButton = widget.NewButton("横转", uv.onLateralJobSelected)

	// 创建努力才能和术书组件
	uv.EffortTalentCheck = widget.NewCheck("努力才能", func(checked bool) {
		uv.onEffortTalentChanged(checked)
	})
	uv.SpellbookCheck = widget.NewCheck("术书", func(checked bool) {
		uv.onSpellbookChanged(checked)
	})

	uv.PersonStatsTable = operations.CreatePersonStatsTable()
	uv.HistoryStatsList = createHistoryList(uv, "")

	return uv
}

func (uv *PersonView) SetJobs(jobs map[string]models.Job) {
	uv.Jobs = jobs
}

func (uv *PersonView) GetContent() *fyne.Container {

	sortedPersonNames := uv.getSortedPersonNames()

	// 创建人物选择列表
	personList := widget.NewList(
		func() int {
			return len(sortedPersonNames)
		},
		func() fyne.CanvasObject {
			return widget.NewLabel("Template")
		},
		func(id widget.ListItemID, obj fyne.CanvasObject) {
			obj.(*widget.Label).SetText(sortedPersonNames[id])
		},
	)

	// 设置人物选择列表的事件处理
	personList.OnSelected = func(id widget.ListItemID) {
		uv.onPersonSelected(sortedPersonNames[id])
	}

	// 创建可滚动的人物选择区域
	personSelectContainer := container.NewVScroll(personList)
	personSelectContainer.SetMinSize(fyne.NewSize(200, 600))

	// 创建人物顶部信息区域
	personTopInfo := container.NewGridWithColumns(4,
		uv.PersonIDLabel,
		uv.PersonJidLabel,
		uv.LevelLabel,
		uv.GenderLabel,
	)

	// 创建模拟升级控件
	sliderContainer := container.NewBorder(
		nil, nil,
		widget.NewLabel("等级:"),
		nil,
		uv.LevelSlider,
	)

	levelControl := container.NewVBox(
		container.NewHBox(
			widget.NewLabel(models.LabelLevelSetup),
			uv.LevelUpButton,
			uv.LevelDownButton,
			uv.ResetButton,
		),
		sliderContainer,
	)

	// 创建转职控制区域
	jobPromotionControl := container.NewVBox(
		uv.JobPromotionLabel,
		uv.JobPromotionContainer,
		widget.NewLabel(models.LabelTransfer),
		container.NewHBox(uv.LateralJobSelect, uv.LateralJobButton, uv.EffortTalentCheck, uv.SpellbookCheck),
	)

	controlPanel := container.NewVBox(
		personTopInfo,
		levelControl,
		jobPromotionControl,
	)

	controlCard := widget.NewCard("", "", controlPanel)

	StateSection := container.NewHSplit(uv.PersonStatsTable, uv.HistoryStatsList)
	StateSection.SetOffset(0.57)

	mainSection := container.NewBorder(
		controlCard,
		nil, nil, nil,
		StateSection,
	)

	mainContainer := container.NewBorder(
		nil, nil, personSelectContainer, nil,
		mainSection,
	)

	return mainContainer
}

// 在 onPersonSelected 方法中，确保初始状态正确保存
func (uv *PersonView) onPersonSelected(selected string) {
	// 根据中文名称查找对应的PID
	var selectedPid string
	for pid, name := range models.PersonCanRead {
		if name == selected {
			selectedPid = pid
			break
		}
	}

	// 获取对应的Person属性值
	person, exists := uv.Persons[selectedPid]
	if !exists {
		return
	}

	// 保存当前人物信息
	uv.CurrentPerson = person

	// 初始化状态历史记录
	if _, exists := uv.PersonStateHistory[selectedPid]; !exists {
		uv.PersonStateHistory[selectedPid] = &models.PersonStateHistory{
			Pid:    models.PersonCanRead[selectedPid],
			States: []models.PersonStateHistoryItem{},
		}
	}
	job := uv.Jobs[person.Jid]

	// 获取当前模拟状态
	if _, exists := uv.PersonSimStates[selectedPid]; !exists {
		// 正确计算初始的个人成长积分
		initialTotalPoints := [9]float64{
			float64(person.PersonGrowStats[0]) / 100,
			float64(person.PersonGrowStats[1]) / 100,
			float64(person.PersonGrowStats[2]) / 100,
			float64(person.PersonGrowStats[3]) / 100,
			float64(person.PersonGrowStats[4]) / 100,
			float64(person.PersonGrowStats[5]) / 100,
			float64(person.PersonGrowStats[6]) / 100,
			float64(person.PersonGrowStats[7]) / 100,
			float64(person.PersonGrowStats[8]) / 100,
		}

		uv.PersonSimStates[selectedPid] = &models.PersonSimState{
			Pid:           selectedPid,
			JobID:         person.Jid,
			Level:         person.Level,
			MaxLevel:      job.MaxLevel,
			MinLevel:      person.Level,
			TotalPoints:   initialTotalPoints, // 使用正确计算的初始积分
			GrowthCount:   [9]int{},
			JobBaseStats:  job.JobBaseStats,
			JobGrowStats:  job.JobGrowStats,
			JobLimitStats: job.JobLimitStats,
			EffortTalent:  false,
			Spellbook:     false,
		}
	}

	// 立刻保存当前信息进States: []models.PersonStateHistoryItem{}
	if len(uv.PersonStateHistory[selectedPid].States) == 0 {
		// 正确计算初始的个人成长积分
		initialTotalPoints := [9]float64{
			float64(person.PersonGrowStats[0]) / 100,
			float64(person.PersonGrowStats[1]) / 100,
			float64(person.PersonGrowStats[2]) / 100,
			float64(person.PersonGrowStats[3]) / 100,
			float64(person.PersonGrowStats[4]) / 100,
			float64(person.PersonGrowStats[5]) / 100,
			float64(person.PersonGrowStats[6]) / 100,
			float64(person.PersonGrowStats[7]) / 100,
			float64(person.PersonGrowStats[8]) / 100,
		}

		uv.PersonStateHistory[selectedPid].States = append(uv.PersonStateHistory[selectedPid].States, models.PersonStateHistoryItem{
			Level:            uv.PersonSimStates[selectedPid].Level,
			JobID:            uv.PersonSimStates[selectedPid].JobID,
			MaxLevel:         uv.PersonSimStates[selectedPid].MaxLevel,
			MinLevel:         uv.PersonSimStates[selectedPid].Level,
			TotalPoints:      initialTotalPoints, // 使用正确计算的初始积分
			GrowthCount:      uv.PersonSimStates[selectedPid].GrowthCount,
			PersonBaseStats:  person.PersonBaseStats,
			JobBaseStats:     uv.PersonSimStates[selectedPid].JobBaseStats,
			PersonLimitStats: person.PersonLimitStats,
			JobLimitStats:    uv.PersonSimStates[selectedPid].JobLimitStats,
			PersonGrowStats:  person.PersonGrowStats,
			JobGrowStats:     uv.PersonSimStates[selectedPid].JobGrowStats,
			EffortTalent:     false,
			Spellbook:        false,
		})
	}

	// 重置努力才能术书UI状态
	uv.EffortTalentCheck.SetChecked(false)
	uv.SpellbookCheck.SetChecked(false)

	uv.PersonIDLabel.SetText(models.LabelPersonID + models.PersonCanRead[person.Pid])
	uv.PersonJidLabel.SetText(models.LabelPersonJid + models.JobCanRead[uv.PersonSimStates[selectedPid].JobID])
	uv.LevelLabel.SetText(models.LabelPersonLevel + strconv.Itoa(uv.PersonSimStates[selectedPid].Level))

	genderText := models.LabelUnknown
	if person.Gender == models.SEX_MALE {
		genderText = models.LabelMale
	} else if person.Gender == models.SEX_FEMALE {
		genderText = models.LabelFemale
	}
	uv.GenderLabel.SetText(models.LabelPersonGender + genderText)

	// 更新历史记录列表
	uv.updatePersonHistoryList(selectedPid)
}

func (uv *PersonView) updatePersonHistoryList(selectedPid string) {
	// 更新历史记录列表的数据源
	uv.HistoryStatsList.Length = func() int {
		if uv.PersonStateHistory[selectedPid] != nil {
			return len(uv.PersonStateHistory[selectedPid].States)
		}
		return 0
	}

	uv.HistoryStatsList.CreateItem = func() fyne.CanvasObject {
		button := widget.NewButton("Template", func() {})
		return button
	}

	uv.HistoryStatsList.UpdateItem = func(id widget.ListItemID, template fyne.CanvasObject) {
		if uv.PersonStateHistory[selectedPid] != nil &&
			id < len(uv.PersonStateHistory[selectedPid].States) {

			reverseIndex := len(uv.PersonStateHistory[selectedPid].States) - 1 - id

			entry := uv.PersonStateHistory[selectedPid].States[reverseIndex]
			button := template.(*widget.Button)

			// 获取职业中文名
			jobName := entry.JobID
			if name, ok := models.JobCanRead[entry.JobID]; ok {
				jobName = name
			}

			// 设置按钮文本
			button.SetText(fmt.Sprintf("职业: %s, 等级: %d", jobName, entry.Level))

			button.OnTapped = func() {
				// 创建一个模拟状态对象
				simState := &models.PersonSimState{
					Level:         entry.Level,
					JobID:         entry.JobID,
					TotalPoints:   entry.TotalPoints,
					GrowthCount:   entry.GrowthCount,
					MaxLevel:      entry.MaxLevel,
					MinLevel:      entry.MinLevel,
					JobLimitStats: entry.JobLimitStats,
					JobBaseStats:  entry.JobBaseStats,
					JobGrowStats:  entry.JobGrowStats,
					EffortTalent:  entry.EffortTalent, // 恢复努力才能状态
					Spellbook:     entry.Spellbook,    // 恢复术书状态
				}
				// 把simState赋值给uv.PersonSimStates
				uv.PersonSimStates[selectedPid] = simState
				uv.PersonJidLabel.SetText(models.LabelPersonJid + models.JobCanRead[entry.JobID])
				uv.LevelLabel.SetText(models.LabelPersonLevel + strconv.Itoa(entry.Level))
				uv.EffortTalentCheck.SetChecked(entry.EffortTalent) // 更新UI
				uv.SpellbookCheck.SetChecked(entry.Spellbook)       // 更新UI

				// fmt.Println(uv.PersonSimStates[selectedPid])
				operations.UpdatePersonStatsTable(uv.PersonStatsTable, uv.Persons[selectedPid], uv.Jobs[entry.JobID], simState)

				deleteHistory := func(index int) {
					// 也需要相应调整删除逻辑
					uv.PersonStateHistory[selectedPid].States = uv.PersonStateHistory[selectedPid].States[:reverseIndex+1]
				}
				deleteHistory(reverseIndex)

				// 删除后立刻刷新历史表
				uv.updatePersonHistoryList(selectedPid)
			}

		}
	}

	simState := uv.PersonSimStates[selectedPid]
	// 刷新历史记录列表显示
	uv.HistoryStatsList.Refresh()
	operations.UpdatePersonStatsTable(uv.PersonStatsTable, uv.Persons[selectedPid], uv.Jobs[simState.JobID], simState)
	uv.updateJobPromotionButtons(uv.Persons[selectedPid])
	uv.updateLevelButtonsState()
	uv.updateLateralJobSelect(uv.Persons[selectedPid])
}

func (uv *PersonView) onLevelUpClicked() {
	if uv.PersonSelect == nil {
		return
	}

	// 获取当前模拟状态
	simState := uv.PersonSimStates[uv.CurrentPerson.Pid]
	if simState == nil {
		return
	}

	if simState.Level >= simState.MaxLevel {
		uv.LevelUpButton.Disable()
	}

	// 增加模拟等级，但不超过职业上限
	if simState.Level < simState.MaxLevel {
		// 计算从当前等级到下一等级的积分增长
		for i := 0; i < 9; i++ {
			// 计算总成长率
			personGrow := uv.CurrentPerson.PersonGrowStats[i]
			jobGrow := simState.JobGrowStats[i]

			if simState.Spellbook {
				personGrow += 5
			}

			if simState.EffortTalent {
				jobGrow *= 2.0
			}
			totalGrow := personGrow + jobGrow

			// 计算增加的积分
			points := float64(totalGrow) / 100.0

			simState.TotalPoints[i] += points

			threshold := float64(simState.GrowthCount[i]) + 0.99

			if simState.TotalPoints[i] >= threshold {
				// 计算可能的最大增长次数
				totalBase := uv.CurrentPerson.PersonBaseStats[i] + simState.JobBaseStats[i]
				totalLimit := uv.CurrentPerson.PersonLimitStats[i] + simState.JobLimitStats[i]
				currentValue := totalBase + simState.GrowthCount[i]
				// 只有在未达到上限时才增加属性
				if currentValue < totalLimit {
					simState.GrowthCount[i]++
				}
			}
		}
		//fmt.Println(simState)

		simState.Level++
		uv.LevelLabel.SetText(models.LabelPersonLevel + strconv.Itoa(simState.Level))
		// 保存当前等级状态到历史记录
		historyItem := models.PersonStateHistoryItem{
			Level:            simState.Level,
			JobID:            simState.JobID,
			MaxLevel:         simState.MaxLevel,
			MinLevel:         simState.MinLevel,
			TotalPoints:      simState.TotalPoints,
			GrowthCount:      simState.GrowthCount,
			PersonBaseStats:  uv.CurrentPerson.PersonBaseStats,
			PersonLimitStats: uv.CurrentPerson.PersonLimitStats,
			PersonGrowStats:  uv.CurrentPerson.PersonGrowStats,
			JobBaseStats:     uv.Jobs[simState.JobID].JobBaseStats,
			JobLimitStats:    uv.Jobs[simState.JobID].JobLimitStats,
			JobGrowStats:     uv.Jobs[simState.JobID].JobGrowStats,
			EffortTalent:     simState.EffortTalent,
			Spellbook:        simState.Spellbook,
		}
		uv.PersonStateHistory[uv.CurrentPerson.Pid].States = append(
			uv.PersonStateHistory[uv.CurrentPerson.Pid].States,
			historyItem,
		)

		// 更新历史记录列表
		uv.updatePersonHistoryList(uv.CurrentPerson.Pid)

	}
}

func (uv *PersonView) onLevelDownClicked() {
	if uv.PersonSelect == nil {
		return
	}

	// 获取当前模拟状态
	simState := uv.PersonSimStates[uv.CurrentPerson.Pid]
	if simState == nil {
		return
	}

	// 检查是否有足够的历史记录可以回溯
	if len(uv.PersonStateHistory[uv.CurrentPerson.Pid].States) < 2 {
		return
	}

	lastSecondItem := uv.PersonStateHistory[uv.CurrentPerson.Pid].States[len(uv.PersonStateHistory[uv.CurrentPerson.Pid].States)-2]
	// 使用倒数第二条记录恢复模拟状态
	uv.PersonSimStates[uv.CurrentPerson.Pid] = &models.PersonSimState{
		Level:         lastSecondItem.Level,
		JobID:         lastSecondItem.JobID,
		TotalPoints:   lastSecondItem.TotalPoints,
		GrowthCount:   lastSecondItem.GrowthCount,
		MaxLevel:      lastSecondItem.MaxLevel,
		MinLevel:      lastSecondItem.MinLevel,
		JobLimitStats: lastSecondItem.JobLimitStats,
		JobBaseStats:  lastSecondItem.JobBaseStats,
		JobGrowStats:  lastSecondItem.JobGrowStats,
		EffortTalent:  lastSecondItem.EffortTalent,
		Spellbook:     lastSecondItem.Spellbook,
	}

	if lastSecondItem.JobID != simState.JobID {
		uv.PersonJidLabel.SetText(models.LabelPersonJid + models.JobCanRead[lastSecondItem.JobID])
	}
	uv.LevelLabel.SetText(models.LabelPersonLevel + strconv.Itoa(lastSecondItem.Level))

	// 删除历史记录中最后一个状态
	uv.PersonStateHistory[uv.CurrentPerson.Pid].States = uv.PersonStateHistory[uv.CurrentPerson.Pid].States[:len(uv.PersonStateHistory[uv.CurrentPerson.Pid].States)-1]
	uv.updatePersonHistoryList(uv.CurrentPerson.Pid)
}
func (uv *PersonView) onResetClicked() {
	if uv.PersonSelect == nil {
		return
	}

	// 重置模拟状态
	simState := uv.PersonSimStates[uv.CurrentPerson.Pid]

	// 正确计算初始的个人成长积分
	initialTotalPoints := [9]float64{
		float64(uv.CurrentPerson.PersonGrowStats[0]) / 100,
		float64(uv.CurrentPerson.PersonGrowStats[1]) / 100,
		float64(uv.CurrentPerson.PersonGrowStats[2]) / 100,
		float64(uv.CurrentPerson.PersonGrowStats[3]) / 100,
		float64(uv.CurrentPerson.PersonGrowStats[4]) / 100,
		float64(uv.CurrentPerson.PersonGrowStats[5]) / 100,
		float64(uv.CurrentPerson.PersonGrowStats[6]) / 100,
		float64(uv.CurrentPerson.PersonGrowStats[7]) / 100,
		float64(uv.CurrentPerson.PersonGrowStats[8]) / 100,
	}

	simState.GrowthCount = [9]int{}
	simState.TotalPoints = initialTotalPoints // 使用正确计算的初始积分
	simState.Level = uv.CurrentPerson.Level
	simState.JobID = uv.CurrentPerson.Jid
	simState.MaxLevel = uv.Jobs[uv.CurrentPerson.Jid].MaxLevel
	simState.MinLevel = uv.CurrentPerson.Level
	simState.JobLimitStats = uv.Jobs[uv.CurrentPerson.Jid].JobLimitStats
	simState.JobBaseStats = uv.Jobs[uv.CurrentPerson.Jid].JobBaseStats
	simState.JobGrowStats = uv.Jobs[uv.CurrentPerson.Jid].JobGrowStats

	simState.EffortTalent = false // 重置努力才能
	simState.Spellbook = false    // 重置术书

	uv.LevelLabel.SetText(models.LabelPersonLevel + strconv.Itoa(uv.PersonSimStates[uv.CurrentPerson.Pid].Level))
	uv.PersonJidLabel.SetText(models.LabelPersonJid + models.JobCanRead[uv.CurrentPerson.Jid])
	uv.EffortTalentCheck.SetChecked(false) // 重置UI
	uv.SpellbookCheck.SetChecked(false)    // 重置UI

	// 更新历史记录中的最后一个节点
	if len(uv.PersonStateHistory[uv.CurrentPerson.Pid].States) > 0 {
		lastIndex := len(uv.PersonStateHistory[uv.CurrentPerson.Pid].States) - 1
		uv.PersonStateHistory[uv.CurrentPerson.Pid].States[lastIndex].EffortTalent = false
		uv.PersonStateHistory[uv.CurrentPerson.Pid].States[lastIndex].Spellbook = false
		uv.PersonStateHistory[uv.CurrentPerson.Pid].States[lastIndex].TotalPoints = initialTotalPoints // 更新积分
		uv.PersonStateHistory[uv.CurrentPerson.Pid].States[lastIndex].GrowthCount = [9]int{}           // 更新成长计数
	}

	uv.clearPersonHistory(uv.Persons[uv.CurrentPerson.Pid])
	uv.updatePersonHistoryList(uv.CurrentPerson.Pid)
}

// 添加恢复到历史状态函数

func (uv *PersonView) clearPersonHistory(person models.Person) {
	if uv.PersonStateHistory[person.Pid] != nil {
		uv.PersonStateHistory[person.Pid].States = nil
	}

	currentPersonInitialJid := uv.CurrentPerson.Jid
	initialJob := uv.Jobs[currentPersonInitialJid]

	// 正确计算初始的个人成长积分
	initialTotalPoints := [9]float64{
		float64(uv.CurrentPerson.PersonGrowStats[0]) / 100,
		float64(uv.CurrentPerson.PersonGrowStats[1]) / 100,
		float64(uv.CurrentPerson.PersonGrowStats[2]) / 100,
		float64(uv.CurrentPerson.PersonGrowStats[3]) / 100,
		float64(uv.CurrentPerson.PersonGrowStats[4]) / 100,
		float64(uv.CurrentPerson.PersonGrowStats[5]) / 100,
		float64(uv.CurrentPerson.PersonGrowStats[6]) / 100,
		float64(uv.CurrentPerson.PersonGrowStats[7]) / 100,
		float64(uv.CurrentPerson.PersonGrowStats[8]) / 100,
	}

	// 创建初始状态的历史记录项
	initialState := models.PersonStateHistoryItem{
		Level:            uv.CurrentPerson.Level,
		JobID:            uv.CurrentPerson.Jid,
		MaxLevel:         initialJob.MaxLevel,
		MinLevel:         uv.CurrentPerson.Level,
		TotalPoints:      initialTotalPoints,                // 使用正确计算的初始积分
		GrowthCount:      [9]int{},                          // 每个属性成长数
		PersonBaseStats:  uv.CurrentPerson.PersonBaseStats,  // 人物基础属性值
		JobBaseStats:     initialJob.JobBaseStats,           // 职业基础属性值
		PersonLimitStats: uv.CurrentPerson.PersonLimitStats, // 人物属性上限值
		JobLimitStats:    initialJob.JobLimitStats,          // 职业属性上限值
		PersonGrowStats:  uv.CurrentPerson.PersonGrowStats,  // 人物成长率
		JobGrowStats:     initialJob.JobGrowStats,           // 职业成长率
		EffortTalent:     false,                             // 初始努力才能状态
		Spellbook:        false,                             // 初始术书状态
	}

	uv.PersonStateHistory[person.Pid].States = append(uv.PersonStateHistory[person.Pid].States, initialState)

	uv.updatePersonHistoryList(uv.CurrentPerson.Pid)
}
func (uv *PersonView) updateLevelButtonsState() {
	simState := uv.PersonSimStates[uv.CurrentPerson.Pid]
	if simState.Level >= simState.MaxLevel {
		uv.LevelUpButton.Disable()
	} else {
		uv.LevelUpButton.Enable()
	}
	if simState.Level <= simState.MinLevel {
		uv.LevelDownButton.Disable()
	} else {
		uv.LevelDownButton.Enable()
	}

	if len(uv.PersonStateHistory[uv.CurrentPerson.Pid].States) >= 2 {
		lastItem := uv.PersonStateHistory[uv.CurrentPerson.Pid].States[len(uv.PersonStateHistory[uv.CurrentPerson.Pid].States)-1]
		lastSecondItem := uv.PersonStateHistory[uv.CurrentPerson.Pid].States[len(uv.PersonStateHistory[uv.CurrentPerson.Pid].States)-2]
		if lastItem.JobID != lastSecondItem.JobID {
			if lastItem.MaxLevel > lastSecondItem.MaxLevel {
				uv.LevelDownButton.SetText(models.LabelCancelPromotion)
				uv.LevelDownButton.Enable()
			} else {
				uv.LevelDownButton.SetText(models.LabelCancelTransfer)
				uv.LevelDownButton.Enable()
			}
		} else {
			uv.LevelDownButton.SetText(models.LabelLevelDown)
		}
	}

	if uv.LateralJobSelect != nil && uv.LateralJobSelect.Selected != "" {
		uv.LateralJobSelect.Selected = ""
	}

	uv.LevelSlider.Max = float64(simState.MaxLevel)
	uv.LevelSlider.Min = float64(simState.MinLevel)
	uv.LevelSlider.SetValue(float64(simState.Level))
	uv.LevelSlider.Refresh()

	if simState.EffortTalent == true {
		uv.EffortTalentCheck.SetChecked(true)
	} else {
		uv.EffortTalentCheck.SetChecked(false)
	}

	if simState.Spellbook == true {
		uv.SpellbookCheck.SetChecked(true)
	} else {
		uv.SpellbookCheck.SetChecked(false)
	}
}

func createHistoryList(uv *PersonView, pid string) *widget.List {
	// 初始化历史记录列表
	historyList := widget.NewList(
		// 数据长度函数 - 返回历史记录条目数量
		func() int {
			if uv.PersonStateHistory[pid] != nil {
				// 返回历史记录条目数量
				return len(uv.PersonStateHistory[pid].States)
			}
			return 0
		},
		func() fyne.CanvasObject {
			label := widget.NewButton("Template", func() {})
			return label
		},

		func(id widget.ListItemID, obj fyne.CanvasObject) {

			entry := uv.PersonStateHistory[pid].States[id]
			label := obj.(*widget.Button)
			label.SetText(fmt.Sprintf("职业: %d, 等級: %s", models.JobCanRead[entry.JobID], entry.Level))
			label.OnTapped = func() {
				simState := &models.PersonSimState{
					Level:         entry.Level,
					JobID:         entry.JobID,
					TotalPoints:   entry.TotalPoints,
					GrowthCount:   entry.GrowthCount,
					MaxLevel:      entry.MaxLevel,
					MinLevel:      entry.MinLevel,
					JobLimitStats: entry.JobLimitStats,
					JobBaseStats:  entry.JobBaseStats,
					JobGrowStats:  entry.JobGrowStats,
				}
				uv.PersonJidLabel.SetText(models.LabelPersonJid + models.JobCanRead[entry.JobID])
				uv.LevelLabel.SetText(models.LabelPersonLevel + strconv.Itoa(entry.Level))

				//把simState代入uv.PersonSimStates
				uv.PersonSimStates[pid] = simState
				fmt.Println(uv.PersonSimStates[pid])
				operations.UpdatePersonStatsTable(uv.PersonStatsTable, uv.Persons[pid], uv.Jobs[entry.JobID], simState)
				deleteHistory := func(index int) {
					for i := index + 1; i < len(uv.PersonStateHistory[pid].States); i++ {
						uv.PersonStateHistory[pid].States = append(uv.PersonStateHistory[pid].States[:i], uv.PersonStateHistory[pid].States[i+1:]...)
					}
				}
				deleteHistory(id)
			}
		},
	)

	historyList.SetItemHeight(0, 40)

	return historyList
}

// 更新上位转职按钮
func (uv *PersonView) updateJobPromotionButtons(person models.Person) {
	// 清空之前的按钮
	uv.JobPromotionContainer.Objects = nil
	uv.JobPromotionButtons = uv.JobPromotionButtons[:0]
	simState := uv.PersonSimStates[person.Pid]
	currentJob := uv.Jobs[simState.JobID]
	buttonCount := 0 // 记录按钮数量

	// 检查HighJob1 (上位职业1)
	if currentJob.HighJob1 != "" {
		if highJob, exists := uv.Jobs[currentJob.HighJob1]; exists {
			// 检查转职限制
			canPromote := models.CanPromoteToJob(person, highJob)
			// 使用JobCanRead映射表获取中文名称
			jobName := highJob.StyleName
			if name, ok := models.JobCanRead[highJob.Jid]; ok {
				jobName = name
			}
			button := widget.NewButton("转职到 "+jobName, func() {
				uv.onJobPromotionClicked(currentJob.HighJob1) // false表示不是横转
			})

			// 如果不能转职，则禁用按钮
			if !canPromote {
				button.Disable()
			}

			//如果人物等级不足10级 按钮变灰
			if uv.PersonSimStates[person.Pid].Level < 10 {
				button.Disable()
			}

			uv.JobPromotionContainer.Add(button)
			uv.JobPromotionButtons = append(uv.JobPromotionButtons, button)
			buttonCount++
		}
	}

	// 检查HighJob2 (上位职业2)
	if currentJob.HighJob2 != "" {
		if highJob, exists := uv.Jobs[currentJob.HighJob2]; exists {
			// 检查转职限制
			canPromote := models.CanPromoteToJob(person, highJob)
			// 使用JobCanRead映射表获取中文名称
			jobName := highJob.StyleName
			if name, ok := models.JobCanRead[highJob.Jid]; ok {
				jobName = name
			}
			button := widget.NewButton("转职到 "+jobName, func() {
				uv.onJobPromotionClicked(currentJob.HighJob2) // false表示不是横转
			})

			// 如果不能转职，则禁用按钮
			if !canPromote {
				button.Disable()
			}
			//如果人物等级不足10级 按钮变灰
			if uv.PersonSimStates[person.Pid].Level < 10 {
				button.Disable()
			}

			uv.JobPromotionContainer.Add(button)
			uv.JobPromotionButtons = append(uv.JobPromotionButtons, button)
			buttonCount++
		}
	}

	// 刷新容器
	uv.JobPromotionContainer.Refresh()

	//设置一个已经是上位职业灰色按钮
	if buttonCount == 0 {
		button := widget.NewButton("已经是上位职业", nil)
		button.Disable()
		uv.JobPromotionContainer.Add(button)
		uv.JobPromotionButtons = append(uv.JobPromotionButtons, button)
	}
}

// 检查是否可以转职到目标职业
func (uv *PersonView) onJobPromotionClicked(targetJobID string) {

	if targetJobID == "" {
		return
	}

	targetJobMaxLevel := uv.Jobs[targetJobID].MaxLevel
	simState := uv.PersonSimStates[uv.CurrentPerson.Pid]

	//上转
	if targetJobMaxLevel > simState.MaxLevel {
		simState.MaxLevel = 40

		if simState.Level < 20 {
			simState.Level = 20
		}
		simState.MinLevel = 20
	}

	simState.JobID = targetJobID
	simState.JobBaseStats = uv.Jobs[targetJobID].JobBaseStats
	simState.JobGrowStats = uv.Jobs[targetJobID].JobGrowStats
	simState.JobLimitStats = uv.Jobs[targetJobID].JobLimitStats

	uv.PersonStateHistory[uv.CurrentPerson.Pid].States = append(uv.PersonStateHistory[uv.CurrentPerson.Pid].States, models.PersonStateHistoryItem{
		JobID:         simState.JobID,
		Level:         simState.Level,
		MinLevel:      simState.MinLevel,
		MaxLevel:      simState.MaxLevel,
		TotalPoints:   simState.TotalPoints,
		GrowthCount:   simState.GrowthCount,
		JobBaseStats:  simState.JobBaseStats,
		JobGrowStats:  simState.JobGrowStats,
		JobLimitStats: simState.JobLimitStats,
		EffortTalent:  simState.EffortTalent, // 保留努力才能状态
		Spellbook:     simState.Spellbook,    // 保留术书状态
	})

	uv.PersonJidLabel.SetText(models.LabelPersonJid + models.JobCanRead[simState.JobID])
	uv.LevelLabel.SetText(models.LabelPersonLevel + strconv.Itoa(simState.Level))

	uv.updatePersonHistoryList(uv.CurrentPerson.Pid)

}

// 横转职业选择事件处理
func (uv *PersonView) onLateralJobSelected() {

	// 获取选中的职业名称
	selectedJobName := uv.LateralJobSelect.Selected
	///获取目标职业ID
	targetJobID, _ := uv.LateralJobMap[selectedJobName]

	// 执行横转
	uv.onJobPromotionClicked(targetJobID)
}

// 更新横转职业下拉框
func (uv *PersonView) updateLateralJobSelect(person models.Person) {
	// 清空之前的选择
	uv.LateralJobSelect.Options = []string{}

	simState := uv.PersonSimStates[person.Pid]
	currentJob := uv.Jobs[simState.JobID]
	// 清空映射
	uv.LateralJobMap = make(map[string]string)

	// 获取所有相同等级上限的职业
	var lateralJobs []string

	for _, job := range uv.Jobs {
		// 检查是否具有相同的最高等级
		if job.MaxLevel == currentJob.MaxLevel && job.Jid != currentJob.Jid {
			// 检查转职限制
			if models.CanPromoteToJob(person, job) {
				jobName := job.StyleName
				if name, ok := models.JobCanRead[job.Jid]; ok {
					jobName = name
				}
				lateralJobs = append(lateralJobs, jobName)
				uv.LateralJobMap[jobName] = job.Jid
			}
		}
	}

	// 排序职业名称
	sort.Strings(lateralJobs)

	// 更新下拉框选项
	uv.LateralJobSelect.Options = lateralJobs

	// 刷新下拉框
	uv.LateralJobSelect.Refresh()
}

// 滑动条变更处理
func (uv *PersonView) onLevelSliderChanged(value float64) {
	if uv.CurrentPerson.Pid == "" {
		return
	}

	// 确保值在有效范围内
	if value < uv.LevelSlider.Min {
		value = uv.LevelSlider.Min
	}
	if value > uv.LevelSlider.Max {
		value = uv.LevelSlider.Max
	}

	simState := uv.PersonSimStates[uv.CurrentPerson.Pid]
	targetLevel := int(value)
	initialLevel := simState.Level

	// 如果目标等级与当前模拟等级不同，则更新
	if targetLevel != simState.Level {
		// 检查拖动方向：从左往右拉（升级）还是从右往左拉（降级）
		if targetLevel > initialLevel {
			// 升级逻辑保持不变...
			for simState.Level < targetLevel && simState.Level < simState.MaxLevel {
				// 计算从当前等级到下一等级积分増加
				for i := 0; i < 9; i++ {
					// 计算总成长率
					personGrow := uv.CurrentPerson.PersonGrowStats[i]
					jobGrow := simState.JobGrowStats[i]

					// 如果启用术书，个人职业成长全部+5
					if simState.Spellbook {
						personGrow += 5
					}

					if simState.EffortTalent {
						jobGrow *= 2.0
					}
					// 计算总成长率
					totalGrow := personGrow + jobGrow

					points := float64(totalGrow) / 100.0

					// 如果启用努力才能，积分翻倍

					simState.TotalPoints[i] += points
					threshold := float64(simState.GrowthCount[i]) + 0.99
					if simState.TotalPoints[i] >= threshold {
						totalBase := uv.CurrentPerson.PersonBaseStats[i] + simState.JobBaseStats[i]
						totalLimit := uv.CurrentPerson.PersonLimitStats[i] + simState.JobLimitStats[i]
						currentValue := totalBase + simState.GrowthCount[i]
						if currentValue < totalLimit {
							simState.GrowthCount[i]++
						}
					}
				}
				simState.Level++
				uv.LevelLabel.SetText(models.LabelPersonLevel + strconv.Itoa(simState.Level))
				// 保存当前等级状态到历史记录
				historyItem := models.PersonStateHistoryItem{
					Level:            simState.Level,
					JobID:            simState.JobID,
					MaxLevel:         simState.MaxLevel,
					MinLevel:         simState.MinLevel,
					TotalPoints:      simState.TotalPoints,
					GrowthCount:      simState.GrowthCount,
					PersonBaseStats:  uv.CurrentPerson.PersonBaseStats,
					PersonLimitStats: uv.CurrentPerson.PersonLimitStats,
					PersonGrowStats:  uv.CurrentPerson.PersonGrowStats,
					JobBaseStats:     uv.Jobs[simState.JobID].JobBaseStats,
					JobLimitStats:    uv.Jobs[simState.JobID].JobLimitStats,
					JobGrowStats:     uv.Jobs[simState.JobID].JobGrowStats,
					EffortTalent:     simState.EffortTalent, // 保存努力才能状态
					Spellbook:        simState.Spellbook,    // 保存术书状态
				}
				uv.PersonStateHistory[uv.CurrentPerson.Pid].States = append(
					uv.PersonStateHistory[uv.CurrentPerson.Pid].States,
					historyItem,
				)

				// 更新历史记录列表
				uv.updatePersonHistoryList(uv.CurrentPerson.Pid)
			}
		} else if targetLevel < initialLevel {
			// 降级逻辑：添加长度检查防止索引越界
			for simState.Level > targetLevel {
				// 检查是否有足够的历史记录可以回溯
				states := uv.PersonStateHistory[uv.CurrentPerson.Pid].States
				if len(states) < 2 {
					break // 没有足够的历史记录可以回溯
				}

				lastSecondItem := states[len(states)-2]
				lastItem := states[len(states)-1]

				if lastSecondItem.JobID != lastItem.JobID {
					uv.PersonJidLabel.SetText(models.LabelPersonJid + models.JobCanRead[lastSecondItem.JobID])
				}

				// 使用倒数第二条记录恢复模拟状态
				uv.PersonSimStates[uv.CurrentPerson.Pid] = &models.PersonSimState{
					Level:         lastSecondItem.Level,
					JobID:         lastSecondItem.JobID,
					TotalPoints:   lastSecondItem.TotalPoints,
					GrowthCount:   lastSecondItem.GrowthCount,
					MaxLevel:      lastSecondItem.MaxLevel,
					MinLevel:      lastSecondItem.MinLevel,
					JobLimitStats: lastSecondItem.JobLimitStats,
					JobBaseStats:  lastSecondItem.JobBaseStats,
					JobGrowStats:  lastSecondItem.JobGrowStats,
					Spellbook:     lastSecondItem.Spellbook,
					EffortTalent:  lastSecondItem.EffortTalent,
				}

				simState = uv.PersonSimStates[uv.CurrentPerson.Pid]
				uv.LevelLabel.SetText(models.LabelPersonLevel + strconv.Itoa(lastSecondItem.Level))
				uv.PersonStateHistory[uv.CurrentPerson.Pid].States = states[:len(states)-1]
				uv.updatePersonHistoryList(uv.CurrentPerson.Pid)
			}
		}
	}
}

// 努力才能选项变更处理
func (uv *PersonView) onEffortTalentChanged(checked bool) {
	if uv.CurrentPerson.Pid == "" {
		return
	}

	simState := uv.PersonSimStates[uv.CurrentPerson.Pid]
	if simState != nil {
		simState.EffortTalent = checked

		// 更新历史记录中的最后一个节点而不是新增
		if len(uv.PersonStateHistory[uv.CurrentPerson.Pid].States) > 0 {
			lastIndex := len(uv.PersonStateHistory[uv.CurrentPerson.Pid].States) - 1
			uv.PersonStateHistory[uv.CurrentPerson.Pid].States[lastIndex].EffortTalent = checked
		}

		// 只更新表格，不刷新历史记录列表，避免光标位置变化
		operations.UpdatePersonStatsTable(uv.PersonStatsTable, uv.CurrentPerson, uv.Jobs[simState.JobID], simState)
	}
}

// 术书选项变更处理
func (uv *PersonView) onSpellbookChanged(checked bool) {
	if uv.CurrentPerson.Pid == "" {
		return
	}

	simState := uv.PersonSimStates[uv.CurrentPerson.Pid]
	if simState != nil {
		simState.Spellbook = checked

		// 更新历史记录中的最后一个节点而不是新增
		if len(uv.PersonStateHistory[uv.CurrentPerson.Pid].States) > 0 {
			lastIndex := len(uv.PersonStateHistory[uv.CurrentPerson.Pid].States) - 1
			uv.PersonStateHistory[uv.CurrentPerson.Pid].States[lastIndex].Spellbook = checked
		}

		// 只更新表格，不刷新历史记录列表，避免光标位置变化
		operations.UpdatePersonStatsTable(uv.PersonStatsTable, uv.CurrentPerson, uv.Jobs[simState.JobID], simState)
	}
}
