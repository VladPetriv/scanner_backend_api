package utils

func FormatPage(page int) int {
	if page == 0 || page == 1 {
		page = 0
	} else if page == 2 {
		page = 10
	} else {
		page *= 10
		page -= 10
	}

	return page
}
