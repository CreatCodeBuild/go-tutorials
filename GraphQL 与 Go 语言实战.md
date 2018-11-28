## 《GraphQL 与 Go 语言实战》
这本书会专注于分享作者在工作中对 GraphQL 技术结合 Go 语言的使用。本书着重于 GraphQL 技术的使用以及 API 设计，不会讲解如何实现 GraphQL 语言的解析器。本书也不会着重于 Go 语言的讲解。会用到基本的 Go 语言语法，但是很简单，即使没有使用过 Go 的读者也可以轻松理解。

读者范围：希望学习 GraphQL 的职业软件工程师或者计算机科学、软件工程在读学生。有 2 年以上编程经验。

阅读顺序：
- 对于经验较丰富的读者，可以从自己喜欢的章节开始。每一章的知识点是相关的，但是没有必然的前后关系。
- 对于初学者，从第一章开始。

### GraphQL
### 1. 简介
  - 1.1 简介  
        讲述 GraphQL 的由来
  - 1.2 基本语法
  - 1.3 例子  
        在深入讲解之前，我们会简单地看看 GitHub 的 GraphQL API （虽然 GraphQL 是 facebook 开发的，但是对于国内开发者来说，GitHub 更相关）

### 2. Query 语法 （第二章会用 GitHub GraphQL API 举例）
  - 2.1 Source Text
  - 2.2 Document
  - 2.3 Operations
  - 2.4 Selection Sets
  - 2.5 Fields
  - 2.6 Arguments
  - 2.7 Field Alias
  - 2.8 Fragments
  - 2.9 Input Values
  - 2.10 Variables
  - 2.11 Type References
  - 2.12 Directives

### 3. 类型系统 与 Sechema 语法
  - 3.1 Type System Extensions
  - 3.2 Schema
  - 3.3 Descriptions
  - 3.4 Types
  - 3.5 Scalars
  - 3.6 Objects
  - 3.7 Interfaces
  - 3.8 Unions
  - 3.9 Enums
  - 3.10 Input Objects
  - 3.11 List
  - 3.12 Non-Null
  - 3.13 Directives

### 4. 关于设计的讨论 Public API vs Private API. Data API vs BFF. 以及它们和 GraphQL 有什么关系。
  - 4.1. What is a public API? 什么是公开API？
  - 4.2. What is a private API? 什么是私密API？
  - 4.3. What is a data API? 什么是数据API？
  - 4.4. What is a Backend For Frontend? 什么是前端的后端？  

### 5. Go 语言加油站
这一章会讲述 G0 语言的基本使用方法，为之后的代码实战做准备。（可以推荐其他专门讲 Go 语言的书）
- 1. 基本语法
- 2. Interface
- 3. 并行
- 4. 写 Web App 会用到的标准库 (net/http)

### 6. 具体类库介绍  
Go 语言有 3 个 GraphQL 类库，我们来看看每一个的优劣。
- 1. https://github.com/graphql-go/graphql
- 2. https://github.com/graph-gophers/graphql-go
- 3. https://github.com/99designs/gqlgen

### 7. 实战代码  
这一章会展示一个完整的 Web 后端程序。

### 8. 思考总结
- 1. Go 是一个合适的 GraphQL 实现语言吗？
- 2. NodeJS or Go？
- 3. 生态  
这一章会简单谈谈 GraphQL 的开源社区以及生态系统。会讲到非 Go 语言的解决方案。比如 Apollo。但是不会深入。
