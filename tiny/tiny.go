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
	// AppPngquant png quant application
	AppPngquant = "pngquant"
	// AppGuetzli guetzli application
	AppGuetzli = "guetzli"
)
