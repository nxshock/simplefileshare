pkgname=simplefileshare
pkgver=0.1.4
pkgrel=0
pkgdesc="Simple file share"
arch=('x86_64' 'aarch64')
license=('GPL')
url="https://github.com/nxshock/$pkgname"
makedepends=('go' 'git')
options=('!strip')
backup=("etc/$pkgname.conf")
source=(
	"git+https://github.com/nxshock/$pkgname.git"
	'git+https://github.com/dmhendricks/file-icon-vectors')
sha256sums=(
	'SKIP'
	'SKIP')

build() {
	cd "$srcdir/$pkgname"

	mv ../file-icon-vectors/dist/icons/high-contrast icons

	go build -o $pkgname -buildmode=pie -trimpath -ldflags="-linkmode=external -s -w"
}

package() {
	cd "$srcdir/$pkgname"

	install -Dm755 "$pkgname"          "$pkgdir/usr/bin/$pkgname"
	install -Dm644 "$pkgname.conf"     "$pkgdir/etc/$pkgname.conf"
	install -Dm644 "$pkgname.service"  "$pkgdir/usr/lib/systemd/system/$pkgname.service"
	install -Dm644 "$pkgname.sysusers" "$pkgdir/usr/lib/sysusers.d/$pkgname.conf"
	install -Dm644 "$pkgname.tmpfiles" "$pkgdir/usr/lib/tmpfiles.d/$pkgname.conf"
}
