# Certification badge images

Drop badge images here with the filename referenced by `image:` in
`data/certifications.yaml` (e.g. `crto.png`). Hugo resizes/optimises them at
build time (`.Fit "200x200 Contain"` in `layouts/shortcodes/certifications.html`),
so a reasonably large source (400px+) is fine — no need to pre-resize.

Until a file is present, the certification card falls back to a plain shield
icon instead of failing the build, so metadata can be added ahead of the image.
