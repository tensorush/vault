package sqlite

import (
	"database/sql"
	"log"
	"os"
	"reflect"
	"testing"

	"vault-bot/internal/database/queries"
	"vault-bot/internal/secret"
)

var st *Sqlite3
var dbName = "test.db"

func TestMain(m *testing.M) {
	db, err := sql.Open("sqlite", dbName)
	if err != nil {
		log.Fatalf("can't opening the db: %v", err)
	}
	defer cleanUp(dbName)
	defer db.Close()

	st, err = New(db, "file://..//..//..//schema/sqlite")
	if err != nil {
		log.Fatalf("can't creating the db: %v", err)
	}

	err = queries.Prepare(db, "sqlite")
	if err != nil {
		log.Fatalf("error preparing db: %v", err)
	}

	m.Run()

}

func cleanUp(filename string) {
	if err := os.Remove(filename); err != nil {
		log.Fatalf("can't remove the db: %v", err)
	}
}

func TestDB_Delete(t *testing.T) {
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
				chatID:  1,
				service: "vk",
			},
		},
		{
			name: "ok 2",
			args: args{
				chatID:  2,
				service: "yandex",
			},
		},
		{
			name: "not found",
			args: args{
				chatID:  3,
				service: "test",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if !tt.wantErr {
				_, err := st.Exec(
					"INSERT INTO services (service, login, password, owner)  VALUES (?, ?, ?, ?)",
					tt.args.service, "test", "test", tt.args.chatID,
				)
				if err != nil {
					t.Errorf("can't insert the record: %v", err)
					return
				}
			}

			if err := st.Delete(tt.args.chatID, tt.args.service); (err != nil) != tt.wantErr {
				t.Errorf("Delete() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestDB_Get(t *testing.T) {
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
				chatID:  1,
				service: "vk.com",
			},
			want: secret.Credentials{
				Login:    "test",
				Password: "test",
			},
		},
		{
			name: "ok 2",
			args: args{
				chatID:  2,
				service: "yandex.ru",
			},
			want: secret.Credentials{
				Login:    "test",
				Password: "test",
			},
		},
		{
			name: "not found",
			args: args{
				chatID:  3,
				service: "test.com",
			},
			wantErr: true,
			want:    secret.Credentials{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if !tt.wantErr {
				_, err := st.Exec(
					"INSERT INTO services (service, login, password, owner)  VALUES (?, ?, ?, ?)",
					tt.args.service, "test", "test", tt.args.chatID,
				)
				if err != nil {
					t.Errorf("can't insert the record: %v", err)
					return
				}
			}
			got, err := st.Get(tt.args.chatID, tt.args.service)
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

func TestDB_GetLang(t *testing.T) {
	type args struct {
		chatID int64
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
				chatID: 11,
			},
			want: "pt",
		},
		{
			name: "ok 2",
			args: args{
				chatID: 22,
			},
			want: "en",
		},
		{
			name: "not found",
			args: args{
				chatID: 33,
			},
			wantErr: true,
			want:    "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if !tt.wantErr {
				_, err := st.Exec(
					"INSERT INTO chats (chat_id, chat_lang)  VALUES (?, ?)",
					tt.args.chatID, tt.want)
				if err != nil {
					t.Errorf("can't insert the record: %v", err)
				}
			}
			got, err := st.GetLang(tt.args.chatID)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetLang() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetLang() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDB_Save(t *testing.T) {
	type args struct {
		chatID  int64
		service string
		secret  secret.Credentials
	}

	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "ok",
			args: args{
				chatID:  111,
				service: "vk.ru.com",
				secret: secret.Credentials{
					Login:    "test",
					Password: "XXXX",
				},
			},
		},
		{
			name: "ok 2",
			args: args{
				chatID:  222,
				service: "yan2x.ru",
				secret: secret.Credentials{
					Login:    "test",
					Password: "XXXX",
				},
			},
		},
		{
			name: "duplicate",
			args: args{
				chatID:  222,
				service: "yan2x.ru",
				secret: secret.Credentials{
					Login:    "test",
					Password: "XXXX",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := st.Save(tt.args.chatID, tt.args.service, tt.args.secret); (err != nil) != tt.wantErr {
				t.Errorf("Save() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr {
				got, err := st.Get(tt.args.chatID, tt.args.service)
				if err != nil {
					t.Errorf("can't get the record: %v", err)
				}
				if !reflect.DeepEqual(got, tt.args.secret) {
					t.Errorf("Save() got = %v, want %v", got, tt.args.secret)
				}
			}
		})
	}
}

func TestDB_SetLang(t *testing.T) {
	type args struct {
		chatID int64
		lang   string
	}

	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "ok",
			args: args{
				chatID: 111,
				lang:   "pt",
			},
		},
		{
			name: "ok 2",
			args: args{
				chatID: 222,
				lang:   "en",
			},
		},
		{
			name: "duplicate",
			args: args{
				chatID: 222,
				lang:   "en",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := st.SetLang(tt.args.chatID, tt.args.lang); (err != nil) != tt.wantErr {
				t.Errorf("SetLang() error = %v, wantErr %v", err, tt.wantErr)
			}

			if !tt.wantErr {
				got, err := st.GetLang(tt.args.chatID)
				if err != nil {
					t.Errorf("can't get the record: %v", err)
				}
				if got != tt.args.lang {
					t.Errorf("Save() got = %v, want %v", got, tt.args.lang)
				}
			}
		})
	}
}
