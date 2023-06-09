{
  description = "Automate restaurant orders";

  inputs = {
    nixpkgs.url = "github:nixos/nixpkgs/nixos-unstable";
    treefmt-nix = {
      url = "github:numtide/treefmt-nix";
      inputs.nixpkgs.follows = "nixpkgs";
    };
  };

  outputs = { self, nixpkgs, treefmt-nix }:
    let
      pkgs = import nixpkgs {
        system = "x86_64-linux";
      };
      dotinder = pkgs.buildNpmPackage {
        pname = "dotinder-js-server";
        version = "0.0.1";
        src = ./.;
        npmDepsHash = "sha256-9FyRiZRB6IZihUzRQmCX+O3m1qvpN6T8EZ9iV3YSlMM=";
        installPhase = ''
          mkdir -p $out/bin
          mv build/src/* $out/bin
          mv node_modules $out
          mv package.json $out
          mv package-lock.json $out
        '';
      };
      wrap = pkgs.writeScriptBin "dotinder" ''
        ${pkgs.nodejs_18}/bin/node ${self.packages.x86_64-linux.default}/bin/app.js
      '';
    in
    {
      devShells.x86_64-linux.default = pkgs.mkShell {
        packages = with pkgs;[
          nodejs_18
        ];
      };
      formatter.x86_64-linux = treefmt-nix.lib.mkWrapper
        nixpkgs.legacyPackages.x86_64-linux
        {
          projectRootFile = "flake.nix";
          programs.nixpkgs-fmt.enable = true;
          programs.prettier.enable = true;
          settings.formatter.prettier.excludes = [ "*.html" ];
        };
      packages.x86_64-linux.default = dotinder;
      apps.x86_64-linux.default = {
        type = "app";
        program = "${wrap}/bin/dotinder";
      };
    };
}
