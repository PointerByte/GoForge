module github.com/PointerByte/GoForge/component/tinygo

go 1.25.0

require (
	github.com/PointerByte/GoForge/component v0.0.0
	go.bytecodealliance.org/cm v0.3.0
)

require github.com/PointerByte/GoForge/portable v0.0.0 // indirect

replace github.com/PointerByte/GoForge/component => ..

replace github.com/PointerByte/GoForge/portable => ../../portable
