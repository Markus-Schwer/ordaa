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
        omega-star-dev = pkgs.writeScriptBin "omega-star-dev" ''${omega-star.bin}/bin/omega-star -address 127.0.0.1:8080 -v'';
        # omega-star = omega-star.container;
        galactus-bin = galactus.bin;
        galactus-dev = pkgs.writeScriptBin "galactus-dev" ''${galactus.bin}/bin/galactus -address 127.0.0.1:8081 -omega-star http://localhost:8080 -v'';
        # galactus = galactus.container;
        malenalau-bin = malenalau.bin;
        malenalau-dev = pkgs.writeScriptBin "malenalau-dev" ''${malenalau.bin}/bin/malenalau -v -room '!hvJGXMkjcyzxtSNNsx:matrix.org' -user order-bot-aa -password-file ./secret.txt'';
        # malenalau = malenalau.container;
      };
    };
}
