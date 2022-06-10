class Tubekit < Formula
  desc "Tool that helps you to operate Kubernetes clusters more effectively"
  homepage "https://github.com/reconquest/tubekit"
  url "https://github.com/reconquest/tubekit/releases/download/v4/tubekit_4_Darwin_x86_64.tar.gz"
  sha256 "10256142d95e0fe8879a1c46cf85d87d7943499ba3878afc30779a4125ec53d5"

  depends_on "kubernetes-cli"

  def install
    bin.install "tubectl"
  end

  test do
    system "#{bin}/tubectl", "--tube-version"
    assert_match /usage/i, shell_output("#{bin}/tubectl --tube-help")
  end
end
