package main

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

func GetDeviceFormMac() (MacconfData MacConfDataRespStu) {
	Debug("GetDeviceFormMac()----通过mac地址查询预配置信息")
	Debug("Using Config.toml Conf")
	MacconfData.Mac = Config.CommonKeys.RegistrationId
	MacconfData.Idscope = Config.DeviceProvisioning.IdScope
	MacconfData.Primarykey = Config.CommonKeys.PrimaryKey
	MacconfData.Secondarykey = Config.CommonKeys.SecondaryKey
	MacconfData.Devendpoint = Config.DeviceProvisioning.DevEndPoint
	MacconfData.Orderid = Config.DeviceProvisioning.OrderId
	return
}

// DeviceRegToIot 注册设备
func DeviceRegToIot() (macConfData MacConfDataRespStu, devRegistStatus bool, err error) {
	Debug("DeviceRegToIot()----设备注册到IotHub")
	devRegistStatus = false
	macConfData = GetDeviceFormMac()
	registrationId, reqUrl := GetRegDpsReqUrl(macConfData)
	Debug("设备信息, Reg请求Url", registrationId, reqUrl)
	var reqBody DevProvisionStu
	reqBody.RegistrationId = registrationId
	reqBytes, _ := json.Marshal(reqBody)
	header := make(map[string]string)
	header["Content-Type"] = "application/json"
	header["Authorization"] = geneDeviceSign(macConfData)
	resp, err := HttpRequest("PUT", reqUrl, reqBytes, header)
	if err != nil {
		Error(err.Error())
		return
	}
	defer resp.Body.Close()
	respStr, err := ParseResponseString(resp)
	var registerResp RegistRespStu
	err = json.NewDecoder(strings.NewReader(respStr)).Decode(&registerResp)
	if err != nil {
		Error("RegistRespStu Json To Struct Err！err>>> ", err)
	}
	Debug("Response", resp)
	Debug("ResponseCode", resp.StatusCode)
	Debug("ResponseStr", respStr)
	if registerResp.Status == "assigning" && registerResp.OperationId != "" {
		Debug("Device Regist operationId----", registerResp.OperationId)
		Debug("Device Regist Status----", registerResp.Status)
		devRegistStatus = true
		return
	}
	return
}

// GetRegInfoFromIot 从IotHub获取设备注册信息
func GetRegInfoFromIot() (assignedHub string, err error) {
	Debug("GetRegInfoFromIot()----从IotHub获取设备注册信息")
	var macConfData MacConfDataRespStu
	macConfData = GetDeviceFormMac()
	Debug("设备配置信息:----", macConfData)
	registrationId, reqUrl := GetRegStatusDpsReqUrl(macConfData)
	Debug("设备信息, Reg请求Url", registrationId, reqUrl)
	var reqBody DevProvisionStu
	reqBody.RegistrationId = registrationId
	reqBytes, _ := json.Marshal(reqBody)
	header := make(map[string]string)
	header["Content-Type"] = "application/json"
	header["Authorization"] = geneDeviceSign(macConfData)
	resp, err := HttpRequest("POST", reqUrl, reqBytes, header)
	if err != nil {
		Error(err.Error())
		return
	}
	defer resp.Body.Close()
	respStr, err := ParseResponseString(resp)
	var deviceRegStatusRespData DeviceRegStatusRespStu
	err = json.NewDecoder(strings.NewReader(respStr)).Decode(&deviceRegStatusRespData)
	if err != nil {
		Error("RegistRespStu Json To Struct Err！err>>> ", err)
	}
	Debug("Response", resp)
	Debug("ResponseCode", resp.StatusCode)
	Debug("ResponseStr", respStr)
	if (deviceRegStatusRespData.Status == "assigned" || deviceRegStatusRespData.Status == "assigning") && deviceRegStatusRespData.RegistrationId == macConfData.Mac && deviceRegStatusRespData.AssignedHub != "" {
		Debug("Device RegistInfo assignedHub----", deviceRegStatusRespData.AssignedHub)
		Debug("Device RegistInfo registrationId----", deviceRegStatusRespData.RegistrationId)
		Debug("Device RegistInfo Status----", deviceRegStatusRespData.Status)
		Debug("Device RegistInfo deviceId----", deviceRegStatusRespData.DeviceId)
		Debug("Device RegistInfo createdDateTimeUtc----", deviceRegStatusRespData.CreatedDateTimeUtc)
		assignedHub = deviceRegStatusRespData.AssignedHub
		return
	}
	return
}

// DeviceUploadDataToIot 设备上报数据到iothub
func DeviceUploadDataToIot(macConfData MacConfDataRespStu, writeChan chan []byte) {
	Debug("DeviceUploadDataToIot()---设备数据上报至IotHub")
	payloadJson := readJsonFile(Config.DirectMethodConf.MethodJsonPath)
	uploadDate := time.Now().AddDate(0, 0, 1).Format("2006-01-02:15:04:05")
	msgArrayData, err := setDataUploadStu(macConfData.Mac, uploadDate, payloadJson.Payload)
	if err != nil {
		Error(err)
	}
	writeChan <- msgArrayData
	Debug("DeviceUploadDataToIot----", string(msgArrayData))
}

//TimeDeviceUploadDataToIot 定时上报
func TimeDeviceUploadDataToIot(macConfData MacConfDataRespStu) {
	Debugf("start %s dataUpload", "Device2Iot数据上报")
	ticker := time.NewTicker(time.Second * Config.DataUploadConf.TimeVal)
	for {
		select {
		case <-ticker.C:
			DeviceUploadDataToIot(macConfData, writeChan)
		}
	}
}

//DeviceDirectMethod  设备直接方法调用(控制设备) receive direct method invoke message
func DeviceDirectMethod(macConfData MacConfDataRespStu, assignedHub string, writeChan chan []byte) {
	Debug("InvokeDirectMethod()---设备方法调用(控制设备)")
	//Connect IotDevice Client
	dc, err := GetConnectedDevice(macConfData, assignedHub)
	Debug("Device Connect conf:", dc)
	if err != nil {
		Error(err)
	}
	if err = dc.Connect(context.Background()); err != nil {
		Error(err)
	}

	//Device RegisterMethod
	err = dc.RegisterMethod(
		context.Background(),
		Config.DirectMethodConf.MethodName,
		DirectMethodFunc,
	)
	if err != nil {
		Error(err)
	}
	for {
		msgArrayData := <-writeChan
		go func() {
			fmt.Println("UploadData----End", string(msgArrayData))
			if err := dc.SendEvent(context.Background(), msgArrayData); err != nil {
				Error(err)
			}
		}()
	}
	select {}
}
