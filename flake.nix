{
  description = "Automate restaurant orders";

  inputs = {
    nixpkgs.url = "github:nixos/nixpkgs/nixos-unstable";
    treefmt-nix = {
      url = "github:numtide/treefmt-nix";
      inputs.nixpkgs.follows = "nixpkgs";
    };
    rust-overlay = {
      url = "github:oxalica/rust-overlay";
      inputs.nixpkgs.follows = "nixpkgs";
    };
  };

  outputs = { self, nixpkgs, treefmt-nix, rust-overlay }:
    let
      overlays = [ (import rust-overlay) ];
      pkgs = import nixpkgs {
        inherit overlays;
        system = "x86_64-linux";
      };
      rustVersion = pkgs.rust-bin.stable.latest.default;
      rustPlatform = pkgs.makeRustPlatform {
        cargo = rustVersion;
        rustc = rustVersion;
      };
      dotinder = rustPlatform.buildRustPackage {
        pname = "dotinder";
        version = "0.0.1";
        src = ./.;
        cargoLock.lockFile = ./Cargo.lock;
      };
      treefmtEval = treefmt-nix.lib.evalModule pkgs ./treefmt.nix;
    in
    {
      devShells.x86_64-linux.default = pkgs.mkShell {
        packages = with pkgs;[
          rustc
          cargo
        ];
      };
      formatter.x86_64-linux = treefmtEval.config.build.wrapper;
      checks.x86_64-linux.formatter = treefmtEval.config.build.check self;
      packages.x86_64-linux.default = dotinder;
    };
}
