package main

import (
	"fmt"
	"time"
)

// (设备控制)直接方法调用示例

// MethodName 方法名称
// SuccessStatusCode 状态码
var (
	//DirectMethodFuncName = Config.DirectMethodConf.MethodName
	SuccessStatusCode = Config.DirectMethodConf.SuccessStatusCode
)

// DirectMethodFunc 统一回调返回
func DirectMethodFunc(payload map[string]interface{}) (int, map[string]interface{}, error) {
	Debug("设备控制成功:----", payload)
	payloadJson := readJsonFile(Config.DirectMethodConf.MethodJsonPath)
	for k, _ := range payloadJson.Payload {
		if v, ok := payload[k]; ok {
			payloadJson.Payload[k] = v
		}
	}
	//reportedData := map[string]interface{}{
	//	"control": payload,
	//}
	uploadDate := time.Now().AddDate(0, 0, 1).Format("2006-01-02:15:04:05")
	msgArrayData, err := setDataUploadStu(macConfData.Mac, uploadDate, payload)
	if err != nil {
		Error(err)
	}
	fmt.Println("UploadData----Before", string(msgArrayData))
	writeChan <- msgArrayData
	return SuccessStatusCode, map[string]interface{}{"result": payload}, nil
}

/*
type MethodResult struct {
	Status  int                    `json:"status,omitempty"`
	Payload map[string]interface{} `json:"payload,omitempty"`
}

// DirectMethodFunc 统一回调返回
func DirectMethodFunc(v map[string]interface{}) (MethodResult, error) {
	mResult := MethodResult{
		Status: 200,
		Payload: map[string]interface{}{
			"result": v["a"].(float64) + v["b"].(float64),
		},
	}
	return mResult, nil
}
*/
