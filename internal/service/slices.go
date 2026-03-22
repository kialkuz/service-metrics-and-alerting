package service

func CheckExistValue(value string, listvalues []string) bool {
	for _, v := range listvalues {
		if v == value {
			return true
		}
	}

	return false
}
