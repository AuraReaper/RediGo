package core

var KeyspaceStat [4]map[string]int

func UpdateDBStat(num, value int, metric string) {
	KeyspaceStat[num][metric] = value
}
