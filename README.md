<div align="center">

# mygo-ui-sdk

**专为 `egoist/mygo` Native GPU UI 打造的 Go 可组合交互行为库、无障碍/Agent 可达性测试套件与 Agent Skill**

[![Go](https://img.shields.io/badge/Go-1.27.4-00ADD8?style=flat-square)](#-快速开始)
[![MyGo](https://img.shields.io/badge/MyGo-v0.10.0-0A84FF?style=flat-square)](#-项目架构)
[![License](https://img.shields.io/badge/License-MIT-10B981?style=flat-square)](LICENSE)
[![CI](https://img.shields.io/github/actions/workflow/status/Lumiaqian/mygo-ui-sdk/ci.yml?branch=main&style=flat-square&label=CI)](.github/workflows/ci.yml)
[![Release](https://img.shields.io/github/v/release/Lumiaqian/mygo-ui-sdk?style=flat-square&label=Release)](https://github.com/Lumiaqian/mygo-ui-sdk/releases)

Go 即时模式行为原语 · 零额外运行时依赖 · 公开 API 零侵入 · `ui.Tester` Agent 体检 · 随仓库分发 Skill

[核心功能](#-核心功能一览) · [API 速查](#-api-速查) · [快速开始](#-快速开始) · [项目架构](#-项目架构) · [本地开发与 CI](#-本地开发与持续集成)

</div>

---

`mygo-ui-sdk` 借鉴 `vercel-labs/ui` 的行为分层与 Agent-ready 校验理念，面向 `github.com/egoist/mygo/ui` 的即时模式（immediate-mode）原生界面设计。项目不替代 `mygo/ui` 的渲染与内置 `*Base` 控件，也不依赖任何 `mygo/internal/*` 内部包，而是补齐自定义视觉控件（如 Liquid Glass、浮动标注工具栏）的交互行为封装、无障碍与键盘自动化检验，以及面向 AI 编程代理的 `mygo-ui` Skill。

## ✨ 核心功能一览

### 1. 可组合交互行为原语（`behaviors`）
- **一维拖拽与吸附（`TrackDrag` / `SwitchTrack`）**：基于 `ui.Local` 封装拖拽阈值判定、指针跟随、边界钳制与松手帧调度（`c.After`），搭配 `ui.SwitchBase` 保留原生开关语义与键盘操作。
- **分段控制行为（`SegmentedTrack` / `SegmentedGroup`）**：`SegmentedTrack` 面向连续滑块/液态玻璃透镜式分段控件，内置拖拽吸附、点击选段、`Left`/`Right`（自动适配 RTL）、`Home`/`End` 与 `Label`/`Value` 语义；`SegmentedGroup` 封装 `ui.SegmentedBase`，提供预命名离散分段。
- **无样式浮动工具栏（`Toolbar` / `ToolbarParts`）**：基于 `RoleToolbar` + `FocusGroup` 提供单次 Tab 进入与方向键漫游，遵循 MyGo `*Parts` 惯例先样式化 `tb.Root` 再构建子项，内置 `Button`、`Toggle`（`RoleToggleButton` + `Checked`）与 `RadioItem`（`RoleRadio` + `Checked`，适合色块/线宽选择器）。
- **侧栏/顶部导航列表（`NavList` / `NavListParts`）**：基于 `RoleTabList` + `RoleTab` 封装垂直/水平（含 RTL 镜像）方向键与 `Home`/`End` 切换，在文本输入聚焦时自动让出方向键，支持不回绕夹紧或循环切换，适合自定义 Liquid Glass 侧栏导航。
- **自定义单选组（`RadioGroup` / `RadioGroupParts`）**：基于 `RoleRadioGroup` + `FocusGroup` + `RoleRadio` 封装色块、主题圆点与模板卡片等任意自定义单选组，保留完整键盘漫游与选中项优先聚焦。
- **视口安全定位与轻关闭（`PlaceFloating` / `ShouldDismiss`）**：支持首选方向溢出翻转、锚点内回退、边缘钳制与 RTL 对齐镜像，以及统一的 `Escape` / 外部点击关闭判定。

### 2. 无障碍与 Agent 可达性测试（`uitest`）
- **一键体检（`Audit` / `AssertAgentReady`）**：扫描 `ui.Tester.AccessAll()`，自动检出未命名交互/可聚焦控件（`FindingUnnamedControl`）、重名交互控件冲突（`FindingDuplicateName`）、缺失必选控件（`FindingMissingRequired`）、非交互容器伪装点击（`FindingNonInteractiveRole`，配合 `RequiredInteractive`）以及不可键盘聚焦的交互控件（`FindingNotFocusable`）。
- **主题、方向与键盘契约（`AssertThemeParity` / `AssertDirectionParity` / `AssertKeyboardTabOrder` / `AssertRovingFocus`）**：在无窗口 `go test` 中覆盖浅深色、LTR/RTL、Tab 顺序与 `FocusGroup` 方向键漫游。
- **结构化快照（`Snapshot` / `FindByRole`）**：导出紧凑的无障碍节点列表用于断言或排障。

### 3. 面向 AI Agent 的 `mygo-ui` Skill（`skills/mygo-ui`）
- **渐进式索引与规则**：通过 `skills.sh.json`、`skills/mygo-ui/SKILL.md` 与 `llms.txt` 分页索引，固化 `*Base` 选型、`Context.Key`/`ui.Local` 生命周期、`Label`/`Value` 规范与 `uitest` 必测闭环。

---

## 📖 API 速查

<details open>
<summary><strong>点击展开 / 收起全部对照</strong></summary>

### `behaviors` 行为原语
| 函数 / 方法 | 用途 |
|---|---|
| `behaviors.TrackDrag(c, el, key, restPos, disabled, opts)` | 一维指针拖拽状态机（阈值、钳制、帧调度） |
| `behaviors.SwitchTrack(c, &on, opts)` | 封装 `ui.SwitchBase` 的可拖拽开关 |
| `behaviors.SegmentedTrack(c, &selected, opts)` | 连续拖拽透镜式分段控件（含方向键 / RTL / Home / End） |
| `behaviors.SegmentedGroup(c, &selected, opts)` | 封装 `ui.SegmentedBase` 的离散自定义分段控件 |
| `behaviors.Toolbar(c, opts)` | 无样式 `RoleToolbar` + `FocusGroup` 漫游焦点容器（返回 `ToolbarParts`） |
| `tb.Button` / `tb.Toggle` / `tb.RadioItem` | 工具栏内的可聚焦按钮、开关按钮与单选项 |
| `behaviors.NavList(c, &selected, count, opts)` | 无样式 `RoleTabList` + `RoleTab` 侧栏/标签导航（支持方向键 / RTL / Home / End / 输入态避让） |
| `behaviors.RadioGroup(c, opts)` / `rg.Item(...)` | 无样式 `RoleRadioGroup` + `FocusGroup` + `RoleRadio` 自定义单选组（色块/卡片） |
| `behaviors.PlaceFloating(viewport, anchor, size, opts)` | 浮动面板定位、溢出翻转、内回退与 RTL 镜像 |
| `behaviors.ShouldDismiss(c, surface, opts)` | `Escape` 与 `PressedOutside` 轻关闭判定 |

### `uitest` 验收工具
| 函数 | 用途 |
|---|---|
| `uitest.Audit(tester, opts)` | 返回无障碍与 Agent 可达性问题列表 `[]Finding` |
| `uitest.AssertAgentReady(t, tester, opts)` | 发现任何问题即报错并输出结构化摘要 |
| `uitest.AssertThemeParity(t, tester, opts)` | 浅色 / 深色主题双跑体检 |
| `uitest.AssertDirectionParity(t, tester, opts)` | LTR / RTL 双跑体检并校验交互标签集合一致 |
| `uitest.AssertKeyboardTabOrder(t, tester, names)` | 校验顺序 `Tab` 焦点落点 |
| `uitest.AssertRovingFocus(t, tester, orient, names)` | 校验 `FocusGroup` 方向键、循环与 `Home`/`End` |
| `uitest.Snapshot(tester)` / `uitest.FindByRole(tester, role)` | 读取当前帧无障碍树快照或按角色过滤节点 |

</details>

---

## 🚀 快速开始

1. 在你的 MyGo 项目中安装依赖：
   ```sh
   go get github.com/Lumiaqian/mygo-ui-sdk@v0.1.1
   ```
2. 在自定义控件中组合行为，并在 `go test` 里执行 Agent 可达性检查（完整示例见 [`examples/demo/main_test.go`](examples/demo/main_test.go)）：
   ```go
   tb := behaviors.Toolbar(c.Key("tools"), behaviors.ToolbarOptions{Label: "Tools"})
   tb.Root.Gap(4).Children(func() {
       if tb.Toggle(c.Key("rect"), "Rect", active == "Rect", behaviors.ToolbarItemOptions{}).Size(32, 28).Clicked() {
           active = "Rect"
       }
   })
   ```
   ```go
   tester := ui.NewTester(view, 640, 480)
   uitest.AssertAgentReady(t, tester, uitest.AuditOptions{
       RequiredLabels:   []string{"Tools", "Rect"},
       RequireFocusable: true,
   })
   ```
3. 将 [`skills/mygo-ui`](skills/mygo-ui) 提供给编程 Agent：
   ```sh
   npx skills add Lumiaqian/mygo-ui-sdk --skill mygo-ui
   ```

---

## 🏗 项目架构

核心模块职责与目录结构如下：

| 模块 / 文件 | 职责 |
|---|---|
| [`behaviors`](behaviors) | 面向 `mygo/ui` 即时模式的可组合交互行为原语与单元测试 |
| [`uitest`](uitest) | 基于公开 `ui.Tester` 的无障碍、键盘、方向、主题与 Agent 可达性断言 |
| [`skills/mygo-ui`](skills/mygo-ui) | 面向 AI Agent 编写 MyGo Native UI 的 `SKILL.md`、`llms.txt` 与参考文档 |
| [`skills.sh.json`](skills.sh.json) | 声明 Skill 发现路径（`./skills`） |
| [`examples/demo`](examples/demo) | 组合 `behaviors` 与 `uitest` 的最小可测试示例 |

---

## 🛠 本地开发与持续集成

### 本地常用命令
```sh
# 1. 整理依赖
go mod tidy

# 2. 运行全部包与示例测试（含竞态检测）
go test -race -count=1 ./...

# 3. 静态检查
go vet ./...
```

### GitHub Actions 流水线
- **持续集成（`.github/workflows/ci.yml`）**：在推送到 `main` 或发起针对 `main` 的 Pull Request 时，于 `ubuntu-latest` 与 `macos-latest` 矩阵运行 `go vet ./...` 与 `go test -race -count=1 ./...`。
