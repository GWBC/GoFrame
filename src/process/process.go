package process

var processList = []Process{
	&LogProcess{},
	&DBProcess{},
}

type Process interface {
	PreProcess() error //前置处理器
	PostProcess()      //后置处理器
}

func PreProcess() error {
	for _, p := range processList {
		err := p.PreProcess()
		if err != nil {
			return err
		}
	}

	return nil
}

func PostProcess() {
	for _, p := range processList {
		p.PostProcess()
	}
}
