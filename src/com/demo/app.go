package demo

import (
	"fmt"
	"com/models"
)

func AppDemo(disInstance DISInstance) {
	fmt.Println("---------------------------AppDemo start----------------------------");
	CreateApp(disInstance)
	DescApp(disInstance)
	ListApps(disInstance)
	DeleteApp(disInstance)
	fmt.Println("----------------------------AppDemo end-----------------------------");
}

//创建APP
func CreateApp(disInstance DISInstance) {
	dis := disInstance.Dis
	result := dis.CreateApp(disInstance.AppName)
	if result.Err != nil {
		fmt.Printf("Failed to create app [%s].\n", disInstance.AppName);
	} else {
		fmt.Printf("Success to create app [%s].\n", disInstance.AppName);
	}
}

//查询APP详情
func DescApp(disInstance DISInstance) {
	dis := disInstance.Dis
	result, output := dis.DescApp(disInstance.AppName)
	if result.Err != nil {
		fmt.Printf("Failed to desc app [%s].\n", disInstance.AppName);
	} else {
		fmt.Printf("App name [%s], Id [%s], createTime [%d].\n", output.AppName, output.AppId, output.CreateTime);
	}
}

//查询APP列表
func ListApps(disInstance DISInstance) {
	input := &models.ListAppsRequest{
		Limit:10,
	}
	dis := disInstance.Dis
	result, output := dis.ListApps(input)
	if result.Err != nil {
		fmt.Printf("Failed to list app.\n");
	} else {
		fmt.Printf("ListApps, size [%d], hasMoreApp [%t].\n", len(output.Apps), output.HasMoreApp);
		for i, app := range output.Apps {
			fmt.Printf("App%d name [%s], Id [%s], createTime [%d].\n", i, app.AppName, app.AppId, app.CreateTime);
		}
	}
}

//删除APP
func DeleteApp(disInstance DISInstance) {
	dis := disInstance.Dis
	result := dis.DeleteApp(disInstance.AppName)
	if result.Err != nil {
		fmt.Printf("Failed to delete app [%s].\n", disInstance.AppName);
	} else {
		fmt.Printf("Success to delete app [%s].\n", disInstance.AppName);
	}
}

