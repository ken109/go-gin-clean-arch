package driver

import (
	"golang.org/x/oauth2/jwt"
)

var conf *jwt.Config

// 現状GCPを利用する実装がないのでコメントアウト
// func init() {
// 	credBytes, err := ioutil.ReadFile(config.Env.Gcp.CredentialPath)
// 	if err != nil {
// 		panic(err)
// 	}
//
// 	conf, err = google.JWTConfigFromJSON(credBytes)
// 	if err != nil {
// 		panic(err)
// 	}
// }

func GcpEmail() string {
	return conf.Email
}

func GcpPrivateKey() []byte {
	return conf.PrivateKey
}
