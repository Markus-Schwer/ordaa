{ pkgs }:
with pkgs; rec {
  bin = buildGoModule {
    pname = "malenalau";
    version = "1.0.0";
    # run with fake sha first and then copy actual sha from error message
    # vendorSha256 = lib.fakeSha256;
    vendorSha256 = "sha256-HKbSuuMRSbcyZ8zn5OXbPFpwT+MMo2MF9wcHYJIv4hE=";
    src = ./.;
  };
  container = dockerTools.buildLayeredImage {
    name = "malenalau";
    tag = "latest";
    # contents = pkgs.cacert;
    config = {
      Cmd = [ "${bin}/bin/malenalau" ];
      Env = [
        "GALACTUS_ADDRESS=0.0.0.0"
        "GALACTUS_PORT=80"
      ];
    };
  };

}
