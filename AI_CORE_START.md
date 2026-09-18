# AI Core Prototype — Start Document

## 1. Цель

Построить быстрое, понятное и расширяемое AI-ядро на Go, которое станет общей основой для разных кабинетов и продуктов.

Главный принцип:

> **Кабинеты — это композиции возможностей единого ядра, а не отдельные системы.**

Ядро должно позволять подключать новые продукты без переписывания Core.

---

## 2. Базовый цикл ядра

Первая версия строится вокруг замкнутого цикла:

```text
Project
↓
Hypothesis
↓
Test
↓
Evidence
↓
Decision
↓
Next Loop
↺
```

На первом этапе нам нужен только этот работающий цикл.

---

## 3. Архитектурный подход

Используем:

```text
Modular Monolith
+ clear module boundaries
+ ports only at real external boundaries
+ plugins for product-specific logic
+ cabinets as compositions
```

Не используем full DDD и тяжёлую Hexagonal Architecture с первого дня.

Не создаём абстракции ради абстракций.

Правило:

```text
REUSE → EXTEND → ADAPTER → NEW
```

---

## 4. Core не должен знать о конкретных продуктах

В Core живут только универсальные сущности:

```text
Project
Goal
Loop
Hypothesis
Evidence
Decision
Artifact
Risk
Capability
Workflow
Memory
Context
```

В Core не должно быть:

```text
SoloFounderDashboard
DevOpsDashboard
TestingDashboard
CRMPage
```

Продукты подключаются через plugins.

---

## 5. Целевая архитектура

```text
PRODUCT / CABINETS
        │
plugins / workflows
        │
┌───────▼────────┐
│  ORCHESTRATOR  │
└───────┬────────┘
        │
┌───────┼─────────────────────────────┐
│       │                             │
Context Engine                 Capability Registry
│                                     │
Memory Core                        PredictGuard
│                                     │
Representation Layer            Evaluation
│
├── Symbolic
├── Vector / Tensor
└── Graph
        │
     Ports
        │
Adapters
```

---

## 6. Представления и память

Нужно заложить места для нескольких видов представления знаний.

### Symbolic

Будущие сущности:

```text
Fact
Entity
Relation
Rule
Constraint
Claim
```

Пример:

```text
subject: project-123
predicate: requires
object: rollback-plan
```

---

### Vector / Tensor

Core не должен зависеть от конкретной embedding-модели или GPU-runtime.

На уровне Core будут только структуры и порты:

```text
Vector
TensorRef
EmbeddingPort
TensorPort
```

Позже можно подключить:

```text
OpenAI embeddings
pgvector
Qdrant
Milvus
local embedding model
```

без переписывания Core.

---

### Graph

Минимальные понятия:

```text
Node
Edge
GraphStore
```

Пример:

```text
Project
├── HAS_GOAL → Goal
├── HAS_LOOP → Loop
├── PRODUCED → Artifact
├── HAS_RISK → Risk
└── DECIDED → Decision
```

Первая реализация может быть InMemoryGraph.

Graph DB подключается позже через adapter.

---

## 7. Типы памяти

Memory Core должен в будущем поддерживать четыре типа памяти.

### Working Memory

Текущее состояние:

```text
current goal
current loop
active task
current context
```

### Episodic Memory

Что происходило:

```text
prompts
decisions
test results
failed hypotheses
deployments
incidents
```

### Semantic Memory

Что система знает:

```text
facts
entities
relationships
architecture decisions
business knowledge
```

### Procedural Memory

Как система должна действовать:

```text
methodologies
prompt patterns
workflows
checklists
capabilities
```

---

## 8. Гибридная память

Один Memory Record в будущем сможет объединять:

```text
TEXT
+
SYMBOLIC FACTS
+
VECTOR REPRESENTATION
+
GRAPH CONNECTIONS
```

Не нужно выбирать только один подход.

---

## 9. Context Engine

LLM не должен каждый раз получать всю историю чата.

Context Engine должен собирать компактный Context Pack:

```text
Current Project State
+
Relevant Decisions
+
Relevant Evidence
+
Relevant Memories
+
Symbolic Constraints
+
Vector Similarity
+
Graph Neighbors
+
Available Capabilities
```

Результат:

```json
{
  "goal": "prepare MVP",
  "stage": "validation",
  "facts": [],
  "decisions": [],
  "relevant_memories": [],
  "risks": [],
  "capabilities": [],
  "constraints": []
}
```

---

## 10. Capability Registry

AI не должен иметь неограниченные скрытые действия.

Каждое действие — capability с контрактом.

Примеры:

```text
project.read
project.update
research.web
repository.inspect
document.create
architecture.review
testing.plan
loadtest.design
deploy.review
```

Будущий контракт:

```go
type Capability interface {
    Name() string
    Execute(ctx context.Context, input Input) (Output, error)
}
```

---

## 11. PredictGuard

PredictGuard — policy layer ядра.

Любое потенциально опасное действие проходит через:

```text
Proposed Action
↓
Risk Analysis
↓
ALLOW
WARN
REQUIRE APPROVAL
BLOCK
```

Пример:

```text
read repository       → ALLOW
create markdown       → ALLOW
modify code           → WARN
migration             → APPROVAL
production deploy     → APPROVAL
delete database       → BLOCK
expose secret         → BLOCK
```

