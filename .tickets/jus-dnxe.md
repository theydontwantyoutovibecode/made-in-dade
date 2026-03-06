---
id: jus-dnxe
status: closed
deps: []
links: []
created: 2026-03-01T22:07:13Z
type: feature
priority: 1
assignee: Alex Cabrera
tags: [windows, manifest, cross-platform]
---
# Add setup_windows field to dade.toml manifest schema

Extend the manifest schema to support Windows-specific setup scripts and commands. Add setup_windows, dev_windows, prod_windows fields.

## Acceptance Criteria

- Add setup_windows to [scaffold] section
- Add dev_windows to [dev] section
- Add prod_windows to [prod] section
- Update manifest parser to read these fields
- Use Windows-specific commands when runtime.GOOS == windows

