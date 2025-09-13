package controllers

import (
	"EngageCalculator/models"
	"EngageCalculator/utils"
	"EngageCalculator/views"
	"encoding/xml"
	"fmt"
	"io"
	"strings"

	"fyne.io/fyne/v2/widget"
)

type PersonController struct {
	Persons map[string]models.Person
	Jobs    map[string]models.Job
}

func NewPersonController() *PersonController {
	return &PersonController{
		Persons: make(map[string]models.Person),
		Jobs:    make(map[string]models.Job),
	}
}

// LoadPersons 从文件路径加载人物数据
func (pc *PersonController) LoadPersons(filePath string, window interface{}) {
	decoder, file, err := utils.ReadXMLFile(filePath)
	if err != nil {
		fmt.Println("Error opening person file:", err)
		return
	}
	defer utils.CloseXMLFile(file)

	pc.parsePersons(decoder)
}

// LoadPersonsFromEmbedded 从内嵌数据加载人物数据
func (pc *PersonController) LoadPersonsFromEmbedded() {
	fmt.Println("Loading persons from embedded data...")
	fmt.Printf("Job XML data length: %d bytes\n", len(models.JobXMLData))
	if len(models.PersonXMLData) == 0 {
		fmt.Println("ERROR: Person XML data is empty!")
		return
	}
	decoder := xml.NewDecoder(strings.NewReader(string(models.PersonXMLData)))
	fmt.Println("Finished loading jobs from embedded data.")
	pc.parsePersons(decoder)
}

func (pc *PersonController) parsePersons(decoder *xml.Decoder) {
	fmt.Println("Starting to parse persons...")
	count := 0
	for {
		token, err := decoder.Token()
		if err != nil {
			if err == io.EOF {
				fmt.Println("Reached end of person XML file")
				break
			}
			fmt.Println("Error reading person XML:", err)
			return
		}

		count++

		if se, ok := token.(xml.StartElement); ok {
			if se.Name.Local == "Param" {
				// 读取开始标签和结束标签之间的所有内容
				var jobElement strings.Builder
				jobElement.WriteString("<Param")
				// 添加属性
				for _, attr := range se.Attr {
					jobElement.WriteString(fmt.Sprintf(" %s=\"%s\"", attr.Name.Local, attr.Value))
				}
				jobElement.WriteString(">")

				// 读取元素内容
				for {
					token, err := decoder.Token()
					if err != nil {
						break
					}
					switch t := token.(type) {
					case xml.CharData:
						jobElement.Write(t)
					case xml.StartElement:
						jobElement.WriteString("<" + t.Name.Local)
						for _, attr := range t.Attr {
							jobElement.WriteString(fmt.Sprintf(" %s=\"%s\"", attr.Name.Local, attr.Value))
						}
						jobElement.WriteString(">")
					case xml.EndElement:
						jobElement.WriteString("</" + t.Name.Local + ">")
					}
				}

				jobStr := jobElement.String()
				lines := strings.FieldsFunc(jobStr, func(c rune) bool {
					return c == '\n' || c == '\r'
				})

				// 打印每一行以检查内容
				for _, line := range lines {
					if strings.Contains(line, "<Param Out=\"\" Pid=") {
						person := models.ProcessPersonLine(line)
						//person的Pid在PersonCanRead中时才添加
						if models.PersonCanRead[person.Pid] != "" {
							pc.Persons[person.Pid] = person
						}
					}
				}
			}
		}
	}
	fmt.Printf("Finished parsing jobs. Total tokens processed: %d, Total persons found: %d\n", count, len(pc.Persons))
}

// SetJobs 设置职业数据
func (pc *PersonController) SetJobs(jobs map[string]models.Job) {
	pc.Jobs = jobs
}

// GetPersonByID 根据ID获取人物
func (pc *PersonController) GetPersonByID(id string) (models.Person, bool) {
	person, exists := pc.Persons[id]
	return person, exists
}

// GetPersonList 获取人物列表
func (pc *PersonController) GetPersonList() []models.Person {
	var persons []models.Person
	for _, person := range pc.Persons {
		persons = append(persons, person)
	}
	return persons
}

// CreatePersonSelection 创建人物选择组件
func (pc *PersonController) CreatePersonSelection(onSelected func(string)) *widget.Select {
	var personNames []string
	personNameToID := make(map[string]string)

	for _, person := range pc.Persons {
		personNames = append(personNames, person.Pid)
		personNameToID[person.Pid] = person.Pid
	}

	personSelection := widget.NewSelect(personNames, func(selected string) {
		if personID, exists := personNameToID[selected]; exists {
			onSelected(personID)
		}
	})

	return personSelection
}

func (pc *PersonController) CreatePersonView() *views.PersonView {
	// 创建人物名称列表
	var personNames []string
	for _, person := range pc.Persons {
		personNames = append(personNames, person.Pid)
	}

	// 创建并返回PersonView
	personView := views.NewPersonView(pc.Persons, personNames, pc.Jobs)

	return personView
}

func (pc *PersonController) GetPersonAndJob(selectedPid string) (models.Person, models.Job, bool) {
	person, exists := pc.Persons[selectedPid]
	if !exists {
		return models.Person{}, models.Job{}, false
	}

	job, jobExists := pc.Jobs[person.Jid]
	if !jobExists {
		job = models.Job{}
	}

	return person, job, true
}
