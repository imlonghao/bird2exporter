package bird

func (b *Bird) Exec(command string) parsed {
	b.write(command)
	return parsed{
		"raw": b.read(),
	}
}
