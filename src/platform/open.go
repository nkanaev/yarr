package platform

func Open(input string) error {
	return open(input).Run()
}
