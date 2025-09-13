package models

import (
	"EngageCalculator/utils"
	"fmt"
	"strconv"
)

// 模拟状态结构体
type PersonSimState struct {
	Pid           string
	Level         int
	TotalPoints   [9]float64 // 每个属性的总积分
	GrowthCount   [9]int     // 每个属性的增长次数
	JobID         string     // 当前职业ID
	MaxLevel      int        // 当前职业的最大等级
	MinLevel      int
	JobBaseStats  [9]int // 职业基础属性值
	JobLimitStats [9]int // 职业属性上限值
	JobGrowStats  [9]int // 职业成长率

	EffortTalent bool // 努力才能：使职业成长翻倍计算
	Spellbook    bool // 术书：使个人成长全部+5
}

// 人物状态历史记录项
type PersonStateHistoryItem struct {
	Level            int
	JobID            string // 职业ID
	MaxLevel         int
	MinLevel         int
	TotalPoints      [9]float64 // 每个属性的总积分
	GrowthCount      [9]int     // 每个属性的增长次数
	PersonBaseStats  [9]int     // 人物基础属性值
	JobBaseStats     [9]int     // 职业基础属性值
	PersonLimitStats [9]int     // 人物属性上限值
	JobLimitStats    [9]int     // 职业属性上限值
	PersonGrowStats  [9]int     // 人物成长率
	JobGrowStats     [9]int     // 职业成长率

	EffortTalent bool // 努力才能：使职业成长翻倍计算
	Spellbook    bool // 术书：使个人成长全部+5
}

// 人物状态历史记录
type PersonStateHistory struct {
	Pid    string
	States []PersonStateHistoryItem // 按时间顺序存储的状态信息
}

type Person struct {
	Pid              string
	Jid              string
	Gender           int
	Level            int
	SubAptitude      int
	Sid              string
	PersonBaseStats  [9]int // 对应基础属性值
	PersonLimitStats [9]int // 对应属性上限值
	PersonGrowStats  [9]int // 对应属性成长值
}

// 可视化打印
func (p Person) String() string {
	result := fmt.Sprintf("Person ID: %s\n", p.Pid)
	result += fmt.Sprintf("Job ID: %s\n", p.Jid)
	result += fmt.Sprintf("Gender: %d\n", p.Gender)
	result += fmt.Sprintf("Level: %d\n", p.Level)
	result += fmt.Sprintf("Sub Aptitude: %d\n", p.SubAptitude)
	result += fmt.Sprintf("SID: %s\n", p.Sid)

	result += "Base Stats:\n"
	for i, stat := range p.PersonBaseStats {
		if name, ok := StatNames[i]; ok {
			result += fmt.Sprintf("  %s: %d\n", name, stat)
		}
	}

	result += "Limit Stats:\n"
	for i, stat := range p.PersonLimitStats {
		if name, ok := StatNames[i]; ok {
			result += fmt.Sprintf("  %s: %d\n", name, stat)
		}
	}

	result += "Grow Stats:\n"
	for i, stat := range p.PersonGrowStats {
		if name, ok := StatNames[i]; ok {
			result += fmt.Sprintf("  %s: %d\n", name, stat)
		}
	}

	return result
}

func ProcessPersonLine(line string) Person {
	pid := utils.ExtractAttribute(line, PidRegex)
	jid := utils.ExtractAttribute(line, JidRegex)
	sid := utils.ExtractAttribute(line, SidRegex)

	// 转换整数属性
	genderStr := utils.ExtractAttribute(line, GenderRegex)
	gender, _ := strconv.Atoi(genderStr)

	levelStr := utils.ExtractAttribute(line, LevelRegex)
	level, _ := strconv.Atoi(levelStr)

	subAptitudeStr := utils.ExtractAttribute(line, SubAptitudeRegex)
	subAptitude, _ := strconv.Atoi(subAptitudeStr)

	// 初始化基础属性数组
	var baseStats [9]int
	baseStats[StatHP] = utils.ParseIntAttribute(line, PersonBaseHPRegex)
	baseStats[StatStr] = utils.ParseIntAttribute(line, PersonBaseStrRegex)
	baseStats[StatTech] = utils.ParseIntAttribute(line, PersonBaseTechRegex)
	baseStats[StatQuick] = utils.ParseIntAttribute(line, PersonBaseQuickRegex)
	baseStats[StatLuck] = utils.ParseIntAttribute(line, PersonBaseLuckRegex)
	baseStats[StatMagic] = utils.ParseIntAttribute(line, PersonBaseMagicRegex)
	baseStats[StatDef] = utils.ParseIntAttribute(line, PersonBaseDefRegex)
	baseStats[StatMdef] = utils.ParseIntAttribute(line, PersonBaseMdefRegex)
	baseStats[StatMove] = utils.ParseIntAttribute(line, PersonBaseMoveRegex)

	// 初始化属性上限数组
	var limitStats [9]int
	limitStats[StatHP] = utils.ParseIntAttribute(line, PersonLimitHPRegex)
	limitStats[StatStr] = utils.ParseIntAttribute(line, PersonLimitStrRegex)
	limitStats[StatTech] = utils.ParseIntAttribute(line, PersonLimitTechRegex)
	limitStats[StatQuick] = utils.ParseIntAttribute(line, PersonLimitQuickRegex)
	limitStats[StatLuck] = utils.ParseIntAttribute(line, PersonLimitLuckRegex)
	limitStats[StatMagic] = utils.ParseIntAttribute(line, PersonLimitMagicRegex)
	limitStats[StatDef] = utils.ParseIntAttribute(line, PersonLimitDefRegex)
	limitStats[StatMdef] = utils.ParseIntAttribute(line, PersonLimitMdefRegex)
	limitStats[StatMove] = utils.ParseIntAttribute(line, PersonLimitMoveRegex)

	// 初始化属性成长数组
	var growStats [9]int
	growStats[StatHP] = utils.ParseIntAttribute(line, PersonGrowHPRegex)
	growStats[StatStr] = utils.ParseIntAttribute(line, PersonGrowStrRegex)
	growStats[StatTech] = utils.ParseIntAttribute(line, PersonGrowTechRegex)
	growStats[StatQuick] = utils.ParseIntAttribute(line, PersonGrowQuickRegex)
	growStats[StatLuck] = utils.ParseIntAttribute(line, PersonGrowLuckRegex)
	growStats[StatMagic] = utils.ParseIntAttribute(line, PersonGrowMagicRegex)
	growStats[StatDef] = utils.ParseIntAttribute(line, PersonGrowDefRegex)
	growStats[StatMdef] = utils.ParseIntAttribute(line, PersonGrowMdefRegex)
	growStats[StatMove] = utils.ParseIntAttribute(line, PersonGrowMoveRegex)

	return Person{
		Pid:              pid,
		Jid:              jid,
		Gender:           gender,
		Level:            level,
		SubAptitude:      subAptitude,
		Sid:              sid,
		PersonBaseStats:  baseStats,
		PersonLimitStats: limitStats,
		PersonGrowStats:  growStats,
	}
}
