#!/usr/bin/env bash
# note: the `lightningcss` binary is available in our PATH because our Makefile
# passes down its path to this subshell. We set our Makefile's PATH to include
# the location of our local GOBIN and the Lightning CSS binary in our
# node_modules folder for easier use in the build scripts

# the directory our static site is generated to
outputDir=""
# these will be the default flags for both "dev" and "prod" mode
flags="--bundle --custom-media --targets defaults"
# should be "dev" or "prod", "prod" by default
mode="prod"

# loop through all the passed arguments and set the $outputDir variable equal
# to the value passed to the "-o", and the $mode variable equal to the value
# passed to the "-m" flag
while getopts "o:m:" arg; do
  case $arg in
    o) outputDir=$OPTARG;; # -o "output"
    m) mode=$OPTARG;; # -m "mode"
    *) ;;
  esac
done

# set the appropriate flags based upon whether $mode is "dev" or "prod"
if [ "$mode" == "dev" ]; then
  # generate a sourcemap for all our CSS files in dev mode and tell lightningcss
  # to not throw errors if it runs into any parsing issues when in dev mode.
  # we want errors before we build for production but for fast iterations in dev
  # mode it's nicer to not have it error out
  flags="$flags --sourcemap --error-recovery"
else
  # minify our CSS in production
  flags="$flags --minify"
fi

# process and bundle our main CSS file
lightningcss $flags styles/index.css -o "$outputDir/styles.css" \
  && echo "Bundled $outputDir/styles.css"

# create the output directory of route-specific CSS files
mkdir -p "$outputDir/styles/routes"

# process all of the CSS files in the styles/routes directory individually
for file in styles/routes/*.css; do
  if [ ! -f "$file" ]; then
    continue
  fi

  lightningcss $flags "$file" -o "$outputDir/$file" \
    && echo "Bundled $outputDir/$file"
done
