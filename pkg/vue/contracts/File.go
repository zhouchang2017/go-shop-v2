package contracts

type StaticFile interface {
	GetKey() string
	GetName() string
	GetDrive() string
	PreviewUrl() string
}
