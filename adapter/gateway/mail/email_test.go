package mail

import (
	"testing"

	"go-gin-clean-arch/resource/mail_body"
)

func Test_email_Send(t *testing.T) {
	type args struct {
		to   string
		body mail_body.MailBody
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "default",
			args: args{
				to: "test@test.com",
				body: mail_body.Default{
					Title: "テスト",
					Body:  "メール送信テスト",
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := email{}
			if err := e.Send(tt.args.to, tt.args.body); (err != nil) != tt.wantErr {
				t.Errorf("Send() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
