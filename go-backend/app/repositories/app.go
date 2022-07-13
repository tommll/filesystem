package repositories

type AppRepositories struct {
	ItemRepo ItemRepository
}

var appRepositories AppRepositories

func InitAppRepositories(itemRepo ItemRepository) {
	appRepositories = AppRepositories{itemRepo}
}

func GetAppRepo() AppRepositories {
	return appRepositories
}
