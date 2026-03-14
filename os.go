package gosdk

import (
	"os"
	"os/signal"
	"syscall"
)

// 等待退出信号
func WaitExitSignal(fn func()) {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, os.Interrupt, syscall.SIGTERM)
	<-sigs
	fn()
}
