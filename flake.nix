{
  inputs = {
    nixpkgs.url = "github:nixos/nixpkgs/nixos-24.05";
    utils.url = "github:numtide/flake-utils";
  };
  outputs =
    { nixpkgs, utils, ... }:
    utils.lib.eachDefaultSystem (
      system:
      let
        pkgs = import nixpkgs { inherit system; };
        pkgs-bop = {
          client-cli = pkgs.callPackage ./client-cli/default.nix { };
          server-api = pkgs.callPackage ./server-api/default.nix { };
        };
      in
      {
        devShell = pkgs.mkShell {
          buildInputs = (
            with pkgs;
            [
              buf-language-server
              fzf
              git
              go
              gopls
              mpv
              protobuf
              protoc-gen-go
              protoc-gen-go-grpc
            ]
          );
        };
        packages = pkgs-bop // {
          default = pkgs-bop.client-cli;
        };
      }
    );
}
