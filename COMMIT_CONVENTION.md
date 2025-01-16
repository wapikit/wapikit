# Commit and Branch Conventions

At `wapikit`, we follow a strict set of conventions for commit messages and branch names to ensure that our repository stays organized, and our commit history remains crystal clear. Here's a guide on our conventions:

---

## Conventional Commits

Conventional commits ensure our commit messages are clear and useful.

**Benefits**:

- Automatic changelog generation
- Simplified versioning
- Improved readability of commit messages

### Commit Message Format

Each commit message should follow this format:

```
<type>(<scope>): <short summary>
<BLANK LINE>
<optional body>
<BLANK LINE>
<optional footer>
```

- **`<type>`**: Describes the purpose of the commit:
  - `feat`: Introduces a new feature
  - `fix`: Addresses a bug
  - `docs`: Updates documentation
  - `style`: Code that relates to styling, not affecting logic
  - `refactor`: Refactoring existing code
  - `perf`: Improving performance
  - `test`: All about tests
  - `chore`: Maintenance tasks
- **`<scope>`**: (Optional) Specifies which part of the codebase is affected.

- **`<short summary>`**: A concise summary of the changes made.

### Examples:

1. Introducing a new feature:

```
feat(auth): implement social login
```

2. Addressing a bug:

```
fix(button): resolve alignment issue
```

---

## Conventional Branching

A standardized naming system for branches helps everyone quickly understand a branch's purpose.

### Branch Naming Format

Branch names should adhere to:

```
<type>/<short-description>
```

- **`<type>`**: Purpose of the branch, common ones being:

  - `feature`: Developing a new feature.
  - `fix`: Addressing a bug.
  - `chore`: Regular maintenance tasks.
  - `hotfix`: Immediate fixes, often tied to production issues.
  - `docs`: Documentation enhancements.

- **`<short-description>`**: A brief, kebab-cased (words separated by hyphens) description of the branch's objective.

### Examples:

1. Developing a new user dashboard:

```
feature/user-dashboard
```

2. Resolving a login issue:

```
fix/login-issue
```

---

🙌 Thanks for contributing to `wapikit`! By adhering to these conventions, we're making our repository a better place. If you're new, welcome aboard, and if you've been here, thanks for sticking around!

---
