#!/usr/bin/env python3
"""Copy the NuGet README into the generated .NET SDK and pack it.

pulumi package gen-sdk rewrites sdk/dotnet, including the .csproj and a
one-line README.md, so this runs after generation (see the Makefile).
"""
from __future__ import annotations

import shutil
import sys
from pathlib import Path

ROOT = Path(__file__).resolve().parents[1]
SRC_README = ROOT / "docs" / "nuget-README.md"
DOTNET = ROOT / "sdk" / "dotnet"
DEST_README = DOTNET / "README.md"
CSPROJ = DOTNET / "NodeRings.Pulumi.csproj"

README_PROPERTY = "    <PackageReadmeFile>README.md</PackageReadmeFile>\n"
README_ITEM = """
  <ItemGroup>
    <None Include="README.md" Pack="true" PackagePath="\\" />
  </ItemGroup>
"""


def apply_csproj(text: str) -> str:
    if "<PackageReadmeFile>" not in text:
        needle = "    <PackageIcon>logo.png</PackageIcon>\n"
        if needle not in text:
            raise SystemExit("NodeRings.Pulumi.csproj is missing <PackageIcon>logo.png</PackageIcon>; cannot insert PackageReadmeFile")
        text = text.replace(needle, needle + README_PROPERTY, 1)
    if 'Include="README.md"' not in text:
        closing = "\n</Project>\n"
        if not text.endswith(closing) and not text.rstrip().endswith("</Project>"):
            raise SystemExit("NodeRings.Pulumi.csproj does not end with </Project>")
        # Insert the item group immediately before </Project>, keeping a trailing newline.
        text = text.rstrip()[: -len("</Project>")].rstrip() + "\n" + README_ITEM + "</Project>\n"
    return text


def main() -> int:
    if not SRC_README.is_file():
        raise SystemExit(f"missing {SRC_README}")
    if not CSPROJ.is_file():
        raise SystemExit(f"missing {CSPROJ}")
    DEST_README.parent.mkdir(parents=True, exist_ok=True)
    shutil.copyfile(SRC_README, DEST_README)
    updated = apply_csproj(CSPROJ.read_text())
    CSPROJ.write_text(updated)
    return 0


if __name__ == "__main__":
    sys.exit(main())
