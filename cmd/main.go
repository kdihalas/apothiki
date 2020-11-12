package main

func main() {
	if err := cmd.Execute(); err != nil {
		er(err.Error())
	}
}