#!type:1
#!target:{{.global.pkg}}/proto/{{.table.Name}}_service.proto
{{$comment := .table.Comment}}\
{{$title := .table.Title}}\
{{$shortTitle := .table.ShortTitle}}
{{$pkType := join .table.Title .table.PkProp ""}}
syntax = "proto3";
option go_package = "./;proto";
option java_package = "{{replace .global.pkg "/" "."}}.rpc";

import "global.proto";

// {{$comment}}服务
service {{$title}}Service {
    // 保存{{$comment}}
    rpc Save{{$shortTitle}} (Save{{$shortTitle}}Request) returns (Save{{$shortTitle}}Response) {
    }
    // 获取{{$comment}}
    rpc Get{{$shortTitle}} ({{$pkType}}) returns (S{{$shortTitle}}) {
    }
    // 获取{{$comment}}列表
    rpc Query{{$shortTitle}}List (Query{{$shortTitle}}Request) returns (Query{{$shortTitle}}Response) {
    }
    // 删除{{$comment}}
    rpc Delete{{$shortTitle}} ({{$pkType}}) returns (Result) {
    }
    // 获取{{$comment}}分页数据
    rpc Paging{{$shortTitle}} ({{$shortTitle}}PagingRequest) returns ({{$shortTitle}}PagingResponse);
}

message Save{{$shortTitle}}Request{
    {{range $i,$c := .columns}}
    /** {{$c.Comment}} */
    {{type "protobuf" $c.Type}} {{$c.Prop}} = {{plus $c.Ordinal 1}};{{end}}
}

message Save{{$shortTitle}}Response{
    int64 ErrCode = 1;
    string ErrMsg = 2;
    {{type "protobuf" .table.PkType}} {{.table.PkProp}} = 3;
}

message {{$pkType}}{
   {{type "protobuf" .table.PkType}} Value = 1;
}

message S{{$shortTitle}}{
    {{range $i,$c := .columns}}
    /** {{$c.Comment}} */
    {{type "protobuf" $c.Type}} {{$c.Prop}} = {{plus $c.Ordinal 1}};{{end}}
}

message Query{{$shortTitle}}Request{
    /** 自定义参数 */
}

message Query{{$shortTitle}}Response{
    repeated S{{$shortTitle}} Value = 1;
}

message Paging{{$shortTitle}}{
    {{range $i,$c := .columns}}
    /** {{$c.Comment}} */
    {{type "protobuf" $c.Type}} {{$c.Prop}} = {{plus $c.Ordinal 1}};{{end}}
}

message {{$shortTitle}}PagingRequest{
    // 分页参数
    SPagingParams Params = 1;
}

message {{$shortTitle}}PagingResponse {
   int64 Total = 1;
   repeated Paging{{$shortTitle}} Value = 2;
}
