package usecase

type AppUseCase struct {
	TreeUC TreeUseCase
}

var appUseCase AppUseCase

func InitAppUseCase(treeUC TreeUseCase) {
	appUseCase = AppUseCase{treeUC}
}

func GetAppUseCase() AppUseCase {
	return appUseCase
}
