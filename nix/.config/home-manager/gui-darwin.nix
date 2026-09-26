{ ... }:

{
  homebrew = {
    enable = true;
    casks = [
      "bitwarden"
      "discord"
      "ghostty"
      "google-chrome"
      "slack"
      "spotify"
      "tailscale"
    ];
    cleanup = false; # do not delete unmanaged brew casks
    enableShellIntegration = false; # shell package already handles brew environment
  };
}
