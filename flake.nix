{
  description = "Automate restaurant orders";

  inputs = {
    nixpkgs.url = "github:nixos/nixpkgs/nixos-unstable";
    treefmt-nix = {
      url = "github:numtide/treefmt-nix";
      inputs.nixpkgs.follows = "nixpkgs";
    };
    devenv.url = "github:cachix/devenv";
    fenix = {
      url = "github:nix-community/fenix";
      inputs.nixpkgs.follows = "nixpkgs";
    };
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
            languages.rust.enable = true;
            # https://devenv.sh/reference/options/#languagesrustchannel
            languages.rust.channel = "stable";
            packages = with pkgs; [ sqlx-cli ];
            env.DATABASE_URL = "postgresql:///dotinder";

            services.postgres = {
              enable = true;
              package = pkgs.postgresql_15;
              initialDatabases = [{ name = "dotinder"; }];
            };
          }
        ];
      };
      formatter.${system} = treefmtEval.config.build.wrapper;
      checks.${system}.formatter = treefmtEval.config.build.check self;
      packages.${system}.devenv-up = self.devShells.${system}.default.config.procfileScript;
    };
}
