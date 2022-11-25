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

/** {{$comment}}服务 */
service {{$title}}Service {
    /** 保存{{$comment}} */
    rpc Save{{$shortTitle}} (Save{{$shortTitle}}Request) returns (Save{{$shortTitle}}Response) {
    }
    /** 获取{{$comment}} */
    rpc Get{{$shortTitle}} ({{$pkType}}) returns (S{{$shortTitle}}) {
    }
    /** 获取{{$comment}}列表 */
    rpc Query{{$shortTitle}}List (Query{{$shortTitle}}Request) returns (Query{{$shortTitle}}Response) {
    }
    /** 删除{{$comment}} */
    rpc Delete{{$shortTitle}} ({{$pkType}}) returns (Result) {
    }
    /** 获取{{$comment}}分页数据 */
    rpc Paging{{$shortTitle}} ({{$shortTitle}}PagingRequest) returns ({{$shortTitle}}PagingResponse);
}

/** 保存{{$comment}}请求 */
message Save{{$shortTitle}}Request{
    {{range $i,$c := .columns}}
    /** {{$c.Comment}} */
    {{type "protobuf" $c.Type}} {{lower_title $c.Prop}} = {{plus $c.Ordinal 1}};{{end}}
}

/** 保存{{$comment}}响应 */
message Save{{$shortTitle}}Response{
    int64 errCode = 1;
    string errMsg = 2;
    {{type "protobuf" .table.PkType}} {{lower_title .table.PkProp}} = 3;
}

/** {{$comment}}编号 */
message {{$pkType}}{
   {{type "protobuf" .table.PkType}} value = 1;
}

/** {{$comment}} */
message S{{$shortTitle}}{
    {{range $i,$c := .columns}}
    /** {{$c.Comment}} */
    {{type "protobuf" $c.Type}} {{lower_title $c.Prop}} = {{plus $c.Ordinal 1}};{{end}}
}

/** 查询{{$comment}}请求 */
message Query{{$shortTitle}}Request{
    /** 自定义参数 */
}

/** 查询{{$comment}}响应 */
message Query{{$shortTitle}}Response{
    repeated S{{$shortTitle}} value = 1;
}

/** {{$comment}}分页数据 */
message Paging{{$shortTitle}}{
    {{range $i,$c := .columns}}
    /** {{$c.Comment}} */
    {{type "protobuf" $c.Type}} {{lower_title $c.Prop}} = {{plus $c.Ordinal 1}};{{end}}
}

/** {{$comment}}分页请求 */
message {{$shortTitle}}PagingRequest{
    // 分页参数
    SPagingParams params = 1;
}

/** {{$comment}}分页响应 */
message {{$shortTitle}}PagingResponse {
   int64 total = 1;
   repeated Paging{{$shortTitle}} value = 2;
}
