{ pkgs }:
with pkgs; rec {
  bin = buildGoModule {
    pname = "galactus";
    version = "1.0.0";
    # run with fake sha first and then copy actual sha from error message
    # vendorSha256 = lib.fakeSha256;
    vendorSha256 = "sha256-l3A2/iFFbH8nqVBmmwXXtRA73/O5TdrzOkxtVpuPgbA=";
    src = ./.;
  };
  container = dockerTools.streamLayeredImage {
    name = "galactus";
    tag = "latest";
    # contents = pkgs.cacert;
    config = {
      Cmd = [ "${bin}/bin/galactus" ];
      Env = [
        "ADDRESS=0.0.0.0:80"
      ];
    };
  };

}
