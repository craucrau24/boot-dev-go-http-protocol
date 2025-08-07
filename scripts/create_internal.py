#! /usr/bin/env python

import sys
import os

def write_to_file(name, content):
  with open(name, mode="w") as f:
    f.write(content)


if __name__ == "__main__":
  if len(sys.argv) < 2:
    print("Missing argument (package name)")
    sys.exit(1)
  
  pkg_name = sys.argv[1].lower()

  target_dir = os.path.join(".", "internal", pkg_name)
  if os.path.exists(target_dir):
    print("Package with same name already exists")
    sys.exit(1)

  os.makedirs(target_dir)
  content = f"package {pkg_name}\n\n"
  write_to_file(os.path.join(target_dir, f"{pkg_name}.go"), content)
  write_to_file(os.path.join(target_dir, f"{pkg_name}_test.go"), content)
  