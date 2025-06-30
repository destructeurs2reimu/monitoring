{
  inputs = {
    nixpkgs.url = "github:nixos/nixpkgs?ref=nixos-unstable";
    flake-utils.url = "github:numtide/flake-utils";
  };

  outputs = { self, nixpkgs, flake-utils }: 
    flake-utils.lib.eachDefaultSystem (system:
      let
        pkgs = import nixpkgs { inherit system; };
      in {
        devShell = pkgs.mkShell{
          buildInputs = with pkgs; [
            go_1_24
          ];
        };

        packages.default = pkgs.buildGoModule {
          name = "monitoring";
          src = ./.;
          vendorHash = null;
        };
      }
    );
}
