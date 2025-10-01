package models

import "regexp"

var (
	PidRegex         = regexp.MustCompile(`Pid="([^"]*)"`)
	JidRegex         = regexp.MustCompile(`Jid="([^"]*)"`)
	GenderRegex      = regexp.MustCompile(`Gender="([^"]*)"`)
	StyleNameRegex   = regexp.MustCompile(`StyleName="([^"]*)"`)
	LevelRegex       = regexp.MustCompile(`Level="([^"]*)"`)
	SubAptitudeRegex = regexp.MustCompile(`SubAptitude="([^"]*)"`)
	SidRegex         = regexp.MustCompile(`CommonSids="([^"]*);"`)
	HighJob1Regex    = regexp.MustCompile(`HighJob1="([^"]*)"`)
	HighJob2Regex    = regexp.MustCompile(`HighJob2="([^"]*)"`)

	MaxLevelRegex              = regexp.MustCompile(`MaxLevel="([^"]*)"`)
	MaxWeaponLevelSwordRegex   = regexp.MustCompile(`MaxWeaponLevelSword="([^"]*)"`)
	MaxWeaponLevelLanceRegex   = regexp.MustCompile(`MaxWeaponLevelLance="([^"]*)"`)
	MaxWeaponLevelAxeRegex     = regexp.MustCompile(`MaxWeaponLevelAxe="([^"]*)"`)
	MaxWeaponLevelBowRegex     = regexp.MustCompile(`MaxWeaponLevelBow="([^"]*)"`)
	MaxWeaponLevelDaggerRegex  = regexp.MustCompile(`MaxWeaponLevelDagger="([^"]*)"`)
	MaxWeaponLevelMagicRegex   = regexp.MustCompile(`MaxWeaponLevelMagic="([^"]*)"`)
	MaxWeaponLevelRodRegex     = regexp.MustCompile(`MaxWeaponLevelRod="([^"]*)"`)
	MaxWeaponLevelFistRegex    = regexp.MustCompile(`MaxWeaponLevelFist="([^"]*)"`)
	MaxWeaponLevelSpecialRegex = regexp.MustCompile(`MaxWeaponLevelSpecial="([^"]*)"`)

	PersonBaseHPRegex    = regexp.MustCompile(`OffsetL\.Hp="([^"]*)"`)
	PersonBaseStrRegex   = regexp.MustCompile(`OffsetL\.Str="([^"]*)"`)
	PersonBaseTechRegex  = regexp.MustCompile(`OffsetL\.Tech="([^"]*)"`)
	PersonBaseQuickRegex = regexp.MustCompile(`OffsetL\.Quick="([^"]*)"`)
	PersonBaseLuckRegex  = regexp.MustCompile(`OffsetL\.Luck="([^"]*)"`)
	PersonBaseMagicRegex = regexp.MustCompile(`OffsetL\.Magic="([^"]*)"`)
	PersonBaseDefRegex   = regexp.MustCompile(`OffsetL\.Def="([^"]*)"`)
	PersonBaseMdefRegex  = regexp.MustCompile(`OffsetL\.Mdef="([^"]*)"`)
	PersonBaseMoveRegex  = regexp.MustCompile(`OffsetL\.Move="([^"]*)"`)

	JobBaseHPRegex    = regexp.MustCompile(`Base\.Hp="([^"]*)"`)
	JobBaseStrRegex   = regexp.MustCompile(`Base\.Str="([^"]*)"`)
	JobBaseTechRegex  = regexp.MustCompile(`Base\.Tech="([^"]*)"`)
	JobBaseQuickRegex = regexp.MustCompile(`Base\.Quick="([^"]*)"`)
	JobBaseLuckRegex  = regexp.MustCompile(`Base\.Luck="([^"]*)"`)
	JobBaseMagicRegex = regexp.MustCompile(`Base\.Magic="([^"]*)"`)
	JobBaseDefRegex   = regexp.MustCompile(`Base\.Def="([^"]*)"`)
	JobBaseMdefRegex  = regexp.MustCompile(`Base\.Mdef="([^"]*)"`)
	JobBaseMoveRegex  = regexp.MustCompile(`Base\.Move="([^"]*)"`)

	PersonLimitHPRegex    = regexp.MustCompile(`Limit\.Hp="([^"]*)"`)
	PersonLimitStrRegex   = regexp.MustCompile(`Limit\.Str="([^"]*)"`)
	PersonLimitTechRegex  = regexp.MustCompile(`Limit\.Tech="([^"]*)"`)
	PersonLimitQuickRegex = regexp.MustCompile(`Limit\.Quick="([^"]*)"`)
	PersonLimitLuckRegex  = regexp.MustCompile(`Limit\.Luck="([^"]*)"`)
	PersonLimitMagicRegex = regexp.MustCompile(`Limit\.Magic="([^"]*)"`)
	PersonLimitDefRegex   = regexp.MustCompile(`Limit\.Def="([^"]*)"`)
	PersonLimitMdefRegex  = regexp.MustCompile(`Limit\.Mdef="([^"]*)"`)
	PersonLimitMoveRegex  = regexp.MustCompile(`Limit\.Move="([^"]*)"`)

	JobLimitHPRegex    = regexp.MustCompile(`Limit\.Hp="([^"]*)"`)
	JobLimitStrRegex   = regexp.MustCompile(`Limit\.Str="([^"]*)"`)
	JobLimitTechRegex  = regexp.MustCompile(`Limit\.Tech="([^"]*)"`)
	JobLimitQuickRegex = regexp.MustCompile(`Limit\.Quick="([^"]*)"`)
	JobLimitLuckRegex  = regexp.MustCompile(`Limit\.Luck="([^"]*)"`)
	JobLimitMagicRegex = regexp.MustCompile(`Limit\.Magic="([^"]*)"`)
	JobLimitDefRegex   = regexp.MustCompile(`Limit\.Def="([^"]*)"`)
	JobLimitMdefRegex  = regexp.MustCompile(`Limit\.Mdef="([^"]*)"`)
	JobLimitMoveRegex  = regexp.MustCompile(`Limit\.Move="([^"]*)"`)

	PersonGrowHPRegex    = regexp.MustCompile(`Grow\.Hp="([^"]*)"`)
	PersonGrowStrRegex   = regexp.MustCompile(`Grow\.Str="([^"]*)"`)
	PersonGrowTechRegex  = regexp.MustCompile(`Grow\.Tech="([^"]*)"`)
	PersonGrowQuickRegex = regexp.MustCompile(`Grow\.Quick="([^"]*)"`)
	PersonGrowLuckRegex  = regexp.MustCompile(`Grow\.Luck="([^"]*)"`)
	PersonGrowMagicRegex = regexp.MustCompile(`Grow\.Magic="([^"]*)"`)
	PersonGrowDefRegex   = regexp.MustCompile(`Grow\.Def="([^"]*)"`)
	PersonGrowMdefRegex  = regexp.MustCompile(`Grow\.Mdef="([^"]*)"`)
	PersonGrowMoveRegex  = regexp.MustCompile(`Grow\.Move="([^"]*)"`)

	JobGrowHPRegex    = regexp.MustCompile(`DiffGrow\.Hp="([^"]*)"`)
	JobGrowStrRegex   = regexp.MustCompile(`DiffGrow\.Str="([^"]*)"`)
	JobGrowTechRegex  = regexp.MustCompile(`DiffGrow\.Tech="([^"]*)"`)
	JobGrowQuickRegex = regexp.MustCompile(`DiffGrow\.Quick="([^"]*)"`)
	JobGrowLuckRegex  = regexp.MustCompile(`DiffGrow\.Luck="([^"]*)"`)
	JobGrowMagicRegex = regexp.MustCompile(`DiffGrow\.Magic="([^"]*)"`)
	JobGrowDefRegex   = regexp.MustCompile(`DiffGrow\.Def="([^"]*)"`)
	JobGrowMdefRegex  = regexp.MustCompile(`DiffGrow\.Mdef="([^"]*)"`)
	JobGrowMoveRegex  = regexp.MustCompile(`DiffGrow\.Move="([^"]*)"`)
)

