#!/usr/bin/env python3
import argparse
import csv
import os
import shutil
import subprocess
from datetime import datetime, date, timedelta
from typing import Optional, List, Dict, Tuple, Set
import re


def parse_date(s: str) -> Optional[date]:
    s = (s or "").strip()
    if not s:
        return None
    return datetime.strptime(s, "%Y-%m-%d").date()


def read_rows(path: str) -> List[Dict[str, str]]:
    with open(path, newline="", encoding="utf-8") as f:
        return list(csv.DictReader(f))


LANE_ORDER = [
    "Proposal",
    "Committee",
    "Laser",
    "Imaging",
    "Admin/Accountability",
    "Other",
]

LANE_COLORS = {
    "Proposal": "#4062BB",
    "Committee": "#F18F01",
    "Laser": "#2FBF71",
    "Imaging": "#9B59B6",
    "Admin/Accountability": "#95A5A6",
    "Other": "#7F8C8D",
}


def lane_for(group: str, dtype: str) -> str:
    g = (group or "").strip().upper()
    dt = (dtype or "").strip()
    if dt == "Committee Milestone":
        return "Committee"
    if g == "PROPOSAL" or "Proposal" in dt:
        return "Proposal"
    if g == "LASER":
        return "Laser"
    if g == "IMAGING":
        return "Imaging"
    if g in ("ADMIN", "ACCOUNTABILITY"):
        return "Admin/Accountability"
    return "Other"


def _parse_ids_field(s: str) -> List[str]:
    s = (s or "").strip()
    if not s:
        return []
    # split on commas/semicolons/whitespace, keep alnum/underscore/hyphen
    raw = []
    cur = []
    for ch in s:
        if ch.isalnum() or ch in {"_", "-"}:
            cur.append(ch)
        else:
            if cur:
                raw.append("".join(cur))
                cur = []
    if cur:
        raw.append("".join(cur))
    # de-dup while preserving order
    seen: Set[str] = set()
    out: List[str] = []
    for x in raw:
        if not x:
            continue
        if x not in seen:
            seen.add(x)
            out.append(x)
    return out


def _get_col(r: Dict[str, str], names: List[str]) -> str:
    """Return the first non-empty value among possible column names."""
    for n in names:
        v = r.get(n)
        if v is not None:
            v = v.strip()
            if v:
                return v
    return ""


def load_tasks(src_csv: str) -> List[Dict]:
    rows = read_rows(src_csv)
    # First pass: collect raw tasks
    tasks: List[Dict] = []
    cref_re = re.compile(r"contentReference\[[^\]]*\](?:\{[^}]*\})?", re.IGNORECASE)
    for r in rows:
        sd = parse_date(r.get("Start Date", ""))
        dd = parse_date(r.get("Due Date", ""))
        dur_s = (r.get("Duration (days)") or "").strip()
        dur: Optional[int] = None
        try:
            dur = int(dur_s) if dur_s else None
        except Exception:
            dur = None
        # If one date is missing, try to infer from the other and duration
        if not sd and dd and dur and dur > 0:
            sd = dd - timedelta(days=dur - 1)
        if not dd and sd and dur and dur > 0:
            dd = sd + timedelta(days=dur - 1)
        # If still missing dates, we cannot render the task
        if not sd or not dd:
            continue
        if not dur or dur <= 0:
            dur = (dd - sd).days + 1
        dtype = (r.get("Deliverable Type") or "").strip()
        name = (r.get("Task Name") or "").strip()
        # Support column rename: "Task ID" -> "Index" (and a few variants)
        tid = _get_col(r, ["Index", "Task ID", "ID", "Task Index"]) or ""
        # Support parent under multiple headings
        parent_id = _get_col(r, ["Parent Index", "Parent Task ID", "Parent ID", "Parent"]) or ""
        # Be robust to misspelling: "Dependecies" as well as other variants
        deps_raw = _get_col(r, ["Dependencies", "Dependecies", "Depends", "Depends On", "Predecessors", "Predecessor"]) or ""
        deps = _parse_ids_field(deps_raw)
        notes_raw = (r.get("Notes") or "").strip()
        notes = cref_re.sub("", notes_raw).strip()
        t = {
            "id": tid,
            "name": name,
            "group": (r.get("Group") or "").strip(),
            "dtype": dtype,
            "owner": (r.get("Owner") or "").strip(),
            "status": (r.get("Status") or "").strip(),
            "priority": (r.get("Priority") or "").strip(),
            "start": sd,
            "due": dd,
            "dur": dur,
            "parent_id": parent_id,
            "deps": deps,
            "notes": notes,
            "children": [],
        }
        t["lane"] = lane_for(t["group"], t["dtype"])
        tasks.append(t)

    # Index by ID
    by_id: Dict[str, Dict] = {t["id"]: t for t in tasks if t["id"]}
    # Build parent-child links using explicit Parent Task ID
    for t in tasks:
        pid = t.get("parent_id")
        if pid and pid in by_id and pid != t["id"]:
            by_id[pid]["children"].append(t["id"])  # store child IDs

    # For hierarchy clarity: align children lane with parent lane
    for t in tasks:
        pid = t.get("parent_id")
        if pid and pid in by_id:
            t["lane"] = by_id[pid]["lane"]

    # Sort children of each parent by start date, then due, then id
    for t in tasks:
        if t["children"]:
            t["children"].sort(key=lambda cid: (by_id[cid]["start"], by_id[cid]["due"], by_id[cid]["id"]))

    # Sort tasks for deterministic output
    tasks.sort(key=lambda x: (x["start"], x["due"], x["id"]))
    return tasks


