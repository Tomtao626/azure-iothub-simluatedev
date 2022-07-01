package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/tomtao626/iothub/iotdevice"
	iotmqtt "github.com/tomtao626/iothub/iotdevice/transport/mqtt"
	"io/ioutil"
	"os"
)

func Hmacsha256(data string, key []byte) (token string) {
	h := hmac.New(sha256.New, key)
	h.Write([]byte(data))
	token = hex.EncodeToString(h.Sum(nil))
	return
}

func Base64Encode(src []byte) string {
	return base64.StdEncoding.EncodeToString(src)
}

func Base64Decode(src string) []byte {
	dst, err := base64.StdEncoding.DecodeString(src)
	if err != nil {
		return nil
	}
	return dst
}

func calcSig(data string, key []byte) string {
	h := hmac.New(sha256.New, key)
	h.Write([]byte(data))
	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}

// GetRegDpsReqUrl 获取注册请求url
func GetRegDpsReqUrl(macConfData MacConfDataRespStu) (string, string) {
	baseUrl := fmt.Sprintf("https://%s", macConfData.Devendpoint)
	idScope := macConfData.Idscope
	registrationId := macConfData.Mac
	regBaseUrl := Config.DeviceProvisioning.RegBaseUrl
	return registrationId, fmt.Sprintf(regBaseUrl, baseUrl, idScope, registrationId)
}

// POST https://global.azure-devices-provisioning.net/{idScope}/registrations/{registrationId}?api-version=2021-06-01
func GetRegStatusDpsReqUrl(macConfData MacConfDataRespStu) (string, string) {
	baseUrl := fmt.Sprintf("https://%s", macConfData.Devendpoint)
	registrationId := macConfData.Mac
	idScope := macConfData.Idscope
	regBaseUrl := Config.DeviceProvisioning.RegStatusBaseUrl
	return registrationId, fmt.Sprintf(regBaseUrl, baseUrl, idScope, registrationId)
}

// GetRegDpsUri 获取注册dps校验头uri参数
func GetRegDpsUri(macConfData MacConfDataRespStu) string {
	idscope := macConfData.Idscope
	registrationid := macConfData.Mac
	return fmt.Sprintf("%s/registrations/%s", idscope, registrationid)
}

//GetConnectedDevice 获取已连接的device
func GetConnectedDevice(macConfData MacConfDataRespStu, assignedHub string) (*iotdevice.Client, error) {
	//获取当前设备的连接字符串
	connectString := GetIotHubConnectString(macConfData, assignedHub)
	Debug("Device ConnectString----", connectString)
	//解析connectString
	return iotdevice.NewFromConnectionString(iotmqtt.New(), connectString), nil
}

func setDataUploadStu(macAddr, uploadDate string, reportedData map[string]interface{}) ([]byte, error) {
	var dataUpload DataUploadStu
	dataUpload.IdentityId = Config.DataUploadConf.IdentityId
	dataUpload.Mac = macAddr
	dataUpload.Pid = Config.DataUploadConf.Pid
	dataUpload.MsgType = Config.DataUploadConf.MsgType
	dataUpload.FirmwareVersion = Config.DataUploadConf.FirmwareVersion
	dataUpload.OccuTime = uploadDate
	dataUpload.Reported = reportedData
	var msgArray []DataUploadStu
	msgArray = append(msgArray, dataUpload)
	msgArrayData, err := json.Marshal(msgArray)
	if err != nil {
		Error(err)
	}
	Debug("Device MacAdress:----", macAddr)
	Debug("DataUpload Date:----", uploadDate)
	Debug("dataUpload Struct:----", msgArray)
	Debug("dataUpload Json:----", string(msgArrayData))
	return msgArrayData, nil
}

func readJsonFile(filePath string) (payloadStu MethodPayload) {
	Debug("设备直接调用配置文件路径----MethodFile Path:----", filePath)
	filePtr, err := os.Open(filePath)
	if err != nil {
		Errorf("Open file faile, [Err:%s]", err.Error())
		return
	}
	defer filePtr.Close()

	jsonData, err := ioutil.ReadAll(filePtr)
	if err != nil {
		Error(err)
		return
	}
	var methodData MethodPayload
	err = json.Unmarshal(jsonData, &methodData)
	body, _ := json.Marshal(methodData)
	if err != nil {
		Debug("MethodPayload JsonFile Decode Fail", err.Error())
	} else {
		Debug("MethodPayload JsonFile Decode Success", string(body))
		return methodData
	}
	return
}

func writeJsonFile(payload MethodPayload) {
	filePath := Config.DirectMethodConf.MethodJsonPath
	filePtr, err := os.Create(filePath)
	if err != nil {
		fmt.Println("JsonFile Create Fail", err.Error())
		return
	}
	defer filePtr.Close()
	// 创建Json编码器
	encoder := json.NewEncoder(filePtr)
	err = encoder.Encode(payload)
	if err != nil {
		fmt.Println("Encode Fail", err.Error())
	} else {
		fmt.Println("Encode Success")
	}
}