// 武器类型索引常量
const (
	WeaponSword = iota
	WeaponLance
	WeaponAxe
	WeaponBow
	WeaponDagger
	WeaponMagic
	WeaponRod
	WeaponFist
	WeaponSpecial
)

// 属性类型索引常量
const (
	StatHP = iota
	StatStr
	StatTech
	StatQuick
	StatLuck
	StatMagic
	StatDef
	StatMdef
	StatMove
)

// 武器名称映射
var WeaponNames = map[int]string{
	WeaponSword:   "Sword",
	WeaponLance:   "Lance",
	WeaponAxe:     "Axe",
	WeaponBow:     "Bow",
	WeaponDagger:  "Dagger",
	WeaponMagic:   "Magic",
	WeaponRod:     "Rod",
	WeaponFist:    "Fist",
	WeaponSpecial: "Special",
}

// 属性名称映射
var StatNames = map[int]string{
	StatHP:    "HP",
	StatStr:   "Str",
	StatTech:  "Tech",
	StatQuick: "Quick",
	StatLuck:  "Luck",
	StatMagic: "Magic",
	StatDef:   "Def",
	StatMdef:  "Mdef",
	StatMove:  "Move",
}

// PersonCanRead 定义可读的Person列表
// key: PID, value: 人物名称
var PersonCanRead = map[string]string{
	"PID_リュール":   "琉尔",
	"PID_ヴァンドレ":  "凡德雷",
	"PID_クラン":    "柯岚",
	"PID_フラン":    "芙兰",
	"PID_アルフレッド": "阿尔弗雷德",
	"PID_エーティエ":  "艾提耶",
	"PID_ブシュロン":  "布修隆",
	"PID_セリーヌ":   "锡莉奴",
	"PID_クロエ":    "克罗艾",
	"PID_ルイ":     "路易",
	"PID_ユナカ":    "尤娜卡",
	"PID_スタルーク":  "史塔卢克",
	"PID_シトリニカ":  "希特莉妮卡",
	"PID_ラピス":    "拉琵思",
	"PID_ディアマンド": "迪亚曼德",
	"PID_アンバー":   "安巴",
	"PID_ジェーデ":   "洁德",
	"PID_アイビー":   "艾比",
	"PID_カゲツ":    "花月",
	"PID_ゼルコバ":   "杰尔柯巴",
	"PID_フォガート":  "佛贾特",
	"PID_パンドロ":   "庞德罗",
	"PID_ボネ":     "波聂",
	"PID_ミスティラ":  "蜜丝提拉",
	"PID_パネトネ":   "帕涅托涅",
	"PID_メリン":    "玫琳",
	"PID_オルテンシア": "奥尔坦西亚",
	"PID_セアダス":   "赛安达斯",
	"PID_ロサード":   "罗萨德",
	"PID_ゴルドマリー": "戈尔德玛丽",
	"PID_リンデン":   "霖丹",
	"PID_ザフィーア":  "扎菲亚",
	"PID_ヴェイル":   "贝珥",
	"PID_モーヴ":    "莫布",
	"PID_アンナ":    "安娜",
	"PID_ジャン":    "贾恩",
	"PID_マデリーン":  "玛德琳",
	"PID_セレスティア": "塞勒斯提亚",
	"PID_グレゴリー":  "孤雷葛里",
}

