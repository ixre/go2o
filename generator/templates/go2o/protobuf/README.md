global.proto

```
syntax = "proto3";
option go_package = "./;proto";
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