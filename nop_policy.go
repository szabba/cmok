package cmok

type NopPolicy struct{}

var _ AccessPolicy = NopPolicy{}

func (_ NopPolicy) Protect(storage Storage, _ User) Storage {
	return storage
}
