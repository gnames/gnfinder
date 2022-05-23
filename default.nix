{ lib, buildGoModule, fetchFromGitHub, stdenv, glibc }:

buildGoModule rec {
  pname = "gnfinder";
  version = "v1.0.0-RC1";
  date = "2022-05-23";

  src = ./.;

  vendorSha256 = "sha256-+EpoTLIGTdWfEhgpwj/QHnj+HPnKBbZaw5Shy+4jxzs=";

  buildInputs = [
    stdenv
    glibc.static
  ];

  doChecks = false;

  subPackages = "gnfinder";

  ldflags = [
    "-s"
    "-w"
    "-linkmode external"
    "-extldflags"
    "-static"
    "-X github.com/gnames/gnfinder.Version=${version}"
    "-X github.com/gnames/gnfinder.Build=${date}"
  ];

  meta = with lib; {
    description = "Scientific names detection in documents.";
    homepage = "https://github.com/gnames/gnfinder";
    license = licenses.mit;
    maintainers = with maintainers; [ "dimus" ];
  };
}
