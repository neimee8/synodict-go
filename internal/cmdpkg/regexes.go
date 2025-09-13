package cmdpkg

var cmdRegexes = []string{
	`^add(?:\s+"[A-Za-zÀ-ɏЀ-ӿ\- ]+")+$`,
	`^add-words(?:\s+"[A-Za-zÀ-ɏЀ-ӿ\- ]+")+$`,
	`^remove(?:\s+"[A-Za-zÀ-ɏЀ-ӿ\- ]+")+$`,
	`^unlink\s+"[A-Za-zÀ-ɏЀ-ӿ\- ]+"\s+"[A-Za-zÀ-ɏЀ-ӿ\- ]+"$`,
	`^unlink-clean\s+"[A-Za-zÀ-ɏЀ-ӿ\- ]+"\s+"[A-Za-zÀ-ɏЀ-ӿ\- ]+"$`,
	`^count\s+"[A-Za-zÀ-ɏЀ-ӿ\- ]+"$`,
	`^synonyms\s+"[A-Za-zÀ-ɏЀ-ӿ\- ]+"$`,
	`^groups$`,
	`^words$`,
	`^cleanup$`,
	`^clear$`,
	`^help$`,
}
