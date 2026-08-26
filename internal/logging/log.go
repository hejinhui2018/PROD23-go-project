package logging

import (
	"encoding/json"
	"log"
	"os"
	"time"
)

type Logger struct{ l *log.Logger }

func New() *Logger { return &Logger{log.New(os.Stdout, "", 0)} }
func (x *Logger) Event(level, msg string, fields map[string]any) {
	m := map[string]any{"time": time.Now().UTC(), "level": level, "message": msg}
	for k, v := range fields {
		m[k] = v
	}
	b, _ := json.Marshal(m)
	x.l.Print(string(b))
}
