package stdlib

import (
	"context"
	"fmt"
	"log"
	"os"

	liblog "github.com/slcjordan/library/log"
)

func init() {
	flags := log.Ldate | log.Ltime | log.Lmicroseconds | log.Lshortfile | log.LUTC
	liblog.RegisterInfo(withDepth{
		Logger: log.New(os.Stdout, " [INFO] ", flags),
		Depth:  3,
	})
	liblog.RegisterError(withDepth{
		Logger: log.New(os.Stderr, " [ERROR] ", flags),
		Depth:  3,
	})
}

type withDepth struct {
	Logger *log.Logger
	Depth  int
}

func (w withDepth) Printf(ctx context.Context, format string, a ...interface{}) {
	err := w.Logger.Output(w.Depth, fmt.Sprintf(format, a...))
	if err != nil {
		panic(err)
	}
}
