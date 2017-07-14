package kit

// RPC服务
var RPC = NewRpcToolkit()

// 模板包含函数
var TInc = &templateIncludeToolkit{}

// 模板包含函数包装
var TIncWrapper = &templateIncludeKitWrapper{
	FuncMap:    TInc.getFuncMap(),
	Middleware: TInc.includeMiddle,
}
