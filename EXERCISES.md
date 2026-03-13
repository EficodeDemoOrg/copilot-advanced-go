# Advanced GitHub Copilot Exercises

These exercises are split across two workshops. The weather app codebase is the substrate — you will build tooling around it, not fix it.

- **Workshop 1** (~2 hours): Learn the building blocks — custom agents, skills, and hooks.
- **Workshop 2** (~1 hour): Assemble them into a multi-agent orchestration workflow.

**Approach:** These exercises are intentionally open-ended. Discuss design decisions with your group and use Copilot to explore options. There is no single correct solution.

**Prerequisites:** VS Code with GitHub Copilot, agent mode available, enterprise subscription.

**Two principles to carry through every exercise:**

1. **Human-in-the-loop.** Design your agents so they stop and report back at key decision points. The human reviews, comments, and approves before the agent continues. An orchestrator that runs end-to-end without checkpoints is a liability, not an asset.
2. **Brevity.** Agent definitions, skill instructions, and protocols that grow beyond ~150 lines start to lose focus. The LLM cannot reliably follow a wall of text. Keep instructions tight, specific, and structured.

**Design checklist** — review this before you consider any exercise done:

- [ ] Does the agent stop for human approval at key decision points?
- [ ] Could any LLM guesswork be replaced with a deterministic tool (script, skill, hook) to get facts into context?
- [ ] Are the instructions under ~150 lines and clearly structured?
- [ ] Is the tool list minimal — only what this agent actually needs?
- [ ] Does each subagent have a defined input/output contract?
- [ ] Have you tested the agent with a concrete task and observed its behavior?

**Getting help:** This repo includes an **Exercise Tutor** agent. Switch to it in the agent picker whenever you need guidance on concepts, design decisions, or debugging. Open separate chat threads for different topics to keep conversations focused — don't pile everything into one thread.

---

# Workshop 1: Building Blocks

---

## Exercise 0: Setup Verification

Before starting, confirm your environment works.

**macOS / Linux:**
```bash
make install
make test
```

**Windows (PowerShell):**
```powershell
go test ./internal/...
```

All tests should pass. You do not need a real OpenWeatherMap API key — the tests mock all external calls.

Verify that VS Code agent mode is functional: open the Chat view, check that the agent picker shows both built-in agents (Agent, Plan, Ask, Edit) and that you can switch between them.

Confirm the **Exercise Tutor** agent is available: open the agent picker and look for it in the list. Switch to it and ask a question to verify it responds. This agent is your workshop coach — use it throughout the exercises when you need guidance on concepts, design decisions, or debugging.

Create the directories you will use:

**macOS / Linux:**
```bash
mkdir -p .github/agents .github/skills .github/hooks
```

**Windows (PowerShell):**
```powershell
New-Item -ItemType Directory -Force .github/agents, .github/skills, .github/hooks
```

---

## Exercise 1: Custom Agent — Project Manager

**Goal:** Create a custom agent that acts as a project manager for this codebase. It should be able to assess the project, create structured backlog items, and plan features.

**What to build:** An `.agent.md` file in `.github/agents/`.

A custom agent is a Markdown file with YAML frontmatter that defines the agent's name, description, available [tools](https://code.visualstudio.com/docs/copilot/agents/agent-tools), and behavioral instructions. Key frontmatter properties: `description`, `tools`, `model`, `agents`, `handoffs`. See the [file structure reference](https://code.visualstudio.com/docs/copilot/customization/custom-agents#_custom-agent-file-structure).

**Key decisions to make:**
- What tools does the PM need? It mostly analyzes, but it may also need to write files (plans, backlog documents). Think about where the boundaries are — what should it be allowed to touch and what should it not?
- What instructions make the PM produce consistent, structured output?
- What should a backlog item look like? (Think: title, description, acceptance criteria, TDD requirements, definition of done.) For this codebase, TDD means `testing` + `testify`, with tests placed alongside source files in `internal/<package>/`.
- What model would you choose for this agent? The PM needs to reason about architecture but doesn't need to generate code. Consult the [model comparison reference](https://docs.github.com/en/copilot/reference/ai-models/model-comparison) and think about cost, capability, and context window trade-offs.
- How would you test that the agent behaves consistently?

