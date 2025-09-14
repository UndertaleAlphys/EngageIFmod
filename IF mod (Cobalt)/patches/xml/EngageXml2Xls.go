package xml

import (
	"encoding/xml"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/tealeg/xlsx"
)

// Book XML结构体定义
type Book struct {
	XMLName xml.Name `xml:"Book"`
	Count   string   `xml:"Count,attr"`
	Sheets  []Sheet  `xml:"Sheet"`
}

// Sheet XML结构体定义
type Sheet struct {
	Name  string  `xml:"Name,attr"`
	Count string  `xml:"Count,attr"`
	Header Header  `xml:"Header"`
	Data   Data    `xml:"Data"`
}

// Header XML结构体定义
type Header struct {
	Params []Param `xml:"Param"`
}

// Data XML结构体定义
type Data struct {
	Params []Param `xml:"Param"`
}

// Param XML结构体定义
type Param struct {
	XMLName xml.Name
	Attrs   []xml.Attr `xml:",any,attr"`
}

// UnmarshalXML 实现自定义XML解组
func (p *Param) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	p.XMLName = start.Name
	p.Attrs = start.Attr
	// 消费元素的结束标签
	_, err := d.Token()
	return err
}

// Xls2Xml 将XLS文件转换为XML文件
func Xls2Xml(xlsFilePath string, xmlFilePath string) error {
	// 打开XLSX文件
	xlFile, err := xlsx.OpenFile(xlsFilePath)
	if err != nil {
		return err
	}

	// 创建XML文件
	xmlFile, err := os.Create(xmlFilePath)
	if err != nil {
		return err
	}
	defer xmlFile.Close()

	// 写入XML头部
	xmlFile.WriteString(`<?xml version="1.0" encoding="utf-8"?>` + "\n")
	xmlFile.WriteString(`<Book Count="` + fmt.Sprintf("%d", len(xlFile.Sheets)) + `">` + "\n")

	// 遍历所有工作表
	for sheetName, sheet := range xlFile.Sheets {
		xmlSheet := Sheet{
			Name:  sheetName,
			Count: fmt.Sprintf("%d", sheet.MaxRow-1), // 减去标题行
		}

		if sheet.MaxRow > 0 {
			// 处理Header
			headerRow := sheet.Rows[0]
			for _, cell := range headerRow.Cells {
				value, _ := cell.String()
				xmlSheet.Header.Params = append(xmlSheet.Header.Params, Param{
					XMLName: xml.Name{Local: "Param"},
					Attrs: []xml.Attr{
						{Name: xml.Name{Local: "Name"}, Value: value},
						{Name: xml.Name{Local: "Ident"}, Value: value},
					},
				})
			}

			// 处理Data
			for i := 1; i < sheet.MaxRow; i++ {
				dataRow := sheet.Rows[i]
				dataParam := Param{
					XMLName: xml.Name{Local: "Param"},
					Attrs:   make([]xml.Attr, 0),
				}
				
				for j, cell := range dataRow.Cells {
					if j < len(xmlSheet.Header.Params) {
						value, _ := cell.String()
						ident := xmlSheet.Header.Params[j].Attrs[1].Value // 获取Ident值
						dataParam.Attrs = append(dataParam.Attrs, xml.Attr{
							Name:  xml.Name{Local: ident},
							Value: value,
						})
					}
				}
				
				xmlSheet.Data.Params = append(xmlSheet.Data.Params, dataParam)
			}
		}

		book.Sheets = append(book.Sheets, xmlSheet)
	}

	// 手动写入Sheet内容
	for _, sheet := range book.Sheets {
		// 写入Sheet头部
		xmlFile.WriteString("  <Sheet Name=\"" + sheet.Name + "\" Count=\"" + sheet.Count + "\">" + "\n")
		
		// 写入Header
		xmlFile.WriteString("    <Header>" + "\n")
		for _, header := range sheet.Header.Params {
			xmlFile.WriteString("      <Param Name=\"" + getAttrValue(header.Attrs, "Name") + "\" Ident=\"" + getAttrValue(header.Attrs, "Ident") + "\">" + "\n")
		}
		xmlFile.WriteString("    </Header>" + "\n")
		
		// 写入Data
		xmlFile.WriteString("    <Data>" + "\n")
		for _, data := range sheet.Data.Params {
			xmlFile.WriteString("      <Param")
			for _, attr := range data.Attrs {
				xmlFile.WriteString(" " + attr.Name.Local + "=\"" + attr.Value + "\"")
			}
			xmlFile.WriteString("/>" + "\n")
		}
		xmlFile.WriteString("    </Data>" + "\n")
		
		// 写入Sheet尾部
		xmlFile.WriteString("  </Sheet>" + "\n")
	}
	
	// 写入Book尾部
	xmlFile.WriteString("</Book>" + "\n")

	return nil
}

