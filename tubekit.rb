class Tubekit < Formula
  desc "Tool that helps you to operate Kubernetes clusters more effectively"
  homepage "https://github.com/reconquest/tubekit"

  if OS.mac?
    url "https://github.com/reconquest/tubekit/releases/download/v4/tubekit_4_Darwin_x86_64.tar.gz"
    sha256 "10256142d95e0fe8879a1c46cf85d87d7943499ba3878afc30779a4125ec53d5"
  end

  if OS.linux?
    url "https://github.com/reconquest/tubekit/releases/download/v4/tubekit_4_Linux_x86_64.tar.gz"
    sha256 "40891525fd0f799b54b66ed8ff683ddb26a78d742271ad523fca3679592678a7"
  end

  depends_on "kubernetes-cli"

  def install
    bin.install "tubectl"
  end

  test do
    system "#{bin}/tubectl", "--tube-version"
    assert_match /usage/i, shell_output("#{bin}/tubectl --tube-help")
  end
end
