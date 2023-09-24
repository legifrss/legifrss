{
  description = "A basic gomod2nix flake";

  inputs.nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
  inputs.flake-utils.url = "github:numtide/flake-utils";
  inputs.gomod2nix.url = "github:nix-community/gomod2nix";
  inputs.gomod2nix.inputs.nixpkgs.follows = "nixpkgs";
  inputs.gomod2nix.inputs.flake-utils.follows = "flake-utils";

  outputs = { self, nixpkgs, flake-utils, gomod2nix }:
    (flake-utils.lib.eachDefaultSystem
      (system:
        let
          pkgs = nixpkgs.legacyPackages.${system};
          inherit (pkgs) stdenv lib;
        in
        rec {

          packages.legifrss = pkgs.callPackage ./. {
            inherit (gomod2nix.legacyPackages.${system}) buildGoApplication;
          };
          packages.default = packages.legifrss;
          devShells.default = pkgs.callPackage ./shell.nix {
            inherit (gomod2nix.legacyPackages.${system}) mkGoEnv gomod2nix;
          };
          nixosModules.default = { config, lib, pkgs, ... }:
            with lib;
            let

            in
            {
              services.nginx.virtualHosts."legifrss.org" = {
                enableACME = true;
                forceSSL = true;
                root = "${packages.doc}";
              };
            };
        })
    );
}