---

## 12. Cabinets

Кабинет не является отдельной системой.

Он состоит из:

```text
manifest
+
workflow
+
capabilities
+
projections
```

Пример:

```text
AI Core
├── SoloStep
├── LoadLab
├── Architect
├── DevOps
└── Support
```

Все используют одно ядро.

Архитектурный инвариант:

> **Plugins extend Core; they never fork Core.**

---

## 13. Целевая структура Go-проекта

Не создавать все папки заранее. Папка появляется только вместе с реальным кодом.

Целевая карта:

```text
ai-core-prototype/
│
├── cmd/
│   └── api/
│
├── internal/
│
│   ├── core/
│   │   ├── project/
│   │   ├── loop/
│   │   ├── evidence/
│   │   ├── decision/
│   │   ├── representation/
│   │   │   ├── symbolic/
│   │   │   ├── tensor/
│   │   │   └── graph/
│   │   ├── memory/
│   │   ├── context/
│   │   ├── capability/
│   │   ├── workflow/
│   │   ├── orchestration/
│   │   ├── guard/
│   │   └── evaluation/
│   │
│   ├── ports/
│   │   ├── llm/
│   │   ├── embedding/
│   │   ├── memory/
│   │   ├── graph/
│   │   └── storage/
│   │
│   ├── adapters/
│   │   ├── memory/
│   │   ├── graph/
│   │   ├── llm/
│   │   └── http/
│   │
│   ├── plugins/
│   │   └── solostep/
│   │
│   └── cabinets/
│
├── examples/
├── docs/
└── go.mod
```

---

## 14. Порядок разработки через Codex

Строим ядро кольцами.

### Loop 1 — Core Skeleton

Реализуем:

```text
Project
Loop
Evidence
Decision
Engine
InMemory Store
```

Цель:

> получить один реально работающий замкнутый цикл.

---

### Loop 2 — Symbolic

Добавляем:

```text
Fact
Entity
Relation
Constraint
```

---

### Loop 3 — Memory

Добавляем:

```text
Working
Episodic
Semantic
Procedural
```

Первая реализация — in-memory.

---

### Loop 4 — Graph

Добавляем:

```text
Node
Edge
GraphStore
InMemoryGraph
```

---

### Loop 5 — Vector / Tensor

Добавляем только контракты:

```text
Vector
TensorRef
EmbeddingPort
MockEmbeddingAdapter
```

Без vector DB и тяжёлого runtime.

---

### Loop 6 — Context Engine

Собираем:

```text
symbols
+
memory
+
graph
+
vectors
+
project state
```

в единый:

```text
ContextPack
```

---

### Loop 7 — LLM Port

Добавляем:

```text
LLM Port
↓
OpenAI Adapter
```

Core не зависит напрямую от OpenAI.

---

### Loop 8 — Orchestrator

Цикл:

```text
Observe
→ Build Context
→ Plan
→ Select Capability
→ Execute
→ Collect Evidence
→ Evaluate
→ Decide
```

---

### Loop 9 — PredictGuard

Добавляем:

```text
ALLOW
WARN
APPROVAL
BLOCK
```

---

### Loop 10 — Первый Plugin

Подключаем SoloStep как первый реальный продукт поверх ядра.

---

## 15. Что делаем сейчас

Создаём локальную папку:

```powershell
cd C:\Users\User\Desktop
mkdir ai-core-prototype
cd ai-core-prototype
code .
```

Инициализируем Go:

```powershell
go mod init example.com/aicore
```

GitHub пока не нужен.

Название продукта выберем позже.

---

## 16. Первый промт для Codex

```text
F-001 — AI Core foundation.

GOAL

Build the smallest working foundation for a generic project loop engine in Go.

The engine must support:

Project
Loop
Evidence
Decision

Lifecycle:

Project
→ Hypothesis
→ Test
→ Evidence
→ Decision
→ Next Loop

ARCHITECTURE

Use a modular monolith.

Do not implement full DDD or full hexagonal architecture.

Create abstractions only for real boundaries.

Core must not depend on future products or cabinets.

DO NOT IMPLEMENT YET

- LLM
- OpenAI
- embeddings
- vector database
- tensors
- graph database
- symbolic reasoning
- persistence database
- UI
- auth
- deployment
- plugins

However, do not make architectural decisions that would prevent these layers from being added later.

FIRST INSPECT

Show the proposed minimal package structure before implementation.

Keep it small.

TARGET

A working Go program with:

- Project
- Loop
- Evidence
- Decision
- Engine
- in-memory storage
- tests proving one complete loop

TEST

go test ./...

REPORT

Explain:
1. package structure
2. responsibilities
3. dependency direction
4. exact files
5. tests
6. where memory/symbolic/vector/graph layers will connect later
```

---

## 17. Правило для первого шага

Codex сначала показывает минимальную структуру.

Мы проверяем её.

Только после этого разрешаем реализацию.

Не строим всё ядро одним запросом.

Главный принцип разработки:

```text
Loop 1
→ proof
→ Loop 2
→ proof
→ Loop 3
→ proof
```

Каждый новый слой появляется только после того, как предыдущий цикл доказан работающим.
