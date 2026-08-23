package domain

func allowFromActive(next string) bool { return next != "closed" }
