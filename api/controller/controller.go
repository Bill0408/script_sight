package controller

type Controller struct {
	Hash string
}

func New(h string) *Controller {
	return &Controller{
		Hash: h,
	}
}
