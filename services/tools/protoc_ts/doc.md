# 生成ts代码

本工具可以根据指定proto文件，生成对应的ts代码

## 使用

安装：`make install`
使用：`make protoc`

输入参数介绍：`-h`
使用介绍：参考相应powershell脚本

## 数据结构

```go 
type pbFile struct {
	name     string
	imports  []*dependence
	payloads []*payload

	payloadStack *payloadStack
}

type dependence struct {
	from       string   // e.g. "common", 'common.proto' remove extension name
	structures []string // e.g. []string{"Error","Pagination"}
}

type payload struct {
	name        string // e.g. "Note", from 'message Note {'
	typ         string // e.g. "namespace"/"class"/"enum"
	indentation int    // payload indentation, field indentation is 'v'+2
	fields      []*field
	children    []*payload
}

// field include 'message' and 'enum' field
type field struct {
	name    string // e.g. "note_id" from 'string note_id = 1;'
	typ     string // e.g. "string"/"ListRule"/"enum"
	isArray bool   // e.g. "repeated Data data = 1;"
}
```

## 实现

将一个`.proto`文件中的message，翻译成ts代码

步骤：

1. 初始化：工作路径、缩进、用于解析proto文件的正则表达式
2. 遍历路径下的每一个`.proto`文件，将其抽象成一般格式(`*pbFile`)
3. 根据一般格式的实例，生成ts代码，写入文件

解析一个proto文件：

1. 按行读取，使用正则表达式判断当前行内容，并决定行为

根据一般格式的实例生成ts代码：

1. 基本格式：生成文件注释、import、已翻译成ts的message

## 优化

1. 目前采用按行读取、正则表达式匹配行为的方式，考虑调整为按byte读取，使用有限状态机(FSM)解析

## 问题

1. 自用工具，与proto文件编写规范双向奔(tuo)赴(xie)，如果你用到了我没有用到的proto语法，请谨慎使用本工具
2. 第一版仅包含根据message生成结构体的功能
3. 要求至少所有被引用的proto文件，包名与文件名一致，这样做是为了import的时候，可以立刻获知一个结构应该从哪个文件中引入
    1. 该问题主要受限于**是否允许导入范围外的proto文件**，即导入文件`a.proto`，生成工具都不知道它在哪，  
       也就没办法在生成的ts文件中import相应结构    
       如果不要求proto文件的包名，我能想到的办法是根据proto文件的import路径去找导入文件，然后解析它的包名  
       不过由于该功能比较边缘，且可以通过约定解决，暂不处理
4. 删除了一个proto文件，需要手动删除该文件曾经的ts生成文件
