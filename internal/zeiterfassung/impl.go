package zeiterfassung

func blub() {
	var c Kunde
	c = Privatkunde("asd")
	c = Firmenkunde("asdf")
	MatchKunde[int](c, func(privatkunde Privatkunde) int {
		return 1
	}, func(firmenkunde Firmenkunde) int {
		return 2
	})
}
