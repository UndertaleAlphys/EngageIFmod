// XmlReader.go
package utils

import (
	"encoding/xml"
	"io"
	"os"
	"regexp"
	"strconv"
)

func ReadXMLFile(filePath string) (*xml.Decoder, *os.File, error) {
	// 打开文件
	file, err := os.Open(filePath)
	if err != nil {
		return nil, nil, err
	}

	// 创建XML解码器
	decoder := xml.NewDecoder(file)

	return decoder, file, nil
}

func ParseXMLElements(decoder *xml.Decoder, handler func(xml.StartElement) error) error {
	for {
		token, err := decoder.Token()
		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}

		if se, ok := token.(xml.StartElement); ok {
			if err := handler(se); err != nil {
				return err
			}
		}
	}
	return nil
}

func CloseXMLFile(file *os.File) error {
	return file.Close()
}

func ExtractAttribute(line string, re *regexp.Regexp) string {
	matches := re.FindStringSubmatch(line)
	if len(matches) > 1 {
		return matches[1]
	}
	return ""
}

func ParseIntAttribute(line string, regex *regexp.Regexp) int {
	str := ExtractAttribute(line, regex)
	val, _ := strconv.Atoi(str)
	return val
}

func GetMapKeys(m map[string]string) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}

func GetOrderedPersonKeys(personCanRead map[string]string) []string {
	// 定义PID的预设顺序
	orderedPIDs := []string{
		"PID_リュール", "PID_ヴァンドレ", "PID_クラン", "PID_フラン",
		"PID_アルフレッド", "PID_エーティエ", "PID_ブシュロン", "PID_セリーヌ",
		"PID_クロエ", "PID_ルイ", "PID_ユナカ", "PID_スタルーク",
		"PID_シトリニカ", "PID_ラピス", "PID_ディアマンド", "PID_アンバー",
		"PID_ジェーデ", "PID_アイビー", "PID_カゲツ", "PID_ゼルコバ",
		"PID_フォガート", "PID_パンドロ", "PID_ボネ", "PID_ミスティラ",
		"PID_パネトネ", "PID_メリン", "PID_オルテンシア", "PID_セアダス",
		"PID_ロサード", "PID_ゴルドマリー", "PID_リンデン", "PID_ザフィーア",
		"PID_ヴェイル", "PID_モーヴ", "PID_アンナ", "PID_ジャン",
	}

	// 过滤出实际存在于personCanRead中的PID
	var result []string
	for _, pid := range orderedPIDs {
		if _, exists := personCanRead[pid]; exists {
			result = append(result, pid)
		}
	}

	return result
}
