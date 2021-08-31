pkgname=simplefileshare
pkgver=0.1.0
pkgrel=0
pkgdesc="Simple file share"
arch=('x86_64' 'aarch64')
license=('GPL')
url="https://github.com/nxshock/$pkgname"
makedepends=('go' 'git')
options=('!strip')
backup=("etc/$pkgname.toml")
source=("git+https://github.com/nxshock/$pkgname.git")
sha256sums=('SKIP')

build() {
	cd "$srcdir/$pkgname"

	export CGO_CPPFLAGS="${CPPFLAGS}"
	export CGO_CFLAGS="${CFLAGS}"
	export CGO_CXXFLAGS="${CXXFLAGS}"
	export CGO_LDFLAGS="${LDFLAGS}"

	go build -o $pkgname -buildmode=pie -trimpath -ldflags="-linkmode=external -s -w"
}

package() {
	cd "$srcdir/$pkgname"

	install -Dm755 "$pkgname"          "$pkgdir/usr/bin/$pkgname"
	install -Dm644 "$pkgname.conf"     "$pkgdir/etc/$pkgname.conf"
	install -Dm644 "$pkgname.service"  "$pkgdir/usr/lib/systemd/system/$pkgname.service"
	install -Dm644 "$pkgname.sysusers" "$pkgdir/usr/lib/sysusers.d/$pkgname.conf"
}
