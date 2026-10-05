{ ... }:

{
  homebrew = {
    enable = true;
    casks = [
      "antigravity"
      "bitwarden"
      "discord"
      "drawio"
      "ghostty"
      "gimp"
      "google-chrome"
      "google-drive"
      "macdown"
      "monitorcontrol"
      "plex"
      "rectangle"
      "slack"
      "spotify"
      "tunnelblick"
      "utm"
      "visual-studio-code"
      "wifiman"
      "windows-app"
      "xquartz"
      "zen"
    ];
    cleanup = false; # do not delete unmanaged brew casks
    enableShellIntegration = false; # shell package already handles brew environment
  };
}
