package models

import (
	_ "embed"
)

//go:embed Job.xml
var JobXMLData []byte

//go:embed Person.xml
var PersonXMLData []byte

const (
	JobXMLFilePath    = ".\\IF mod (Cobalt)\\patches\\xml\\Job.xml"
	PersonXMLFilePath = ".\\IF mod (Cobalt)\\patches\\xml\\Person.xml"
)

// PersonView相关的标签常量
const (
	LabelPersonID        = "人物ID: "
	LabelPersonJid       = "职业: "
	LabelPersonGender    = "性别: "
	LabelPersonLevel     = "等级: "
	LabelUnknown         = "未知"
	LabelMale            = "男"
	LabelFemale          = "女"
	LabelLevelUp         = "升级"
	LabelLevelDown       = "降级"
	LabelReset           = "重置"
	LabelLevelSetup      = "等级设定"
	LabelAttr            = "属性"
	LabelBaseStats       = "基础值"
	LabelGrowStats       = "成长率"
	LabelLimitStats      = "上限値"
	LabelSumStats        = "总和"
	LabelSumExcludeHP    = "不含HP总和"
	LabelPromotion       = "上位转职"
	LabelTransfer        = "横转"
	LabelCancelPromotion = "取消上转"
	LabelCancelTransfer  = "取消横转"
)

const (
	SEX_UNKNOWN = iota
	SEX_MALE
	SEX_FEMALE
)

func CanPromoteToJob(person Person, job Job) bool {
	switch job.Jid {
	case "JID_神竜ノ子", "JID_神竜ノ王": // 神龙之子/神龙之王只能由琉尔转职
		if person.Pid != "PID_リュール" {
			return false
		}
	case "JID_邪竜ノ娘": // 邪龙之女只能由贝尔转职
		if person.Pid != "PID_ヴェイル" {
			return false
		}
	case "JID_アヴニール下級", "JID_アヴニール": // 王室贵族(アルフレッド)
		if person.Pid != "PID_アルフレッド" {
			return false
		}
	case "JID_フロラージュ下級", "JID_フロラージュ": // 王室贵族(锡莉奴)
		if person.Pid != "PID_セリーヌ" {
			return false
		}
	case "JID_ティラユール下級", "JID_ティラユール": // 领主(史塔卢克)
		if person.Pid != "PID_スタルーク" {
			return false
		}
	case "JID_スュクセサール下級", "JID_スュクセサール": // 领主(迪亚曼德)
		if person.Pid != "PID_ディアマンド" {
			return false
		}
	case "JID_ピッチフォーク下級", "JID_ピッチフォーク": // 义警队(蜜丝提拉)
		if person.Pid != "PID_ミスティラ" {
			return false
		}
	case "JID_クピードー下級", "JID_クピードー": // 义警队(佛贾特)
		if person.Pid != "PID_フォガート" {
			return false
		}
	case "JID_リンドブルム下級", "JID_リンドブルム": // 驯兽师(艾比)
		if person.Pid != "PID_アイビー" {
			return false
		}
	case "JID_スレイプニル下級", "JID_スレイプニル": // 驯兽师(奥尔坦西亚)
		if person.Pid != "PID_オルテンシア" {
			return false
		}
	case "JID_メリュジーヌ_味方": // 驯兽师(奥尔坦西亚)
		if person.Pid != "PID_セレスティア" {
			return false
		}
	case "JID_ダンサー":
		if person.Pid != "PID_セアダス" {
			return false
		}
	case "JID_メイド":
		if person.Gender != SEX_FEMALE {
			return false
		}
	case "JID_バトラー":
		if person.Gender != SEX_MALE {
			return false
		}
	}

	return true
}
