package html

import "embed"

//go:embed assets/gohtml.js assets/idiomorph.js
var Assets embed.FS

//go:embed assets/tailwind.js
var Tailwind embed.FS
