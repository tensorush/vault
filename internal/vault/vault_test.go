package vault

import (
	"log"
	"reflect"
	"testing"

	"go.uber.org/zap"

	"vault-bot/internal/database"
	"vault-bot/internal/secret"
)

func newVault(t *testing.T) *Vault {
	s, err := database.New("test", "file::memory:?cache=shared")
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatalf("zap error: %s", err)
	}

	v, err := New(s, "1234567890123456", logger)
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}
	return v
}

func getEncrypted(text string, err error) string {
	if err != nil {
		log.Fatalf("Encrypt() error = %v", err)
	}
	return text
}

func TestVault_Decrypt(t *testing.T) {
	v := newVault(t)

	type args struct {
		text string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "too small text",
			args: args{
				text: "test1",
			},
			want: "test1",
		},
		{
			name: "ok",
			args: args{
				text: getEncrypted(v.Encrypt("test")),
			},
			want: "test",
		},
		{
			name: "ok #2",
			args: args{
				text: getEncrypted(v.Encrypt("dqwfqwedfefqfqfhkqfjqjfgqwdqwfgqwefhqvdjvqwvf")),
			},
			want: "dqwfqwedfefqfqfhkqfjqjfgqwdqwfgqwefhqvdjvqwvf",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := v.Decrypt(tt.args.text)
			if err != nil {
				t.Errorf("Decrypt() error = %v", err)
			}
			if got != tt.want {
				t.Errorf("Decrypt() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestVault_Delete(t *testing.T) {
	v := newVault(t)

	type args struct {
		chatID  int64
		service string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "ok",
			args: args{
				chatID:  190,
				service: "test",
			},
		},
		{
			name: "ok #2",
			args: args{
				chatID:  191,
				service: "teqdwqwdqdst",
			},
		},
		{
			name: "not found",
			args: args{
				chatID:  197,
				service: "teqdwqwdqdst",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if !tt.wantErr {
				err := v.Save(tt.args.chatID, tt.args.service, "test", "test")
				if err != nil {
					t.Errorf("Save() error = %v", err)
				}
			}
			if err := v.Delete(tt.args.chatID, tt.args.service); (err != nil) != tt.wantErr {
				t.Errorf("Delete() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr {
				if _, err := v.Get(tt.args.chatID, tt.args.service); err == nil {
					t.Errorf("Delete() error = %v, wantErr %v", err, tt.wantErr)
				}
			}
		})
	}
}

func TestVault_Encrypt(t *testing.T) {
	v := newVault(t)

	type args struct {
		text string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "small text",
			args: args{
				text: "test1",
			},
		},
		{
			name: "big text",
			args: args{
				text: "dqwfqwedfefqfqfhkqfjqjfgqwdqwfgqwefhqvdjvqwvf,qwld,qlfmkqmqfqjdnqf",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := v.Encrypt(tt.args.text)
			if err != nil {
				t.Errorf("Encrypt() error = %v", err)
			}

			got, err = v.Decrypt(got)
			if err != nil {
				t.Errorf("Decrypt() error = %v", err)
			}

			if got != tt.args.text {
				t.Errorf("Decrypt() = %v, want %v", got, tt.args.text)
			}
		})
	}
}

func TestVault_Get(t *testing.T) {
	v := newVault(t)

	type args struct {
		chatID  int64
		service string
	}
	tests := []struct {
		name    string
		args    args
		want    secret.Credentials
		wantErr bool
	}{
		{
			name: "ok",
			args: args{
				chatID:  290,
				service: "test",
			},
			want: secret.Credentials{
				Login:    "test login",
				Password: "test password",
			},
		},
		{
			name: "ok #2",
			args: args{
				chatID:  291,
				service: "teqdwqwdqdst",
			},
			want: secret.Credentials{
				Login:    "teqdwqwdqdst",
				Password: "XXXXXXXXXXXXXXXXXXXXX",
			},
		},
		{
			name: "not found",
			args: args{
				chatID:  291,
				service: "not found",
			},
			want:    secret.Credentials{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if !tt.wantErr {
				if err := v.Save(tt.args.chatID, tt.args.service, tt.want.Login, tt.want.Password); err != nil {
					t.Errorf("Set() error = %v, wantErr %v", err, tt.wantErr)
				}
			}

			got, err := v.Get(tt.args.chatID, tt.args.service)
			if (err != nil) != tt.wantErr {
				t.Errorf("Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Get() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestVault_GetLang(t *testing.T) {
	v := newVault(t)

	type args struct {
		chatID int64
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "ok",
			args: args{
				chatID: 390,
			},
			want: "pt",
		},
		{
			name: "ok  #2",
			args: args{
				chatID: 391,
			},
			want: "en",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v.SetLang(tt.args.chatID, tt.want)
			if got := v.GetLang(tt.args.chatID); got != tt.want {
				t.Errorf("GetLang() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestVault_Hash(t *testing.T) {
	v := newVault(t)

	type args struct {
		text string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "ok",
			args: args{
				text: "test",
			},
			want:    "n4bQgYhMfWWaL+qgxVrQFaO/TxsrC4Is0V1sFbDwCgg",
			wantErr: false,
		},
		{
			name: "ok #2",
			args: args{
				text: "fnqjlnfjkqndfjkqnfjqndfqjnj",
			},
			want:    "iSRSDMOs2PAGQ5hS4CBvaU53UrRpfF8TXSOoONJpv0w",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := v.Hash(tt.args.text)
			if (err != nil) != tt.wantErr {
				t.Errorf("Hash() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Hash() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestVault_Save(t *testing.T) {
	v := newVault(t)
	type args struct {
		chatID   int64
		service  string
		login    string
		password string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "ok",
			args: args{
				chatID:   490,
				service:  "test",
				login:    "test login",
				password: "XXXXXXXXXXXXX",
			},
		},
		{
			name: "ok #2",
			args: args{
				chatID:   491,
				service:  "test2",
				login:    "test login2",
				password: "XXXXXXXXXXXXX",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := v.Save(tt.args.chatID, tt.args.service, tt.args.login, tt.args.password); (err != nil) != tt.wantErr {
				t.Errorf("Save() error = %v, wantErr %v", err, tt.wantErr)
			}

			if !tt.wantErr {
				if p, err := v.Get(tt.args.chatID, tt.args.service); err != nil {
					t.Errorf("Get() error = %v, wantErr %v", err, tt.wantErr)
				} else if p.Login != tt.args.login || p.Password != tt.args.password {
					t.Errorf("Get() got = %v, want %v", p, tt.args)
				}
			}
		})
	}
}

func TestVault_SetLang(t *testing.T) {
	v := newVault(t)
	type args struct {
		chatID int64
		lang   string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "ok",
			args: args{
				chatID: 590,
				lang:   "pt",
			},
		},
		{
			name: "ok #2",
			args: args{
				chatID: 590,
				lang:   "en",
			},
		},
		{
			name: "ok #3",
			args: args{
				chatID: 591,
				lang:   "pt",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v.SetLang(tt.args.chatID, tt.args.lang)
			if got := v.GetLang(tt.args.chatID); got != tt.args.lang {
				t.Errorf("GetLang() = %v, want %v", got, tt.args.lang)
			}
		})
	}
}
