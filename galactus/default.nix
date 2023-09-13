{ pkgs }:
with pkgs; rec {
  bin = buildGoModule {
    pname = "galactus";
    version = "1.0.0";
    # run with fake sha first and then copy actual sha from error message
    # vendorSha256 = lib.fakeSha256;
    vendorSha256 = "sha256-sGSBPztiLa0Ngq8zHIZUoqeQJ0CYivDcJl5/fnhZ/+0=";
    src = ./.;
  };
  container = dockerTools.buildLayeredImage {
    name = "galactus";
    tag = "latest";
    # contents = pkgs.cacert;
    config = {
      Cmd = [ "${bin}/bin/galactus" ];
      Env = [
        "GALACTUS_ADDRESS=0.0.0.0"
        "GALACTUS_PORT=80"
      ];
    };
  };

}
