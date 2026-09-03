---
name: git-ops
description: >-
  Use this skill when staging files, crafting commit messages, verifying git status,
  and managing repository hygiene.
---

# Git Operations & Repository Hygiene

This skill enforces version control conventions and repository isolation.

## 1. Commit Message Convention (Conventional Commits)
Commit messages must be written in English following Conventional Commits format:
`<type>(<optional scope>): <description>`

### Types
- `feat`: A new feature
- `fix`: A bug fix
- `docs`: Documentation only changes
- `style`: Changes that do not affect the meaning of code (formatting, whitespace)
- `refactor`: A code change that neither fixes a bug nor adds a feature
- `test`: Adding missing tests or correcting existing tests
- `chore`: Build process, dependencies, tooling, or agent configuration updates

### Rules
- Write summary in imperative mood (e.g., `feat: add division by zero validation`).
- Keep the subject line under 72 characters.
- Do not commit generated build artifacts, dependencies (`node_modules`, `.venv`), or secrets.

## 2. Repository Scope Restriction
- All operations, creations, and edits are strictly confined to this repository (`vitoresende/calculator`).
- Never touch or modify files outside the repository root.
