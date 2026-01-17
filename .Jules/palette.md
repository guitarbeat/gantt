## 2026-01-17 - LaTeX Accessibility with accsupp
**Learning:** LaTeX PDFs are often inaccessible to screen readers because decorative symbols (like bullets or icons) are read aloud as "bullet" or "star", creating noise.
**Action:** Use the `accsupp` package and `\BeginAccSupp{ActualText={}}...\EndAccSupp{}` to hide decorative elements or provide alternative text for icons.
