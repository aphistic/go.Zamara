package main

func handleStdout(flags zamaraFlags) {
	switch flags.runType {
	case "mpq":
		mpqStdout(flags)
		break
	}
}
