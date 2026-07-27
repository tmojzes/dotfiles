{
  description = "Home Manager configuration (dev tools)";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
    home-manager = {
      url = "github:nix-community/home-manager/master";
      inputs.nixpkgs.follows = "nixpkgs";
    };
  };

  outputs =
    { nixpkgs, home-manager, ... }:
    let
      mkHome =
        system: homeDirectory:
        home-manager.lib.homeManagerConfiguration {
          pkgs = import nixpkgs {
            inherit system;
            config.allowUnfreePredicate = pkg: builtins.elem (nixpkgs.lib.getName pkg) [ "terraform" ];
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
