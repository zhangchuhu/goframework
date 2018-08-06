package servant

import (
	"code.yy.com/yytars/goframework/kissgo/appzaplog"
	"code.yy.com/yytars/goframework/kissgo/appzaplog/zap"
)

func checkPanic() {
	if r := recover(); r != nil {
		//path, _ := filepath.Abs(filepath.Dir(os.Args[0]))
		//os.Chdir(path)
		//file, _ := os.Create(fmt.Sprintf("panic.%s", time.Now().Format("20060102-150405")))
		//file.WriteString(string(debug.Stack()))
		//file.Close()
		appzaplog.DPanic("panic recovered ", zap.Any("panic", r))
		//os.Exit(-1)
	}
}
