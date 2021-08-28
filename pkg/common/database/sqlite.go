package database

import (
	"context"
	"os"
)

type Sqlite struct {
	File     *os.File
	FullPath string
	Driver   string
	Ctx      context.Context
}
