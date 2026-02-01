{
  description = "A beautiful TUI for searching NixOS packages with vim keybindings";

  inputs = {
    nixpkgs.url = "github:nixos/nixpkgs?ref=nixos-unstable";
    flake-utils.url = "github:numtide/flake-utils";
  };

  outputs = { self, nixpkgs, flake-utils }:
    flake-utils.lib.eachDefaultSystem (system:
      let
        pkgs = import nixpkgs { inherit system; };
        version = "0.1.0";
      in
      {
        packages = {
          default = pkgs.buildGoModule {
            pname = "ns-tui";
            inherit version;
            src = ./.;

            # vendorHash needs to be updated after go.mod changes
            # Set to null first, then nix will tell you the correct hash
            vendorHash = "sha256-47V2qVlRacRSqNKpsnp9LL2nSlEgSFTWJUKX/jUs0d8=";

            # Exclude vendor directory
            proxyVendor = true;

            subPackages = [ "cmd/ns-tui" ];

            ldflags = [
              "-s"
              "-w"
              "-X main.version=${version}"
            ];

            meta = with pkgs.lib; {
              description = "A beautiful TUI for searching NixOS packages with vim keybindings";
              homepage = "https://github.com/briheet/ns-tui";
              license = licenses.mit;
              maintainers = [ ];
              mainProgram = "ns-tui";
            };
          };

          ns-tui = self.packages.${system}.default;
        };

        apps = {
          default = {
            type = "app";
            program = "${self.packages.${system}.default}/bin/ns-tui";
          };

          ns-tui = self.apps.${system}.default;
        };

        devShells.default = pkgs.mkShell {
          buildInputs = with pkgs; [
            go
            gopls
            gotools
            go-tools
            golangci-lint
            delve
            go-task
          ];

          shellHook = ''
            echo "🚀 ns-tui development environment"
            echo "Go version: $(go version)"
            echo ""
            echo "Available commands:"
            echo "  task build  - Build the application"
            echo "  task run    - Run the application"
            echo "  task dev    - Run in development mode"
            echo "  task test   - Run tests"
            echo ""
          '';
        };
      }
    );
}
