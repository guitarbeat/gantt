# First Contribution Checklist

**Purpose:** Checklist for making your first contribution. Delete after your first PR is merged.

---

## ✅ Before You Start

- [ ] Read [Developer Guide](../tasks/DEVELOPER_GUIDE.md)
- [ ] Complete [Initial Setup](INITIAL_SETUP_CHECKLIST.md)
- [ ] Join GitHub Discussions (if available)
- [ ] Review open issues for "good first issue" label

---

## ✅ Choose Your Contribution

Pick one:
- [ ] Fix a bug
- [ ] Add a feature
- [ ] Improve documentation
- [ ] Add tests
- [ ] Refactor code

**Issue #:** _____________  
**Description:** _____________

---

## ✅ Development Process

- [ ] Create a new branch
  ```bash
  git checkout -b feature/my-feature
  # or
  git checkout -b fix/bug-description
  ```

- [ ] Make your changes
  - [ ] Write code
  - [ ] Add/update tests
  - [ ] Update documentation
  - [ ] Run tests locally

- [ ] Follow code style
  ```bash
  # Format code
  gofmt -w .
  # or
  make fmt
  ```

- [ ] Run linters
  ```bash
  golangci-lint run ./...
  # or
  make lint
  ```

- [ ] Run all tests
  ```bash
  go test ./...
  # or
  make test
  ```

---

## ✅ Commit Your Changes

- [ ] Stage your changes
  ```bash
  git add .
  ```

- [ ] Write a good commit message
  ```bash
  git commit -m "feat: add new feature"
  # or
  git commit -m "fix: resolve bug in X"
  # or
  git commit -m "docs: update user guide"
  ```

- [ ] Follow commit message convention
  - `feat:` - New feature
  - `fix:` - Bug fix
  - `docs:` - Documentation
  - `test:` - Tests
  - `refactor:` - Code refactoring
  - `chore:` - Maintenance

---

## ✅ Push and Create PR

- [ ] Push your branch
  ```bash
  git push origin feature/my-feature
  ```

- [ ] Create Pull Request on GitHub
  - [ ] Use descriptive title
  - [ ] Fill out PR template
  - [ ] Link related issue
  - [ ] Add screenshots (if UI changes)
  - [ ] List breaking changes (if any)

---

## ✅ PR Checklist

- [ ] All tests pass
- [ ] Code is formatted
- [ ] Documentation updated
- [ ] No merge conflicts
- [ ] Commit messages follow convention
- [ ] PR description is clear

---

## ✅ After PR Submission

- [ ] Respond to review comments
- [ ] Make requested changes
- [ ] Push updates
  ```bash
  git add .
  git commit -m "fix: address review comments"
  git push origin feature/my-feature
  ```

- [ ] Wait for approval
- [ ] Celebrate when merged! 🎉

---

## ✅ After First PR is Merged

- [ ] Update your local main branch
  ```bash
  git checkout main
  git pull origin main
  ```

- [ ] Delete your feature branch
  ```bash
  git branch -d feature/my-feature
  git push origin --delete feature/my-feature
  ```

- [ ] **Delete this file** - You're now a contributor!

---

## 💡 Tips for Success

- Start small - don't try to change everything at once
- Ask questions in PR comments or discussions
- Be patient - reviews take time
- Learn from feedback
- Have fun!

---

## 📚 Resources

- [Developer Guide](../tasks/DEVELOPER_GUIDE.md)
- [Troubleshooting](../tasks/TROUBLESHOOTING.md)
- [GitHub Flow Guide](https://guides.github.com/introduction/flow/)
- [Writing Good Commit Messages](https://chris.beams.io/posts/git-commit/)

---

**First Contribution Date:** _____________  
**PR Number:** _____________  
**What I Learned:** _____________
