package moduleservice

type ModuleService interface {
	EnableModule(id uint) (string, error)
	DisableModule(id uint) (string, error)
}
