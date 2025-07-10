class Howdoi < Formula
  desc "Command-line assistant for Linux, macOS, and Windows powered by Anthropic Claude"
  homepage "https://github.com/aktagon/llmkit"
  url "https://github.com/aktagon/llmkit/archive/refs/heads/master.tar.gz"
  version "0.1.0"
  sha256 "0000000000000000000000000000000000000000000000000000000000000000"
  license "MIT"

  depends_on "go" => :build

  def install
    cd "examples/anthropic/clify" do
      system "go", "build", *std_go_args(ldflags: "-s -w"), "-o", bin/"clify"
    end
  end

  def caveats
    <<~EOS
      clify requires an Anthropic API key to function.
      Set your API key as an environment variable:
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