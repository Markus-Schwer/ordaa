{ pkgs }:
with pkgs; rec {
  bin = buildGoModule {
    pname = "omega-star";
    version = "1.0.0";
    # run with fake sha first and then copy actual sha from error message
    /* vendorSha256 = lib.fakeSha256; */
    vendorSha256 = "sha256-cnrL8Y6AFkKPykNfIlidwvSjic/psNBimbqqhD7gKlg=";
    src = ./.;
  };
  container = dockerTools.buildLayeredImage {
    name = "omega-star";
    tag = "latest";
    contents = pkgs.cacert;
    config = {
      Cmd = [ "${bin}/bin/omega-star" ];
      Env = [
        "OMEGA_STAR_ADDRESS=0.0.0.0"
        "OMEGA_STAR_PORT=80"
      ];
    };
  };

}
