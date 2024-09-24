# magic-url

## WORK IN PROGRESS

### Sebuah aplikasi Backend REST API untuk layanan aplikasi pemendek URL atau tautan.

## How to build

### release build
```sh
# linux
pkgname="github.com/binarstrike/magic-url"

go build -o magic-url -ldflags="-s -w -X $pkgname/config.APP_ENV=production -X $pkgname/config.APP_VERSION=0.0.3-release" -trimpath $pkgname/cmd/magic-url/...

# windows powershell
$pkgname = "github.com/binarstrike/magic-url"

go build -o magic-url.exe -ldflags="-s -w -X $pkgname/config.APP_ENV=production -X $pkgname/config.APP_VERSION=0.0.3-release" -trimpath $pkgname/cmd/magic-url/...
```

### debug build
```sh
# linux
go build -o magic-url ./cmd/magic-url

# windows
go build -o magic-url.exe .\cmd\magic-url
```

## How to run

```sh
# linux
./magic-url

# windows
.\magic-url.exe
```