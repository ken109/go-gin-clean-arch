package driver

import (
	"cloud.google.com/go/storage"
)

var client *storage.Client

// 現状GCSを利用する実装がないのでコメントアウト
// func init() {
// 	var err error
//
// 	ctx := context.Background()
// 	client, err = storage.NewClient(ctx, option.WithCredentialsFile(config.Env.Gcp.CredentialPath))
// 	if err != nil {
// 		panic(err)
// 	}
// }

func GcsClient() *storage.Client {
	return client
}
