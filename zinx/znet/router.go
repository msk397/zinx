package znet

import "zinx/zinx/ziface"

// 实现 Router 时，先嵌入这个基类，然后根据需要对这个基类的方法进行重写就行了
type BaseRouter struct{}

// 这里之所以 BaseRouter 的方法都为空
// 是因为有的 Router 不希望有 PreHandle 或 PostHandle
// 所以 Router 全部继承 BaseRouter 的好处就是，不需要实现 PreHandle 和 PostHandle 也可以将IRouter实例化
func (b *BaseRouter) PreHandle(request ziface.IRequest) {}

func (b *BaseRouter) Handle(request ziface.IRequest) {}

func (b *BaseRouter) PostHandle(request ziface.IRequest) {}
