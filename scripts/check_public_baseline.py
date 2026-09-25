#!/usr/bin/env python3
from pathlib import Path
import re
import sys

ROOT = Path(__file__).resolve().parent.parent
CODE_EXTENSIONS = {".go", ".ini", ".ps1", ".cmd", ".sh"}
PATTERNS = [
    (re.compile(r"(?i)[A-Z]:\\\\Users\\\\"), "hardcoded Windows user-profile path"),
    (re.compile(r"(?i)[A-Z]:\\\\(?!Windows\\\\)"), "hardcoded drive-root path"),
]
errors = []
for path in ROOT.rglob("*"):
    if not path.is_file() or path.suffix.lower() not in CODE_EXTENSIONS:
        continue
    if "build" in path.parts or "dist" in path.parts:
        continue
    text = path.read_text(encoding="utf-8", errors="ignore")
    for pattern, label in PATTERNS:
        if pattern.search(text):
            errors.append(f"{path.relative_to(ROOT)} contains {label}")
if errors:
    print("\n".join(errors))
    sys.exit(1)
print("Public baseline check: OK")
