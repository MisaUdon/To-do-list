# Todo List 项目架构与技术栈

## 1. 项目简介

这是一个使用 Go、SQLite 和原生 JavaScript 开发的简易待办事项项目，支持任务的创建、查询、更新、完成和删除。

整体结构如下：

```text
浏览器页面
    ↓ HTTP / JSON
Go Web 服务
    ↓ GORM
SQLite 数据库
```

## 2. 目录结构

```text
To-do-list/
├── main.go                 # 后端入口、路由和接口处理
├── go.mod                  # Go 模块及依赖声明
├── go.sum                  # 依赖版本校验信息
├── test.db                 # SQLite 数据库文件
├── db/
│   └── model.go            # Todo 模型及数据库增删改查
├── static/
│   ├── index.html          # 页面结构
│   ├── style.css           # 页面样式
│   └── script.js           # 页面交互及后端接口调用
└── PROJECT_ARCHITECTURE.md # 项目架构说明
```

## 3. 技术栈

### 3.1 后端

- **Go**：项目的主要开发语言。
- **net/http**：Go 标准库，用于启动 HTTP 服务、注册路由和处理请求。
- **encoding/json**：解析前端发送的 JSON，并将查询结果编码为 JSON。
- **GORM**：Go ORM 框架，用于通过结构体操作数据库。
- **google/uuid**：为每条 Todo 生成唯一 ID。

### 3.2 数据库

- **SQLite**：轻量级文件数据库，不需要单独启动数据库服务。
- 数据保存在项目根目录的 `test.db` 中。
- GORM 的 `AutoMigrate` 根据 `Todo` 结构体自动创建或更新数据表。

### 3.3 前端

- **HTML5**：定义页面结构和表单。
- **CSS3**：负责页面布局、按钮及任务卡片样式。
- **原生 JavaScript**：处理页面事件和 DOM 更新。
- **Fetch API**：向 Go 后端发送 HTTP 请求。

项目没有使用 Vue、React 等前端框架，适合学习基础的前后端交互过程。

## 4. 后端架构

### 4.1 服务入口

`main.go` 主要负责：

1. 调用 `db.Init()` 初始化数据库。
2. 注册创建、查询、更新和删除路由。
3. 配置跨域请求响应头。
4. 提供 `static` 目录中的网页资源。
5. 在 `8080` 端口启动 HTTP 服务。

### 4.2 API 路由

| 路径 | 用途 | 主要请求数据 |
| --- | --- | --- |
| `/getall` | 查询全部任务 | 无 |
| `/create` | 创建任务 | `name`、`description` |
| `/update` | 更新任务 | `id`、`name`、`description`、`completed` |
| `/delete` | 删除任务 | `id` |
| `/` | 加载前端静态网页 | 无 |

目前前端主要通过 `POST` 请求完成创建、更新和删除操作。

### 4.3 请求处理器

`main.go` 中的 Handler 是 HTTP 层，负责：

1. 接收浏览器请求。
2. 解析请求体中的 JSON。
3. 调用 `db` 包执行数据库操作。
4. 返回 HTTP 状态码或 JSON 数据。

例如，创建任务的请求数据为：

```json
{
  "name": "学习 Go",
  "description": "完成 HTTP 服务章节"
}
```

### 4.4 数据模型

`db/model.go` 中定义了 Todo 数据模型：

```go
type Todo struct {
    ID          string `json:"id"`
    Name        string `json:"name"`
    Description string `json:"description"`
    Completed   bool   `json:"completed"`
}
```

各字段含义：

| 字段 | 类型 | 作用 |
| --- | --- | --- |
| `ID` | `string` | Todo 的唯一标识，由 UUID 生成 |
| `Name` | `string` | 任务名称 |
| `Description` | `string` | 任务描述 |
| `Completed` | `bool` | 任务是否完成 |

### 4.5 数据库操作

`db/model.go` 提供以下函数：

- `Init()`：连接数据库并自动迁移数据表。
- `CreateTodo()`：创建任务。
- `GetAllTodos()`：查询全部任务。
- `UpdateTodo()`：更新任务信息和完成状态。
- `DeleteTodo()`：根据 ID 删除任务。

## 5. 前端架构

### 5.1 HTML

`static/index.html` 包含：

- 项目标题。
- 任务名称和任务描述输入框。
- 创建按钮。
- 用于展示任务的无序列表。

### 5.2 JavaScript

`static/script.js` 主要负责：

- 页面加载时查询全部任务。
- 将任务数据转换为页面元素。
- 创建新任务。
- 修改任务名称和描述。
- 切换任务完成状态。
- 删除任务。
- 每次操作后重新获取最新任务列表。

页面加载过程为：

```text
打开 index.html
    ↓
触发 window.onload
    ↓
调用 GET /getall
    ↓
后端查询 SQLite
    ↓
返回 JSON 数组
    ↓
JavaScript 生成任务列表
```

### 5.3 CSS

`static/style.css` 负责：

- 页面居中布局。
- 表单和输入框样式。
- Todo 卡片布局。
- 创建、更新、完成和删除按钮样式。
- 已完成任务的删除线效果。

## 6. 一次完整请求的执行过程

以创建任务为例：

1. 用户在网页中填写任务名称和描述。
2. JavaScript 阻止表单默认刷新行为。
3. `fetch` 将任务数据以 JSON 形式发送到 `/create`。
4. Go Handler 使用 `json.Decoder` 解析请求体。
5. 后端使用 UUID 生成任务 ID。
6. `db.CreateTodo()` 通过 GORM 写入 SQLite。
7. 后端返回成功状态码。
8. 前端重新调用 `/getall` 并刷新任务列表。

## 7. 核心概念

- **路由（Route）**：把 URL 映射到对应的 Go 处理函数。
- **Handler**：处理一次 HTTP 请求并生成响应。
- **JSON**：前后端交换数据的格式。
- **ORM**：使用 Go 结构体操作数据库，减少手写 SQL。
- **CRUD**：Create、Read、Update、Delete，即增删改查。
- **CORS**：控制网页是否可以跨来源访问后端接口。
- **静态文件服务**：由 Go 返回 HTML、CSS 和 JavaScript 文件。

## 8. 当前架构可改进点

1. 数据库连接可以在初始化时创建一次并复用，避免每次请求重新连接。
2. Handler 中不应使用 `log.Fatal`，否则一次请求出错可能导致整个服务退出。
3. 每个接口应限制允许的 HTTP 方法。
4. 应增加任务名称、ID 和请求参数校验。
5. 前端可使用 `/getall` 等相对地址，避免写死 `localhost:8080`。
6. 可以将项目进一步拆分为 Handler、Service、Repository 三层。
7. 可以增加统一的 JSON 错误响应和自动化测试。

## 9. 启动方式

在项目根目录运行：

```bash
go run main.go
```

然后在浏览器访问：

```text
http://localhost:8080
```

如果出现 `address already in use`，说明 `8080` 端口已经被其他进程占用，需要先停止旧的服务进程。
