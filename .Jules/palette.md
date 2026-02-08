## 2026-01-14 - CLI Log Noise vs UI Output
**Learning:** Default loggers writing to `stderr` often interrupt formatted CLI UI output on `stdout`, creating a messy user experience.
**Action:** In CLI tools with rich UI, default logging to `Discard` or a file, or use `Debug` level for operational logs to keep the standard output clean.

## 2026-01-14 - Go Template Receiver Types
**Learning:** `text/template` fails to call methods defined on pointer receivers if the data passed is a value type.
**Action:** Ensure data passed to templates is a pointer, or access struct fields directly if public, to avoid runtime template execution errors.

## 2026-01-14 - Structured CLI Error Messages
**Learning:** Plain text error lists (like "Row 5, Field 'Status', Value 'X': Error") are hard to scan.
**Action:** Structure CLI validation errors visually using dim/bold/color (e.g., `Row 5` • **Field** • 'Value': Message) to separate metadata from the message content.

## 2026-02-08 - Accessible LaTeX Icons
**Learning:** Purely visual LaTeX symbols (like `$\star$` or `$\checkmark$`) are read literally by screen readers. The `accsupp` package provides a way to substitute them with semantic text.
**Action:** Wrap decorative icons with `\BeginAccSupp{method=pdfstringdef,unicode,ActualText={Semantic Text} } ... \EndAccSupp{}`. Use braces `{}` around the text to preserve spaces.
