{
  description = "Legifrss server flake";

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
              cfg = config.services.legifrss;
              pkg = self.packages.${system}.default;
            in
            {
              options.services.legifrss = {
                enable = mkEnableOption "Enable legifrss service";
                envFile = mkOption { type = types.str; };
              };
              config = mkIf cfg.enable {
                services.nginx.virtualHosts."legifrss.org" = {
                  enableACME = true;
                  forceSSL = true;
                  locations."/" = {
                    proxyPass = "http://127.0.0.1:8080/";
                  };
                };
                users.groups = { legifrss = { }; };

                users.users.legifrss = {
                  group = "legifrss";
                  isNormalUser = true;
                };

                systemd.services.legifrss = {
                  description = "Legifrss server";
                  wantedBy = [ "multi-user.target" ];
                  environment = { };
                  serviceConfig = {
                    User = "legifrss";
                    Group = "legifrss";
                    #  DynamicUser = "yes";
                    ExecStart = "${pkg}/bin/server";
                    Restart = "on-failure";
                    RestartSec = "10s";
                    WorkingDirectory = "/home/legifrss";
                    ReadWritePaths = [ "/home/legifrss" ];
                  };
                };

                systemd.services.legifrss-batch = {
                  description = "Legifrss batch";
                  wantedBy = [ "multi-user.target" ];
                  environment = {
                    ENV_FILE = "${config.services.legifrss.envFile}";
                  };
                  serviceConfig = {
                    User = "legifrss";
                    Group = "legifrss";
                    #  DynamicUser = "yes";
                    ExecStart = "${pkg}/bin/batch";
                    Restart = "no";
                    WorkingDirectory = "/home/legifrss";
                    ReadWritePaths = [ "/home/legifrss" ];
                  };
                };

                systemd.timers = {
                  legifrss-batch-timer = {
                    # Unit = {
                    #   Description = "Fetch Legifrance updates";
                    #   After = [ "network.target" ];
                    # };
                    timerConfig = {
                      OnBootSec = "5 min";
                      OnUnitInactiveSec = "60 min";
                      Unit = "legifrss-batch.service";
                    };
                  };
                };
              };
            };
        })
    );
}