// EngageXml2Xls 将XML文件转换为XLS文件
func EngageXml2Xls(xmlFilePath string, xlsFilePath string) error {
	// 打开XML文件
	xmlFile, err := os.Open(xmlFilePath)
	if err != nil {
		return err
	}
	defer xmlFile.Close()

	// 解析XML
	var book Book
	decoder := xml.NewDecoder(xmlFile)
	err = decoder.Decode(&book)
	if err != nil {
		return err
	}

	// 创建XLSX文件
	xlsxFile := xlsx.NewFile()

	// 遍历所有Sheet
	for _, sheet := range book.Sheets {
		// 创建工作表
		xlSheet, err := xlsxFile.AddSheet(sheet.Name)
		if err != nil {
			return err
		}

		// 创建标题行
		if len(sheet.Header.Params) > 0 {
			headerRow := xlSheet.AddRow()
			for _, param := range sheet.Header.Params {
				// 查找Ident属性
				ident := ""
				for _, attr := range param.Attrs {
					if attr.Name.Local == "Ident" {
						ident = attr.Value
						break
					}
				}
				headerRow.AddCell().SetValue(ident)
			}
		}

		// 添加数据行
		for _, param := range sheet.Data.Params {
			dataRow := xlSheet.AddRow()
			// 为每个标题创建单元格
			for _, headerParam := range sheet.Header.Params {
				// 查找Ident属性
				ident := ""
				for _, attr := range headerParam.Attrs {
					if attr.Name.Local == "Ident" {
						ident = attr.Value
						break
					}
				}

				// 在数据参数中查找匹配的属性
				value := ""
				for _, attr := range param.Attrs {
					if attr.Name.Local == ident {
						value = attr.Value
						break
					}
				}
				dataRow.AddCell().SetValue(value)
			}
		}
	}

	// 保存XLSX文件
	return xlsxFile.Save(xlsFilePath)
}

// getAttrValue 获取指定属性的值
func getAttrValue(attrs []xml.Attr, name string) string {
	for _, attr := range attrs {
		if attr.Name.Local == name {
			return attr.Value
		}
	}
	return ""
}

// Main方法实现
func main() {
	// 获取当前目录
	currentDir, err := os.Getwd()
	if err != nil {
		fmt.Printf("获取当前目录失败: %v\n", err)
		return
	}

	// 查找当前目录下的所有XML文件
	files, err := os.ReadDir(currentDir)
	if err != nil {
		fmt.Printf("读取目录失败: %v\n", err)
		return
	}

	// 处理每个XML文件
	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(file.Name(), ".xml") {
			xmlFilePath := filepath.Join(currentDir, file.Name())
			xlsFilePath := strings.TrimSuffix(xmlFilePath, ".xml") + ".xlsx"
			
			fmt.Printf("正在转换 %s 到 %s...\n", file.Name(), xlsFilePath)
			err := EngageXml2Xls(xmlFilePath, xlsFilePath)
			if err != nil {
				fmt.Printf("转换失败 %s: %v\n", file.Name(), err)
			} else {
				fmt.Printf("转换成功 %s\n", file.Name())
			}
		}
	}

	fmt.Println("所有XML文件已转换完成。")
}