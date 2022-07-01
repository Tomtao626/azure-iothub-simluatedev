package main

import (
	"runtime"
	"time"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	// 加载配置
	LoadCfg()
	var err error
	var registerStatus bool
	// 设备注册
	macConfData, registerStatus, err = DeviceRegToIot()
	if err != nil {
		Error(err)
	}
	Debug(macConfData, registerStatus)
	var assignedHub string
	time.Sleep(time.Second * 1)
	assignedHub, err = GetRegInfoFromIot()
	if err != nil {
		Error(err)
	}
	if registerStatus == true {
		//数据上报
		go TimeDeviceUploadDataToIot(macConfData)
		//设备控制
		DeviceDirectMethod(macConfData, assignedHub, writeChan)
	}
}
