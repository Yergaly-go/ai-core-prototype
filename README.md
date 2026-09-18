# AI Core Prototype

Experimental modular AI core written in Go.

The goal of this project is to build a reusable foundation for AI-powered products, agents and domain-specific cabinets without rewriting the core for every new product.

## Why this project exists

Most AI applications start from a chat interface and gradually accumulate:

- duplicated logic;
- duplicated state;
- hidden business rules inside prompts;
- multiple sources of truth;
- tightly coupled LLM calls;
- weak memory;
- difficult-to-debug agent behavior.

This project takes a different approach.

The AI system is built as a set of explicit layers with clear responsibilities.

The core stays generic.

Products, vertical domains and user cabinets are connected later through plugins and workflows.

---

## Core idea

The main execution loop is:

```text
Project
→ Hypothesis
→ Test
→ Evidence
→ Decision
→ Next Loop
```

Each completed loop changes the state of the project and creates the next step.

The system is designed around repeated validation cycles rather than one large autonomous AI action.

---

## Architecture vision

```text
                    Products / Cabinets
                           │
                    Plugins / Workflows
                           │
                    ┌──────▼──────┐
                    │ Orchestrator │
                    └──────┬──────┘
                           │
         ┌─────────────────┼─────────────────┐
         │                 │                 │
   Context Engine     Capabilities      PredictGuard
         │
   ┌─────▼─────────────────────────────────────┐
   │                 Memory Core               │
   ├───────────────────────────────────────────┤
   │ Working                                   │
   │ Episodic                                  │
   │ Semantic                                  │
   │ Procedural                                │
   └─────┬─────────────────────────────────────┘
         │
   ┌─────▼─────────────────────────────────────┐
   │          Representation Layer             │
   ├──────────────┬──────────────┬─────────────┤
   │ Symbolic     │ Graph        │ Vector /    │
   │              │              │ Tensor      │
   └──────────────┴──────────────┴─────────────┘
         │
       Ports
         │
      Adapters
```

The important architectural rule is:

> Canonical data and derived representations are separated.

---

## Current implementation

The repository already contains working Go code for three independent layers.

### 1. Deterministic Core

The deterministic core manages project state and validation loops.

```text
Project
→ Loop
→ Evidence
→ Decision
→ Next Loop
```

Responsibilities:

- create projects;
- open validation loops;
- collect evidence;
- close loops through decisions;
- create the next loop;
- protect lifecycle invariants.

The Engine does not call an LLM.

This keeps project state deterministic and testable.

---

### 2. Symbolic Knowledge Layer

The symbolic layer represents explicit knowledge.

Current concepts:

```text
Entity
Fact
Relation
Constraint
```

Semantics:

```text
Fact
Entity → Predicate → Literal Value

Relation
Entity → Predicate → Entity

Constraint
MUST / MUST_NOT
```

Example:

```text
Project
HAS_TARGET_USER
SoloFounder
```

or:

```text
Deployment
REQUIRES
RollbackPlan
```

Symbolic knowledge is project-scoped and append-only.

It is intended to become the canonical semantic layer of the system.

---

### 3. Episodic Memory

The memory layer currently stores an append-only timeline of references to canonical records.

```text
Episode
├── ProjectID
├── Target RecordRef
├── OccurredAt
└── CreatedAt
```

The memory layer does not copy Evidence, Decisions or Symbolic knowledge.

Instead it remembers:

```text
what happened
when it happened
what canonical record represents it
```

This prevents Memory from becoming a second source of truth.

---

## Planned memory model

The long-term memory architecture separates four different concepts.

### Working Memory

Temporary current execution state.

Examples:

```text
current goal
current loop
active task
temporary context
```

Working Memory should not automatically become permanent storage.

### Episodic Memory

What happened over time.

Examples:

```text
evidence collected
decision made
test completed
deployment happened
incident happened
handoff created
```

### Semantic Memory

What the system knows.

Canonical semantic knowledge is represented through:

```text
Entity
Fact
Relation
Constraint
```

### Procedural Memory

How the system knows how to act.

Future examples:

```text
methodologies
workflows
prompt patterns
capabilities
checklists
```

---

## Future Graph Layer

The Graph Layer will be a derived representation of canonical knowledge.

Example:

```text
Project
├── HAS_LOOP → Loop
├── HAS_EPISODE → Episode
├── HAS_RISK → Risk
└── HAS_DECISION → Decision
```

Symbolic relations naturally map to graph edges:

```text
Entity → Node
Relation → Edge
```

The Graph will not become a second source of truth.

