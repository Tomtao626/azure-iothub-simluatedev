package main

type Device struct {
	DeviceID       string          `json:"deviceId,omitempty"`
	Authentication *Authentication `json:"authentication,omitempty"`
}

type DevProvisionStu struct {
	RegistrationId string `json:"registrationId"`
}

type SymmetricKey struct {
	PrimaryKey   string `json:"primaryKey,omitempty"`
	SecondaryKey string `json:"secondaryKey,omitempty"`
}

type Authentication struct {
	SymmetricKey *SymmetricKey `json:"symmetricKey,omitempty"`
}

type DataUploadStu struct {
	IdentityId      string                 `json:"identityid"`
	Mac             string                 `json:"mac"`
	Pid             string                 `json:"pid"`
	MsgType         string                 `json:"msgtype"`
	FirmwareVersion string                 `json:"firmwareversion"`
	OccuTime        string                 `json:"occutime"`
	Reported        map[string]interface{} `json:"reported"`
}

type AssignMacStu struct {
	FactoryId   string `json:"FactoryId"`
	Pid         string `json:"Pid"`
	GenerateId  string `json:"GenerateId"`
	MacNum      int    `json:"MacNum"`
	ExpireAt    int    `json:"ExpireAt"`
	Description string `json:"Description"`
	Iottype     string `json:"iottype"`
}

type IdentityStu struct {
	Type      string `json:"type"`
	AccountId string `json:"accountId"`
}

type MacConfDataStu struct {
	Devs []string `json:"devs"`
}

type MacConfStu struct {
	Message       string         `json:"message"`
	MessageId     string         `json:"messageId"`
	Identity      IdentityStu    `json:"identity"`
	Authorization string         `json:"authorization"`
	Data          MacConfDataStu `json:"data"`
}

type MacConfDataRespStu struct {
	Mac          string `json:"mac"`
	Orderid      string `json:"orderid"`
	Primarykey   string `json:"primarykey"`
	Secondarykey string `json:"secondarykey"`
	Idscope      string `json:"idscope"`
	Devendpoint  string `json:"devendpoint"`
}

type MacConfRespStu struct {
	Message   string               `json:"message"`
	MessageId string               `json:"messageId"`
	Status    int                  `json:"status"`
	Data      []MacConfDataRespStu `json:"data"`
}

type RegistRespStu struct {
	OperationId string `json:"operationId"`
	Status      string `json:"status"`
}

type MethodPayload struct {
	Payload map[string]interface{} `json:"payload"`
}

type DeviceRegStatusRespStu struct {
	RegistrationId         string `json:"registrationId"`
	CreatedDateTimeUtc     string `json:"createdDateTimeUtc"`
	AssignedHub            string `json:"assignedHub"`
	DeviceId               string `json:"deviceId"`
	Status                 string `json:"status"`
	SubStatus              string `json:"substatus"`
	LastUpdatedDateTimeUtc string `json:"lastUpdatedDateTimeUtc"`
	Etag                   string `json:"etag"`
}

var writeChan = make(chan []byte, 1024)
var macConfData MacConfDataRespStu
