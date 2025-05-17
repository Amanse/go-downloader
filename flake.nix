{
  description = "A very basic flake";

  inputs = {
    nixpkgs.url = "github:nixos/nixpkgs?ref=nixos-unstable";
    utils.url = "github:numtide/flake-utils";
    gomod2nix = {
      url = "github:tweag/gomod2nix";
      inputs.nixpkgs.follows = "nixpkgs";
    };
  };

  outputs = { self, nixpkgs, utils, gomod2nix }:
    let
      pkgs = import nixpkgs {
        system = "x86_64-linux";
        overlays = [ gomod2nix.overlays.default ];
      };
    in {
      packages.x86_64-linux.default = pkgs.buildGoApplication {
        name = "go-downloader";
        version = "1";
        modules = ./gomod2nix.toml;
        src = ./.;
      };
      devShells.x86_64-linux.default = pkgs.mkShell {
        buildInputs = with pkgs; [
          go
          gopls
          gomod2nix.packages.${system}.default
        ];
      };
    };
}