It should always be rebuildable from canonical records.

---

## Future Vector / Tensor Layer

Vector and tensor representations will be used for retrieval and multidimensional similarity.

They will not own canonical data.

Future responsibilities may include:

```text
semantic similarity
context retrieval
multidimensional relation strength
ranking
relevance scoring
```

The expected principle is:

```text
Graph
= structure of relationships

Vector / Tensor
= strength and multidimensional characteristics of relationships
```

---

## Context Engine

The future Context Engine will assemble only the context required for the current task.

Instead of sending an entire conversation to an LLM:

```text
Current Project State
+
Relevant Episodes
+
Relevant Symbolic Knowledge
+
Graph Neighborhood
+
Vector Retrieval
+
Constraints
+
Available Capabilities
```

will become:

```text
Context Pack
```

This is intended to reduce prompt noise and make AI behavior easier to explain and reproduce.

---

## Hypothesis-driven reasoning

The project uses cyclic hypothesis validation.

```text
Current State
→ Hypothesis
→ Test
→ Evidence
→ Evaluation
→ Decision
→ Transformation
→ Next State
```

The next loop starts from the result of the previous one.

Therefore the system moves not only in circles, but as a spiral toward a target.

```text
S0
→ S1
→ S2
→ S3
→ Target
```

A future reasoning layer should also remember failed paths so the system does not repeat the same unsuccessful route indefinitely.

---

## Intelligent Orchestrator

The future Orchestrator will sit above the deterministic Engine.

Its role will be:

```text
Observe
→ Retrieve context
→ Form hypothesis
→ Plan
→ Select capability
→ Execute
→ Collect evidence
→ Evaluate
→ Decide next action
```

The deterministic Engine will remain responsible for valid state transitions.

The Orchestrator will not directly own project state.

---

## PredictGuard

PredictGuard is planned as a policy and safety layer.

Before potentially risky actions:

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

Examples:

```text
read repository       → ALLOW
create document       → ALLOW
modify code           → WARN
database migration    → APPROVAL
production deploy     → APPROVAL
delete database       → BLOCK
expose secret         → BLOCK
```

---

## Capabilities

AI actions will eventually be exposed through explicit capabilities.

Examples:

```text
repository.inspect
architecture.review
research.web
testing.plan
loadtest.design
document.create
deployment.review
```

Each capability will have a contract:

```text
name
input
output
permissions
risk level
approval policy
```

The AI should not have hidden unlimited access to the system.

---

## Plugins and Cabinets

The core is designed so multiple products can use the same AI foundation.

```text
AI Core
├── SoloStep
├── Architecture Assistant
├── Load Testing Lab
├── DevOps Assistant
└── Support Assistant
```

A cabinet is not a separate system.

> Cabinets are compositions, not systems.

A future cabinet will combine:

```text
workflow
+
capabilities
+
projections
+
navigation
```

while reusing the same Core.

---

## Design principles

The project follows several engineering principles.

### Deterministic core

LLMs do not own canonical lifecycle state.

### Single source of truth

Derived layers do not duplicate canonical data.

### Append-only knowledge where possible

Evidence, decisions, symbolic knowledge and episodes preserve history.

### Project isolation

Knowledge and memory belong to a specific project.

### Derived representations

Graph and vector/tensor layers remain rebuildable indexes.

### Small verified loops

The architecture is developed in short cycles:

```text
Recon
→ Architecture Decision
→ Implementation
→ Tests
→ Review
→ Commit
→ Next Loop
```

---

## Current repository structure

```text
internal/
├── core/
├── engine/
├── symbolic/
├── memory/
└── store/
    └── memory/
```

The repository intentionally stays small.

New packages are created only when a real responsibility or external boundary appears.

---

## Development roadmap

Completed:

```text
F-001 — deterministic project loop core
F-002 — symbolic knowledge layer
F-003 — episodic memory timeline
```

Next layers:

```text
F-004 — Graph representation
F-005 — Vector / Tensor contracts
F-006 — Context Engine
F-007 — Hypothesis Engine
F-008 — Path / Failure Memory
F-009 — Intelligent Orchestrator
F-010 — LLM Provider Port
F-011 — Retrieval / RAG
F-012 — PredictGuard
F-013 — Plugin / Cabinet system
```

---

## Tests

Run:

```bash
go test ./...
```

Current layers are covered by automated tests.

---

## Status

Active experimental prototype.

The current goal is not to build a large AI framework at once.

The goal is to validate each architectural layer independently and add complexity only when the previous layer is proven.
