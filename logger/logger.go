package logger

import "log"

// ANSI 颜色代码
const (
	ColorReset  = "\033[0m"
	ColorRed    = "\033[31m"
	ColorGreen  = "\033[32m"
	ColorYellow = "\033[33m"
	ColorBlue   = "\033[34m"
	ColorPurple = "\033[35m"
	ColorCyan   = "\033[36m"
	ColorWhite  = "\033[37m"
)

// 图标
const (
	IconSuccess = "✅"
	IconError   = "❌"
	IconInfo    = "ℹ️"
	IconWarning = "⚠️"
	IconStart   = "🚀"
	IconKey     = "🔑"
	IconAddress = "🏠"
	IconNetwork = "🌐"
)

// 自定义日志函数
func LogInfo(format string, v ...interface{}) {
	log.Printf(ColorCyan+IconInfo+" INFO: "+format+ColorReset, v...)
}

func LogSuccess(format string, v ...interface{}) {
	log.Printf(ColorGreen+IconSuccess+" SUCCESS: "+format+ColorReset, v...)
}

func LogWarning(format string, v ...interface{}) {
	log.Printf(ColorYellow+IconWarning+" WARNING: "+format+ColorReset, v...)
}

func LogError(format string, v ...interface{}) {
	log.Printf(ColorRed+IconError+" ERROR: "+format+ColorReset, v...)
}

func LogStart(format string, v ...interface{}) {
	log.Printf(ColorBlue+IconStart+" START: "+format+ColorReset, v...)
}
func LogLoading(msg string) {
	log.Printf("%s[⟳] %s%s", ColorCyan, msg, ColorReset)
}

func LogStep(msg string) {
	log.Printf("%s[➤] %s%s", ColorWhite, msg, ColorReset)
}
