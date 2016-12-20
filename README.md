# ggs
Go Game Server

# 执行流程
ggs首先会在同一个goroutine中按模块注册顺序执行模块的OnInit方法，等到所有模块的OnInit方法执行完后则为每一个模块启动一个goroutine来执行模块的Run方法。最后，游戏服务器关闭时(Ctrl + C 关闭游戏服务器)，将按与模块注册相反顺序在同一个goroutine中执行模块的OnDestroy方法。 

# 消息发送流程 
客户端发送到游戏服务器的消息需要通过gate模块路由，简而言之，gate模块决定了某个消息具体交给内部的哪个模块来处理。
