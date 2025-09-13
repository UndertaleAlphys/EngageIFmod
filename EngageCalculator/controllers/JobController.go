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

type JobController struct {
	Jobs map[string]models.Job
}

func NewJobController() *JobController {
	return &JobController{
		Jobs: make(map[string]models.Job),
	}
}

// LoadJobs 从文件路径加载职业数据
func (jc *JobController) LoadJobs(filePath string, window interface{}) {
	decoder, file, err := utils.ReadXMLFile(filePath)
	if err != nil {
		fmt.Println("Error opening job file:", err)
		return
	}
	defer utils.CloseXMLFile(file)

	jc.parseJobs(decoder)
}

// LoadJobsFromEmbedded 从内嵌数据加载职业数据
func (jc *JobController) LoadJobsFromEmbedded() {
	fmt.Println("Loading jobs from embedded data...")
	fmt.Printf("Job XML data length: %d bytes\n", len(models.JobXMLData))
	if len(models.JobXMLData) == 0 {
		fmt.Println("ERROR: Job XML data is empty!")
		return
	}
	decoder := xml.NewDecoder(strings.NewReader(string(models.JobXMLData)))
	fmt.Println("Finished loading jobs from embedded data.")
	jc.parseJobs(decoder)
}

func (jc *JobController) parseJobs(decoder *xml.Decoder) {
	fmt.Println("Starting to parse jobs...")
	count := 0
	for {
		token, err := decoder.Token()
		if err != nil {
			if err == io.EOF {
				fmt.Println("Reached end of job XML file")
				break
			}
			fmt.Println("Error reading job XML:", err)
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
					if strings.Contains(line, "<Param Out=\"\" Jid=") {
						job := models.ProcessJobLine(line)
						if models.JobCanRead[job.Jid] != "" {
							jc.Jobs[job.Jid] = job
						}
					}
				}
			}
		}
	}
	fmt.Printf("Finished parsing jobs. Total tokens processed: %d, Total jobs found: %d\n", count, len(jc.Jobs))
}

// GetJobByID 根据ID获取职业
func (jc *JobController) GetJobByID(id string) (models.Job, bool) {
	job, exists := jc.Jobs[id]
	return job, exists
}

// GetJobList 获取职业列表
func (jc *JobController) GetJobList() []models.Job {
	var jobs []models.Job
	for _, job := range jc.Jobs {
		jobs = append(jobs, job)
	}
	return jobs
}

// CreateJobSelection 创建职业选择组件
func (jc *JobController) CreateJobSelection(onSelected func(string)) *widget.Select {
	var jobNames []string
	jobNameToID := make(map[string]string)

	for _, job := range jc.Jobs {
		jobNames = append(jobNames, job.StyleName)
		jobNameToID[job.StyleName] = job.Jid
	}

	jobSelection := widget.NewSelect(jobNames, func(selected string) {
		if jobID, exists := jobNameToID[selected]; exists {
			onSelected(jobID)
		}
	})

	return jobSelection
}

func (jc *JobController) CreateJobView() *views.JobView {
	var jobNames []string
	for _, job := range jc.Jobs {
		jobNames = append(jobNames, job.StyleName)
	}

	return views.NewJobView(jc.Jobs, jobNames)
}

func (jc *JobController) GetJob(selectedJid string) (models.Job, bool) {
	job, exists := jc.Jobs[selectedJid]
	return job, exists
}
