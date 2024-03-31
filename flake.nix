{
  description = "Automate restaurant orders";

  inputs = {
    nixpkgs.url = "github:nixos/nixpkgs/nixos-unstable";
    treefmt-nix = {
      url = "github:numtide/treefmt-nix";
      inputs.nixpkgs.follows = "nixpkgs";
    };
    devenv.url = "github:cachix/devenv";
    templ.url = "github:a-h/templ";
  };

  outputs = { self, nixpkgs, treefmt-nix, devenv, templ, ... } @ inputs:
    let
      system = "x86_64-linux";
      pkgs = import nixpkgs { inherit system; };
      treefmtEval = treefmt-nix.lib.evalModule pkgs ./treefmt.nix;
      templ = system: inputs.templ.packages.${system}.templ;
    in
    {
      devShells.${system}.default = devenv.lib.mkShell {
        inherit inputs pkgs;
        modules = [
          {
            dotenv.disableHint = true;
            languages.go.enable = true;
            packages = with pkgs; [ go-migrate (templ system) reflex ];
            env.DATABASE_URL = "postgresql:///dotinder";
            env.ADDRESS = "localhost:8080";

            services.postgres = {
              enable = true;
              package = pkgs.postgresql_15;
              initialDatabases = [{ name = "dotinder"; }];
            };

            scripts.dev-server.exec = ''
              reflex -R '_templ.go$' -s -- sh -c 'templ generate && go run main.go'
            '';
          }
        ];
      };
      formatter.${system} = treefmtEval.config.build.wrapper;
      checks.${system}.formatter = treefmtEval.config.build.check self;
      packages.${system} = {
        devenv-up = self.devShells.${system}.default.config.procfileScript;
        default = pkgs.buildGoModule {
          pname = "dotinder";
          version = "1.0.0";
          # run with fake sha first and then copy actual sha from error message
          #vendorSha256 = nixpkgs.lib.fakeSha256;
          vendorSha256 = "sha256-N1jG2uB6K6/T0+jzx6qQw/1S36EgmktV7aBXCKTaKhM=";
          src = ./.;

          preBuild = ''
            ${templ system}/bin/templ generate
          '';
        };
      };
    };
}
