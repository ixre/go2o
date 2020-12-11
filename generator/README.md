# 代码生成器

**TTO是一款使用Go编写的代码生成器,可根据模板定制生成代码.**

特点:
- 支持mysql和postgresql数据库
- 支持Go,JAVA,Kotlin,Thrift,Javascript,Html,C#
- 支持代码模板, 支持模板函数

资源:

- [下载地址](https://github.com/ixre/tto/releases/)
- [Go模板语法-中](http://www.g-var.com/posts/translation/hugo/hugo-21-go-template-primer/)
- [Go模板语法-English](https://golang.org/pkg/text/template/)


_注：您看到的文档有可能已经更新，请参见最新[使用文档](https://github.com/ixre/tto)_

## 快速开始

1. 配置数据源
```
下载安装包,解压修改tto.conf文件进行数据源配置.
```
2. 定制修改模板
```
根据实际需求对模板进行修改, 或创建自己的模板. 模板语法请参考: Go Template
```
模板注释,使用`/** #! 注释 */`的语法,使用`#!`与普通的代码注释区分
```
/** #! 这是模板注释,不会出现在代码中 */
```
3. 运行命令生成代码
```bash
tto -conf tto.conf
Usage of tto:
  -arch string
        program language
  -clean
        clean last generate files
  -conf string
        config path (default "./tto.conf")
  -debug
        debug mode
  -o string
        path of output directory (default "./output")
  -table string
        table name or table prefix
  -t string
        path of code templates directory (default "./templates")
  -v    print version
```


## 预定义语法

预定义语法用来在代码模板中定义一些数据, 在生成代码时预定义语法不输入任何内容.
预定义语法格式为: !预定义参数名:预定义参数值

目前,支持的预定义语法如下:

- \#!kind: 定义模板生成类型,0:普通,1:生成所有表 2:按表名前缀生成,默认为0
- \#!target: 用来定义代码文件存放的目标路径
- \#!append: 是否追加到文件,可选值为:true和false , 默认为false
- \#!format: 是否启用格式化代码，可选值为:true和false，默认开启
- \#!lang: 指定当前生成代码的语言
如:
```
#!target:java/{{.global.pkg}}/pojo/{{.table.Title}}Entity.java
```
多个预定义表达式可以放在一行
```
#!format:true#!target:Entity.java
```

## 函数

大/小写函数: lower和upper
```
{{lower .table.Name}}
{{upper .table.Name}}
```
单词首字大写函数:title
```
{{title .table.Name}}
```
首字母小写函数: lower_title
```
{{lower_title .table.Name}}
```
语言类型函数: type
```
{{type "go" .columns[0].Type}}
```
返回SQL/ORM类型: sql_type
```
{{sql_type "py" .columns[0].Type .columns[0].Length}}
```
包名函数: pkg
```
{{pkg "go" .global.pkg}}
```
默认值函数: default
```
{{default "go" .columns[0].TypeId}}
```
是否相等
```
{{equal (3%2) 1}
```
是否与任意值相等,　如表的主键是否为int类型
```
{{equal_any .table.PkType 3 4 5}}
```
替换, 如将`table_name`替换为:`table-name`
```
{{replace "table_name" "_" "-"}}
```
替换N次, 如将`table_name`替换为:`table-name`
```
{{replace_n "table_name" "_" "-" 1}}
```
截取第N个字符位置后的字符串,如以下语句将输出:user_list
```
{{substr_n "sys_user_list" "_" 1}}
```
截取索引为N的元素
```
{{$first_table := get_n .tables 0}}
```
字符组合,如以下语句将输出:`1,2,3`
```
{{join "," "1","2","3"}}
{{$api := join "/" .global.base_path (name_path .table.Name)}}
```

包含函数
```
{{contain .table.Pk "id"}}
```
是否以指定字符开始
```
{{starts_with .table.Pk "user_"}}
```
是否以指定字符结束
```
{{ends_with .table.Pk "_time"}}
```
是否为表的列(数组)的最后一列
```
{{$columns := .columns}}
{{range $,$v := .columns}}{{if is_last $i .columns}} last column {{end}}{{end}}
```
排除列元素, 组成新的列数组, 如：
```
{{ $columns := exclude .columns "id","create_time" }}
```
尝试获取一个列,返回列及是否存在的Boolean, 如: 
```
{{ $c := try_get .columns "update_time" }}
{{if $c}}prop={{$c.Prop}}{{end}}
```
将名称转为路径,规则： 替换首个"_"为"/"
```
{{$path := name_path .table.Name}}
```


## 模板

使用`go template`作为模板引擎, 可以通过内置的函数和语法, 生成任意代码. 如果在模板行的末尾添加`\`,
将自动合并下一行. 项目中集成了部分语言的模板,当然也可以在/templates创建自己的模板

 
模板主要包含三大对象: 
- global
- table
- columns

按所有表(前缀分组)模板包含对象：
- global
- tables

### global

**用于读取全局变量**

1. 输出生成器的版本号
```
// this file created by generate {{.global.version}}
```
2. 输出包名,包名通过配置文件配置.格式为: com/pkg
```
package {{.global.pkg}}
```
如果是Java或其他语言, 包名以"."分割, 可使用pkg函数,如:
```
// java package
package {{pkg "java" .global.pkg}}
// c# namespace
namespace {{pkg "csharp" .global.pkg}}
```
3. 输出当前时间
```
generate time {{.global.time}}
```
4. 获取数据库驱动
可选值：pgsql | mysql, 可针对不同数据库生成代码
```
{{.global.db}}
```
5. 输出自定义变量
用户可以通过在配置文件的节点`[global]`中进行添加变量,如:
```
[global]
base_path="/api"
```
使用以下语法读取变量
```
{{.global.base_path}}
```

### table 数据表对象

数据表对象对来返回表的信息,包含如下属性:

- Name: 表名
- Prefix: 表前缀
- Pk: 主键,默认为:id
- PkProp: 主键属性, 首字母大写
- PkType: 主键类型编号
- Title: 表名单词首字大写,通常用来表示类型,
  如:user_info对应的Title为UserInfo
- Comment: 表注释
- Engine: 数据库引擎
- Schema: 架构
- Charset: 数据库编码
- Ordinal: 表的序号

### columns 数据列对象

数据列对象存储表的数据列数组, 并且可遍历. 每个数据列都包含如下属性:

- Name: 列名
- Prop: 列名首字大写,　通常用作属性
- IsPk: 是否主键(bool)
- IsAuto:  是否自动生成(bool)
- NotNull: 是否不能为空(bool)
- DbType: 数据库数据类型
- Comment: 注释
- Length: 长度
- Type: 类型编号,使用type函数转换为对应语言的类型
- Ordinal: 列的序号

示例:
```
{{range $i,$c := .columns}}
    列名:$c.Name {{if $c.IsPk}}是主键{{end}}, 类型:{{type "java" $c.Type}}
{{end}}
```

## 模板示例

以下代码用于生成Java的Pojo对象, 更多示例点击[这里](templates)

```
#!target:{{.global.pkg}}/pojo/{{.table.Title}}Entity.java
package {{pkg "java" .global.pkg}}.pojo;

import javax.persistence.Basic;
import javax.persistence.Id;
import javax.persistence.Column;
import javax.persistence.Entity;
import javax.persistence.Table;
import javax.persistence.GenerationType;
import javax.persistence.GeneratedValue;

/** {{.table.Comment}} */
@Entity
@Table(name = "{{.table.Name}}", schema = "{{.table.Schema}}")
public class {{.table.Title}}Entity {
    {{range $i,$c := .columns}}{{$type := type "java" $c.Type}}
    private {{$type}} {{$c.Name}}
    public void set{{$c.Prop}}({{$type}} {{$c.Name}}){
        this.{{$c.Name}} = {{$c.Name}}
    }

    /** {{$c.Comment}} */{{if $c.IsPk}}
    @Id
    @GeneratedValue(strategy = GenerationType.IDENTITY){{else}}
    @Basic{{end}}
    @Column(name = "{{$c.Name}}"
     {{if not $c.NotNull}}, nullable = true{{end}}
     {{if ne $c.Length 0}},length = {{$c.Length}}{{end}})
    public {{$type}} get{{$c.Prop}}() {
        return this.{{$c.Name}};
    }
    {{end}}
}

```


**如果您觉得这个项目不错, 请给个star吧.**


<img src="images/cq-alipay.png" width="320" style="display:inline-block"/><img src="images/cq-wxpay.png" width="354" style="display:inline-block"/>
