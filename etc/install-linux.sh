#!/bin/bash

if [[ ! -d "$HOME/.local/share/applications" ]]; then
  mkdir -p "$HOME/.local/share/applications"
fi

cat >"$HOME/.local/share/applications/yarr.desktop" <<END
[Desktop Entry]
Name=yarr
Exec=$HOME/.local/bin/yarr -open
Icon=yarr
Type=Application
Categories=Internet;
END

if [[ ! -d "$HOME/.local/share/icons" ]]; then
  mkdir -p "$HOME/.local/share/icons"
fi

cat >"$HOME/.local/share/icons/yarr.svg" <<END
<?xml version="1.0" encoding="UTF-8"?>
<svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="feather feather-anchor-favicon">
  <circle cx="12" cy="5" r="3" stroke-width="4" stroke="#ffffff"></circle>
  <line x1="12" y1="22" x2="12" y2="8" stroke-width="4" stroke="#ffffff"></line>
  <path d="M5 12H2a10 10 0 0 0 20 0h-3" stroke-width="4" stroke="#ffffff"></path>

  <circle cx="12" cy="5" r="3"></circle>
  <line x1="12" y1="22" x2="12" y2="8"></line>
  <path d="M5 12H2a10 10 0 0 0 20 0h-3"></path>
</svg>
END
