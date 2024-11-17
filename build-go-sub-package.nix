{
  buildGoModule,
  buildInputs ? [ ],
  buildPackages,
  installShellFiles,
  pname,
  stdenv,
  subPackage,
  version,
}:

buildGoModule {
  inherit pname version buildInputs;
  src = ./.;
  subPackages = [ subPackage ];
  nativeBuildInputs = [ installShellFiles ];
  doCheck = false;
  postInstall =
    let
      emulator = stdenv.hostPlatform.emulator buildPackages;
    in
    ''
      mv $out/bin/${subPackage} $out/bin/${pname}
      installShellCompletion --cmd ${pname} \
        --bash <(${emulator} $out/bin/${pname} completion bash) \
        --fish <(${emulator} $out/bin/${pname} completion fish) \
        --zsh  <(${emulator} $out/bin/${pname} completion zsh)
    '';
  vendorHash = "sha256-1ZiskvF401IQ9qzNHCJG2F3TohNxl3+zEv+cU1kbBWo=";
}
