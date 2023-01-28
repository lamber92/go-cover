# go-cover

[**简体中文**](https://github.com/lamber92/go-cover) **|** **English**

A tool for generating full/differential coverage reports (HTML files) by converting Go-Coverage-Profile.



## Important statement

Most of the source code, implementation ideas and report styles of this tool come from the following two open source projects. Thank them very much:  
(If infringement, please notify me to delete this repository)

- [**axw/gocov**](https://github.com/axw/gocov)：Coverage reporting tool for The Go Programming Language

- [**matm/gocov-html**](https://github.com/matm/gocov-html)：This is a simple helper tool for generating HTML output from [axw/gocov](https://github.com/axw/gocov/)

Functional differences from the above items：

- Support for generating differential code coverage reports
  - Based on the comparison of two different git branches to get the difference
- Highlight the covered code line
- Support coverage report result merging (waiting for implementation)


## Dependencies

- Go Coverage Profile(explained below)
- Go version >= 1.16.x
- Git version >= 2.22


## Usage

#### Step 1: Generate Go Coverage Profile

- (**Recommend**) If you want to get the instantaneous coverage information during the running of the program, use [**qiniu/goc**](https://github.com/qiniu/goc)
- Go official toolchain can be used：`go test -cover`
  - For binary programs, you can use `TestFunc()` to wrap `main()` to achieve; but due to the time limit of the test method to generate the profile, it is not as convenient to use as goc;
  - For versions above Go1.20, new version features are available: https://go.dev/testing/coverage/

#### Step 2: Generate coverage report (HTML)

  ```shell
  go convert <go-profile filepath>
  ```

#### [More examples](https://github.com/lamber92/go-cover-example)



## Detailed command

| Command                                                                                                                                                                                | Option Key                                                                                                                                                    | Option Value                                                                                                                                                                                                                                                                        |
|----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|---------------------------------------------------------------------------------------------------------------------------------------------------------------|-------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| **convert** \<go-coverage-profile filepath\><br>Load and convert Go-Coverage-Profile file<br>and generate HTML report                                                                  | **-o**<br>Output report mode.<br>Optional, default: **\<all\>**                                                                                               | **all**：Output differential & full coverage report<br>**full-only**：Only output the full coverage report (full.html)<br>**diff-only**：Only output the differential coverage report (diff.html)<br>**json-only**：Only output the json information of the intermediate state (stdout) |
|                                                                                                                                                                                        | **-f** \<css-format-filepath\><br>HTML report rendering style file path.<br>Optional, use internal style by default                                           | -                                                                                                                                                                                                                                                                                   |
|                                                                                                                                                                                        | **-d** \<diff-filepath\><br>Branch code diff information file path.<br>Optional, by default, it is obtained according to the combination of -c and -t options | -                                                                                                                                                                                                                                                                                   |
|                                                                                                                                                                                        | **-c**<br>The Git branch name of the current project<br>Optional, by default, call the git command in the program to obtain                                   | -                                                                                                                                                                                                                                                                                   |
|                                                                                                                                                                                        | **-t**<br>The name of the Git branch being compared in the current project<br>Optional, the master branch is used by default                                  | -                                                                                                                                                                                                                                                                                   |
|                                                                                                                                                                                        | **-i**<br>The hash_id interval submitted by the current branch<br>Optional, all submission points are collected by default                                    | 格式：start-hash-id,end-hash-id                                                                                                                                                                                                                                                        |
| **diff** \<diff-filepath\><br>The file path that records the difference information between the current branch and the compared branch<br>                                             | **-c**<br>The Git branch name of the current project<br>Optional, by default, call the git command in the program to obtain                                   | -                                                                                                                                                                                                                                                                                   |
|                                                                                                                                                                                        | **-t**<br>The name of the Git branch being compared in the current project<br>Optional, the master branch is used by default                                  | -                                                                                                                                                                                                                                                                                   |
| **trim** \<go-cover json filepath\><br>Load the intermediate json file generated by go-cover,<br>and cut out the information that needs to be preserved based on the diff file.        | **-d** \<diff-filepath\><br>Branch code diff information file path<br>Required                                                                                | -                                                                                                                                                                                                                                                                                   |
| (waiting for implementation)<br>**report** \<go-cover json filepath\><br>Load the intermediate json file generated by go-cover,<br>and generate the corresponding coverage HTML report | (waiting for implementation)                                                                                                                                  | (waiting for implementation)                                                                                                                                                                                                                                                        |



## TODO List

- Enables multiple coverage reports to be merged
  - Referring to the jacoco implementation scheme, it supports merging coverage reports based on different Git submission points
- Implement \<report\> function
- (busy farming, and the update time is random; if you like it, please help optimize it~ Thanks~)
