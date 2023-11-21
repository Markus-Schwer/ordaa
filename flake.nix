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
      system = "x86_64-linux";
      pkgs = import nixpkgs { inherit system; };
      treefmtEval = treefmt-nix.lib.evalModule pkgs ./treefmt.nix;
      omega-star = import ./omega-star { inherit pkgs; };
      galactus = import ./galactus { inherit pkgs; };
      malenalau = import ./malenalau { inherit pkgs; };
    in
    {
      devShells.${system}.default = pkgs.mkShell {
        packages = with pkgs;[
          go
        ];
      };
      formatter.${system} = treefmtEval.config.build.wrapper;
      checks.${system}.formatter = treefmtEval.config.build.check self;
      packages.${system} = {
        omega-star-bin = omega-star.bin;
        omega-star-dev = pkgs.writeScriptBin "omega-star-dev" ''
          ADDRESS=127.0.0.1:8080 ${omega-star.bin}/bin/omega-star -v
        '';
        omega-star = omega-star.container;
        galactus-bin = galactus.bin;
        galactus-dev = pkgs.writeScriptBin "galactus-dev" ''
          ADDRESS=127.0.0.1:8081 OMEGA_STAR_URL=http://localhost:8080 ${galactus.bin}/bin/galactus -v
        '';
        galactus = galactus.container;
        malenalau-bin = malenalau.bin;
        malenalau-dev = pkgs.writeScriptBin "malenalau-dev" ''
          export ROOM="!hvJGXMkjcyzxtSNNsx:matrix.org"
          export USER=order-bot-aa
          export HOME_SERVER=matrix.org
          export PASSWORD_FILE="./secret.txt"
          export OMEGA_STAR_URL=http://localhost:8080
          export GALACTUS_URL=http://localhost:8081
          ${malenalau.bin}/bin/malenalau -v
        '';
        malenalau = malenalau.container;
      };
    };
}
