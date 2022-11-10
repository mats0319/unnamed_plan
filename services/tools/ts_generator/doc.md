# 生成ts代码

本工具可以根据指定proto文件，生成对应的ts代码

问题：

1. 与我写的proto文件强关联，如果你用到了我没有用到的proto语法，请谨慎使用本工具
2. 第一版仅包含根据message生成结构体的功能

我原本想着用模板填充的方式，可是学到了一个新的名词：抽象语法树(ast)，似乎这些代码生成都是依赖于这样一棵树做的，了解一下这个概念再做决定

抽象语法树，是源码语法结构的一种抽象表示。大致是将源代码使用类似伪代码的方式描述，然后用这个树去做事情，比如生成其他语言的代码。

1. 词法分析：scan，将代码转化为token流
2. 语法分析：parse，将token流转化为结构化的ast

我原本的想法，是做一个proto和ts的语法转换，例如结构定义，proto中使用message，ts使用interface，找到这些语法差异，然后挨个转换，  
而ast提供了一种新思路，是从整体出发，描述整个proto源代码，然后翻译成ts，这让我更想自己实现一个代码生成工具了。

在lc那边尝试了ast的使用，效果基本符合预期，至于这边，只要套规则就好了

序列化：（根据proto文件生成ast结构）

1. 按行读取文件，如果发现import，则表示该文件依赖其他文件，记录依赖的文件，准备生成ts的import
    1. 每个文件都只管自己，可能存在A import B，但B不存在的情况，用户自行处理
2. 按行读取文件，第一次扫描到`message xxx {`/`enum xxx {`行，开始解析
    1. 针对message和enum，扫描field采用不同的正则表达式
    2. 扫描到`}`行结束
3. 如果扫描过程中再次遇到相同结构，item的type调整为namespace，开始解析新遇到的类型
4. 对于扫描到的类型，如果有其他包的包名，则在保存时直接删除
5. 得到结果：file

反序列化：（根据ast结构生成ts代码）

1. 引入依赖
2. 遍历序列化的结果，如果是class或者enum则直接生成代码
3. 如果是namespace，先根据field生成代码（可能没有），然后遍历其children

tips:

1. 通过读取输入参数，知道要处理哪些proto文件，默认要求输入文件名，支持*.proto
2. 保存当前解析的item栈，遇到一个开始解析入口时，将其加入栈，解析完成后移除
3. 要求至少所有被引用的proto文件，包名与文件名一致，这样做是为了import的时候，可以立刻获知一个结构应该从哪个文件中引入

file:
name: "3_note" // '3_note.proto' just record file name
imports: []dependence
items: []item

dependence:
from: "common" // 'common.proto' remove extension name
structures: []string{"Error","Pagination"}

item:
name: "Note"
type: "namespace"/"export class"/"export enum"
indentation: 0 // item indentation, field indentation is 'v'+2, calc when find child
fields: []field
children: []item

field: // include 'message' and 'enum'
name: "note_id"
type: "string"/"number"(for enum, self-define number)
comment: "something after // flag"

