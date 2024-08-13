package utils

type Player struct {
	isFirst bool
	boardId int64
}

func InitializePlayer(id int64, isFirst bool) Player {
	var player Player
	player.boardId = id
	player.isFirst = isFirst
	return player
}
