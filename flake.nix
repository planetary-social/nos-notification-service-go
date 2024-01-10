{
  description = "Flake for developing on nos-notification-service";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
    flake-utils.url = "github:numtide/flake-utils";
    devshell.url = "github:numtide/devshell";
  };

  outputs = { self, nixpkgs, flake-utils, devshell }:
    flake-utils.lib.eachDefaultSystem (system:
    let pkgs = import nixpkgs {
          inherit system;
          overlays = [
            devshell.overlays.default
          ];
        };
    in
      {
        devShells.default =
          pkgs.devshell.mkShell {
            name = "nos-notification-service";
            packages = with pkgs; [
              go
              gopls
              golangci-lint
              gotools
              go-tools
              gosimports
            ];
            commands = [
              {
                name = ''start'';
                help = ''alias for `go run cmd/notification-service/main.go`'';
                category = ''dev helpers'';
                command = ''
                go run cmd/notification-service/main.go
                '';
              }
              {
                name = ''lint'';
                help = ''check style and formatting'';
                category = ''dev helpers'';
                command = ''
                go vet ./...
                golangci-lint run ./...
                '';
              }

            ];
          };
      });
}