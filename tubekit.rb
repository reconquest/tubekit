class Tubekit < Formula
  desc "Tool that helps you to operate Kubernetes clusters more effectively"
  homepage "https://github.com/reconquest/tubekit"
  url "https://github.com/reconquest/tubekit/releases/download/v2/tubekit_2_Darwin_x86_64.tar.gz"
  sha256 "e64ce03124f6f3ab7675d5f4d5f721add6fab0184db318a3dc24da0fd9ed2f92"
  head "https://github.com/reconquest/tubekit.git"

  bottle :unneeded

  depends_on "kubernetes-cli"

  def install
    bin.install "tubectl"
  end
end
