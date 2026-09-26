{ ... }:

{
  services.flatpak = {
    enable = true;
    remotes = [
      {
        name = "flathub";
        location = "https://dl.flathub.org/repo/flathub.flatpakrepo";
      }
    ];
    packages = [
      "com.bitwarden.desktop"
      "com.google.Chrome"
      "com.spotify.Client"
      "dev.vencord.Vesktop"
      "com.slack.Slack"
    ];
  };
}
