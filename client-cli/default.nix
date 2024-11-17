{ callPackage, mpv }:

callPackage ../build-go-sub-package.nix {
  subPackage = "client-cli";
  pname = "bop";
  version = "0.0.1";
  buildInputs = [ mpv ];
}
