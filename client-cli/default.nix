{
  callPackage,
  fzf,
  lib,
  makeBinaryWrapper,
  mpv,
  stdenv,
}:

let
  base = callPackage ../build-go-sub-package.nix {
    subPackage = "client-cli";
    pname = "bop";
    version = "0.0.1";
  };
in
stdenv.mkDerivation {
  pname = "bop";
  version = "0.0.1";
  nativeBuildInputs = [
    base
    makeBinaryWrapper
  ];
  dontUnpack = true;
  installPhase = ''
    mkdir -p $out/bin
    makeBinaryWrapper ${base}/bin/bop $out/bin/bop \
      --prefix PATH : ${
        lib.makeBinPath [
          mpv
          fzf
        ]
      }
  '';
}
