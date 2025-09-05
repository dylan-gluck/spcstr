# Deployment Architecture

## Deployment Strategy

**Binary Distribution:**
- **Platform:** Multi-platform binary releases (macOS, Linux, Windows WSL)
- **Build Command:** `goreleaser build --clean`
- **Output Directories:** `dist/` with platform-specific binaries
- **Distribution:** Package managers (Homebrew, APT, Pacman) + GitHub Releases

**Installation Methods:**
```bash
# Homebrew (macOS/Linux)
brew install spcstr

# APT (Debian/Ubuntu)
curl -s https://api.github.com/repos/username/spcstr/releases/latest | grep "browser_download_url.*deb" | cut -d '"' -f 4 | wget -i -
sudo dpkg -i spcstr_*.deb

# Direct binary download
wget https://github.com/username/spcstr/releases/latest/download/spcstr_linux_amd64.tar.gz
tar -xzf spcstr_linux_amd64.tar.gz
sudo mv spcstr /usr/local/bin/
```

## CI/CD Pipeline
```yaml
# .github/workflows/release.yaml
name: Release

on:
  push:
    tags: ['v*']

jobs:
  goreleaser:
    runs-on: ubuntu-latest
    steps:
    - uses: actions/checkout@v4
      with:
        fetch-depth: 0
    
    - uses: actions/setup-go@v4
      with:
        go-version: '1.21'
    
    - name: Run GoReleaser
      uses: goreleaser/goreleaser-action@v4
      with:
        version: latest
        args: release --clean
      env:
        GITHUB_TOKEN: ${{ secrets.GITHUB_TOKEN }}
```

## Environments

| Environment | Purpose | Binary Location |
|-------------|---------|-----------------|
| Development | Local development and testing | `./dist/spcstr_<platform>` |
| CI | Automated testing and validation | GitHub Actions runners |
| Production | End user installations | Package managers + GitHub Releases |
