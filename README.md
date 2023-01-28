# go-cover

**简体中文** **|** [**English**](https://github.com/lamber92/go-cover/blob/main/README_en.md)

通过转换Go-Coverage-Profile，生成全量/增量覆盖率报告(HTML文件)的工具。



## 重要声明

本工具大部分源码、实现思路及报告样式来源于以下两款开源项目，非常感谢他们：  
(如果侵权，请通知我删除本仓库)

- [**axw/gocov**](https://github.com/axw/gocov)：Coverage reporting tool for The Go Programming Language

- [**matm/gocov-html**](https://github.com/matm/gocov-html)：This is a simple helper tool for generating HTML output from [axw/gocov](https://github.com/axw/gocov/)

与以上项目的功能差异：

- 支持生成增量代码覆盖率报告
  - 基于两个不同 Git 分支对比得到差异
- 对已覆盖代码行高亮标识
- 支持覆盖率报告结果合并(待实现)


## 依赖

- Go Coverage Profile(下文有说明)
- Go 版本 >= 1.16.x
- Git 版本 >= 2.22


## 用法

#### 步骤1: 生成Go Coverage Profile

- **（推荐）**如果想得到程序运行过程中瞬时覆盖率信息，使用[**qiniu/goc**](https://github.com/qiniu/goc)
- 可以使用Go官方工具链：`go test -cover`
  - 对二进制程序，可以使用`TestFunc()`包裹`main()`配合实现；但由于test方式生成profile的时机限制，使用上没有goc方便；
  - 对于Go1.20以上版本，可以使用新版本特性：https://go.dev/testing/coverage/

#### 步骤2: 生成覆盖率报告(HTML文件)

  ```shell
  go convert <go-profile filepath>
  ```

#### [更多示例集](https://github.com/lamber92/go-cover-example)



## 命令详解

| 命令                                                                                         | 选项键                                                              | 选项值                                                       |
|--------------------------------------------------------------------------------------------|------------------------------------------------------------------| ------------------------------------------------------------ |
| **convert** \<go-coverage-profile filepath\><br>加载并转换go-coverage-profile文件<br>并生成HTML报告             | **-o**<br>输出报告的模式<br>选填，缺省时使用**\<all\>**                         | **all**：输出增量&全量覆盖率报告<br>**full-only**：只输出全量覆盖率报告 (full.html)<br>**diff-only**：只输出增量覆盖率报告 (diff.html)<br>**json-only**：只输出中间态的json信息 (stdout) |
|                                                                                            | **-f** \<css-format-filepath\><br>HTML报告渲染样式文件路径<br>选填，缺省时使用内部样式 | -                                                            |
|                                                                                            | **-d** \<diff-filepath\><br>分支代码差异信息文件路径<br>选填，缺省时按-c与-t组合选项获取   | -                                                            |
|                                                                                            | **-c**<br>当前项目，当前git分支名称<br>选填，缺省时程序内调用git命令获取                   | -                                                            |
|                                                                                            | **-t**<br>当前项目，被对比git分支名称<br>选填，缺省时使用master分支                    | -                                                            |
|                                                                                            | **-i**<br>当前分支提交的hash_id区间<br>选填，缺省时采集所有提交点       | 格式：start-hash-id,end-hash-id                              |
| **diff** \<diff-filepath\><br>记录有当前分支与被对比分支的差异信息的文件路径<br>                                  | **-c**<br/>当前项目的git分支名称<br/>选填，缺省时程序内调用git命令获取                   | -                                                            |
|                                                                                            | **-t**<br/>当前项目的git被对比分支名称<br/>选填，缺省时使用master分支                  | -                                                            |
| **trim** \<go-cover json filepath\><br>加载go-cover生成的中间态json文件<br>并以diff文件为依据裁剪出需要<br>保留的信息 | **-d** \<diff-filepath\><br>分支代码差异信息文件路径<br/>必填                    | -                                                            |
| (待实现)<br>**report** \<go-cover json filepath\><br>加载go-cover生成的中间态json文件<br>生成对应的覆盖率HTML报告 | (待实现)                                                            | (待实现)                                                     |



## TODO List

- 实现 覆盖率报告 结果合并
  - 参照jacoco实现方案，支持合并基于不同Git提交点的覆盖率报告
- 实现 go-cover **report** 功能
- (底层码农搬砖中，更新时间随缘；如果对你有用，请帮忙优化它~感谢~)
