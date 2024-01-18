package bridges

func Inshellisense() []string {
	// out, err := execlog.Command("inshellisense", "list").Output()
	// if err != nil {
	// 	return []string{}
	// }

	// var entries []string
	// if err := json.Unmarshal(out, &entries); err != nil {
	// 	return []string{}
	// }
	// TODO hardcoded for now (https://github.com/microsoft/inshellisense/pull/154)
	entries := inshellisenseCompleters

	unique := make(map[string]bool)
	for _, entry := range entries {
		unique[entry] = true
	}

	completers := make([]string, 0)
	for c := range filter(unique, inshellisenseBuiltins) {
		completers = append(completers, c)
	}
	return completers
}
