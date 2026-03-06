---
id: jus-tjpo
status: open
deps: []
links: []
created: 2026-03-02T04:54:24Z
type: epic
priority: 1
assignee: Alex Cabrera
tags: [cli, restructure]
---
# Command restructure: execution order and version bump

Execution order for the command restructure:

1. jus-936u  Remove completion command (quick, independent)
2. jus-ri29  Group template commands under 'template' parent  
3. jus-x2r5  Group project commands under 'project' parent
4. jus-du4p  Absorb refresh into 'proxy reload'
5. jus-air9  Merge tunnel into 'share --attach'
6. jus-6i08  Add cobra command groups to help output
7. jus-7z4w  Add hidden backward-compat aliases
8. jus-vf6p  Keep top-level commands as-is (verify/adjust)
9. jus-xw6m  Update README and docs

Each step should be independently testable. Run full test suite after each. Bump to v1.1.0 when complete (this is a feature change, not a patch).

