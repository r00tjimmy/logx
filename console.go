package logx

import (
	"fmt"
	"io"
	"os"
	"time"
)

// 写出到终端的 logger
type ConsoleLogWriter struct {
	format  string
	writeCh chan *LogRecord // 传输日志的缓冲 channel
}

var (
	consoleOut io.Writer = os.Stdout              // 写入目标
	defaultFmt           = "[%D %T] [%L] (%S) %M" // 默认格式
)

func NewConsoleLogWriter() *ConsoleLogWriter {
	w := &ConsoleLogWriter{
		format:  defaultFmt,
		writeCh: make(chan *LogRecord, LogBufMsgs),
	}
	go w.run(consoleOut)
	return w
}

func (cw *ConsoleLogWriter) SetFormat(format string) {
	cw.format = format
}

func (cw *ConsoleLogWriter) run(out io.Writer) {
	for rec := range cw.writeCh {
		fmt.Fprint(out, FormatLogRecord(cw.format, rec)) // 写入日志
	}
}

// 负责 logger 的日志写入
// 注意：若写入过快导致缓冲 channel 装满此方法会阻塞
func (cw *ConsoleLogWriter) LogWrite(rec *LogRecord) {
	cw.writeCh <- rec
}

// 关闭此 logger
// 负责资源的回收
func (cw *ConsoleLogWriter) Close() {
	close(cw.writeCh)
	time.Sleep(10 * time.Millisecond) // 等待 channel 中的数据全部写回日志
}