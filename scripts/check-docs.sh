#!/usr/bin/env bash
set -euo pipefail

root="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "${root}"

required=(
	README.md CHANGELOG.md COMPATIBILITY.md CONTRIBUTING.md DEPRECATION.md
	LICENSE SECURITY.md SUPPORT.md example_test.go docs/README.md
	docs/specification-decisions.md docs/numeric-model.md docs/precision.md
	docs/conditions.md docs/serialization.md docs/security.md
	docs/performance.md docs/benchmark-baseline.md docs/migration.md
	docs/cookbook.md docs/compatibility.md docs/faq.md
	docs/troubleshooting.md docs/verification.md specification/README.md
	specification/manifest.tsv specification/decisions.json
	specification/conformance.json specification/decision-history.json
	specification/monitoring.json specification/maintained-peers.json
)
for path in "${required[@]}"; do
	if [[ ! -s "${path}" ]]; then
		printf 'missing documentation: %s\n' "${path}" >&2
		exit 1
	fi
done

python3 - "${required[@]}" <<'PY'
from pathlib import Path
from urllib.parse import unquote, urlsplit
import json
import os
import re
import stat
import subprocess
import sys
import tempfile


root = Path.cwd().resolve(strict=True)


def decode_destination(value: str, document: Path, raw_target: str) -> str:
    if re.search(r"%(?![0-9A-Fa-f]{2})", value):
        raise SystemExit(f"malformed percent escape in {document}: {raw_target}")
    decoded = unquote(value, errors="strict")
    if unquote(decoded, errors="strict") != decoded:
        raise SystemExit(f"double-encoded documentation link in {document}: {raw_target}")
    return decoded


def contained_regular(path: Path, *, directory_index: bool = False) -> Path:
    lexical = Path(os.path.abspath(path))
    try:
        relative = lexical.relative_to(root)
    except ValueError as error:
        raise SystemExit(f"documentation target escapes repository: {path}") from error

    current = root
    for part in relative.parts:
        current /= part
        try:
            mode = current.lstat().st_mode
        except FileNotFoundError as error:
            raise SystemExit(f"documentation target does not exist: {path}") from error
        if stat.S_ISLNK(mode):
            raise SystemExit(f"documentation target must not be a symlink: {path}")

    if directory_index and current.is_dir():
        current /= "README.md"
        try:
            mode = current.lstat().st_mode
        except FileNotFoundError as error:
            raise SystemExit(f"documentation directory has no README.md: {path}") from error
        if stat.S_ISLNK(mode):
            raise SystemExit(f"documentation target must not be a symlink: {path}")
    else:
        mode = current.lstat().st_mode

    if not stat.S_ISREG(mode):
        raise SystemExit(f"documentation target must be a regular file: {path}")
    physical = current.resolve(strict=True)
    try:
        physical.relative_to(root)
    except ValueError as error:
        raise SystemExit(f"documentation target escapes repository: {path}") from error
    return physical


for required_path in sys.argv[1:]:
    contained_regular(Path(required_path))


