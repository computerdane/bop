{ callPackage }:

callPackage ../build-go-sub-package.nix {
  subPackage = "server-api";
  pname = "bop-api";
  version = "0.0.1";
}
