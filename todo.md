# 待完善的功能 (Polish Needed)
这些地方目前可用，但可以做得更好：

- HTML 转文本的质量: 题目描述的 htmlToText 实现（虽然我没直接看到代码，但推测比较基础）可能无法完美处理复杂的数学公式 (LaTeX) 或表格，导致生成的注释阅读体验一般。
- 分页支持: ltgo list 目前硬编码获取前 50 题，无法查看后面的题目。
- 搜索范围: ltgo gen 目前搜索范围限制在前 2000 题，对于新出的题目可能搜不到（除非直接用精确的 Slug 或 ID）。
- 文件名强依赖: run 和 submit 命令强依赖文件名格式 ID_slug.go 来提取 slug。如果用户重命名了文件，命令会失败。
- 建议：可以在生成文件的注释头里写入 @lc app=leetcode.cn id=1 lang=golang 这样的元数据，读取时优先解析元数据，而不是依赖文件名。
- 硬编码语言: 目前仅支持 Go 语言 (lang: "golang" 硬编码在 Client 和 Generator 中)。
# 建议补充的功能 (Future Features)
如果你想让这个工具更上一层楼，可以考虑以下功能：

- 每日一题 (ltgo daily):
    - 增加一个命令自动获取当天的每日一题并生成文件，这是刷题用户最高频的需求之一。
- 自定义测试用例:
    - ltgo run 目前只跑官方的 Sample Case。建议支持 ltgo run file.go --test "1,2,3"，允许用户输入自定义数据进行调试。
- 打开浏览器 (ltgo open):
    - 在终端看题目描述有时比较累，增加一个命令直接调用系统浏览器打开当前题目的网页。
- 显示更多题目:
    - 给 list 增加参数 ltgo list --page 2 或 ltgo list --limit 100。
