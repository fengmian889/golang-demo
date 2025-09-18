# 练习需求
用 cobra 从零建一个名为 memo 的 CLI，要求支持以下子命令与选项，且所有数据只保存在本地 JSON 文件（无需数据库）
# 必须实现的命令与功能
- 新增一条备忘录；若同标题已存在，提示冲突并拒绝
```
memo add -t <title> -c "<content>"
```
- 默认只显示未归档条目；-a 同时显示已归档；-k 按标题或内容关键词过滤
```
memo list [-a] [-k <keyword>]
```
- 把指定标题的条目标记为“已完成”（即归档）
```
memo done <title>
```
- 删除指定标题的条目，存在才删，不存在给出提示
```
memo del <title>
```
- 删除所有已归档条目，并告知删掉了多少条
```
memo clean
```
# 全局可选 flag
- --data string  指定数据文件路径（默认 ~/.memo.json）
- --no-color     关闭所有彩色输出，方便脚本调用
