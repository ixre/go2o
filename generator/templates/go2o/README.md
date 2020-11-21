# protobuf模板

global.proto
```
syntax = "proto3";
option go_package = ".;proto";
option java_package="com.github.go2o.rpc";

message Empty {
}
message String {
  string Value = 1;
}
message Int64 {
  sint64 Value = 1;
}
message Int32 {
  sint32 Value = 1;
}
message Bool {
  bool Value = 1;
}
message StringMap {
  map<string, string> Value = 1;
}
message StringArray {
  repeated string Value = 1;
}
//传输结果对象
message Result {
  /* 状态码,如为0表示成功 */
  sint32 ErrCode = 1;
  /* 消息 */
  string ErrMsg = 2;
  /** 数据字典 */
  map<string, string> Data = 3;
}

/** 分页参数 */
message SPagingParams {
  // 开始记录数
  sint64 Begin = 1;
  // 结束记录数
  sint64 End = 2;
  // 条件
  string Where = 3;
  // 排序字段
  string SortBy = 4;
  // 参数
  map<string, string> Parameters = 5;
}

/** 分页结果 */
message SPagingResult {
  /** 代码 */
  sint32 ErrCode = 1;
  /** 消息 */
  string ErrMsg = 2;
  /** 总数 */
  sint32 Count = 3;
  /** 数据 */
  string Data = 4;
  /** 额外的数据 */
  map<string, string> Extras = 5;
}


message Id {
  int64 Value = 1;
}
```

serviceUtil
```

// 服务工具类，实现的服务组合此类,可直接调用其方法
type serviceUtil struct{}

// 返回失败的结果
func (s serviceUtil) failResult(msg string) *proto.Result {
	return s.failCodeResult(1, msg)
}

// 返回错误的结果
func (s serviceUtil) error(err error) *proto.Result {
	if err == nil {
		return s.success(nil)
	}
	return s.failResult(err.Error())
}

// 返回结果
func (s serviceUtil) result(err error) *proto.Result {
	if err == nil {
		return s.success(nil)
	}
	return s.error(err)
}

// 返回自定义编码的结果
func (s serviceUtil) resultWithCode(code int, message string) *proto.Result {
	return &proto.Result{ErrCode: int32(code), ErrMsg: message, Data: map[string]string{}}
}

// 返回失败的结果
func (s serviceUtil) errorCodeResult(code int, err error) *proto.Result {
	return &proto.Result{ErrCode: int32(code), ErrMsg: err.Error(), Data: map[string]string{}}
}

// 返回失败的结果
func (s serviceUtil) failCodeResult(code int, msg string) *proto.Result {
	return &proto.Result{ErrCode: int32(code), ErrMsg: msg, Data: map[string]string{}}
}

// 返回成功的结果
func (s serviceUtil) success(data map[string]string) *proto.Result {
	if data == nil {
		data = map[string]string{}
	}
	return &proto.Result{ErrCode: 0, ErrMsg: "", Data: data}
}

// 将int32数组装换为int数组
func (s serviceUtil) intArray(values []int32) []int {
	arr := make([]int, len(values))
	for i, v := range values {
		arr[i] = int(v)
	}
	return arr
}

// 转换为JSON
func (s serviceUtil) json(data interface{}) string {
	if data == nil {
		return "{}"
	}
	r, err := json.Marshal(data)
	if err != nil {
		return "{\"error\":\"parse error:" + err.Error() + "\"}"
	}
	return string(r)
}

// 分页响应结果
func (s serviceUtil) pagingResult(total int, data interface{}) *proto.SPagingResult {
	switch data.(type) {
	case string:
		return &proto.SPagingResult{
			Count:  int32(total),
			Data:   data.(string),
			Extras: map[string]string{},
		}
	}
	r, _ := json.Marshal(data)
	return &proto.SPagingResult{
		ErrCode: 0,
		ErrMsg:  "",
		Count:   int32(total),
		Data:    string(r),
		Extras:  map[string]string{},
	}
}
```