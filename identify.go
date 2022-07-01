package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"time"
)

func geneDeviceSign(macConfData MacConfDataRespStu) string {
	sharekey := macConfData.Primarykey
	policyName := Config.DeviceProvisioning.PolicyName
	uri := GetRegDpsUri(macConfData)
	expires := time.Now().Unix() + 3600

	signstr := fmt.Sprintf("%s\n%d", url.QueryEscape(uri), expires)
	signkey := Base64Decode(sharekey)

	h := hmac.New(sha256.New, signkey)
	h.Write([]byte(signstr))
	token := base64.StdEncoding.EncodeToString(h.Sum(nil))

	params := url.Values{}
	params.Add("sr", uri)
	params.Add("sig", token)
	params.Add("se", strconv.FormatInt(expires, 10))
	params.Add("skn", policyName)
	Debug("Authorization:----", fmt.Sprintf("SharedAccessSignature %s", params.Encode()))
	return fmt.Sprintf("SharedAccessSignature %s", params.Encode())
}

// GetIotHubConnectString HostName=fornanjing.azure-devices.cn;DeviceId=dca6321ac5fa;SharedAccessKey=6ZOVJ15hjrKFtZSNv1S6N9Vr7nx1Is0HwmTaKQO+GwM2F4oaL5sfgEWthkMvLFypPSAdaMU4ZkkLajUDVOxncA==
func GetIotHubConnectString(macConfData MacConfDataRespStu, assignedHub string) string {
	return fmt.Sprintf("HostName=%s;DeviceId=%s;SharedAccessKey=%s", assignedHub, macConfData.Mac, macConfData.Primarykey)
}

//accessKey
func accessKey(auth *Authentication, secondary bool) (string, error) {

	if secondary {
		return auth.SymmetricKey.SecondaryKey, nil
	}
	return auth.SymmetricKey.PrimaryKey, nil
}

// DeviceConnectionString builds up a connection string for the given device.
func DeviceConnectionString(device *Device, secondary bool, assignedHub string) (string, error) {
	key, err := accessKey(device.Authentication, secondary)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("HostName=%s;DeviceId=%s;SharedAccessKey=%s",
		assignedHub, device.DeviceID, key,
	), nil
}

//DeviceSAS 生成sas-token
func DeviceSAS(device *Device, assignedHub, resource string, duration time.Duration, secondary bool) (string, error) {
	key, err := accessKey(device.Authentication, secondary)
	if err != nil {
		return "", err
	}
	sas, _ := NewSharedAccessSignature(
		assignedHub+"/"+strings.TrimLeft(resource, "/"),
		"iothubowner",
		key,
		time.Now().Add(duration),
	)
	if err != nil {
		return "", err
	}
	return sas.String(), nil
}
