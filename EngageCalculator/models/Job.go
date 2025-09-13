package models

import (
	"EngageCalculator/utils"
	"fmt"
	"strconv"
)

type Job struct {
	Jid             string
	StyleName       string
	MaxLevel        int
	HighJob1        string
	HighJob2        string
	MaxWeaponLevels [11]string // 对应各种武器的最大等级
	JobBaseStats    [9]int     // 对应基础属性值
	JobLimitStats   [9]int     // 对应属性上限值
	JobGrowStats    [9]int     // 对应属性成长值
}

// 可视化打印
func (j Job) String() string {
	result := fmt.Sprintf("Job ID: %s\n", j.Jid)
	result += fmt.Sprintf("Style Name: %s\n", j.StyleName)
	result += fmt.Sprintf("High Job 1: %s\n", j.HighJob1)
	result += fmt.Sprintf("High Job 2: %s\n", j.HighJob2)
	result += fmt.Sprintf("Max Level: %d\n", j.MaxLevel)

	result += "Max Weapon Levels:\n"
	for i, level := range j.MaxWeaponLevels {
		if name, ok := WeaponNames[i]; ok {
			result += fmt.Sprintf("  %s: %s\n", name, level)
		}
	}

	result += "Base Stats:\n"
	for i, stat := range j.JobBaseStats {
		if name, ok := StatNames[i]; ok {
			result += fmt.Sprintf("  %s: %d\n", name, stat)
		}
	}

	result += "JobLimit Stats:\n"
	for i, stat := range j.JobLimitStats {
		if name, ok := StatNames[i]; ok {
			result += fmt.Sprintf("  %s: %d\n", name, stat)
		}
	}

	result += "JobGrow Stats:\n"
	for i, stat := range j.JobGrowStats {
		if name, ok := StatNames[i]; ok {
			result += fmt.Sprintf("  %s: %d\n", name, stat)
		}
	}

	return result
}

func ProcessJobLine(line string) Job {
	jid := utils.ExtractAttribute(line, JidRegex)
	styleName := utils.ExtractAttribute(line, StyleNameRegex)
	highJob1 := utils.ExtractAttribute(line, HighJob1Regex)
	highJob2 := utils.ExtractAttribute(line, HighJob2Regex)

	maxLevelStr := utils.ExtractAttribute(line, MaxLevelRegex)
	maxLevel, _ := strconv.Atoi(maxLevelStr)

	// 初始化武器等级数组
	var maxWeaponLevels [11]string
	maxWeaponLevels[WeaponSword] = utils.ExtractAttribute(line, MaxWeaponLevelSwordRegex)
	maxWeaponLevels[WeaponLance] = utils.ExtractAttribute(line, MaxWeaponLevelLanceRegex)
	maxWeaponLevels[WeaponAxe] = utils.ExtractAttribute(line, MaxWeaponLevelAxeRegex)
	maxWeaponLevels[WeaponBow] = utils.ExtractAttribute(line, MaxWeaponLevelBowRegex)
	maxWeaponLevels[WeaponDagger] = utils.ExtractAttribute(line, MaxWeaponLevelDaggerRegex)
	maxWeaponLevels[WeaponMagic] = utils.ExtractAttribute(line, MaxWeaponLevelMagicRegex)
	maxWeaponLevels[WeaponRod] = utils.ExtractAttribute(line, MaxWeaponLevelRodRegex)
	maxWeaponLevels[WeaponFist] = utils.ExtractAttribute(line, MaxWeaponLevelFistRegex)
	maxWeaponLevels[WeaponSpecial] = utils.ExtractAttribute(line, MaxWeaponLevelSpecialRegex)

	// 初始化基础属性数组
	var JobBaseStats [9]int
	JobBaseStats[StatHP] = utils.ParseIntAttribute(line, JobBaseHPRegex)
	JobBaseStats[StatStr] = utils.ParseIntAttribute(line, JobBaseStrRegex)
	JobBaseStats[StatTech] = utils.ParseIntAttribute(line, JobBaseTechRegex)
	JobBaseStats[StatQuick] = utils.ParseIntAttribute(line, JobBaseQuickRegex)
	JobBaseStats[StatLuck] = utils.ParseIntAttribute(line, JobBaseLuckRegex)
	JobBaseStats[StatMagic] = utils.ParseIntAttribute(line, JobBaseMagicRegex)
	JobBaseStats[StatDef] = utils.ParseIntAttribute(line, JobBaseDefRegex)
	JobBaseStats[StatMdef] = utils.ParseIntAttribute(line, JobBaseMdefRegex)
	JobBaseStats[StatMove] = utils.ParseIntAttribute(line, JobBaseMoveRegex)

	// 初始化属性上限数组
	var JobLimitStats [9]int
	JobLimitStats[StatHP] = utils.ParseIntAttribute(line, JobLimitHPRegex)
	JobLimitStats[StatStr] = utils.ParseIntAttribute(line, JobLimitStrRegex)
	JobLimitStats[StatTech] = utils.ParseIntAttribute(line, JobLimitTechRegex)
	JobLimitStats[StatQuick] = utils.ParseIntAttribute(line, JobLimitQuickRegex)
	JobLimitStats[StatLuck] = utils.ParseIntAttribute(line, JobLimitLuckRegex)
	JobLimitStats[StatMagic] = utils.ParseIntAttribute(line, JobLimitMagicRegex)
	JobLimitStats[StatDef] = utils.ParseIntAttribute(line, JobLimitDefRegex)
	JobLimitStats[StatMdef] = utils.ParseIntAttribute(line, JobLimitMdefRegex)
	JobLimitStats[StatMove] = utils.ParseIntAttribute(line, JobLimitMoveRegex)

	// 初始化属性成长数组
	var JobGrowStats [9]int
	JobGrowStats[StatHP] = utils.ParseIntAttribute(line, JobGrowHPRegex)
	JobGrowStats[StatStr] = utils.ParseIntAttribute(line, JobGrowStrRegex)
	JobGrowStats[StatTech] = utils.ParseIntAttribute(line, JobGrowTechRegex)
	JobGrowStats[StatQuick] = utils.ParseIntAttribute(line, JobGrowQuickRegex)
	JobGrowStats[StatLuck] = utils.ParseIntAttribute(line, JobGrowLuckRegex)
	JobGrowStats[StatMagic] = utils.ParseIntAttribute(line, JobGrowMagicRegex)
	JobGrowStats[StatDef] = utils.ParseIntAttribute(line, JobGrowDefRegex)
	JobGrowStats[StatMdef] = utils.ParseIntAttribute(line, JobGrowMdefRegex)
	JobGrowStats[StatMove] = utils.ParseIntAttribute(line, JobGrowMoveRegex)

	return Job{
		Jid:             jid,
		StyleName:       styleName,
		HighJob1:        highJob1,
		HighJob2:        highJob2,
		MaxLevel:        maxLevel,
		MaxWeaponLevels: maxWeaponLevels,
		JobBaseStats:    JobBaseStats,
		JobLimitStats:   JobLimitStats,
		JobGrowStats:    JobGrowStats,
	}
}
