package oicq

import (
	"crypto/md5"
	"fmt"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/str"
	mobile "github.com/SinKingCloud/sinking-oicq/binary-mobile/client"
	pc "github.com/SinKingCloud/sinking-oicq/binary-pc/client"
	"math/rand"
	"strconv"
	"strings"
)

type Oicq struct {
	Mobile     *mobile.QQClient
	LoginToken []byte

	Pc          *pc.QQClient
	PcLoginTime int64
	PcHeartTime int64
}

type Proxy struct {
	Addr string `json:"addr"`
	Port string `json:"port"`
	User string `json:"user"`
	Pwd  string `json:"pwd"`
}

// NewOicq 实例化OICQ
func NewOicq(uin int64, password [16]byte, proxy *Proxy, device int, loginToken []byte, proto mobile.ClientProtocol) *Oicq {
	oicq := &Oicq{
		Mobile: mobile.NewClientMd5(uin, password).UseDevice(getMobileDevice(uin, device, proto)),
		Pc:     pc.NewClientMd5(uin, password[:]).UseDevice(getPcDevice(uin)),
	}
	if loginToken != nil && len(loginToken) > 0 {
		oicq.LoginToken = loginToken
	}
	if proxy != nil && proxy.Addr != "" && proxy.Port != "" {
		oicq.Mobile.UseProxy(proxy.Addr, proxy.Port, proxy.User, proxy.Pwd)
		oicq.Pc.UseProxy(proxy.Addr, proxy.Port, proxy.User, proxy.Pwd)
	}
	return oicq
}

// getMobileDevice 获取手机设备信息
func getMobileDevice(uin int64, deviceId int, proto mobile.ClientProtocol) *mobile.DeviceInfo {
	device := mobile.SystemDeviceInfo
	s := str.NewStringTool("SK-" + strconv.Itoa(deviceId) + strconv.FormatInt(uin, 10)).Md5()
	md5Lower := strings.ToLower(s)
	s2 := str.NewStringTool("SK-" + strconv.Itoa(deviceId) + md5Lower).Md5()
	md52Lower := strings.ToLower(s2)
	md5Upper := strings.ToUpper(s)
	md52Upper := strings.ToUpper(s2)
	md5Twice := md5.Sum([]byte(s))
	device.Display = []byte("388" + md5Lower[0:13])
	device.Bootloader = []byte("G9550" + md5Upper[0:8])
	device.FingerPrint = []byte("samsung/dream2qltezh/dream2qltechn:7.1/" + md5Upper[6:12] + "/" + string(device.Bootloader) + ":user/release-keys")
	device.ProcVersion = []byte("Linux version 4.4.88-" + md5Lower[13:26] + " (android-build@wpiu5.hot.corp.google.com) (Android clang version 5.0.300080 (based on LLVM 5.0.300080)) #1 SMP PREEMPT Thu May 17 20:37:50 UTC 2018")
	device.BootId = []byte(fmt.Sprintf("%s-%s-%s-%s-%s", md52Lower[0:8], md5Lower[8:12], md52Lower[12:16], md5Lower[16:20], md52Lower[20:32]))
	device.MacAddress = []byte(fmt.Sprintf("%s:%s:%s:%s:%s:%s", md5Upper[1:3], md52Upper[3:5], md5Upper[5:7], md52Upper[7:9], md5Upper[9:11], md52Upper[11:13]))
	device.IpAddress = []byte{172, 16, byte(rand.Intn(255)), byte(rand.Intn(255))}
	device.WifiBSSID = md5Twice[:]
	uuid := md52Lower[20:28]
	device.WifiSSID = []byte(md52Lower[29:32] + uuid + md5Lower[8:16])
	device.IMEI = strconv.FormatInt(395650000000000+uin*6+int64(deviceId), 10)
	device.AndroidId = []byte(uuid + md5Lower[8:16])
	device.QImei = []byte(md52Lower[8:16] + md52Lower[17:25])
	device.QImei36 = []byte("5aac" + md52Lower)
	device.Protocol = proto
	return device
}

func getPcDevice(uin int64) *pc.DeviceInfo {
	device := pc.SystemDeviceInfo
	device.Name = []byte("PC-SK" + strconv.FormatInt(uin, 10))
	return device
}

// AutoLogin 自动登录
func (o *Oicq) AutoLogin() (error, *mobile.LoginResponse) {
	if !o.Mobile.Online.Load() {
		if o.LoginToken == nil || len(o.LoginToken) <= 0 || o.TokenLogin() != nil {
			loginStatus, err := o.Login()
			if err != nil {
				return err, nil
			}
			if !loginStatus.Success {
				return nil, loginStatus
			}
		}
	}
	return nil, &mobile.LoginResponse{
		Success: true,
	}
}

// Login 登陆QQ
func (o *Oicq) Login() (*mobile.LoginResponse, error) {
	resp, e := o.Mobile.Login()
	if e == nil && resp.Success {
		e = o.Mobile.Init(false, false)
		if e == nil {
			o.LoginToken = o.Mobile.GenToken()
		} else {
			resp.Success = false
		}
	}
	return resp, e
}

// TokenLogin token登陆QQ
func (o *Oicq) TokenLogin() error {
	e := o.Mobile.TokenLogin(o.LoginToken)
	if e == nil {
		e = o.Mobile.Init(true, false)
	}
	if e == nil {
		o.LoginToken = o.Mobile.GenToken()
	}
	return e
}
