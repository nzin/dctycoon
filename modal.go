package dctycoon

import(
	"github.com/nzin/sws"
)



func ShowModalError(root *sws.RootWidget,title, desc string, callback func()) {
	modal:=sws.NewMainWidget(500, 200, title, false, false)
	modal.SetCloseCallback(func() {
		root.RemoveChild(modal)
		if callback!=nil {
			callback()
		}
	})
	
	icon:=sws.NewLabelWidget(32,32,"")
	icon.SetImage("resources/icon-triangular-big.png")
	icon.Move(20,50)
	modal.AddChild(icon)
	
	textarea:=sws.NewTextAreaWidget(400,70,desc)
	textarea.Move(70,40)
	textarea.SetReadonly(true)
	modal.AddChild(textarea)
	
	ok:=sws.NewButtonWidget(100,25,"Ok")
	ok.Move(370,120)
	ok.SetClicked(func() {
		root.RemoveChild(modal)
		if callback!=nil {
			callback()
		}
	})
	modal.AddChild(ok)
	modal.Move((root.Width()-500)/2,(root.Height()-200)/2)
	
	root.AddChild(modal)
	root.SetModal(modal)
}
