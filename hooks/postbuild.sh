#!/bin/bash

mkdir -p cdnmanager.app/Contents/Licenses
cp ../../frontend/src/assets/fonts/IBM_Plex_Mono/license.txt cdnmanager.app/Contents/Licenses/IBMPlexMono\ License.txt
cp ../../frontend/src/assets/fonts/glyphicons_halflings/license.txt cdnmanager.app/Contents/Licenses/Glyphicons\ Halflings\ License.txt
rm ../appicon.png