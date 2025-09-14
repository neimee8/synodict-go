package cmdpkg

var cmdRegexes = []string{
	`^add(?:\s+"[A-Za-zÀ-ɏЀ-ӿ\- ]+")+$`,
	`^add-words(?:\s+"[A-Za-zÀ-ɏЀ-ӿ\- ]+")+$`,
	`^remove(?:\s+"[A-Za-zÀ-ɏЀ-ӿ\- ]+")+$`,
	`^unlink\s+"[A-Za-zÀ-ɏЀ-ӿ\- ]+"\s+"[A-Za-zÀ-ɏЀ-ӿ\- ]+"$`,
	`^unlink-clean\s+"[A-Za-zÀ-ɏЀ-ӿ\- ]+"\s+"[A-Za-zÀ-ɏЀ-ӿ\- ]+"$`,
	`^check\s+"[A-Za-zÀ-ɏЀ-ӿ\- ]+"\s+"[A-Za-zÀ-ɏЀ-ӿ\- ]+"$`,
	`^check-direct\s+"[A-Za-zÀ-ɏЀ-ӿ\- ]+"\s+"[A-Za-zÀ-ɏЀ-ӿ\- ]+"$`,
	`^exists\s+"[A-Za-zÀ-ɏЀ-ӿ\- ]+"$`,
	`^count\s+"[A-Za-zÀ-ɏЀ-ӿ\- ]+"$`,
	`^synonyms\s+"[A-Za-zÀ-ɏЀ-ӿ\- ]+"$`,
	`^direct-synonyms\s+"[A-Za-zÀ-ɏЀ-ӿ\- ]+"$`,
	`^count-groups$`,
	`^groups$`,
	`^count-words$`,
	`^words$`,
	`^cleanup$`,
	`^clear$`,
	`^import$`,
	`^export$`,
	`^help$`,
}
