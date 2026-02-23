#!/bin/bash

outputDirectory="build"
targetsList=("linux/amd64" "linux/arm64" "windows/amd64" "windows/arm64")
componentsList=("server" "client")

# Clean and create output directory
rm -rf "$outputDirectory"
mkdir -p "$outputDirectory"

echo "🚀 Starting build process..."

for targetItem in "${targetsList[@]}"; do
  IFS="/" read -r operatingSystem architecture <<< "$targetItem"
  
  # Map architecture names for consistency
  archLabel="$architecture"
  if [ "$architecture" == "amd64" ]; then
    archLabel="x64"
  fi

  echo "----------------------------------------"
  echo "📦 Target: $operatingSystem / $archLabel"
  
  for componentItem in "${componentsList[@]}"; do
    # New naming convention: component-os-arch
    binaryName="${componentItem}-${operatingSystem}-${archLabel}"
    if [ "$operatingSystem" == "windows" ]; then
      binaryName="${binaryName}.exe"
    fi
    
    echo "  🛠️  Building $componentItem..."
    
    buildDirectory="./$componentItem"
    outputPath="${outputDirectory}/${binaryName}"
    
    # Build
    # -trimpath: removes local file system paths from the binary
    # -ldflags="-s -w": strips symbol table and debug information
    GOOS=$operatingSystem GOARCH=$architecture go build -trimpath -ldflags="-s -w" -o "$outputPath" "$buildDirectory/cmd/main.go"
    
    if [ $? -eq 0 ]; then
      echo "  ✅ Success: $outputPath"
    else
      echo "  ❌ Failed: $componentItem for $operatingSystem/$architecture"
    fi
  done
done

echo "----------------------------------------"
echo "✨ Build process complete! Check the '$outputDirectory' directory."
