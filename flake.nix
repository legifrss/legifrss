{
  inputs = {
    crane = {
      url = "github:ipetkov/crane";
    };
    fenix = {
      url = "github:nix-community/fenix";
      inputs.nixpkgs.follows = "nixpkgs";
    };
    flake-utils.url = "github:numtide/flake-utils";
    nixpkgs.url = "nixpkgs/nixos-unstable";
  };

  outputs =
    {
      self,
      crane,
      fenix,
      flake-utils,
      nixpkgs,
    }:
    flake-utils.lib.eachDefaultSystem (
      system:
      let
        pkgs = nixpkgs.legacyPackages.${system};
        lib = pkgs.lib;
        fenixToolchain = (
          with fenix.packages.${system};
          combine [
            stable.toolchain
            stable.cargo
            stable.rustc
            stable.llvm-tools
            stable.rust-std
            stable.rust-src
            stable.clippy
          ]
        );
        craneLib = (crane.mkLib nixpkgs.legacyPackages.${system}).overrideToolchain fenixToolchain;
        commonArgs = {
          src = lib.cleanSourceWith {
            src = ./.;
            filter =
              path: type:
              (lib.hasInfix "tests/" path)
              || (lib.hasInfix "assets/" path)
              || (lib.hasInfix "sql/" path)
              || (lib.hasInfix "migrations/" path)
              || (craneLib.filterCargoSources path type);
          };
          buildInputs = with pkgs; [
            pkg-config
            openssl
          ];
          LD_LIBRARY_PATH = lib.makeLibraryPath [ pkgs.openssl ];
          RUSTFLAGS = "-C linker-features=-lld";
        };
        cargoArtifacts = craneLib.buildDepsOnly (
          commonArgs
          // {
            pname = "legifrss";
          }
        );

        legifrss = craneLib.buildPackage (
          commonArgs
          // {
            pname = "legifrss";
            inherit cargoArtifacts;
            # https://crane.dev/examples/sqlx.html
            preBuild = ''
              export DATABASE_URL=postgresql://postgres:postgres@localhost/test-db
              ${pkgs.postgresql}/bin/initdb -D .tmp/test-db
              ${pkgs.postgresql}/bin/pg_ctl -D .tmp/test-db -l .tmp/test-db.log -o "--unix_socket_directories='$PWD'" start
              ${pkgs.postgresql}/bin/createuser postgres -d -s -h $PWD
              ${pkgs.postgresql}/bin/createdb test-db -h $PWD --owner=postgres
              ${pkgs.sqlx-cli}/bin/sqlx migrate run
            '';

            postInstall = ''
              cp -R migrations $out/
            '';

          }
        );
      in
      {
        packages.default = legifrss;

        devShells.default = pkgs.mkShell {
          RUST_SRC_PATH = "${fenix.packages.${system}.stable.rust-src}/lib/rustlib/src/rust/library";
          LD_LIBRARY_PATH = lib.makeLibraryPath [ pkgs.openssl ];
          buildInputs = with pkgs; [
            pkg-config
            fenixToolchain
            openssl
            poppler
            poppler-utils
            postgresql
            gnumeric
            tailwindcss_4
            nodejs
            sqlx-cli
            (writeShellScriptBin "stop_all" ''
              #!/usr/bin/env bash
              set -e
              docker-compose down -v
            '')
            (writeShellScriptBin "restart_all" ''
              #!/usr/bin/env bash
              set -e
              docker-compose down -v
              docker-compose up -d
              sleep 2
              ${sqlx-cli}/bin/sqlx migrate run
            '')
          ];

          env = {
            DATABASE_URL = "postgres://legifrss:legifrss@localhost/legifrss";
          };

        };
      }
    );

}
