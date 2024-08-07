package config

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"os"
)

var IsInit bool

type Config struct {
	LogLevel             string            `json:"logLevel"` // 日志级别
	Domain               string            `json:"domain"`
	Domains              []string          `json:"domains"` //多域名设置，把所有收信域名都填进去
	WebDomain            string            `json:"webDomain"`
	DkimPrivateKeyPath   string            `json:"dkimPrivateKeyPath"`
	SSLType              string            `json:"sslType"` // 0表示自动生成证书，HTTP挑战模式，1表示用户上传证书，2表示自动-DNS挑战模式
	SSLPrivateKeyPath    string            `json:"SSLPrivateKeyPath"`
	SSLPublicKeyPath     string            `json:"SSLPublicKeyPath"`
	DbDSN                string            `json:"dbDSN"`
	DbType               string            `json:"dbType"`
	HttpsEnabled         int               `json:"httpsEnabled"`    //后台页面是否启用https，0默认（启用），1启用，2不启用
	SpamFilterLevel      int               `json:"spamFilterLevel"` //垃圾邮件过滤级别，0不过滤、1 spf dkim 校验均失败时过滤，2 spf校验不通过时过滤
	HttpPort             int               `json:"httpPort"`        //http服务端口设置，默认80
	HttpsPort            int               `json:"httpsPort"`       //https服务端口，默认443
	WeChatPushAppId      string            `json:"weChatPushAppId"`
	WeChatPushSecret     string            `json:"weChatPushSecret"`
	WeChatPushTemplateId string            `json:"weChatPushTemplateId"`
	WeChatPushUserId     string            `json:"weChatPushUserId"`
	TgBotToken           string            `json:"tgBotToken"`
	TgChatId             string            `json:"tgChatId"`
	IsInit               bool              `json:"isInit"`
	WebPushUrl           string            `json:"webPushUrl"`
	WebPushToken         string            `json:"webPushToken"`
	Tables               map[string]string `json:"-"`
	TablesInitData       map[string]string `json:"-"`
	setupPort            int               // 初始化阶段端口
}

func (c *Config) GetSetupPort() int {
	return c.setupPort
}

func (c *Config) SetSetupPort(setupPort int) {
	c.setupPort = setupPort
}

const DBTypeMySQL = "mysql"
const DBTypeSQLite = "sqlite"
const DBTypePostgres = "postgres"
const SSLTypeAutoHTTP = "0" //自动生成证书
const SSLTypeAutoDNS = "2"  //自动生成证书，DNS api验证
const SSLTypeUser = "1"     //用户上传证书

var DBTypes []string = []string{DBTypeMySQL, DBTypeSQLite, DBTypePostgres}

var Instance *Config = &Config{}

func Init() {
	var cfgData []byte
	var err error
	args := os.Args

	if len(args) >= 2 && args[len(args)-1] == "dev" {
		cfgData, err = os.ReadFile("./config/config.dev.json")
		if err != nil {
			return
		}
	} else {
		cfgData, err = os.ReadFile("./config/config.json")
		if err != nil {
			return
		}
	}

	err = json.Unmarshal(cfgData, &Instance)
	if err != nil {
		return
	}

	if len(Instance.Domains) == 0 && Instance.Domain != "" {
		Instance.Domains = []string{Instance.Domain}
	}

	if Instance.Domain != "" && Instance.IsInit {
		IsInit = true
	}

}

func ReadPrivateKey() (*ecdsa.PrivateKey, bool) {
	key, err := os.ReadFile("./config/ssl/account_private.pem")
	if err != nil {
		return createNewPrivateKey(), true
	}

	block, _ := pem.Decode(key)
	x509Encoded := block.Bytes
	privateKey, _ := x509.ParseECPrivateKey(x509Encoded)

	return privateKey, false
}

func createNewPrivateKey() *ecdsa.PrivateKey {
	// Create a user. New accounts need an email and private key to start.
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		panic(err)
	}
	x509Encoded, _ := x509.MarshalECPrivateKey(privateKey)

	// 将ec 密钥写入到 pem文件里
	keypem, _ := os.OpenFile("./config/ssl/account_private.pem", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	err = pem.Encode(keypem, &pem.Block{Type: "EC PRIVATE KEY", Bytes: x509Encoded})
	if err != nil {
		panic(err)
	}
	return privateKey
}
