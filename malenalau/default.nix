{ pkgs }:
with pkgs; rec {
  bin = buildGoModule {
    pname = "malenalau";
    version = "1.0.0";
    # run with fake sha first and then copy actual sha from error message
    # vendorSha256 = lib.fakeSha256;
    vendorSha256 = "sha256-k+5wnl6DK5VW/CbqVOG/hknyTt0Q7dygovm+v6rXIN8=";
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