**Things to try once the agent exists:**
- Switch to the PM agent and ask it to assess the project and identify improvement areas.
- Ask it to plan a feature it identifies as a good candidate.
- Ask it to review one of its own backlog items for completeness and proper sizing.
- If you're unsure about any design decision, switch to the **Exercise Tutor** agent in a separate thread and ask.

**Discussion points:**
- How specific should the instructions be vs. how much should you rely on the model's judgment?
- The codebase has a strict layered architecture (`handlers` → `services` → `repository`). Should the PM agent understand this? How would you encode that knowledge?

**References:**
- [Custom agents in VS Code](https://code.visualstudio.com/docs/copilot/customization/custom-agents)
- [Creating custom agents (GitHub)](https://docs.github.com/en/copilot/how-tos/use-copilot-agents/coding-agent/create-custom-agents)
- [Model comparison](https://docs.github.com/en/copilot/reference/ai-models/model-comparison)

---

## Exercise 2: Agent Skills for the PM

**Goal:** Create skills that give the PM agent (and later, other agents) repeatable, deterministic capabilities. Skills are not just instructions — they include scripts and resources.

**Key concept:** A skill is a directory under `.github/skills/<skill-name>/` with a `SKILL.md` file and optional scripts/resources. The agent loads the skill when it judges the skill relevant, and follows the instructions, which can reference the included scripts. See [SKILL.md file format](https://code.visualstudio.com/docs/copilot/customization/agent-skills#_skillmd-file-format) and [how Copilot uses skills](https://code.visualstudio.com/docs/copilot/customization/agent-skills#_how-copilot-uses-skills).

The difference between a skill and a custom instruction: skills include scripts, examples, and resources alongside instructions. They are loaded on-demand based on relevance, not always-on. They are portable across VS Code, Copilot CLI, and the coding agent. See [skills vs. custom instructions](https://code.visualstudio.com/docs/copilot/customization/agent-skills#_agent-skills-vs-custom-instructions).

**Getting started:** Type `/skills` in the chat input to open the Configure Skills menu, where you can create a new skill. You can also describe the skill you want to Copilot and ask it to scaffold the directory and `SKILL.md` for you.

**Task:** Think about what the PM agent does repeatedly, and where deterministic scripts would produce more reliable results than the LLM guessing. Discuss with your group and with Copilot. Build at least two skills.

What kind of deterministic capabilities would help a PM agent working with this Go codebase? Consider:
- Running `go test ./internal/... -v` and parsing test counts and failures
- Running `golangci-lint run` and collecting lint issues
- Counting source files per layer (`internal/handlers/`, `internal/services/`, `internal/repository/`, etc.)
- Mapping `*_test.go` files to their source pairs to identify untested packages
- Listing module dependencies from `go.mod` and checking for known vulnerabilities with `go list -m all`

The key idea is that these are things a shell script can do reliably and the LLM should not be guessing at.

**For each skill:**
- Create the directory structure: `.github/skills/<name>/SKILL.md` plus any scripts.
- The `SKILL.md` must have `name` and `description` in the YAML frontmatter. The `name` must match the directory name.
- Keep skills focused and concise — a skill that tries to do too much loses effectiveness.
- Write scripts that produce structured, parseable output. The agent interprets the output; the script ensures deterministic data collection.
- Test the skill by invoking it as a slash command (`/<skill-name>`) or by prompting the PM agent in a scenario where the skill should activate.
- Control visibility with `user-invocable` and `disable-model-invocation` frontmatter properties. See [slash command control](https://code.visualstudio.com/docs/copilot/customization/agent-skills#_use-skills-as-slash-commands).
- To nudge the agent toward using a skill, reference it explicitly in the agent's instructions (e.g., "When assessing the project, use the `<skill-name>` skill to gather metrics").

**Discussion points:**
- What makes a good boundary between "instruction for the agent" and "script that runs deterministically"? (Ask the **Exercise Tutor** if you want to think this through.)
- How do you ensure the agent actually uses the skill vs. trying to do it in its own way?
- Which of these skills would be useful beyond the PM agent?

**References:**
- [Agent skills in VS Code](https://code.visualstudio.com/docs/copilot/customization/agent-skills)
- [Creating agent skills (GitHub)](https://docs.github.com/en/copilot/how-tos/use-copilot-agents/coding-agent/create-skills)

---

## Exercise 3: Hooks

**Goal:** Add lifecycle hooks that enforce guarantees that instructions alone cannot provide.

**Key concept:** Hooks are shell commands that run deterministically at specific points in the agent lifecycle. An instruction that says "always run the linter" is a suggestion the LLM may ignore. A hook that runs the linter is a guarantee. See [hook lifecycle events](https://code.visualstudio.com/docs/copilot/customization/hooks#_hook-lifecycle-events).

Create `.github/hooks/` configuration files. Start with one or two hooks and observe their behavior.

**Starting point:** Think about the PM agent and skills you built in Exercises 1–2. Where did the agent skip steps or produce inconsistent results? Ask Copilot (or the **Exercise Tutor**) in a separate thread to review your current agent definition and suggest where hooks would add value.

Think about the full lifecycle of an agent session and where you want guarantees that instructions alone cannot provide. The available [hook events](https://code.visualstudio.com/docs/copilot/customization/hooks#_hook-lifecycle-events) are: `SessionStart`, `UserPromptSubmit`, `PreToolUse`, `PostToolUse`, `PreCompact`, `SubagentStart`, `SubagentStop`, `Stop`.

Hook configuration is JSON in `.github/hooks/`. See [hook configuration format](https://code.visualstudio.com/docs/copilot/customization/hooks#_hook-configuration-format) and [hook input/output](https://code.visualstudio.com/docs/copilot/customization/hooks#_hook-input-and-output) for details on how hooks communicate with the agent.

**Key points:**
- Hooks run deterministically — they are shell commands, not LLM suggestions.
- A `Stop` hook that blocks must check `stop_hook_active` to prevent infinite loops.
- [Agent-scoped hooks](https://code.visualstudio.com/docs/copilot/customization/hooks#_agent-scoped-hooks) (defined in the agent's frontmatter `hooks` field) only run when that agent is active. Requires `chat.useCustomAgentHooks` to be enabled.
- `PreToolUse` hooks can [control tool approval](https://code.visualstudio.com/docs/copilot/customization/hooks#_pretooluse-output): allow, deny, or ask.

**Things to try:**
- A `PostToolUse` hook that runs `golangci-lint run` after any `.go` file is edited.
- A `PreToolUse` hook that requires confirmation before any terminal command is executed.
- An agent-scoped hook on the PM agent that logs every assessment to a file under `tmp/`.
- A `Stop` hook that runs `go test ./internal/...` and blocks the agent from finishing if tests are failing.

**Discussion:**
- Where is the line between "the agent should decide" and "the system must enforce"?
- What are the risks of hooks that block agent operations?
- How do hooks, skills, and instructions each fit into the control spectrum? (Instructions = guidance, skills = on-demand capabilities, hooks = guarantees.)

**References:**
- [Hooks in VS Code](https://code.visualstudio.com/docs/copilot/customization/hooks)
- [Using hooks (GitHub)](https://docs.github.com/en/copilot/how-tos/use-copilot-agents/coding-agent/use-hooks)

---

## Exercise 4: Feature Implementation with Your Tools

**Goal:** Put the building blocks together. Use the PM agent to plan a feature and then implement it using Agent mode, with your skills and hooks active.

This exercise is the Workshop 1 payoff. You'll experience the full cycle manually — which will make the "why orchestrate?" question concrete when you reach Workshop 2.

**Steps:**
1. Switch to the PM agent. Ask it to assess the project and propose a feature. Observe whether your skills fire and whether the hooks behave as expected.
2. Pick a feature from the PM's output. Ask the PM to produce a detailed plan with acceptance criteria and TDD requirements.
3. Switch to the built-in **Agent** mode (not your PM). Give it the plan and ask it to implement the feature using TDD — tests first, then implementation.
4. Review the result. Run the tests to verify:
   - **macOS / Linux:** `make test`
   - **Windows:** `go test ./internal/...`

**Things to observe:**
- Did the PM's skills activate and produce useful data, or did it ignore them?
- Did hooks enforce the guarantees you set up? (Linting, confirmation gates, etc.)
- Did the Agent mode follow the layered architecture — handler delegates to service, service returns domain errors, handler calls `apperrors.HandleError`?
- How much manual coordination did you have to do between the PM step and the implementation step? (Copy-pasting plans, switching chat threads, re-explaining context...)
- **That manual coordination is exactly what orchestration automates in Workshop 2.**

**Iterate:** If skills didn't fire or hooks misbehaved, go back and fix them. Use the **Exercise Tutor** in a separate thread for debugging.

---

## Exercise 4S: MCP Server Integration (Stretch)

**Goal:** Extend your agents with external tools via Model Context Protocol servers.

MCP servers expose external tools to agents via a standardized protocol. Browse the [GitHub MCP registry](https://github.com/mcp) to see what's available.

**Discussion prompts:**
- The PM agent creates backlog items as text. What if it could create GitHub Issues directly? (The GitHub MCP server can do this.)
- What if test results from `go test -json ./internal/...` were available as a structured MCP tool with parsed output?
- What would a "Go project metrics" MCP server look like — one that exposes test coverage (`go test -cover`), lint results (`golangci-lint run --out-format json`), and module audit results as structured tools?

**If your organization's policies allow:** Configure an existing MCP server from the [GitHub MCP registry](https://github.com/mcp) in one of your agents using the `mcp-servers` property in the agent frontmatter, and try using it.

---

# Workshop 2: Orchestration

You have the building blocks from Workshop 1: a custom agent, skills, and hooks. Now assemble them into a multi-agent workflow.

**Time estimate:** Exercise 5 fits in ~1 hour. Exercises 6–7 are stretch goals.

---

## Exercise 5: Subagent Orchestration — Feature Implementation Workflow

This is the main exercise. You will design and implement a multi-agent workflow that can take a feature from plan to implementation using TDD. Use the PM agent from Workshop 1 to identify a feature, then build the orchestration to implement it.

Remember Exercise 4? The manual coordination you did there — switching agents, copy-pasting plans, re-explaining context — is what orchestration automates.

### 5a: Design the Orchestration

Before writing any agent files, design the workflow as a group. Use the PM agent from Workshop 1 to identify and plan a feature to implement. The **Exercise Tutor** agent can help you think through orchestration patterns — open a separate thread for that conversation.

Subagents are independent agents that perform focused work and report results back to a main agent. Each runs in its own context window. The main agent waits for results before continuing. Multiple subagents can run in parallel. See [how subagent execution works](https://code.visualstudio.com/docs/copilot/agents/subagents#_how-subagent-execution-works).

**What roles does the workflow need?** Think hard about separation of concerns. The goal is not just "who does what" but also "what does each agent receive as input and what exactly does it return." Define the input/output contract for each role explicitly.

Consider:
- Which roles are read-only vs. which need to edit files?
- Which roles can run in parallel?
- The TDD cycle for this codebase: someone writes failing tests in `internal/<package>/<file>_test.go` following the `TestSubject_Condition` naming convention; someone else makes them pass while respecting the layered architecture.
- **Where does the coordinator stop and wait for human approval?** At minimum, after research and before committing to an implementation plan.

A critical design goal: **the coordinator should not read files itself.** It delegates all research (the heavy context ingestion) to researcher subagents that return concise, structured summaries. For example, a researcher given a task to find relevant code should return which files and which line ranges contain the relevant information — not the file contents. This keeps the coordinator's context window lean and focused on decision-making.

Use the [coordinator and worker pattern](https://code.visualstudio.com/docs/copilot/agents/subagents#_coordinator-and-worker-pattern) from the docs as a starting point, but adapt it to your needs.

**Design questions to resolve:**
- What information does the coordinator pass to each subagent? What does it expect back? Be specific about the format.
- How do you instruct researchers to return structured, concise summaries instead of dumping file contents?
- What happens when new code conflicts with existing tests or violates the layer rules (e.g., a handler that contains business logic)? Define the escalation protocol.
- Do all roles need to be separate agents, or could some be combined?

**Model selection for each role:** Not every agent needs the most expensive model. Think about what each role actually does:
- Does it need to reason about architecture and make judgment calls? Strong reasoning model.
- Does it need to ingest a lot of code? Needs a large context window.
- Does it follow straightforward instructions without high abstraction? A mid-tier model works.
- Does it just collect and format information? Cheapest/fastest model.

Consult the [model comparison reference](https://docs.github.com/en/copilot/reference/ai-models/model-comparison) when choosing.

### 5b: Implement the Agents

Create the agent files in `.github/agents/` based on your design.

Key technical details:
- Agents that should only be invoked as subagents: set [`user-invocable: false`](https://code.visualstudio.com/docs/copilot/agents/subagents#_control-subagent-invocation).
- The coordinator should restrict which subagents it can use via the [`agents` property](https://code.visualstudio.com/docs/copilot/agents/subagents#_restrict-which-subagents-can-be-used-experimental).
- Each agent defines its own `tools` list. Read-only agents should not have edit/terminal tools.
- Worker agents can specify a `model` — consider cheaper/faster models for narrow tasks.
- Wire in the skills and hooks you built in Workshop 1 where they apply.

The coordinator's instructions should define the workflow sequence explicitly: what it delegates first, what it does with the results, what triggers the next step, where it stops for human review, and what happens on failure.

For the implementer agent, make the architectural rules explicit: handlers call services, services return domain errors from `internal/apperrors/`, handlers call `apperrors.HandleError()`. New endpoints need Swagger annotations regenerated via `make swagger` (or `swag init -g cmd/server/main.go --output docs` on Windows).

### 5c: Debug and Validate

Before running a full feature through the workflow, do a dry run. Give the coordinator a small, well-defined task and observe.

**Debugging agent behavior:** Open the agent debug panel (right-click in the Chat view and select "Diagnostics") to understand why agents make the decisions they do. Key things to investigate:
- Why did the coordinator choose to call (or skip) a particular subagent?
- Did a skill get loaded? If not, why not — was the description too vague?
- What context did the subagent actually receive?

**Practical tip:** Use separate chat threads for different concerns — one for running the agent, one for tweaking definitions, one for asking the **Exercise Tutor** about debugging strategies. Changes to `.agent.md` files take effect in new threads, not the currently running one.

### 5d: Run a Feature Through It

Take the feature you planned in 5a, switch to the coordinator agent, and give it the feature spec. Observe the full workflow.

**Things to watch for:**
- Does the coordinator actually delegate, or does it try to do everything itself?
- Does the coordinator stop at the designated human checkpoints?
- Are the research summaries concise enough, or do they bloat the coordinator's context?
- Does the TDD cycle work — tests written first, then implementation?
- Does the implementer iterate until `go test ./internal/...` passes?
- Does the final code pass `golangci-lint run`?

**Iterate on the agent instructions based on what you observe.** This is the real work — the first version will not be perfect.

**References:**
- [Subagents in VS Code](https://code.visualstudio.com/docs/copilot/agents/subagents)
- [Orchestration patterns](https://code.visualstudio.com/docs/copilot/agents/subagents#_orchestration-patterns)

---

## Exercise 6: Refactor the PM Agent (Stretch)

**Goal:** Upgrade the PM agent to use the researcher subagents from Exercise 5 instead of doing its own analysis.

Update the PM agent:
- Add your researcher agents to the `agents` property, and `agent` to the tools list.
- Update the instructions to delegate research to subagents and synthesize the results.
- Optionally add a [handoff](https://code.visualstudio.com/docs/copilot/customization/custom-agents#_handoffs) from the PM to the coordinator: after planning a feature, the user can hand off to the coordinator for implementation.

Test by asking the PM to assess the project. It should delegate the research and produce a more focused assessment.

**Discussion:** How does the PM's output quality change when it gets structured research summaries of the Go codebase vs. doing its own ad-hoc file reading?

---

## Exercise 7: Skills Across the Workflow (Stretch)

**Goal:** Identify repeatable operations in the subagent workflow and extract them into skills.

Look at the agents you built in Exercise 5. What do they do repeatedly? Where would a deterministic script produce better results than the LLM improvising?

**Starting point:** If you're unsure what skills to build, open a new chat thread and ask Copilot (or the **Exercise Tutor**) to assess your current agent definitions and propose skills that would make them more reliable. Use that as a starting point, then refine. Remember that skills should be reusable across agents and deterministic — the value is in scripts that produce reliable, structured output, not in more prose instructions.

For Go-specific ideas:
- A skill that runs `go test ./internal/... -json` and returns a structured pass/fail summary per package.
- A skill that checks whether all exported symbols have docstrings (using `go doc` or a simple script).
- A skill that verifies the layered architecture is intact — e.g., no direct HTTP imports in `internal/services/`.

For each skill you identify, create the directory, `SKILL.md`, and any scripts. Then update the relevant agent instructions to reference the skill explicitly (e.g., "Always use the `<skill-name>` skill after modifying source files").

**Discussion:** Skills are loaded based on relevance, not guaranteed to run. How do you increase the likelihood that agents use them? What's the difference between a skill and a hook for enforcing behavior?
