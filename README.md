# ltgo

A fast and simple CLI tool for LeetCode in Go.

## Introduction

**ltgo** is a command-line interface tool designed to streamline your LeetCode problem-solving workflow. It allows you to browse problems, generate boilerplate code, test solutions remotely, and submit directly from your terminal—all without leaving your development environment.

## Features

- **🔐 Init**: Set up your LeetCode session by configuring your cookie for authentication
- **📋 List**: Browse and view LeetCode problems with their status, difficulty, and titles
- **📝 Gen**: Generate Go solution files with boilerplate code for any problem
- **▶️ Run**: Test your code remotely on LeetCode's servers with sample test cases
- **🚀 Submit**: Submit your solution and get instant feedback on acceptance and performance

## Installation

### Option 1: Install via `go install`

```bash
go install github.com/X-for/ltgo/cmd/ltgo@latest
```

Make sure your `$GOPATH/bin` (or `$GOBIN`) is in your system's `PATH`.

### Option 2: Build from Source

```bash
# Clone the repository
git clone https://github.com/X-for/ltgo.git
cd ltgo

# Build the binary
go build -o ltgo ./cmd/ltgo

# Optionally, move to your PATH
sudo mv ltgo /usr/local/bin/
```

## Usage

### `ltgo init` - Initialize Configuration

Set up your LeetCode credentials to enable authentication for all commands.

```bash
ltgo init
```

**Interactive prompts:**
1. Choose your LeetCode site (`cn` for leetcode.cn or `com` for leetcode.com)
2. Paste your LeetCode cookie from browser developer tools

**Getting your cookie:**
1. Log in to LeetCode in your browser
2. Open Developer Tools (F12)
3. Go to Application/Storage → Cookies
4. Copy the entire cookie string (including `LEETCODE_SESSION` and `csrftoken`)

**Example:**
```
Choose site (cn/com) [default: cn]: com
Please paste your LeetCode Cookie (from browser developer tools):
(Include LEETCODE_SESSION and csrftoken)
> LEETCODE_SESSION=your-session-token; csrftoken=your-csrf-token
```

### `ltgo list` - Browse Problems

View a list of LeetCode problems with their status, ID, title, and difficulty.

```bash
ltgo list
```

**Output:**
```
Fetching questions...
Status  ID    Title                              Difficulty
------  --    -----                              ----------
[✓]     1     两数之和 (Two Sum)                  Easy
[ ]     2     两数相加 (Add Two Numbers)          Medium
[✓]     3     无重复字符的最长子串 (Longest...)   Medium
```

### `ltgo gen` - Generate Solution File

Generate a Go file with boilerplate code for a specific problem. You can search by:
- Problem ID (e.g., `1`, `42`)
- Problem slug (e.g., `two-sum`, `add-two-numbers`)
- Keywords from the title

```bash
ltgo gen <id-or-slug-or-keyword>
```

**Examples:**
```bash
# Using problem ID
ltgo gen 1

# Using slug
ltgo gen two-sum

# Using keywords (must match exactly one problem)
ltgo gen "longest substring"
```

**Output:**
- Creates a file in `./questions/` directory
- Filename format: `<ID>_<slug>.go` (e.g., `0001_two-sum.go`)
- Includes problem description as comments and function signature

### `ltgo run` - Test Code Remotely

Run your solution against LeetCode's test cases without submitting.

```bash
ltgo run <file-path>
```

**Example:**
```bash
ltgo run questions/0001_two-sum.go
```

**Output:**
```
Fetching question info for 'two-sum'...
🚀 Sending code to LeetCode...
Waiting for result...

✅ Accepted

Case 1:
  Input:    [2,7,11,15] 9
  Output:   [0,1]
  Expected: [0,1]
  ------------------------
```

**Note:** The file must follow the naming convention `<ID>_<slug>.go` for the tool to identify the problem.

### `ltgo submit` - Submit Solution

Submit your solution to LeetCode for final judgment.

```bash
ltgo submit <file-path>
```

**Example:**
```bash
ltgo submit questions/0001_two-sum.go
```

**Output on success:**
```
Fetching question info for 'two-sum'...
🚀 Submitting to LeetCode...
Submission ID: 123456
Waiting for result...

✅ Accepted!
Runtime: 4 ms (Beats 95.23%)
Memory:  3.5 MB (Beats 87.45%)
```

**Output on failure:**
```
❌ Wrong Answer
Passed:   15/57 cases
Input:    [3,2,4] 6
Output:   [0,2]
Expected: [1,2]
```

## Quick Start

Here's a complete workflow example:

```bash
# 1. Initialize your configuration (one-time setup)
ltgo init

# 2. Browse available problems
ltgo list

# 3. Generate solution file for problem "Two Sum"
ltgo gen two-sum

# 4. Open and edit the generated file
# vim questions/0001_two-sum.go
# (Write your solution)

# 5. Test your solution remotely
ltgo run questions/0001_two-sum.go

# 6. If tests pass, submit your solution
ltgo submit questions/0001_two-sum.go
```

## Configuration

After running `ltgo init`, your configuration is saved to:
- **Linux/macOS**: `~/.config/ltgo/config.json`
- **Windows**: `%APPDATA%\ltgo\config.json`

## License

This project is open source and available under the MIT License.

## Contributing

Contributions are welcome! Feel free to open issues or submit pull requests.
