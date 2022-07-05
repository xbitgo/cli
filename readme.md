# XbitGo 命令行工具

## 安装

    go install github.com/xbitgo/cli@latest

## 使用帮助

    xbit help

## 配置文件

    项目根目录下的xbit.yaml文件

## 代码生成工具注释标识

    @IGNORE     忽略 用于generate忽略对应实体
    @IMPL[...]  **显示声明实现接口** 中括号内为必须配置的接口名称 格式为domain下的{包名.接口名称}
    @DI[...]    自动注册DI  默认注册为{包名.结构体名称} 可在中括号后面配置自定义注册名

## 代码生成工具字段解析标签

    sdi 注册di 用于配置 需要配置对象有CreateInstance() 方法 
    di 依赖注入 
    db 生成数据库层相关
    pb 生成pb文件相关

## 代码范围

    默认generate命令只会解析domain目录,repo_impl和conf目录