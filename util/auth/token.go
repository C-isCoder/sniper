package auth

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/bilibili/twirp"
	"sniper/util/conf"
	"strings"
	"time"
)

// NotTokenError 未传 Token
var NotTokenError = twirp.NewError(twirp.PermissionDenied, "headers must has token")
// FailTokenError Token格式错误
var FailTokenError = twirp.NewError(twirp.InvalidArgument, "token formatter error")
// ExpiredTokenError Token 过期
var ExpiredTokenError = twirp.NewError(twirp.Unauthenticated, "token exp need login")
// ChangeTokenError Token 被窜改
var ChangeTokenError = twirp.NewError(twirp.Unauthenticated, "token has change")

type JWT string

type Token struct {
	header  Header
	payload Payload
}

/*
iss：JWT Token 的签发者
sub：主题
exp：JWT Token 过期时间
aud：接收 JWT Token 的一方
iat：JWT Token 签发时间
nbf：JWT Token 生效时间
jti：JWT Token ID
*/
type Payload struct {
	Uid     string        `json:"uid"` // 用户ID
	Iss     string        `json:"iss"`
	Exp     time.Duration `json:"exp"`
	Nbf     time.Duration `json:"nbf"`
	Name    string        `json:"name"`
	IsAdmin bool          `json:"is_admin"`
}

/**
Token 的类型
Token 所使用的加密算法
{
  "typ": "JWT",
  "alg": "HS256"
}
*/
type Header struct {
	Typ string `json:"typ"`
	Alg string `json:"alg"`
}

/*
SIGNATURE
HMACSHA256(
  base64UrlEncode(header) + "." +
  base64UrlEncode(payload),
  secret
)
*/

// 过期间隔 2 hour
var expInterval = time.Duration(2 * (60 /*s*/ * 60 /*m*/))
// test 1 min
//var expInterval = time.Duration(1 * 60)

func NewToken(uid, name string) JWT {
	head := header()
	payload := payload(uid, name)
	secret := conf.GetString("jwt_secret")
	head64, payload64, secret265 := hs265(secret, head, payload)
	jwt := head64 + "." + payload64 + "." + secret265
	return JWT(jwt)
}

func hs265(secret string, head string, payload string) (string, string, string) {
	hm := hmac.New(sha256.New, []byte(secret))
	head64 := base64.URLEncoding.EncodeToString([]byte(head))
	payload64 := base64.URLEncoding.EncodeToString([]byte(payload))
	hm.Write([]byte(head64 + "." + payload64 + "."))
	secret265 := hex.EncodeToString(hm.Sum(nil))
	return head64, payload64, secret265
}

func (jwt JWT) String() string {
	return string(jwt)
}

func (jwt JWT) VerifyToken() (bool, twirp.Error) {
	header, payload, secret265 := jwt.parse()
	if now() > payload.Exp {
		return false, ExpiredTokenError
	}
	_, _, secret := hs265(conf.GetString("jwt_secret"), header.string(), payload.string())
	if secret265 != secret {
		return false, ChangeTokenError
	}
	return true, nil
}

func (jwt JWT) GetUID() string {
	_, payload, _ := jwt.parse()
	return payload.Uid
}

func header() string {
	header := Header{Typ: "JWT", Alg: "HS256"}
	bytes, err := json.Marshal(header)
	if err != nil {
		fmt.Println("JWT header() json error", err)
	}
	return string(bytes)
}

func payload(uid, name string) string {
	payload := Payload{}
	payload.Uid = uid
	payload.Name = name
	payload.Nbf = now()
	payload.Exp = payload.Nbf + expInterval

	bytes, err := json.Marshal(payload)
	if err != nil {
		fmt.Println("JWT payload() json error", err)
	}
	//log.Get(ctx).Debugln("JWT payload-> ", string(bytes))
	return string(bytes)
}

func now() time.Duration {
	return time.Duration(time.Now().Unix())
}

func (jwt JWT) parse() (Header, Payload, string) {
	sps := strings.Split(jwt.String(), ".")
	var header Header
	var payload Payload
	var secret265 string
	for i := 0; i < 3; i++ {
		by, err := base64.URLEncoding.DecodeString(sps[i])
		if i == 0 {
			err = json.Unmarshal(by, &header)
		} else if i == 1 {
			err = json.Unmarshal(by, &payload)
		} else {
			secret265 = sps[i]
		}
		if err != nil {
			fmt.Println("JWT parse() base64 decode error", err)
		}
	}
	return header, payload, secret265
}

func (header Header) string() string {
	bytes, err := json.Marshal(header)
	if err != nil {
		fmt.Println("Header str() error", err)
	}
	return string(bytes)
}

func (payload Payload) string() string {
	bytes, err := json.Marshal(payload)
	if err != nil {
		fmt.Println("Payload str() error", err)
	}
	return string(bytes)
}
