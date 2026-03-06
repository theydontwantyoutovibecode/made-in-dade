---
id: jus-ni16
status: open
deps: [jus-1p4n]
links: []
created: 2026-02-12T20:01:05Z
type: task
priority: 1
assignee: Alex Cabrera
parent: jus-k78p
tags: [homebrew, distribution]
---
# Update Homebrew formula for dade v1.0.0

Update the Homebrew formula for the new unified dade.

## Repository

theydontwantyoutovibecode/homebrew-tap

## Formula Location

Formula/dade.rb

## Changes Needed

### 1. Update Version and SHA

After tagging v1.0.0:
```ruby
class Justvibin < Formula
  desc "CLI for scaffolding and serving web projects with automatic HTTPS"
  homepage "https://github.com/theydontwantyoutovibecode/dade"
  url "https://github.com/theydontwantyoutovibecode/dade/archive/refs/tags/v1.0.0.tar.gz"
  sha256 "<new-sha256>"
  license "MIT"
  
  depends_on "bash" => "4.0"
  depends_on "jq"
  depends_on "caddy"
  
  def install
    bin.install "dade"
  end
  
  def caveats
    <<~EOS
      To complete setup, run:
        dade setup
        
      This will:
        - Configure the HTTPS proxy
        - Trust the local CA certificate
        - Offer to install official templates
        
      For a better UI experience, install gum:
        brew install gum
    EOS
  end
  
  test do
    assert_match "dade v#{version}", shell_output("#{bin}/dade --version")
  end
end
```

### 2. Add Dependencies

New dependencies to add:
- jq (JSON processing)
- caddy (web server)

### 3. Update Caveats

Point users to 'dade setup' for first-time configuration.

### 4. Handle srv Formula

Options:
A) Remove Formula/srv.rb entirely
B) Keep srv.rb but make it install dade with deprecation warning
C) Make srv a dependency of dade (confusing)

Recommend: Option A - remove srv.rb after archiving repository

## Release Process

1. Complete all dade code changes
2. Tag v1.0.0 on dade repo
3. Get SHA256 of tarball
4. Update formula
5. Test: brew install theydontwantyoutovibecode/tap/dade
6. Push formula changes

## Acceptance Criteria

- [ ] Formula updated with new version
- [ ] SHA256 calculated and added
- [ ] Dependencies (jq, caddy) added
- [ ] Caveats point to setup command
- [ ] Test block works
- [ ] srv formula removed or deprecated
- [ ] Installation tested successfully

