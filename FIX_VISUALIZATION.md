# Visual Explanation of Task Stacking Fix

## Before Fix (Overlapping Problem)

```
┌─────────────────────────────────────────────────────────┐
│ Week 45 - November 2025                                 │
├──────┬──────┬──────┬──────┬──────┬──────┬──────────────┤
│ Mon  │ Tue  │ Wed  │ Thu  │ Fri  │ Sat  │ Sun          │
│  10  │  11  │  12  │  13  │  14  │  15  │  16          │
├──────┼──────┼──────┼──────┼──────┼──────┼──────────────┤
│      │ ┌────┴──────┴──────┴──────┴──────┴────┐         │
│      │ │ Write Methods & Timeline            │         │
│      │ └────────────────────────────────────┬┘         │
│      │      ┌───────────────────────────────┴────────┐ │
│      │      │ Send Proposal ← OVERLAPS!              │ │  ❌
│      │      └────────────────────────────────────────┘ │
└──────┴──────┴──────┴──────┴──────┴──────┴──────────────┘
```

**Problem**: "Send Proposal" draws over "Write Methods" because it doesn't know about the earlier task.

## After Fix (Proper Stacking)

```
┌─────────────────────────────────────────────────────────┐
│ Week 45 - November 2025                                 │
├──────┬──────┬──────┬──────┬──────┬──────┬──────────────┤
│ Mon  │ Tue  │ Wed  │ Thu  │ Fri  │ Sat  │ Sun          │
│  10  │  11  │  12  │  13  │  14  │  15  │  16          │
├──────┼──────┼──────┼──────┼──────┼──────┼──────────────┤
│      │ ┌────┴──────┴──────┴──────┴──────┴────┐         │
│      │ │ Write Methods & Timeline (bottom)   │         │
│      │ └──┬──────────────────────────────────┘         │
│      │    ┌┴──────────────────────────────────────────┐│
│      │    │ Send Proposal (stacked on top)            ││  ✅
│      │    └───────────────────────────────────────────┘│
└──────┴──────┴──────┴──────┴──────┴──────┴──────────────┘
```

**Solution**: "Send Proposal" knows "Write Methods" is active and positions itself above it.

## How It Works Internally

### Day-by-Day Processing

#### November 11 (Tuesday)
```
findActiveTasks(Nov 11):
  ├─ Write Methods (starts Nov 11) ← STARTS TODAY
  └─ Active tasks: [Write Methods]

renderSpanningTaskOverlay(Nov 11):
  ├─ Write Methods: stackPos=0 (bottom)
  └─ Use: \TaskOverlayBox (with \vfill)
```

#### November 12 (Wednesday)
```
findActiveTasks(Nov 12):
  ├─ Write Methods (continuing from Nov 11)
  ├─ Send Proposal (starts Nov 12) ← STARTS TODAY
  └─ Active tasks: [Write Methods, Send Proposal]
                    (sorted by start date)

renderSpanningTaskOverlay(Nov 12):
  ├─ Write Methods: stackPos=0 (continuing, DON'T RENDER)
  ├─ Send Proposal: stackPos=1 (stacked)
  └─ Use: \TaskOverlayBoxNoOffset (no \vfill, stacks above)
```

## LaTeX Macro Behavior

### TaskOverlayBox (Bottom Task)
```latex
\vfill                    ← Pushes to bottom
\begin{tcolorbox}
  Task content
\end{tcolorbox}
```

### TaskOverlayBoxNoOffset (Stacked Task)
```latex
% No \vfill              ← Flows naturally
\begin{tcolorbox}
  Task content
\end{tcolorbox}
```

## Key Insight

The fix recognizes that **vertical stacking requires knowing about ALL active tasks**, not just the ones that start on a specific day. By maintaining awareness of continuing tasks, each new task can calculate its proper position in the vertical stack.

## Code Flow

```
User opens November 2025 page
    ↓
For each day (Nov 11, 12, 13, ...):
    ↓
    findActiveTasks(day)
        ├─ Find tasks starting today
        ├─ Find tasks continuing from earlier
        └─ Return sorted list (by start date)
    ↓
    renderSpanningTaskOverlay(day)
        ├─ Only render tasks starting TODAY
        ├─ But use full active list for stackPos
        └─ Choose macro based on stackPos:
            ├─ stackPos=0 → TaskOverlayBox
            └─ stackPos>0 → TaskOverlayBoxNoOffset
    ↓
PDF rendered with proper stacking
```