var JobCanRead = map[string]string{
	"JID_神竜ノ子":      "神龙之子",
	"JID_神竜ノ王":      "神龙之王",
	"JID_邪竜ノ娘":      "邪龙之女",
	"JID_アヴニール下級":   "王室贵族(阿尔弗雷德)",
	"JID_アヴニール":     "阿巴尼尔",
	"JID_フロラージュ下級":  "王室贵族(锡莉奴)",
	"JID_フロラージュ":    "弗罗刺祝",
	"JID_スュクセサール下級": "领主(迪亚曼德)",
	"JID_スュクセサール":   "继承者",
	"JID_ティラユール下級":  "领主(史塔卢克)",
	"JID_ティラユール":    "神射手",
	"JID_ピッチフォーク下級": "义警队(蜜思提拉)",
	"JID_ピッチフォーク":   "婆娑罗",
	"JID_クピードー下級":   "义警队(佛贾特)",
	"JID_クピードー":     "邱比特",
	"JID_リンドブルム下級":  "驯兽师(艾比)",
	"JID_リンドブルム":    "林德沃姆尔",
	"JID_スレイプニル下級":  "驯兽师(奥尔坦西亚)",
	"JID_スレイプニル":    "斯雷普尼尔",
	"JID_ソードマスター":   "剑圣",
	"JID_ソードファイター":  "武士",
	"JID_上忍":        "上忍",
	"JID_忍":         "忍者",
	"JID_ウルフナイト":    "狼骑士",
	"JID_ベルセルク":     "狂战士",
	"JID_アクスファイター":  "战士",
	"JID_ブレイブヒーロー":  "勇者",
	"JID_傭兵":        "佣兵",
	"JID_ボウナイト":     "弓骑士",
	"JID_シーフ":       "盗贼",
	"JID_アドベンチャラー":  "冒险家",
	"JID_スナイパー":     "弓圣",
	"JID_アーチャー":     "弓箭手",
	"JID_金鵄武者":      "金鸱武者",
	"JID_ソードペガサス":   "剑飞马骑士",
	"JID_ランスペガサス":   "枪飞马骑士",
	"JID_アクスペガサス":   "斧飞马骑士",
	"JID_グリフォンナイト":  "格里芬骑士",
	"JID_パラディン":     "圣骑士",
	"JID_ソードナイト":    "剑骑士",
	"JID_ランスナイト":    "枪骑士",
	"JID_アクスナイト":    "斧骑士",
	"JID_グレートナイト":   "大骑士",
	"JID_ソードアーマー":   "剑重甲骑士",
	"JID_ランスアーマー":   "枪重甲骑士",
	"JID_アクスアーマー":   "斧重甲骑士",
	"JID_ジェネラル":     "将军",
	"JID_マージナイト":    "黑暗骑士",
	"JID_ダークマージ":    "暗法师",
	"JID_セイジ":       "巫师",
	"JID_マージ":       "魔法师",
	"JID_ハイプリースト":   "主教",
	"JID_ハルバーディア":   "枪圣",
	"JID_ランスファイター":  "枪战士",
	"JID_ロイヤルナイト":   "皇家骑士",
	"JID_レヴナントナイト":  "魔龙骑士",
	"JID_ドラゴンナイト":   "龙骑士",
	"JID_ドラゴンマスター":  "龙骑大师",
	"JID_ストラテジスト":   "指挥骑士",
	"JID_ロッドナイト":    "神官骑士",
	"JID_メイド":       "女仆",
	"JID_バトラー":      "管家",
	"JID_修羅":        "修罗",
	"JID_モンク":       "僧侣",
	"JID_マスターモンク":   "僧侣大师",
	"JID_ダンサー":      "舞蹈家",
	//"JID_メリュジーヌ_味方": "美露莘",
}