documents = [
    contained_regular(path)
    for path in Path(".").rglob("*.md")
    if not {".golib-tooling", ".verification", "node_modules"}.intersection(path.parts)
]
markdown_parser = r'''
import { readFileSync } from "node:fs";
import { marked } from "marked";
import GithubSlugger from "github-slugger";

function plainText(tokens = []) {
  return tokens.map((token) => {
    if (["text", "escape", "codespan"].includes(token.type)) {
      return token.text ?? token.raw ?? "";
    }
    if (token.type === "image") {
      return token.text ?? "";
    }
    if (Array.isArray(token.tokens)) {
      return plainText(token.tokens);
    }
    return token.text ?? "";
  }).join("");
}

const records = [];
for (const file of process.argv.slice(1)) {
  const tokens = marked.lexer(readFileSync(file, "utf8"), { gfm: true });
  const slugger = new GithubSlugger();
  const record = { file, anchors: [], links: [], commands: [] };
  marked.walkTokens(tokens, (token) => {
    if (token.type === "heading") {
      record.anchors.push(slugger.slug(plainText(token.tokens)));
    }
    if (token.type === "link" || token.type === "image") {
      record.links.push(token.href);
    }
    if (token.type === "codespan") {
      record.commands.push(token.text);
    }
    if (
      token.type === "code"
      && /^(?:(?:sh|bash|shell)(?:\s|$)|$)/i.test(token.lang ?? "")
    ) {
      record.commands.push(token.text);
    }
  });
  records.push(record);
}
process.stdout.write(JSON.stringify(records));
'''
parsed_markdown = subprocess.run(
    ["node", "--input-type=module", "--eval", markdown_parser, *map(str, documents)],
    check=True,
    capture_output=True,
    text=True,
)
markdown_records = json.loads(parsed_markdown.stdout)
anchor_cache = {
    Path(record["file"]): set(record["anchors"])
    for record in markdown_records
}
for record in markdown_records:
    document = Path(record["file"])
    for raw_target in record["links"]:
        target = raw_target.strip()
        parsed = urlsplit(target)
        scheme = parsed.scheme.lower()
        if scheme in {"http", "https", "mailto"}:
            continue
        if scheme:
            raise SystemExit(f"unsupported documentation link scheme in {document}: {raw_target}")
        if parsed.netloc or target.startswith("//"):
            raise SystemExit(f"network-path documentation link in {document}: {raw_target}")
        relative = decode_destination(parsed.path, document, raw_target)
        fragment = decode_destination(parsed.fragment, document, raw_target).lower()
        if (
            "\\" in relative
            or relative.startswith("//")
            or re.match(r"^[A-Za-z]:", relative)
        ):
            raise SystemExit(f"unsafe local documentation link in {document}: {raw_target}")
        if Path(relative).is_absolute():
            raise SystemExit(f"absolute local link in {document}: {raw_target}")
        linked = contained_regular(
            document.parent / relative if relative else document,
            directory_index=True,
        )
        if fragment and linked.suffix.lower() == ".md":
            linked_anchors = anchor_cache.get(linked)
            if linked_anchors is None:
                raise SystemExit(f"linked Markdown target was not parsed: {linked}")
            if fragment not in linked_anchors:
                raise SystemExit(f"broken local anchor in {document}: {raw_target}")
print("documentation links and anchors resolve")

makefile_targets: set[str] = set()
for line in Path("Makefile").read_text(encoding="utf-8").splitlines():
    if line[:1].isspace():
        continue
    match = re.match(r"^([A-Za-z0-9_.-]+(?:[ \t]+[A-Za-z0-9_.-]+)*):(?:[^=]|$)", line)
    if match:
        makefile_targets.update(match.group(1).split())

documented_make_targets: set[str] = set()
for record in markdown_records:
    for fragment in record["commands"]:
        documented_make_targets.update(
            match.group(1)
            for match in re.finditer(
                r"(?m)(?:^|[;&|])[ \t]*(?:\$[ \t]+)?make[ \t]+([A-Za-z0-9_.-]+)",
                fragment,
            )
        )
unknown_make_targets = sorted(documented_make_targets - makefile_targets)
if unknown_make_targets:
    raise SystemExit(f"documented make targets do not exist: {unknown_make_targets}")
print("documented make targets exist")

package_manifest = json.loads(Path("packages.json").read_text(encoding="utf-8"))
module_manifest = json.loads(Path("modules.json").read_text(encoding="utf-8"))
package_records = sorted(package_manifest["packages"], key=lambda item: item["import_path"])
module_package_records = sorted(
    (package for module in module_manifest["modules"] for package in module["packages"]),
    key=lambda item: item["import_path"],
)
if package_records != module_package_records:
    raise SystemExit("packages.json and modules.json package inventories differ")
module_by_import: dict[str, dict] = {}
for module in module_manifest["modules"]:
    for package in module["packages"]:
        import_path = package["import_path"]
        if import_path in module_by_import:
            raise SystemExit(f"modules.json contains duplicate package ownership: {import_path}")
        module_by_import[import_path] = module

listed = subprocess.run(
    ["go", "list", "-mod=readonly", "-json", "./..."],
    check=True,
    capture_output=True,
    text=True,
)
decoder = json.JSONDecoder()
listed_packages: list[dict] = []
cursor = 0
while cursor < len(listed.stdout):
    while cursor < len(listed.stdout) and listed.stdout[cursor].isspace():
        cursor += 1
    if cursor == len(listed.stdout):
        break
    package, cursor = decoder.raw_decode(listed.stdout, cursor)
    listed_packages.append(package)
listed_packages = [
    package
    for package in listed_packages
    if not {".golib-tooling", ".verification", "node_modules"}.intersection(
        Path(package["Dir"]).resolve(strict=True).relative_to(root).parts
    )
]

