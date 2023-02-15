tag配置说明：

格式为：**\`bor:function(param...)\`**，各function及其param用法如下：

- `const(param)`，固定为param
- `incrByN(start,n)`，从start开始每次自增n
- `rand(start,end)`，在start到end之间随机，仅支持数字或字母，参数为字母时，start和end的长度必须相等
- `time(format)`，按format格式生成时间字符串
- `foreign(first.second[index].third)`，first对象的second属性的数组中，索引值为index对象的third属性的值
- `this#field`，当前对象的field属性的值
- `sum(param1,param2...)`，param1、param2...等的累加
- `toStrf(format,number...)`，把数字number按format格式化后的字符串
- `[n]struct`，结构体数组有n个元素

==注意：字段名必须遵循大驼峰原则；支持函数嵌套，但参数中只能有一个是嵌套的函数==
