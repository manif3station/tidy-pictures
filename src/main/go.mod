module michaelpc.com/tidy_pictures

go 1.19

replace michaelpc.com/shared_lib => ../shared

replace michaelpc.com/exif => ../exif

require (
	github.com/barasher/go-exiftool v1.10.0 // indirect
	michaelpc.com/exif v0.0.0-00010101000000-000000000000 // indirect
	michaelpc.com/shared_lib v0.0.0-00010101000000-000000000000 // indirect
)
