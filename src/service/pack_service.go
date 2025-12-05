package service

type PackService struct {
}

func (s *PackService) Init() error {
	
	return nil
}
func (s *PackService) Uninit() {
}

func (s *PackService) Name() string {
	return "打包"
}

func (s *PackService) Proc() {

}

func (s *PackService) SubMessage() []int {
	return []int{}
}

func (s *PackService) ProcMessage(id int, args ...any) {

}
