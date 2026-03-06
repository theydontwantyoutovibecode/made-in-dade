---
id: jus-uu3b
status: open
deps: [jus-8hqk]
links: []
created: 2026-02-12T20:36:42Z
type: task
priority: 1
assignee: Alex Cabrera
parent: jus-ni16
tags: [homebrew, distribution]
---
# Update Homebrew formula

Update the Homebrew formula in theydontwantyoutovibecode/homebrew-tap.

## File

/Users/acabrera/Code/homebrew-tap/Formula/dade.rb

## Updated Formula

```ruby
class Justvibin < Formula
  desc "CLI for scaffolding and serving web projects with automatic HTTPS"
  homepage "https://github.com/theydontwantyoutovibecode/dade"
  url "https://github.com/theydontwantyoutovibecode/dade/archive/refs/tags/v1.0.0.tar.gz"
  sha256 "<calculate-after-release>"
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

      This will configure the HTTPS proxy and trust certificates.

      For a better UI experience:
        brew install gum

      For public tunnels:
        brew install cloudflared
    EOS
  end

  test do
    assert_match "dade v#{version}", shell_output("#{bin}/dade --version")
  end
end
```

## Steps

1. Wait for v1.0.0 release
2. Download tarball and calculate SHA256
3. Update formula
4. Test: brew install --build-from-source
5. Push to homebrew-tap

## Acceptance Criteria

- [ ] Formula updated with v1.0.0 URL
- [ ] SHA256 calculated correctly
- [ ] Dependencies added (jq, caddy)
- [ ] Caveats updated
- [ ] Installation tested

