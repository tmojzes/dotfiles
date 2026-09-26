{
  description = "Home Manager configuration (dev tools)";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixpkgs-unstable";
    home-manager = {
      url = "github:nix-community/home-manager/master";
      inputs.nixpkgs.follows = "nixpkgs";
    };
    nix-flatpak.url = "github:gmodena/nix-flatpak";
    home-manager-brew = {
      url = "github:koalalorenzo/home-manager-brew";
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
      nix-flatpak,
      home-manager-brew,
      opencode,
      ...
    }:
    let
      mkHome =
        system: homeDirectory: extraModules:
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
          ] ++ extraModules;
        };
    in
    {
      homeConfigurations."tmojzes" = mkHome "aarch64-linux" "/home/tmojzes" [
        nix-flatpak.homeManagerModules.nix-flatpak
        ./gui.nix
      ];
      homeConfigurations."tmojzes-x86_64-linux" = mkHome "x86_64-linux" "/home/tmojzes" [
        nix-flatpak.homeManagerModules.nix-flatpak
        ./gui.nix
      ];
      homeConfigurations."tmojzes-mac" = mkHome "aarch64-darwin" "/Users/tmojzes" [
        home-manager-brew.homeManagerModules.default
        ./gui-darwin.nix
      ];
    };
}
