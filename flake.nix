{
  description = "Automate restaurant orders";

  inputs = {
    nixpkgs.url = "github:nixos/nixpkgs/nixos-unstable";
    treefmt-nix = {
      url = "github:numtide/treefmt-nix";
      inputs.nixpkgs.follows = "nixpkgs";
    };
    devenv.url = "github:cachix/devenv";
  };

  outputs = { self, nixpkgs, treefmt-nix, devenv, ... } @ inputs:
    let
      system = "x86_64-linux";
      pkgs = import nixpkgs { inherit system; };
      treefmtEval = treefmt-nix.lib.evalModule pkgs ./treefmt.nix;
    in
    {
      devShells.${system}.default = devenv.lib.mkShell {
        inherit inputs pkgs;
        modules = [
          {
            dotenv.disableHint = true;
            languages.go.enable = true;
            languages.javascript.enable = true;
            packages = with pkgs; [ go-migrate reflex gcc nodePackages.svelte-language-server delve ];

            env.DATABASE_URL = "postgresql:///ordaa";
            env.ADDRESS = "localhost:8080";
            env.CGO_ENABLED = "0";

            services.postgres = {
              enable = true;
              package = pkgs.postgresql_15;
              initialDatabases = [{ name = "ordaa"; }];
            };

            scripts.dev-server.exec = ''
              reflex -r '\.go$' -s go run main.go
            '';
          }
        ];
      };
      formatter.${system} = treefmtEval.config.build.wrapper;
      checks.${system}.formatter = treefmtEval.config.build.check self;
      packages.${system} = {
        devenv-up = self.devShells.${system}.default.config.procfileScript;
        default = pkgs.buildGoModule {
          pname = "ordaa";
          version = "1.0.0";
          # run with fake hash first and then copy actual hash from error message
          #vendorHash = nixpkgs.lib.fakeHash;
          vendorHash = "sha256-EEF+WoyClJjLTCeLwpRKX1GJ+wLSW/ShvzhXShRhBNs=";
          src = ./.;
        };
        html = pkgs.buildNpmPackage {
          pname = "fontend";
          version = "1.0.0";
          src = html/.;
          npmDepsHash = "sha256-wiBI0HfLlddZsVduJgy5ax3RCS1lzL3o6Q1ccK3+HEI=";
          installPhase = ''
            mkdir $out
            cp -r build/* $out/
          '';
        };
      };
    };
}
