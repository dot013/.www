{
  description = "learning.rs";
  inputs = {
    nixpkgs.url = "github:nixos/nixpkgs/nixos-unstable";
    flake-utils.url = "github:numtide/flake-utils";
    templ.url = "github:a-h/templ?ref=tags/v0.2.697";
  };
  outputs =
    { self
    , nixpkgs
    , flake-utils
    , ...
    } @ inputs:
    flake-utils.lib.eachDefaultSystem (system:
    let
      pkgs = import nixpkgs {
        inherit system;
      };
      templ = inputs.templ.packages.${system}.templ;
    in
    {
      devShells.default = pkgs.mkShell {
        buildInputs = with pkgs; [
          air
          go
          golangci-lint
          templ
          nodePackages_latest.vercel
        ];
      };
    });
}
