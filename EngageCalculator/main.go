package main

import (
	"EngageCalculator/controllers"
	"EngageCalculator/utils"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
)

func main() {
	// 创建Fyne应用
	myApp := app.New()
	myWindow := myApp.NewWindow("Fire Emblem Engage-IFMOD-Calculator")
	myWindow.Resize(fyne.NewSize(1280, 840))

	// 创建控制器
	jobController := controllers.NewJobController()
	personController := controllers.NewPersonController()

	// 加载数据(从内嵌数据)
	jobController.LoadJobsFromEmbedded()
	personController.LoadPersonsFromEmbedded()

	// 设置职业数据到人物控制器
	personController.SetJobs(jobController.Jobs)

	// 创建视图
	personView := personController.CreatePersonView()

	// 创建标签页
	tabs := container.NewAppTabs(
		container.NewTabItem("人物信息", personView.GetContent()),
	)

	// 创建背景
	background := utils.CreateBackground("./resources/background.png", 0.85)

	// 创建主容器，将背景和内容叠放
	content := container.NewStack(background, tabs)

	// 设置窗口内容
	myWindow.SetContent(content)

	// 显示并运行
	myWindow.ShowAndRun()
}
