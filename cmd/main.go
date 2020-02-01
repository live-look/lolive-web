package main

import (
	_ "github.com/jackc/pgx/v4/stdlib"
	_ "github.com/volatiletech/authboss/auth"
	_ "github.com/volatiletech/authboss/confirm"
	_ "github.com/volatiletech/authboss/logout"
	_ "github.com/volatiletech/authboss/register"

	"gitlab.com/isqad/camforchat"
)

func main() {
	app := &camforchat.Application{}

	app.Run()
}