manifest_by_import = {package["import_path"]: package for package in package_records}
listed_by_import = {package["ImportPath"]: package for package in listed_packages}
if len(manifest_by_import) != len(package_records):
    raise SystemExit("package manifests contain duplicate import paths")
if len(listed_by_import) != len(listed_packages):
    raise SystemExit("go list returned duplicate import paths")
if manifest_by_import.keys() != listed_by_import.keys():
    missing = sorted(listed_by_import.keys() - manifest_by_import.keys())
    stale = sorted(manifest_by_import.keys() - listed_by_import.keys())
    raise SystemExit(f"Go package inventory mismatch; missing={missing}, stale={stale}")
for import_path, package in listed_by_import.items():
    relative_directory = Path(os.path.relpath(package["Dir"], root)).as_posix()
    declared = manifest_by_import[import_path]
    if declared["directory"] != relative_directory or declared["name"] != package["Name"]:
        raise SystemExit(f"Go package identity differs from manifests: {import_path}")
    owner = module_by_import[import_path]
    listed_module = package.get("Module")
    if not listed_module:
        raise SystemExit(f"Go package has no module ownership: {import_path}")
    listed_module_directory = Path(
        os.path.relpath(listed_module["Dir"], root)
    ).as_posix()
    if (
        declared["module_directory"] != owner["directory"]
        or owner["directory"] != listed_module_directory
        or owner["module_path"] != listed_module["Path"]
    ):
        raise SystemExit(f"Go package module ownership differs from manifests: {import_path}")
print("Go package inventories reconcile")

public_packages = sorted(
    package["import_path"]
    for package in package_records
    if package["kind"] == "public"
)
if not public_packages:
    raise SystemExit("package inventory contains no public packages")
for import_path in public_packages:
    metadata = listed_by_import[import_path]
    if not metadata.get("Doc", "").strip():
        raise SystemExit(f"public package has no package documentation: {import_path}")
    subprocess.run(["go", "doc", import_path], check=True, stdout=subprocess.DEVNULL)
print("public package documentation renders")

helper = r'''
package main

import (
	"encoding/json"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"strings"
)

func main() {
	names := make([]string, 0)
	for _, path := range os.Args[2:] {
		file, err := parser.ParseFile(token.NewFileSet(), path, nil, 0)
		if err != nil {
			panic(err)
		}
		for _, declaration := range file.Decls {
			function, ok := declaration.(*ast.FuncDecl)
			if !ok || function.Recv != nil || !strings.HasPrefix(function.Name.Name, "Example") {
				continue
			}
			names = append(names, function.Name.Name)
		}
	}
	if err := json.NewEncoder(os.Stdout).Encode(names); err != nil {
		panic(err)
	}
}
'''
with tempfile.TemporaryDirectory(prefix="go-math-doc-examples-") as temporary:
    helper_path = Path(temporary) / "main.go"
    helper_path.write_text(helper, encoding="utf-8")
    declaration_count = 0
    for import_path in sorted(listed_by_import):
        directory = Path(listed_by_import[import_path]["Dir"])
        test_files = sorted(contained_regular(path) for path in directory.glob("*_test.go"))
        parsed = subprocess.run(
            ["go", "run", helper_path, "--", *map(str, test_files)],
            check=True,
            capture_output=True,
            text=True,
        )
        declarations = sorted(json.loads(parsed.stdout))
        declaration_count += len(declarations)
        listed_examples = subprocess.run(
            ["go", "test", "-mod=readonly", "-list", "^Example", import_path],
            check=True,
            capture_output=True,
            text=True,
        )
        registered = sorted(
            line for line in listed_examples.stdout.splitlines() if line.startswith("Example")
        )
        if declarations != registered:
            missing = sorted(set(declarations) - set(registered))
            extra = sorted(set(registered) - set(declarations))
            raise SystemExit(
                f"every exact example declaration in {import_path} must have an attached "
                f"Output or Unordered output oracle; missing={missing}, extra={extra}"
            )

if declaration_count == 0:
    raise SystemExit("repository contains no Example declarations")
print("executable examples have attached output oracles")
PY

packages=()
while IFS= read -r package; do
	packages+=("${package}")
done < <(jq -r '.packages[].import_path' packages.json)
go test -mod=readonly "${packages[@]}" -run '^Example' -count=1
printf 'documentation contract passed\n'
