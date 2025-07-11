class Clify < Formula
  desc "Command-line assistant for Linux, macOS, and Windows powered by Anthropic Claude"
  homepage "https://github.com/aktagon/clify"
  # NOTE: The url, version, and sha256 are updated by the github action (.github/workflows/release.yml) automatically
  url "https://github.com/aktagon/clify/archive/refs/tags/v0.1.2.tar.gz"
  version "v0.1.2"
  sha256 "4416ce05aa8d239af9006d8342905dbeafa223ea395693b2728d2c9afed3e615"
  license "MIT"

  depends_on "go" => :build

  def install
    system "go", "build", *std_go_args(ldflags: "-s -w"), "-o", bin/"clify"
  end

  def caveats
    <<~EOS
      clify requires an Anthropic API key to function.

      Run clify setup or set your API key as an environment variable:
        export ANTHROPIC_API_KEY="your-api-key-here"

      You can add this to your shell profile (~/.zshrc, ~/.bashrc, etc.)
      to make it permanent.
    EOS
  end

  test do
    # Test that the binary was installed and can display help
    assert_match "clify", shell_output("#{bin}/clify --help 2>&1", 1)
  end
end
