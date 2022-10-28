package logic

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"net/http"
	"os"
	"strings"
	"time"
)

const (
	PemBegin    = "-----BEGIN RSA PRIVATE KEY-----\n"
	PemEnd      = "\n-----END RSA PRIVATE KEY-----"
	PemBeginPub = "-----BEGIN RSA PUBLIC KEY-----\n"
	PemEndPub   = "\n-----END RSA PUBLIC KEY-----"
)

// GenerateRSAKey 生成RSA私钥和公钥，保存到文件中 PKCS#1
func GenerateRSAKey(bits int) {
	//GenerateKey函数使用随机数据生成器random生成一对具有指定字位数的RSA密钥
	//Reader是一个全局、共享的密码用强随机数生成器
	privateKey, err := rsa.GenerateKey(rand.Reader, bits)
	if err != nil {
		panic(err)
	}
	//保存私钥
	//通过x509标准将得到的ras私钥序列化为ASN.1 的 DER编码字符串
	X509PrivateKey := x509.MarshalPKCS1PrivateKey(privateKey)
	//x509.MarshalPKCS8PrivateKey(privateKey)
	//使用pem格式对x509输出的内容进行编码
	//创建文件保存私钥
	privateFile, err := os.Create("private.pem")
	if err != nil {
		panic(err)
	}
	defer func(privateFile *os.File) {
		err = privateFile.Close()
		if err != nil {
			fmt.Println(err.Error())
		}
	}(privateFile)
	//构建一个pem.Block结构体对象
	privateBlock := pem.Block{Type: "RSA Private Key", Bytes: X509PrivateKey}
	//将数据保存到文件
	err = pem.Encode(privateFile, &privateBlock)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	//保存公钥
	//获取公钥的数据
	publicKey := privateKey.PublicKey
	//X509对公钥编码
	X509PublicKey, err := x509.MarshalPKIXPublicKey(&publicKey)
	if err != nil {
		panic(err)
	}
	//pem格式编码
	//创建用于保存公钥的文件
	publicFile, err := os.Create("public.pem")
	if err != nil {
		panic(err)
	}
	defer func(publicFile *os.File) {
		err = publicFile.Close()
		if err != nil {
			fmt.Println(err.Error())
		}
	}(publicFile)
	//创建一个pem.Block结构体对象
	publicBlock := pem.Block{Type: "RSA Public Key", Bytes: X509PublicKey}
	//保存到文件
	err = pem.Encode(publicFile, &publicBlock)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
}

// RsaEncryptWithSha1Base64 （1）加密：采用sha1算法加密后转base64格式
func RsaEncryptWithSha1Base64(originalData, publicKey string) (string, error) {
	key, _ := base64.StdEncoding.DecodeString(publicKey)
	pubKey, _ := x509.ParsePKIXPublicKey(key)
	encryptedData, err := rsa.EncryptPKCS1v15(rand.Reader, pubKey.(*rsa.PublicKey), []byte(originalData))
	return base64.StdEncoding.EncodeToString(encryptedData), err
}

// RsaDecryptWithSha1Base64 （2）解密：对采用sha1算法加密后转base64格式的数据进行解密（私钥PKCS1格式）
func RsaDecryptWithSha1Base64(encryptedData, privateKey string) (string, error) {
	encryptedDecodeBytes, err := base64.StdEncoding.DecodeString(encryptedData)
	if err != nil {
		return "", err
	}
	key, _ := base64.StdEncoding.DecodeString(privateKey)
	prvKey, _ := x509.ParsePKCS1PrivateKey(key)
	originalData, err := rsa.DecryptPKCS1v15(rand.Reader, prvKey, encryptedDecodeBytes)
	return string(originalData), err
}

// Sign 签名
func Sign(msg, privateKey string) (signature string, err error) {
	msgHash := sha256.New()
	_, err = msgHash.Write([]byte(msg))
	if err != nil {
		return
	}
	msgHashSum := msgHash.Sum(nil)

	var key *rsa.PrivateKey
	key, err = ParsePrivateKey(privateKey)
	if err != nil {
		return
	}
	var sign []byte
	sign, err = rsa.SignPSS(rand.Reader, key, crypto.SHA256, msgHashSum, nil)
	if err != nil {
		return
	}
	signature = base64.StdEncoding.EncodeToString(sign)

	return
}

