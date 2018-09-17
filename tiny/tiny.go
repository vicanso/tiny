package tiny

const (
	// GZIP gzip type
	GZIP = iota
	// BROTLI brotli type
	BROTLI
	// JPEG jepg type
	JPEG
	// PNG png type
	PNG
	// WEBP webp type
	WEBP
	// GUETZLI guetzli jpeg type
	GUETZLI
)

const (
	// ClipNone none clip
	ClipNone = iota
	// ClipCenter clip center
	ClipCenter
	// ClipLeftTop clip left top
	ClipLeftTop
	// ClipTopCenter clip top center
	ClipTopCenter
)

const (
	// AppPngquant png quant application
	AppPngquant = "pngquant"
	// AppGuetzli guetzli application
	AppGuetzli = "guetzli"
)
