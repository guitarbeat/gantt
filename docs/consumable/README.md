# Consumable Tasks

**Purpose:** This folder contains one-time checklists and tasks that should be **deleted after completion**.

---

## 📋 What Goes Here?

Consumable tasks are:
- ✅ One-time setup checklists
- ✅ Onboarding tasks
- ✅ Migration guides (after migration is complete)
- ✅ Temporary project-specific tasks
- ✅ "Do this once and delete" instructions

---

## 📁 Current Consumable Tasks

### For New Users
- **[Initial Setup Checklist](INITIAL_SETUP_CHECKLIST.md)** - First-time setup
  - Delete after: Initial setup is complete

### For New Contributors
- **[First Contribution Checklist](FIRST_CONTRIBUTION_CHECKLIST.md)** - Making your first PR
  - Delete after: First PR is merged

---

## 🗑️ When to Delete

Delete a consumable task file when:
1. All checklist items are completed
2. The task is no longer relevant
3. The migration/setup is done
4. You've successfully completed the one-time action

---

## 📝 How to Use

1. **Open** the relevant checklist
2. **Check off** items as you complete them
3. **Delete** the file when all items are done
4. **Commit** the deletion to keep the repo clean

Example:
```bash
# After completing initial setup
git rm docs/consumable/INITIAL_SETUP_CHECKLIST.md
git commit -m "chore: remove completed initial setup checklist"
git push
```

---

## ✨ Creating New Consumable Tasks

If you need to create a new consumable task:

1. Create a markdown file in this folder
2. Use a clear, descriptive name (e.g., `MIGRATION_TO_V2_CHECKLIST.md`)
3. Include:
   - Purpose statement
   - Checklist items
   - Deletion instructions
   - Completion date field
4. Add it to this README
5. Delete it when done!

---

## 🔄 Difference from Other Docs

| Folder | Purpose | Lifespan |
|--------|---------|----------|
| **consumable/** | One-time tasks | Delete after completion |
| **tasks/** | How-to guides | Keep forever (reference) |
| **fyi/** | Background info | Keep forever (historical) |

---

## 🎯 Keep This Folder Clean

- Don't let completed checklists accumulate
- Delete files promptly after completion
- This keeps the repo focused and relevant
- Future contributors will thank you!

---

**Remember:** If you're done with it, delete it! 🗑️