// Verify 验签
func Verify(publicKey, data, sign string) (err error) {
	var signature []byte
	signature, err = base64.StdEncoding.DecodeString(sign)
	if err != nil {
		return err
	}
	msgHash := sha256.New()
	_, err = msgHash.Write([]byte(data))
	if err != nil {
		return
	}
	msgHashSum := msgHash.Sum(nil)
	var key *rsa.PublicKey
	key, err = ParsePublicKey(publicKey)
	if err != nil {
		return
	}
	err = rsa.VerifyPSS(key, crypto.SHA256, msgHashSum, signature, nil)
	if err != nil {
		return
	}
	return
}

func ParsePrivateKey(privateKey string) (*rsa.PrivateKey, error) {
	privateKey = FormatPrivateKey(privateKey)
	// 2、解码私钥字节，生成加密对象
	block, _ := pem.Decode([]byte(privateKey))
	if block == nil {
		return nil, errors.New("私钥信息错误！")
	}
	// 3、解析DER编码的私钥，生成私钥对象
	priKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	return priKey, nil
}

func ParsePublicKey(publicKey string) (*rsa.PublicKey, error) {
	publicKey = FormatPublicKey(publicKey)
	// 2、解码公钥字节，生成加密对象
	block, _ := pem.Decode([]byte(publicKey))
	if block == nil {
		return nil, errors.New("公钥信息错误！")
	}
	// 3、解析DER编码的公钥，生成公钥对象
	pubKey, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	return pubKey.(*rsa.PublicKey), nil
}

func FormatPrivateKey(privateKey string) string {
	if !strings.HasPrefix(privateKey, PemBegin) {
		privateKey = PemBegin + privateKey
	}
	if !strings.HasSuffix(privateKey, PemEnd) {
		privateKey = privateKey + PemEnd
	}
	return privateKey
}

func FormatPublicKey(publicKey string) string {
	if !strings.HasPrefix(publicKey, PemBeginPub) {
		publicKey = PemBeginPub + publicKey
	}
	if !strings.HasSuffix(publicKey, PemEndPub) {
		publicKey = publicKey + PemEndPub
	}
	return publicKey
}

type AuthServer struct {
	License string `json:"license" binding:"required"`
}
type Auth struct {
	Project   string `json:"project,omitempty" mapstructure:"project"`
	Num       string `json:"num,omitempty" binding:"required" mapstructure:"num"`
	Name      string `json:"name,omitempty" mapstructure:"name"`
	Days      uint   `json:"days,omitempty" mapstructure:"days"`
	ApplyTime uint   `json:"apply_time,omitempty" mapstructure:"apply_time"`
}

func CheckLicense() (err error) {
	license := viper.GetString("auth_server.license")
	if license == "" {
		return errors.New("no license")
	}

	var a Auth
	err, a = ParseLicense(license)
	if err != nil {
		return err
	}

	if a.Project != viper.GetString("auth.project") {
		return errors.New("license error3, please reauthorize")
	}

	if a.Name != viper.GetString("auth.name") {
		return errors.New("license error4, please reauthorize")
	}

	if a.Days != uint(viper.GetInt("auth.days")) {
		return errors.New("license error5, please reauthorize")
	}

	if a.ApplyTime != uint(viper.GetInt("auth.apply_time")) {
		return errors.New("license error6, please reauthorize")
	}

	return
}

func ParseLicense(license string) (err error, a Auth) {
	arr := strings.Split(license, ",")
	if len(arr) != 2 {
		return errors.New("license error1, please reauthorize"), a
	}

	err = Verify(viper.GetString("auth_server.public_key"), arr[0], arr[1])
	if err != nil {
		fmt.Println(arr[0])
		return
	}

	obj, err2 := RsaDecryptWithSha1Base64(arr[0], viper.GetString("auth.private_key"))
	if err2 != nil {
		return err2, a
	}

	err = json.Unmarshal([]byte(obj), &a)
	if err != nil {
		return
	}

	if a.Num != viper.GetString("auth.num") {
		return errors.New("license error2, please reauthorize"), a
	}

	// 检查时间是否在有效期内
	if time.Now().Unix() > int64(a.ApplyTime+a.Days*60*60*24) {
		return errors.New("license expired, please reauthorize"), a
	}

	// todo: 有网络时尝试联网检查 当前license是否被取消 状态是否正常等等

	return
}

func CheckLicenseMiddleware(c *gin.Context) {
	err := CheckLicense()
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    http.StatusUnauthorized,
			"message": http.StatusText(http.StatusUnauthorized) + ": " + err.Error(),
		})
		c.Abort()
	}
}

func SHA1(s string) string {
	o := sha1.New()
	o.Write([]byte(s))
	return hex.EncodeToString(o.Sum(nil))
}
