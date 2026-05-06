from __future__ import annotations

import sys
from pathlib import Path

from PIL import Image


def build_square_icon(source: Path, target: Path) -> None:
    image = Image.open(source).convert("RGBA")

    size = max(image.size)
    canvas = Image.new("RGBA", (size, size), (0, 0, 0, 0))
    offset = ((size - image.width) // 2, (size - image.height) // 2)
    canvas.paste(image, offset)

    target.parent.mkdir(parents=True, exist_ok=True)
    canvas.save(target, format="ICO", sizes=[(256, 256), (128, 128), (64, 64), (48, 48), (32, 32), (16, 16)])


def main() -> int:
    if len(sys.argv) != 3:
        print("usage: python scripts/make_windows_icon.py <source-image> <target-ico>", file=sys.stderr)
        return 1

    source = Path(sys.argv[1]).resolve()
    target = Path(sys.argv[2]).resolve()

    if not source.is_file():
        print(f"source image not found: {source}", file=sys.stderr)
        return 1

    build_square_icon(source, target)
    print(target)
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
