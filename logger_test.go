package glogger

import (
	"bytes"
	"io"
	"log"
	"os"
	"strings"
	"sync"
	"testing"
)

func Test_checkLogLevel(t *testing.T) {
	tests := []struct {
		name  string
		level string
		want  LogLevel
	}{
		{name: "Test-0", level: "", want: LogLevel{name: "DEBUG", level: 1, formatterFlag: log.Ldate | log.Ltime | log.Lshortfile}},
		{name: "Test-1", level: "DEBUG", want: LogLevel{name: "DEBUG", level: 1, formatterFlag: log.Ldate | log.Ltime | log.Lshortfile}},
		{name: "Test-2", level: "INFO", want: LogLevel{name: "INFO", level: 2, formatterFlag: log.Ldate | log.Ltime | log.Lshortfile}},
		{name: "Test-3", level: "WARN", want: LogLevel{name: "WARN", level: 3, formatterFlag: log.Ldate | log.Ltime | log.Lshortfile}},
		{name: "Test-4", level: "ERROR", want: LogLevel{name: "ERROR", level: 4, formatterFlag: log.Ldate | log.Ltime | log.Lshortfile}},
		{name: "Test-5", level: "TEST", want: LogLevel{name: "DEBUG", level: 1, formatterFlag: log.Ldate | log.Ltime | log.Lshortfile}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			key := "G_LOG_FORMAT"
			if len(tt.level) > 0 {
				key = key + "_" + tt.level
			}
			os.Setenv("G_LOG_LEVEL", tt.level)
			os.Setenv(key, "Lshortfile|LstdFlags")
			if got := checkLogLevel(); got != tt.want {
				t.Errorf("checkLogLevel() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetLogLevel(t *testing.T) {
	tests := []struct {
		name  string
		level string
		want  LogLevel
	}{
		{name: "Test-0", level: "", want: LogLevel{name: "DEBUG", level: 1, formatterFlag: log.Ldate | log.Ltime | log.Lshortfile}},
		{name: "Test-1", level: "DEBUG", want: LogLevel{name: "DEBUG", level: 1, formatterFlag: log.Ldate | log.Ltime | log.Lshortfile}},
		{name: "Test-2", level: "INFO", want: LogLevel{name: "INFO", level: 2, formatterFlag: log.Ldate | log.Ltime | log.Lshortfile}},
		{name: "Test-3", level: "WARN", want: LogLevel{name: "WARN", level: 3, formatterFlag: log.Ldate | log.Ltime | log.Lshortfile}},
		{name: "Test-4", level: "ERROR", want: LogLevel{name: "ERROR", level: 4, formatterFlag: log.Ldate | log.Ltime | log.Lshortfile}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			os.Setenv("G_LOG_LEVEL", tt.level)
			UpdateLogLevel()
			if got := GetLogLevel(); got != tt.want {
				t.Errorf("GetLogLevel() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_isLowerOrEqualThan(t *testing.T) {
	tests := []struct {
		name     string
		l        *LogLevel
		logLevel LogLevel
		want     bool
	}{
		{name: "Test-1", l: &logLevelDebug, logLevel: logLevelDebug, want: true},
		{name: "Test-2", l: &logLevelDebug, logLevel: logLevelInfo, want: true},
		{name: "Test-3", l: &logLevelDebug, logLevel: logLevelWarn, want: true},
		{name: "Test-4", l: &logLevelDebug, logLevel: logLevelError, want: true},
		{name: "Test-5", l: &logLevelError, logLevel: logLevelDebug, want: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.l.isLowerOrEqualThan(tt.logLevel); got != tt.want {
				t.Errorf("LogLevel.isLowerOrEqualThan() = %v, want %v", got, tt.want)
			}
		})
	}
}

func captureOutput(f func()) string {
	reader, writer, err := os.Pipe()
	if err != nil {
		panic(err)
	}
	stdout := os.Stdout
	stderr := os.Stderr
	defer func() {
		os.Stdout = stdout
		os.Stderr = stderr
		log.SetOutput(os.Stderr)
	}()
	os.Stdout = writer
	os.Stderr = writer
	log.SetOutput(writer)
	out := make(chan string)
	wg := new(sync.WaitGroup)
	wg.Add(1)
	go func() {
		var buf bytes.Buffer
		wg.Done()
		io.Copy(&buf, reader)
		out <- buf.String()
	}()
	wg.Wait()
	f()
	writer.Close()
	return <-out
}

func TestDebug(t *testing.T) {
	type args struct {
		format string
		v      []interface{}
	}
	tests := []struct {
		name  string
		level string
		args  args
		want  string
	}{
		{name: "Debug Test-1", level: "DEBUG", args: args{format: "Test-1: %v", v: []interface{}{"DEBUG"}}, want: "Test-1: DEBUG"},
		{name: "Debug Test-2", level: "INFO", args: args{format: "Test-2: %v", v: []interface{}{"INFO"}}, want: ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			os.Setenv("G_LOG_LEVEL", tt.level)
			UpdateLogLevel()
			if got := captureOutput(func() {
				Debug(tt.args.format, tt.args.v...)
			}); !strings.Contains(got, tt.want) {
				t.Errorf("Debug() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInfo(t *testing.T) {
	type args struct {
		format string
		v      []interface{}
	}
	tests := []struct {
		name  string
		level string
		args  args
		want  string
	}{
		{name: "Info Test-1", level: "INFO", args: args{format: "Test-1: %v", v: []interface{}{"INFO"}}, want: "Test-1: INFO"},
		{name: "Info Test-2", level: "WARN", args: args{format: "Test-2: %v", v: []interface{}{"WARN"}}, want: ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			os.Setenv("G_LOG_LEVEL", tt.level)
			UpdateLogLevel()
			if got := captureOutput(func() {
				Info(tt.args.format, tt.args.v...)
			}); !strings.Contains(got, tt.want) {
				t.Errorf("Info() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWarn(t *testing.T) {
	type args struct {
		format string
		v      []interface{}
	}
	tests := []struct {
		name  string
		level string
		args  args
		want  string
	}{
		{name: "Warn Test-1", level: "WARN", args: args{format: "Test-1: %v", v: []interface{}{"WARN"}}, want: "Test-1: WARN"},
		{name: "Warn Test-2", level: "ERROR", args: args{format: "Test-2: %v", v: []interface{}{"ERROR"}}, want: ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			os.Setenv("G_LOG_LEVEL", tt.level)
			UpdateLogLevel()
			if got := captureOutput(func() {
				Warn(tt.args.format, tt.args.v...)
			}); !strings.Contains(got, tt.want) {
				t.Errorf("Warn() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestError(t *testing.T) {
	type args struct {
		format string
		v      []interface{}
	}
	tests := []struct {
		name  string
		level string
		args  args
		want  string
	}{
		{name: "Error Test-1", level: "ERROR", args: args{format: "Test-1: %v", v: []interface{}{"ERROR"}}, want: "Test-1: ERROR"},
		{name: "Error Test-2", level: "", args: args{format: "Test-2: %v", v: []interface{}{"ERROR"}}, want: "Test-2: ERROR"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			os.Setenv("G_LOG_LEVEL", tt.level)
			UpdateLogLevel()
			if got := captureOutput(func() {
				Error(tt.args.format, tt.args.v...)
			}); !strings.Contains(got, tt.want) {
				t.Errorf("Error() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_formatter(t *testing.T) {
	tests := []struct {
		name   string
		level  string
		format string
		want   int
	}{
		{name: "Test-01", level: "DEBUG", format: "Ldate|Ltime|Lmicroseconds|Llongfile|Lshortfile|LUTC|Lmsgprefix|LstdFlags", want: 130},
		{name: "Test-02", level: "DEBUG", format: "Ldate|Ltime|Lmicroseconds|Llongfile|Lshortfile|LUTC|Lmsgprefix", want: 127},
		{name: "Test-03", level: "DEBUG", format: "Ldate|Ltime|Lmicroseconds|Llongfile|Lshortfile|LUTC", want: 63},
		{name: "Test-04", level: "DEBUG", format: "Ldate|Ltime|Lmicroseconds|Llongfile|Lshortfile", want: 31},
		{name: "Test-05", level: "DEBUG", format: "Ldate|Ltime|Lmicroseconds|Llongfile", want: 15},
		{name: "Test-06", level: "DEBUG", format: "Ldate|Ltime|Lmicroseconds", want: 7},
		{name: "Test-07", level: "DEBUG", format: "Ldate|Ltime", want: 3},
		{name: "Test-08", level: "DEBUG", format: "Ldate", want: 1},
		{name: "Test-09", level: "INFO", format: "Ltime|Lmicroseconds|Llongfile|Lshortfile|LUTC|Lmsgprefix", want: 126},
		{name: "Test-10", level: "INFO", format: "Ldate|Ltime|Lmicroseconds|Llongfile|LUTC|Lmsgprefix|LstdFlags", want: 114},
		{name: "Test-11", level: "INFO", format: "Ldate|Ltime|Llongfile|Lshortfile|LUTC|LstdFlags", want: 62},
		{name: "Test-12", level: "INFO", format: "Lshortfile|LstdFlags", want: 19},
		{name: "Test-13", level: "WARN", format: "Ltime|LMICROSECONDS|Llongfile|LSHORTFILE|LUTC|Lmsgprefix", want: 106},
		{name: "Test-14", level: "WARN", format: "LDATE|LTIME|LMICROSECONDS|LLONGFILE|LUTC|LMSGPREFIX|LSTDFLAGS", want: 32},
		{name: "Test-15", level: "ERROR", format: "ldate|Ltime|Llongfile|lshortfile|lutc|LstdFlags", want: 13},
		{name: "Test-16", level: "ERROR", format: "shortfile|stdFlags", want: 0},
		{name: "Test-17", level: "", format: "Lshortfile|LstdFlags", want: 19},
		{name: "Test-17", level: "", want: 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			key := "G_LOG_FORMAT"
			if len(tt.level) > 0 {
				key = key + "_" + tt.level
			}
			os.Setenv("G_LOG_LEVEL", tt.level)
			os.Setenv(key, tt.format)
			if got := formatter(tt.level); got != tt.want {
				t.Errorf("formatter() = %v, want %v", got, tt.want)
			}
		})
	}
}
