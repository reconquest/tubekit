class Tubekit < Formula
  desc "Tool that helps you to operate Kubernetes clusters more effectively"
  homepage "https://github.com/reconquest/tubekit"
  head "https://github.com/reconquest/tubekit.git"
  url "https://github.com/reconquest/tubekit/releases/download/v3/tubekit_3_Darwin_x86_64.tar.gz"
  sha256 "b966e7b014e0d16e22f0d15dbd7c80a084160c6851bf7ad767d0085ac3ac10d1"

  bottle :unneeded

  depends_on "kubernetes-cli"

  def install
    bin.install "tubectl"
  end
end
