# GitHub Configuration Index

This directory contains AI agent configurations and documentation for the showheaders project.

## Contents

### Core Documentation
- [copilot-instructions.md](copilot-instructions.md) - Project conventions and patterns for GitHub Copilot

### Skills & Agents
- [skills/](skills/) - Domain expertise and knowledge bases
  - [DevOps Expert](skills/devops.md) - Complete DevOps lifecycle guidance
- [agents/](agents/) - Executable agent configurations
  - [DevOps Expert Agent](agents/devops-expert.md) - `@devops-expert` - Operational DevOps agent

## Quick Start

### For AI Coding Agents
1. Read [copilot-instructions.md](copilot-instructions.md) first for project-specific patterns
2. Reference [skills/](skills/) for domain knowledge
3. Invoke [agents/](agents/) for specialized, actionable guidance

### For Developers
- Check [copilot-instructions.md](copilot-instructions.md) for architecture and workflows
- Use agents by mentioning them (e.g., `@devops-expert`) when you need specialized help

## Agent Invocation

**DevOps Expert**: `@devops-expert`
- CI/CD pipeline setup and optimization
- Infrastructure as Code
- Monitoring and observability
- Incident response
- Performance optimization
- Deployment strategies

## Adding New Agents

1. Create skill documentation in `skills/`
2. Create agent configuration in `agents/`
3. Define clear invocation pattern
4. Update this index

---

See individual files for detailed documentation.
