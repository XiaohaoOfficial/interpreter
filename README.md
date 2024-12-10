# 基于go语言实现的简易解释器

### 1.支持数据的初始化操作

`let a = 10;`

`let b = "hello world"`

`let c = [1,2,3,4,5]`

`let d = {"key": "value"}`

### 2.支持数据的赋值操作

`a= 20;`

`b = "ni hao zhongguo"`

### 3.定义完成了表达式求值的顺序，并配有完整的测试函数

```bash
haomata@MacBookPro interpreter % go test ./...
?       interpreter     [no test files]
?       interpreter/object      [no test files]
?       interpreter/repl        [no test files]
?       interpreter/token       [no test files]
ok      interpreter/ast (cached)
ok      interpreter/evaluator   (cached)
ok      interpreter/lexer       (cached)
ok      interpreter/parser      (cached)
```

### 4.配有内置函数

```bash
haomata@MacBookPro interpreter % go run main.go
Hello haomata!This is the Mata Programming language
Feel free to type in commands
enter message "quit" to quit
>>let a = 10;
>>let arr = [1,2,3,4,5];
>>puts(arr);
[1, 2, 3, 4, 5]
>>first(arr);
1
>>rest(arr);
[2, 3, 4, 5]
>>len(arr);
5
>>quit
```