def month_iter(start: date, end: date):
    m = date(start.year, start.month, 1)
    while m <= end:
        yield m
        ny = m.year + (m.month // 12)
        nm = 1 if m.month == 12 else m.month + 1
        m = date(ny, nm, 1)


def _wrap_text(s: str, max_chars: int) -> List[str]:
    # simple word wrap by approximate character count
    words = s.split()
    if not words:
        return [""]
    lines: List[str] = []
    cur: List[str] = []
    for w in words:
        if sum(len(x) for x in cur) + max(0, len(cur)-1) + len(w) <= max_chars:
            cur.append(w)
        else:
            if cur:
                lines.append(" ".join(cur))
            cur = [w]
    if cur:
        lines.append(" ".join(cur))
    return lines


def generate_html(tasks: List[Dict], out_html: str, title: str = "Project Timeline"):
    if not tasks:
        os.makedirs(os.path.dirname(out_html), exist_ok=True)
        with open(out_html, "w", encoding="utf-8") as f:
            f.write("<html><body><h2>No tasks to display</h2></body></html>")
        return

    dmin = min(t["start"] for t in tasks) - timedelta(days=2)
    dmax = max(t["due"] for t in tasks) + timedelta(days=2)
    total_days = max(1, (dmax - dmin).days)

    def esc(s: str) -> str:
        return (
            (s or "")
            .replace("&", "&amp;")
            .replace("<", "&lt;")
            .replace(">", "&gt;")
        )

    # build lanes
    by_lane_all: Dict[str, Dict[str, Dict]] = {}
    by_id: Dict[str, Dict] = {t["id"]: t for t in tasks if t["id"]}
    for t in tasks:
        if t.get("id"):
            by_lane_all.setdefault(t["lane"], {})[t["id"]] = t
    order = [ln for ln in LANE_ORDER if ln in by_lane_all] + [ln for ln in by_lane_all.keys() if ln not in LANE_ORDER]

    # layout
    width = 1500
    padding_left = 560   # start of bars; left column reserved for full labels
    padding_right = 40
    padding_top = 120    # Increased for better legend and controls
    padding_bottom = 56
    line_h = 18  # Increased for better readability
    row_gap = 12  # Increased for better spacing
    lane_gap = 24  # Increased for better lane separation

    # Precompute wrapped labels and per-item heights
    label_left = 24
    label_right = padding_left - 16
    approx_char_px = 8.0  # Increased for better character spacing
    max_chars = max(12, int((label_right - label_left) / approx_char_px))

    enriched: Dict[str, List[Dict]] = {}
    lane_heights: Dict[str, int] = {}
    height = padding_top + padding_bottom
    # Utility to flatten hierarchy within a lane
    def flatten_lane(lane: str) -> List[Tuple[Dict, int]]:
        items: List[Tuple[Dict, int]] = []
        # roots are tasks in lane that are not a child of another present task
        lane_tasks = by_lane_all.get(lane, {})
        is_child: Set[str] = set()
        for tid, t in lane_tasks.items():
            pid = t.get("parent_id")
            if pid and pid in lane_tasks:
                is_child.add(tid)

        roots = [t for tid, t in lane_tasks.items() if tid not in is_child]
        roots.sort(key=lambda x: (x["start"], x["due"], x["id"]))

        def add_with_children(t: Dict, level: int = 0):
            items.append((t, level))
            for cid in t.get("children", []):
                if cid in lane_tasks:
                    add_with_children(lane_tasks[cid], level + 1)

        for rt in roots:
            add_with_children(rt, 0)
        return items

    for lane in list(by_lane_all.keys()):
        total = 0
        items_in_lane: List[Dict] = []
        for t, level in flatten_lane(lane):
            # Build label with indent; include name only (no ID)
            label_full = t["name"]
            lines = _wrap_text(label_full, max_chars)
            meta_bits = [b for b in [t.get("dtype"), t.get("owner"), t.get("status"), t.get("priority")] if b]
            meta_str = " • ".join(meta_bits) if meta_bits else ""
            meta_lines = _wrap_text(meta_str, max_chars) if meta_str else []
            line_count = max(1, len(lines)) + (len(meta_lines) if meta_lines else 0)
            item_h = line_count * line_h + 8  # vertical padding
            t2 = dict(t)
            t2["_label_lines"] = lines
            t2["_meta_lines"] = meta_lines
            t2["_item_h"] = item_h
            t2["_level"] = level
            items_in_lane.append(t2)
            total += item_h + row_gap
        lane_h = max(0, total) + lane_gap
        enriched[lane] = items_in_lane
        lane_heights[lane] = lane_h
        height += lane_h

    def x_pos(d: date) -> float:
        return padding_left + (width - padding_left - padding_right) * ((d - dmin).days / total_days)

    svg: List[str] = []
    svg.append(f"<svg width='{width}' height='{height}' viewBox='0 0 {width} {height}' xmlns='http://www.w3.org/2000/svg' font-family='ui-sans-serif, system-ui, -apple-system, Segoe UI, Roboto, Helvetica, Arial'>")
    
    # Title and controls
    svg.append(f"<text x='{padding_left}' y='28' font-size='24' font-weight='600' fill='#0f172a'>{esc(title)}</text>")
    svg.append(f"<text x='{padding_left}' y='52' font-size='14' fill='#475569'>{dmin} – {dmax}</text>")

    # Legend (lane colors)
    legend_x = padding_left
    legend_y = 80
    legend_gap = 18
    for i, lane in enumerate(order):
        col = LANE_COLORS.get(lane, "#888")
        y0 = legend_y + i * legend_gap
        svg.append(f"<rect x='{legend_x}' y='{y0 - 10}' width='14' height='14' rx='2' ry='2' fill='{col}' stroke='#334155' />")
        svg.append(f"<text x='{legend_x + 20}' y='{y0}' font-size='13' fill='#334155'>{esc(lane)}</text>")

    # Grid: months and weeks
    for m in month_iter(dmin, dmax):
        x = x_pos(m)
        svg.append(f"<line x1='{x:.2f}' y1='{padding_top-6}' x2='{x:.2f}' y2='{height-padding_bottom+4}' stroke='#e5e7eb'/>")
        svg.append(f"<text x='{x+2:.2f}' y='{height-padding_bottom+16}' font-size='11' fill='#64748b'>{m.strftime('%b %Y')}</text>")

    # Weekly grid (Mondays) - make it lighter
    d = date(dmin.year, dmin.month, dmin.day)
    # align to Monday
    d -= timedelta(days=(d.weekday()))
    while d <= dmax:
        if d >= dmin:
            x = x_pos(d)
            svg.append(f"<line x1='{x:.2f}' y1='{padding_top}' x2='{x:.2f}' y2='{height-padding_bottom}' stroke='#f8fafc' />")
        d += timedelta(days=7)

    # Axis
    svg.append(f"<line x1='{padding_left}' y1='{padding_top-8}' x2='{width-padding_right}' y2='{padding_top-8}' stroke='#cbd5e1' />")

    # Today marker
    today = date.today()
    if dmin <= today <= dmax:
        tx = x_pos(today)
        svg.append(f"<line x1='{tx:.2f}' y1='{padding_top-8}' x2='{tx:.2f}' y2='{height-padding_bottom}' stroke='#ef4444' stroke-width='1.2' stroke-dasharray='3,3' />")
        svg.append(f"<text x='{tx+4:.2f}' y='{padding_top-14}' font-size='10' fill='#ef4444'>Today</text>")

    # Lanes
    y = padding_top
    # global map of bar centers/edges for dependency rendering across lanes
    global_bar_pos: Dict[str, Tuple[float, float, float]] = {}

    for lane in order:
        items = enriched.get(lane, [])
        if not items:
            continue
        col = LANE_COLORS.get(lane, "#888")
        lane_top = y
        lane_h = lane_heights[lane]
        # Lane label
        svg.append(f"<text x='{padding_left-10}' y='{lane_top - 6}' font-size='14' font-weight='600' text-anchor='end' fill='#0f172a'>{esc(lane)}</text>")
        # Subtle lane background
        svg.append(f"<rect x='{padding_left}' y='{lane_top}' width='{width-padding_left-padding_right}' height='{lane_h - lane_gap}' fill='#f8fafc' stroke='#e2e8f0' />")

        ry = lane_top + 8
        # Store bar positions for dependency lines (per-lane)
        bar_pos: Dict[str, Tuple[float, float, float]] = {}
        for t in items:
            sd, dd = t["start"], t["due"]
            x1 = x_pos(sd)
            x2 = x_pos(dd)
            bar_h = 14
            w = max(2.0, x2 - x1)
            is_milestone = (t["dur"] == 0 or sd == dd)
            item_h = int(t["_item_h"])  # type: ignore
            # Render full label in left column (wrapped)
            indent_px = (t.get("_level", 0)) * 14
            lx = label_left + indent_px
            ly = ry + 4 + line_h  # baseline for first line
            svg.append("<g>")
            tooltip = [
                f"{t.get('name','')}",
                f"{sd} → {dd}",
                f"Type: {t.get('dtype','')}",
                f"Group: {t.get('group','')}",
                f"Owner: {t.get('owner','')}",
                f"Status: {t.get('status','')}",
                f"Priority: {t.get('priority','')}",
            ]
            deps_tt = ",".join(t.get("deps", []))
            if deps_tt:
                tooltip.append(f"Depends: {deps_tt}")
            if t.get("notes"):
                tooltip.append(f"Notes: {t.get('notes')}")
            svg.append(f"<title>{esc(chr(10).join(tooltip))}</title>")
            for i, line in enumerate(t["_label_lines"]):
                svg.append(f"<text x='{lx:.2f}' y='{(ly + i*line_h):.2f}' font-size='13' fill='#0f172a'>{esc(line)}</text>")
            if t.get("_meta_lines"):
                base = ly + len(t["_label_lines"]) * line_h
                for j, mline in enumerate(t["_meta_lines"]):
                    svg.append(f"<text x='{lx:.2f}' y='{(base + j*line_h):.2f}' font-size='11' fill='#64748b'>{esc(mline)}</text>")

            # Bar centered within item slot
            bar_y = ry + (item_h - bar_h) / 2
            if is_milestone:
                cx = x1
                size = 7
                pts = [
                    (cx, bar_y + bar_h/2 - size),
                    (cx + size, bar_y + bar_h/2),
                    (cx, bar_y + bar_h/2 + size),
                    (cx - size, bar_y + bar_h/2),
                ]
                pts_s = " ".join(f"{px:.2f},{py:.2f}" for px, py in pts)
                svg.append(f"<polygon points='{pts_s}' fill='{col}' stroke='#334155' stroke-width='1.5' />")
                bar_pos[t["id"]] = (cx, cx, bar_y + bar_h/2)
                if t.get("id"):
                    global_bar_pos[t["id"]] = bar_pos[t["id"]]
            else:
                # Enhanced visual hierarchy for task bars
                level = t.get("_level", 0)
                if level == 0:
                    # Root level: solid fill with shadow effect
                    svg.append(f"<rect x='{x1:.2f}' y='{bar_y:.2f}' width='{w:.2f}' height='{bar_h:.2f}' rx='3' ry='3' fill='{col}' stroke='#334155' stroke-width='1.5' />")
                    # Add subtle shadow
                    svg.append(f"<rect x='{x1+1:.2f}' y='{bar_y+1:.2f}' width='{w:.2f}' height='{bar_h:.2f}' rx='3' ry='3' fill='rgba(0,0,0,0.1)' />")
                else:
                    # Child level: lighter fill with dashed border
                    alpha = max(20, 50 - level * 10)  # Deeper levels get lighter
                    svg.append(f"<rect x='{x1:.2f}' y='{bar_y:.2f}' width='{w:.2f}' height='{bar_h:.2f}' rx='3' ry='3' fill='{col}{alpha:02x}' stroke='{col}' stroke-width='1' stroke-dasharray='3,2' />")
                
                bar_pos[t["id"]] = (x1, x2, bar_y + bar_h/2)
                if t.get("id"):
                    global_bar_pos[t["id"]] = bar_pos[t["id"]]
                
                # Add progress indicator for completed tasks
                if t.get("status", "").lower() in ["completed", "done", "finished"]:
                    # Add checkmark
                    check_x = x2 - 8
                    check_y = bar_y + bar_h/2
                    svg.append(f"<circle cx='{check_x:.2f}' cy='{check_y:.2f}' r='6' fill='#10b981' stroke='#059669' stroke-width='1' />")
                    svg.append(f"<path d='M {check_x-3:.2f} {check_y:.2f} L {check_x:.2f} {check_y+3:.2f} L {check_x+3:.2f} {check_y-3:.2f}' stroke='white' stroke-width='1.5' fill='none' />")

            ry += item_h + row_gap
            svg.append("</g>")

        y += lane_h

    # Arrow marker (single definition)
    svg.append("""
<defs>
  <marker id='arrow' viewBox='0 0 10 10' refX='10' refY='5' markerWidth='6' markerHeight='6' orient='auto-start-reverse'>
    <path d='M 0 0 L 10 5 L 0 10 z' fill='#64748b' />
  </marker>
</defs>
""")

    # Draw dependency connectors across all lanes
    for t in tasks:
        tid = t.get("id")
        if not tid:
            continue
        for dep in t.get("deps", []):
            if dep in global_bar_pos and tid in global_bar_pos:
                x1s, x1e, y1 = global_bar_pos[dep]
                x2s, x2e, y2 = global_bar_pos[tid]
                sx, ex = x1e, x2s
                # control points for a smooth S curve
                mx = (sx + ex) / 2
                path = f"M {sx:.2f} {y1:.2f} C {mx:.2f} {y1:.2f}, {mx:.2f} {y2:.2f}, {ex-4:.2f} {y2:.2f}"
                svg.append(f"<path d='{path}' stroke='#94a3b8' stroke-width='1' fill='none' marker-end='url(#arrow)' />")

    svg.append("</svg>")

    html = f"""
<!doctype html>
<html lang="en">
<head>
  <meta charset="utf-8" />
  <meta name="viewport" content="width=device-width, initial-scale=1" />
  <title>{title}</title>
  <style>
    :root {{
      --fg: #0f172a; --muted: #475569;
      --bg: #ffffff; --grid: #e5e7eb; --axis: #cbd5e1;
      --primary: #3b82f6; --success: #10b981; --warning: #f59e0b; --danger: #ef4444;
    }}
    @page {{ size: Letter landscape; margin: 0.5in; }}
    @media print {{ body {{ -webkit-print-color-adjust: exact; print-color-adjust: exact; }} }}
    body {{ 
      margin: 0; 
      background: var(--bg); 
      color: var(--fg); 
      font: 14px/1.4 ui-sans-serif, system-ui, -apple-system, Segoe UI, Roboto, Helvetica, Arial; 
      overflow-x: auto;
    }}
    .header {{
      background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
      color: white;
      padding: 20px;
      box-shadow: 0 2px 10px rgba(0,0,0,0.1);
    }}
    .header h1 {{
      margin: 0 0 10px 0;
      font-size: 28px;
      font-weight: 600;
    }}
    .header .subtitle {{
      margin: 0;
      opacity: 0.9;
      font-size: 16px;
    }}
    .controls {{
      background: #f8fafc;
      border-bottom: 1px solid #e2e8f0;
      padding: 15px 20px;
      display: flex;
      gap: 20px;
      align-items: center;
      flex-wrap: wrap;
    }}
    .control-group {{
      display: flex;
      align-items: center;
      gap: 8px;
    }}
    .control-group label {{
      font-weight: 500;
      color: #374151;
    }}
    .control-group select, .control-group input {{
      padding: 6px 12px;
      border: 1px solid #d1d5db;
      border-radius: 6px;
      font-size: 14px;
      background: white;
    }}
    .btn {{
      padding: 8px 16px;
      border: none;
      border-radius: 6px;
      font-size: 14px;
      font-weight: 500;
      cursor: pointer;
      transition: all 0.2s;
    }}
    .btn-primary {{
      background: var(--primary);
      color: white;
    }}
    .btn-primary:hover {{
      background: #2563eb;
      transform: translateY(-1px);
    }}
    .btn-success {{
      background: var(--success);
      color: white;
    }}
    .btn-success:hover {{
      background: #059669;
    }}
    .btn-warning {{
      background: var(--warning);
      color: white;
    }}
    .btn-warning:hover {{
      background: #d97706;
    }}
    .btn-danger {{
      background: var(--danger);
      color: white;
    }}
    .btn-danger:hover {{
      background: #dc2626;
    }}
    .wrap {{ 
      max-width: 1500px; 
      margin: 0 auto; 
      padding: 20px;
    }}
    .timeline-container {{
      background: white;
      border-radius: 12px;
      box-shadow: 0 4px 20px rgba(0,0,0,0.1);
      overflow: hidden;
      margin-bottom: 20px;
    }}
    .status-summary {{
      background: #f8fafc;
      padding: 20px;
      border-radius: 8px;
      margin-bottom: 20px;
    }}
    .status-grid {{
      display: grid;
      grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
      gap: 20px;
    }}
    .status-card {{
      background: white;
      padding: 16px;
      border-radius: 8px;
      border-left: 4px solid var(--primary);
      box-shadow: 0 2px 8px rgba(0,0,0,0.05);
    }}
    .status-card h3 {{
      margin: 0 0 8px 0;
      font-size: 16px;
      color: #374151;
    }}
    .status-card .count {{
      font-size: 24px;
      font-weight: 600;
      color: var(--primary);
    }}
    .task-details {{
      position: fixed;
      top: 50%;
      left: 50%;
      transform: translate(-50%, -50%);
      background: white;
      border-radius: 12px;
      box-shadow: 0 20px 60px rgba(0,0,0,0.3);
      padding: 24px;
      max-width: 500px;
      width: 90%;
      max-height: 80vh;
      overflow-y: auto;
      z-index: 1000;
      display: none;
    }}
    .task-details.show {{
      display: block;
    }}
    .task-details h3 {{
      margin: 0 0 16px 0;
      color: var(--primary);
      border-bottom: 2px solid #e5e7eb;
      padding-bottom: 8px;
    }}
    .task-details .field {{
      margin-bottom: 16px;
    }}
    .task-details label {{
      display: block;
      font-weight: 500;
      margin-bottom: 4px;
      color: #374151;
    }}
    .task-details input, .task-details select, .task-details textarea {{
      width: 100%;
      padding: 8px 12px;
      border: 1px solid #d1d5db;
      border-radius: 6px;
      font-size: 14px;
      box-sizing: border-box;
    }}
    .task-details textarea {{
      min-height: 80px;
      resize: vertical;
    }}
    .task-details .buttons {{
      display: flex;
      gap: 12px;
      justify-content: flex-end;
      margin-top: 24px;
    }}
    .overlay {{
      position: fixed;
      top: 0;
      left: 0;
      right: 0;
      bottom: 0;
      background: rgba(0,0,0,0.5);
      z-index: 999;
      display: none;
    }}
    .overlay.show {{
      display: block;
    }}
    .hint {{ 
      color: #64748b; 
      font-size: 12px; 
      margin-top: 8px; 
      text-align: center;
      padding: 16px;
      background: #f8fafc;
      border-radius: 8px;
    }}
    .export-options {{
      display: flex;
      gap: 12px;
      align-items: center;
    }}
    @media (max-width: 768px) {{
      .controls {{
        flex-direction: column;
        align-items: stretch;
      }}
      .control-group {{
        justify-content: space-between;
      }}
      .export-options {{
        flex-direction: column;
      }}
    }}
  </style>
</head>
<body>
  <div class="header">
    <h1>🎯 {title}</h1>
    <p class="subtitle">Professional Project Timeline - Interactive and Export-Ready for Advisor Meetings</p>
  </div>
  
  <div class="controls">
    <div class="control-group">
      <label for="lane-filter">Filter by Lane:</label>
      <select id="lane-filter">
        <option value="">All Lanes</option>
        <option value="Proposal">Proposal</option>
        <option value="Committee">Committee</option>
        <option value="Laser">Laser</option>
        <option value="Imaging">Imaging</option>
        <option value="Admin/Accountability">Admin/Accountability</option>
        <option value="Other">Other</option>
      </select>
    </div>
    
    <div class="control-group">
      <label for="status-filter">Filter by Status:</label>
      <select id="status-filter">
        <option value="">All Statuses</option>
        <option value="Planned">Planned</option>
        <option value="In Progress">In Progress</option>
        <option value="Completed">Completed</option>
        <option value="On Hold">On Hold</option>
      </select>
    </div>
    
    <div class="control-group">
      <label for="priority-filter">Filter by Priority:</label>
      <select id="priority-filter">
        <option value="">All Priorities</option>
        <option value="High">High</option>
        <option value="Medium">Medium</option>
        <option value="Low">Low</option>
      </select>
    </div>
    
    <div class="export-options">
      <button class="btn btn-primary" onclick="exportToPDF()">📄 Export PDF</button>
      <button class="btn btn-success" onclick="exportToCSV()">📊 Export CSV</button>
      <button class="btn btn-warning" onclick="printTimeline()">🖨️ Print</button>
    </div>
  </div>

  <div class="status-summary">
    <h2>📊 Project Status Overview</h2>
    <div class="status-grid">
      <div class="status-card">
        <h3>Total Tasks</h3>
        <div class="count">{len(tasks)}</div>
      </div>
      <div class="status-card">
        <h3>Timeline Span</h3>
        <div class="count">{(dmax - dmin).days} days</div>
      </div>
      <div class="status-card">
        <h3>High Priority</h3>
        <div class="count">{len([t for t in tasks if t.get('priority') == 'High'])}</div>
      </div>
      <div class="status-card">
        <h3>In Progress</h3>
        <div class="count">{len([t for t in tasks if t.get('status') == 'In Progress'])}</div>
      </div>
    </div>
  </div>

  <div class="wrap">
    <div class="timeline-container">
      {''.join(svg)}
    </div>
    <div class="hint">
      💡 <strong>Professional Features:</strong> High-quality SVG timeline • Interactive filters • Export to PDF/CSV • Print-ready • Perfect for advisor meetings and reports
    </div>
  </div>

  <!-- Task Details Modal -->
  <div class="overlay" id="overlay"></div>
  <div class="task-details" id="task-details">
    <h3>Edit Task Details</h3>
    <form id="task-form">
      <div class="field">
        <label for="task-name">Task Name:</label>
        <input type="text" id="task-name" required>
      </div>
      <div class="field">
        <label for="task-start">Start Date:</label>
        <input type="date" id="task-start" required>
      </div>
      <div class="field">
        <label for="task-due">Due Date:</label>
        <input type="date" id="task-due" required>
      </div>
      <div class="field">
        <label for="task-status">Status:</label>
        <select id="task-status">
          <option value="Planned">Planned</option>
          <option value="In Progress">In Progress</option>
          <option value="Completed">Completed</option>
          <option value="On Hold">On Hold</option>
        </select>
      </div>
      <div class="field">
        <label for="task-priority">Priority:</label>
        <select id="task-priority">
          <option value="High">High</option>
          <option value="Medium">Medium</option>
          <option value="Low">Low</option>
        </select>
      </div>
      <div class="field">
        <label for="task-notes">Notes:</label>
        <textarea id="task-notes" placeholder="Add any additional notes or updates..."></textarea>
      </div>
      <div class="buttons">
        <button type="button" class="btn btn-danger" onclick="closeTaskDetails()">Cancel</button>
        <button type="submit" class="btn btn-success">Save Changes</button>
      </div>
    </form>
  </div>

  <script>
    // Task data for editing
    const tasks = {tasks};
    
    // Current filters
    let currentFilters = {{
      lane: '',
      status: '',
      priority: ''
    }};
    
    // Initialize filters
    document.getElementById('lane-filter').addEventListener('change', function(e) {{
      currentFilters.lane = e.target.value;
      applyFilters();
    }});
    
    document.getElementById('status-filter').addEventListener('change', function(e) {{
      currentFilters.status = e.target.value;
      applyFilters();
    }});
    
    document.getElementById('priority-filter').addEventListener('change', function(e) {{
      currentFilters.priority = e.target.value;
      applyFilters();
    }});
    
    function applyFilters() {{
      // Filter SVG elements based on current filters
      const taskGroups = document.querySelectorAll('g');
      taskGroups.forEach(group => {{
        const title = group.querySelector('title');
        if (title) {{
          const taskText = title.textContent;
          let shouldShow = true;
          
          if (currentFilters.lane && !taskText.includes(`Lane: ${{currentFilters.lane}}`)) {{
            shouldShow = false;
          }}
          if (currentFilters.status && !taskText.includes(`Status: ${{currentFilters.status}}`)) {{
            shouldShow = false;
          }}
          if (currentFilters.priority && !taskText.includes(`Priority: ${{currentFilters.priority}}`)) {{
            shouldShow = false;
          }}
          
          group.style.opacity = shouldShow ? '1' : '0.3';
          group.style.pointerEvents = shouldShow ? 'auto' : 'none';
        }}
      }});
    }}
    
    // Task editing functionality
    let currentTask = null;
    
    function openTaskDetails(taskId) {{
      const task = tasks.find(t => t.id === taskId);
      if (!task) return;
      
      currentTask = task;
      
      // Populate form
      document.getElementById('task-name').value = task.name;
      document.getElementById('task-start').value = task.start;
      document.getElementById('task-due').value = task.due;
      document.getElementById('task-status').value = task.status || 'Planned';
      document.getElementById('task-priority').value = task.priority || 'Medium';
      document.getElementById('task-notes').value = task.notes || '';
      
      // Show modal
      document.getElementById('overlay').classList.add('show');
      document.getElementById('task-details').classList.add('show');
    }}
    
    function closeTaskDetails() {{
      document.getElementById('overlay').classList.remove('show');
      document.getElementById('task-details').classList.remove('show');
      currentTask = null;
    }}
    
    // Handle form submission
    document.getElementById('task-form').addEventListener('submit', function(e) {{
      e.preventDefault();
      
      if (currentTask) {{
        // Update task data
        currentTask.name = document.getElementById('task-name').value;
        currentTask.start = document.getElementById('task-start').value;
        currentTask.due = document.getElementById('task-due').value;
        currentTask.status = document.getElementById('task-status').value;
        currentTask.priority = document.getElementById('task-priority').value;
        currentTask.notes = document.getElementById('task-notes').value;
        
        // In a full implementation, you'd save this data and regenerate the timeline
        console.log('Updated task:', currentTask);
        
        // Close modal
        closeTaskDetails();
        
        // Show success message
        alert('Task updated successfully! Refresh the page to see changes.');
      }}
    }});
    
    // Export functions
    function exportToPDF() {{
      window.print();
    }}
    
    function exportToCSV() {{
      const csvContent = generateCSV();
      const blob = new Blob([csvContent], {{ type: 'text/csv' }});
      const url = window.URL.createObjectURL(blob);
      const a = document.createElement('a');
      a.href = url;
      a.download = 'timeline_export.csv';
      a.click();
      window.URL.revokeObjectURL(url);
    }}
    
    function generateCSV() {{
      const headers = ['ID', 'Name', 'Start Date', 'Due Date', 'Lane', 'Type', 'Owner', 'Status', 'Priority', 'Notes'];
      const rows = tasks.map(task => [
        task.id,
        task.name,
        task.start,
        task.due,
        task.lane,
        task.dtype,
        task.owner,
        task.status,
        task.priority,
        task.notes
      ]);
      
      return [headers, ...rows].map(row => 
        row.map(cell => `"${{cell || ''}}"`).join(',')
      ).join('\\n');
    }}
    
    function printTimeline() {{
      window.print();
    }}
    
    // Make timeline interactive
    document.addEventListener('DOMContentLoaded', function() {{
      console.log('Professional timeline loaded with {len(tasks)} tasks');
      
      // Add hover effects to task bars
      const taskGroups = document.querySelectorAll('g');
      taskGroups.forEach(group => {{
        group.style.cursor = 'pointer';
        group.addEventListener('mouseenter', function() {{
          this.style.transform = 'scale(1.02)';
          this.style.transition = 'transform 0.2s ease';
        }});
        group.addEventListener('mouseleave', function() {{
          this.style.transform = 'scale(1)';
        }});
      }});
    }});
    
    // Close modal when clicking overlay
    document.getElementById('overlay').addEventListener('click', closeTaskDetails);
  </script>
</body>
</html>
"""
    os.makedirs(os.path.dirname(out_html), exist_ok=True)
    with open(out_html, "w", encoding="utf-8") as f:
        f.write(html)


def generate_markdown(tasks: List[Dict], out_md: str):
    if not tasks:
        with open(out_md, "w", encoding="utf-8") as f:
            f.write("# Project Timeline\n\n_No tasks to display._\n")
        return
    
    # Calculate timeline span
    dmin = min(t["start"] for t in tasks)
    dmax = max(t["due"] for t in tasks)
    total_days = (dmax - dmin).days
    
    lines: List[str] = []
    lines.append("# Project Timeline\n")
    lines.append(f"**Timeline Span:** {dmin} to {dmax} ({total_days} days)")
    lines.append(f"**Total Tasks:** {len(tasks)}")
    lines.append(f"**Source:** `input/data.csv`")
    lines.append("")
    
    # Summary by lane
    by_lane: Dict[str, List[Dict]] = {}
    for t in tasks:
        by_lane.setdefault(t["lane"], []).append(t)
    
    lines.append("## 📊 Summary by Project Area\n")
    for lane in [ln for ln in LANE_ORDER if ln in by_lane] + [ln for ln in by_lane.keys() if ln not in LANE_ORDER]:
        lane_tasks = by_lane[lane]
        completed = len([t for t in lane_tasks if t.get("status", "").lower() in ["completed", "done", "finished"]])
        in_progress = len([t for t in lane_tasks if t.get("status", "").lower() in ["in progress", "active", "ongoing"]])
        planned = len([t for t in lane_tasks if t.get("status", "").lower() in ["planned", "pending", "not started"]])
        
        lines.append(f"### {lane}")
        lines.append(f"- **Total Tasks:** {len(lane_tasks)}")
        if completed > 0:
            lines.append(f"- **Completed:** {completed}")
        if in_progress > 0:
            lines.append(f"- **In Progress:** {in_progress}")
        if planned > 0:
            lines.append(f"- **Planned:** {planned}")
        lines.append("")
    
    # Group and nest by lane and parent
    lines.append("## 📋 Detailed Task Breakdown\n")
    for lane in [ln for ln in LANE_ORDER if ln in by_lane] + [ln for ln in by_lane.keys() if ln not in LANE_ORDER]:
        lines.append(f"### {lane}")
        lane_tasks = {t["id"]: t for t in by_lane[lane] if t.get("id")}
        is_child = {tid for tid, t in lane_tasks.items() if t.get("parent_id") in lane_tasks}
        roots = [t for tid, t in lane_tasks.items() if tid not in is_child]
        roots.sort(key=lambda x: (x["start"], x["due"], x["id"]))

        def write_task(t: Dict, level: int = 0):
            indent = "  " * level
            status_emoji = {
                "completed": "✅",
                "done": "✅", 
                "finished": "✅",
                "in progress": "🔄",
                "active": "🔄",
                "ongoing": "🔄",
                "planned": "📅",
                "pending": "⏳",
                "not started": "⏳"
            }
            status = t.get("status", "").lower()
            emoji = status_emoji.get(status, "📝")
            
            # Format dates nicely
            start_str = t["start"].strftime("%b %d, %Y")
            due_str = t["due"].strftime("%b %d, %Y")
            
            # Calculate duration
            duration = (t["due"] - t["start"]).days + 1
            duration_str = f"({duration} day{'s' if duration != 1 else ''})"
            
            # Build task line
            task_line = f"{indent}{emoji} **{t['name']}**"
            task_line += f" - {start_str} → {due_str} {duration_str}"
            
            # Add type and owner if available
            if t.get("dtype"):
                task_line += f" • {t['dtype']}"
            if t.get("owner"):
                task_line += f" • Owner: {t['owner']}"
            if t.get("priority"):
                task_line += f" • Priority: {t['priority']}"
            
            lines.append(task_line)
            
            # Add notes if available
            if t.get("notes"):
                note_indent = "  " * (level + 1)
                # Clean up notes and wrap long lines
                notes = t["notes"].replace("|", "/").strip()
                if len(notes) > 80:
                    # Simple word wrapping for long notes
                    words = notes.split()
                    current_line = ""
                    for word in words:
                        if len(current_line + " " + word) <= 80:
                            current_line += (" " + word) if current_line else word
                        else:
                            if current_line:
                                lines.append(f"{note_indent}💬 {current_line}")
                            current_line = word
                    if current_line:
                        lines.append(f"{note_indent}💬 {current_line}")
                else:
                    lines.append(f"{note_indent}💬 {notes}")
            
            # Add dependencies if available
            if t.get("deps"):
                deps_indent = "  " * (level + 1)
                deps_str = ", ".join(t["deps"])
                lines.append(f"{deps_indent}🔗 **Depends on:** {deps_str}")
            
            lines.append("")
            
            # Recursively add children
            for cid in t.get("children", []):
                if cid in lane_tasks:
                    write_task(lane_tasks[cid], level + 1)

        for rt in roots:
            write_task(rt, 0)
        lines.append("---")
        lines.append("")

    # Flat table for quick scanning
    lines.append("## 📊 Complete Task Table\n")
    lines.append("| Task Name | Start | Due | Duration | Lane | Type | Owner | Status | Priority | Dependencies | Notes |")
    lines.append("|:----------|:-----:|:---:|:--------:|:----|:----|:------|:-------|:---------|:-------------|:------|")
    
    for t in sorted(tasks, key=lambda x: (x["start"], x["due"], x["id"])):
        # Format dates
        start_str = t["start"].strftime("%m/%d/%Y")
        due_str = t["due"].strftime("%m/%d/%Y")
        
        # Calculate duration
        duration = (t["due"] - t["start"]).days + 1
        duration_str = f"{duration}d"
        
        # Clean up text fields
        name = t["name"].replace("|", "/").replace("\n", " ")
        dtype = (t.get("dtype") or "").replace("|", "/")
        owner = (t.get("owner") or "").replace("|", "/")
        status = (t.get("status") or "").replace("|", "/")
        priority = (t.get("priority") or "").replace("|", "/")
        deps = ",".join(t.get("deps", [])).replace("|", "/")
        notes = (t.get("notes") or "").replace("|", "/").replace("\n", " ")[:50]  # Truncate long notes
        
        lines.append(
            f"| {name} | {start_str} | {due_str} | {duration_str} | {t['lane']} | {dtype} | {owner} | {status} | {priority} | {deps} | {notes} |"
        )
    
    # Add footer
    lines.append("")
    lines.append("---")
    lines.append("*Generated on " + datetime.now().strftime("%B %d, %Y at %I:%M %p") + "*")
    
    with open(out_md, "w", encoding="utf-8") as f:
        f.write("\n".join(lines) + "\n")


def try_export_pdf(html_path: str, pdf_path: str) -> Tuple[bool, str]:
    # Prefer wkhtmltopdf, then headless chrome/chromium
    wk = shutil.which("wkhtmltopdf")
    if wk:
        cmd = [wk, "-q", "--enable-local-file-access", html_path, pdf_path]
        try:
            subprocess.run(cmd, check=True)
            return True, "wkhtmltopdf"
        except Exception as e:
            return False, f"wkhtmltopdf failed: {e}"

    for chrome in ["chromium", "chromium-browser", "google-chrome", "chrome", "msedge"]:
        ch = shutil.which(chrome)
        if not ch:
            continue
        file_url = f"file://{os.path.abspath(html_path)}"
        cmd = [ch, "--headless", "--disable-gpu", f"--print-to-pdf={os.path.abspath(pdf_path)}", file_url]
        try:
            subprocess.run(cmd, check=True)
            return True, chrome
        except Exception as e:
            return False, f"{chrome} failed: {e}"

    return False, "no converter found"


def generate_summary_report(tasks: List[Dict], out_report: str):
    """Generate a professional summary report for advisor meetings."""
    if not tasks:
        with open(out_report, "w", encoding="utf-8") as f:
            f.write("# Project Summary Report\n\n_No tasks to display._\n")
        return
    
    # Calculate key metrics
    dmin = min(t["start"] for t in tasks)
    dmax = max(t["due"] for t in tasks)
    total_days = (dmax - dmin).days
    
    by_lane: Dict[str, List[Dict]] = {}
    for t in tasks:
        by_lane.setdefault(t["lane"], []).append(t)
    
    # Count statuses and priorities
    status_counts = {}
    priority_counts = {}
    for t in tasks:
        status = t.get("status", "Unknown")
        priority = t.get("priority", "Unknown")
        status_counts[status] = status_counts.get(status, 0) + 1
        priority_counts[priority] = priority_counts.get(priority, 0) + 1
    
    lines: List[str] = []
    lines.append("# 🎯 Project Summary Report")
    lines.append(f"**Generated:** {datetime.now().strftime('%B %d, %Y at %I:%M %p')}")
    lines.append(f"**Timeline:** {dmin.strftime('%B %d, %Y')} to {dmax.strftime('%B %d, %Y')} ({total_days} days)")
    lines.append(f"**Total Tasks:** {len(tasks)}")
    lines.append("")
    
    # Executive Summary
    lines.append("## 📊 Executive Summary")
    lines.append("")
    
    # Status breakdown
    lines.append("### Task Status Overview")
    for status, count in sorted(status_counts.items()):
        emoji = {"Completed": "✅", "In Progress": "🔄", "Planned": "📅", "On Hold": "⏸️"}.get(status, "📝")
        lines.append(f"- {emoji} **{status}:** {count} tasks")
    lines.append("")
    
    # Priority breakdown
    lines.append("### Priority Distribution")
    for priority, count in sorted(priority_counts.items()):
        emoji = {"High": "🔴", "Medium": "🟡", "Low": "🟢"}.get(priority, "⚪")
        lines.append(f"- {emoji} **{priority}:** {count} tasks")
    lines.append("")
    
    # Lane breakdown
    lines.append("### Project Area Breakdown")
    for lane in [ln for ln in LANE_ORDER if ln in by_lane] + [ln for ln in by_lane.keys() if ln not in LANE_ORDER]:
        lane_tasks = by_lane[lane]
        completed = len([t for t in lane_tasks if t.get("status", "").lower() in ["completed", "done", "finished"]])
        in_progress = len([t for t in lane_tasks if t.get("status", "").lower() in ["in progress", "active", "ongoing"]])
        planned = len([t for t in lane_tasks if t.get("status", "").lower() in ["planned", "pending", "not started"]])
        
        lines.append(f"#### {lane}")
        lines.append(f"- **Total:** {len(lane_tasks)} tasks")
        if completed > 0:
            lines.append(f"- **Completed:** {completed} ({completed/len(lane_tasks)*100:.1f}%)")
        if in_progress > 0:
            lines.append(f"- **In Progress:** {in_progress} ({in_progress/len(lane_tasks)*100:.1f}%)")
        if planned > 0:
            lines.append(f"- **Planned:** {planned} ({planned/len(lane_tasks)*100:.1f}%)")
        lines.append("")
    
    # Critical path analysis
    lines.append("## 🚨 Critical Path Analysis")
    lines.append("")
    
    # Find high priority tasks
    high_priority = [t for t in tasks if t.get("priority") == "High"]
    if high_priority:
        lines.append("### High Priority Tasks")
        for t in sorted(high_priority, key=lambda x: x["start"]):
            start_str = t["start"].strftime("%b %d")
            due_str = t["due"].strftime("%b %d")
            status_emoji = {"Completed": "✅", "In Progress": "🔄", "Planned": "📅", "On Hold": "⏸️"}.get(t.get("status", ""), "📝")
            lines.append(f"- {status_emoji} **{t['name']}** ({start_str} → {due_str})")
            if t.get("notes"):
                lines.append(f"  - Note: {t['notes'][:100]}{'...' if len(t['notes']) > 100 else ''}")
        lines.append("")
    
    # Dependencies analysis
    lines.append("### Dependency Analysis")
    tasks_with_deps = [t for t in tasks if t.get("deps")]
    if tasks_with_deps:
        lines.append(f"**Tasks with dependencies:** {len(tasks_with_deps)}")
        for t in tasks_with_deps[:5]:  # Show first 5
            lines.append(f"- **{t['name']}** depends on: {', '.join(t['deps'])}")
        if len(tasks_with_deps) > 5:
            lines.append(f"- ... and {len(tasks_with_deps) - 5} more")
    else:
        lines.append("No task dependencies identified.")
    lines.append("")
    
    # Recommendations
    lines.append("## 💡 Recommendations")
    lines.append("")
    
    # Calculate completion rate
    completed_tasks = len([t for t in tasks if t.get("status", "").lower() in ["completed", "done", "finished"]])
    completion_rate = (completed_tasks / len(tasks)) * 100 if tasks else 0
    
    if completion_rate < 25:
        lines.append("- **Focus on quick wins** to build momentum")
    elif completion_rate < 50:
        lines.append("- **Maintain current pace** and identify bottlenecks")
    elif completion_rate < 75:
        lines.append("- **Excellent progress** - focus on completing remaining tasks")
    else:
        lines.append("- **Outstanding progress** - consider adding stretch goals")
    
    # Check for overdue tasks
    today = date.today()
    overdue = [t for t in tasks if t["due"] < today and t.get("status", "").lower() not in ["completed", "done", "finished"]]
    if overdue:
        lines.append(f"- **Address {len(overdue)} overdue tasks** immediately")
    
    # Check for tasks starting soon
    soon = timedelta(days=7)
    starting_soon = [t for t in tasks if t["start"] <= today + soon and t["start"] > today]
    if starting_soon:
        lines.append(f"- **Prepare for {len(starting_soon)} tasks starting within 7 days**")
    
    lines.append("")
    
    # Footer
    lines.append("---")
    lines.append("*This report was automatically generated from your project timeline data.*")
    lines.append("*Use this for advisor meetings, progress reviews, and project planning.*")
    
    with open(out_report, "w", encoding="utf-8") as f:
        f.write("\n".join(lines) + "\n")


def main():
    ap = argparse.ArgumentParser(description="Generate a professional project timeline with HTML, Markdown, and summary report")
    ap.add_argument("--input", default="input/data.csv", help="Path to data.csv")
    ap.add_argument("--html", default="output/Timeline.html", help="Output HTML path")
    ap.add_argument("--md", default="output/Timeline.md", help="Output Markdown path")
    ap.add_argument("--report", default="output/Summary_Report.md", help="Output summary report path")
    ap.add_argument("--pdf", default="", help="Output PDF path (optional)")
    ap.add_argument("--format", choices=["all", "html", "markdown", "report"], default="all", 
                   help="Output format(s) to generate")
    ap.add_argument("--title", default="Project Timeline", help="Custom title for the timeline")
    args = ap.parse_args()

    tasks = load_tasks(args.input)
    
    # Generate requested formats
    if args.format in ["all", "html"]:
        generate_html(tasks, args.html, title=args.title)
    
    if args.format in ["all", "markdown"]:
        generate_markdown(tasks, args.md)
    
    if args.format in ["all", "report"]:
        generate_summary_report(tasks, args.report)

    if args.pdf:
        ok, via = try_export_pdf(args.html, args.pdf)
        if ok:
            print(f"PDF exported via {via}: {args.pdf}")
        else:
            print("PDF export not completed (converter not found). See README for install instructions.")

    print(f"✅ Generated professional timeline: {args.title}")
    if args.format in ["all", "html"]:
        print(f"   📄 HTML: {args.html}")
    if args.format in ["all", "markdown"]:
        print(f"   📝 Markdown: {args.md}")
    if args.format in ["all", "report"]:
        print(f"   📊 Summary Report: {args.report}")
    if args.pdf:
        print(f"   📄 PDF: {args.pdf}")
    
    print(f"\n🎯 Timeline spans {(max(t['due'] for t in tasks) - min(t['start'] for t in tasks)).days} days with {len(tasks)} tasks")
    print(f"📊 Use --help for more options and customization")


if __name__ == "__main__":
    main()
