

kou@MacPro app % go mod init nulab-exam.backlog.jp/KOU/app
go: creating new go.mod: module nulab-exam.backlog.jp/KOU/app
kou@MacPro app %
mkdir -p cmd/app static templates
go get github.com/gin-gonic/gin

go get github.com/google/uuid
go get golang.org/x/oauth2
go get github.com/gin-gonic/gin

lsof -i :8081
kill -9 portno1 portno2

Backlog API で課題を登録する
https://support-ja.backlog.com/hc/ja/articles/360046783973-Backlog-API-%E3%81%A7%E8%AA%B2%E9%A1%8C%E3%82%92%E7%99%BB%E9%8C%B2%E3%81%99%E3%82%8B#Go