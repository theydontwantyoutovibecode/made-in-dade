---
id: jus-jobc
status: closed
deps: [jus-7ww4]
links: []
created: 2026-03-02T01:58:38Z
type: task
priority: 2
assignee: Alex Cabrera
parent: jus-odye
tags: [homebrew, ws-6]
---
# Move homebrew tap to new org

Move the dade formula from theydontwantyoutovibecode/homebrew-tap to theydontwantyoutovibecode/homebrew-tap. Read the existing formula from the old tap, update it to reference the new GitHub org for binary downloads, push to the new tap repo. Test: brew tap theydontwantyoutovibecode/tap && brew install dade.

## Acceptance Criteria

1. Formula in new tap. 2. brew install works. 3. Old tap can be deprecated.

