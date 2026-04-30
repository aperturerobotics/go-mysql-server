#!/usr/bin/env bash
set -euo pipefail

# Reuse one compiled engine test binary and remove it when this run finishes.
test_dir=$(mktemp -d)
trap 'rm -rf "$test_dir"' EXIT

# Keep the large engine suite from exhausting one package-wide timeout.
go list ./... > "$test_dir/packages"
packages=()
while IFS= read -r package; do
  if [[ "$package" != github.com/dolthub/go-mysql-server/enginetest ]]; then
    packages+=("$package")
  fi
done < "$test_dir/packages"

go test -race -timeout=2m -tags=gms_icu_regex "${packages[@]}"

# Give every engine test, example, and fuzz seed its own timeout.
go test -race -c -tags=gms_icu_regex -o "$test_dir/enginetest" ./enginetest
cd enginetest
"$test_dir/enginetest" -test.list . > "$test_dir/tests"
while IFS= read -r test_name; do
  case "$test_name" in
    Test*|Example*|Fuzz*)
      echo "Running $test_name"
      "$test_dir/enginetest" -test.timeout=2m -test.run "^${test_name}$"
      ;;
  esac
done < "$test_dir/tests"
