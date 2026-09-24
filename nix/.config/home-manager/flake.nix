{
  description = "Home Manager configuration (dev tools)";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixpkgs-unstable";
    home-manager = {
      url = "github:nix-community/home-manager/master";
      inputs.nixpkgs.follows = "nixpkgs";
    };
    opencode = {
      url = "github:anomalyco/opencode/v2.0.16";
    };
  };

  outputs =
    {
      nixpkgs,
      home-manager,
      opencode,
      ...
    }:
    let
      mkHome =
        system: homeDirectory:
        home-manager.lib.homeManagerConfiguration {
          pkgs = import nixpkgs {
            inherit system;

            # Allow no OSS packages
            config.allowUnfree = true;

            config.allowUnfreePredicate = pkg: builtins.elem (nixpkgs.lib.getName pkg) [ "terraform" ];

            overlays = [
              (final: prev: {
                opencode = opencode.packages.${system}.opencode.overrideAttrs (old: {
                  postInstall = ''
                    installShellCompletion --cmd opencode \
                      --bash <($out/bin/opencode --completions bash) \
                      --zsh <(SHELL=/bin/zsh $out/bin/opencode --completions zsh)

                    installShellCompletion --cmd opencode2 \
                      --bash <($out/bin/opencode2 --completions bash) \
                      --zsh <(SHELL=/bin/zsh $out/bin/opencode2 --completions zsh)
                  '';
                });
              })
            ];
          };
          modules = [
            ./home.nix
            { home.homeDirectory = homeDirectory; }
          ];
        };
    in
    {
      homeConfigurations."tmojzes" = mkHome "aarch64-linux" "/home/tmojzes";
      homeConfigurations."tmojzes-x86_64-linux" = mkHome "x86_64-linux" "/home/tmojzes";
      homeConfigurations."tmojzes-mac" = mkHome "aarch64-darwin" "/Users/tmojzes";
    };
}
