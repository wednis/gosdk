package gosdk

var debug bool

func OnDebug() {
	debug = true
}

func OffDebug() {
	debug = false
}

func IsDebug() bool {
	return debug
}
