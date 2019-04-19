package main

import (
	"fmt"
	"os"
)

/**
 * 广度优先算法
 */
/**
 * 从文件中读取数据
 */

func readMaze(filename string) [][]int {
	file, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	var cols, rows int
	fmt.Fscanf(file, "%d %d", &rows, &cols)
	maze := make([][]int, rows)
	for i := range maze {
		maze[i] = make([]int, cols)
		for j := range maze[i] {
			fmt.Fscanf(file, "%d", &maze[i][j])
		}
	}
	return maze
}

/**
 * 定义要走的节点下标
 */
type point struct {
	i, j int
}

/**
 * 定义,四个方向 ,上左下右
 */
var dirs = [4]point{point{-1, 0}, point{0, -1}, point{1, 0}, point{0, 1},}

//走的节点
//走迷宫
/**
 * maze 迷宫地址
 * strt 开始位置
 * end 走出迷宫位置
 */
func walk(maze [][]int, start, end point) [][]int {
	steps := make([][]int, len(maze))
	for i := range steps {
		steps[i] = make([]int, len(maze[i]))
	}
	Q := []point{start}
	for len(Q) > 0 {
		cur := Q[0] //取一个第一个位置
		Q = Q[1:]
		//fmt.Printf("Q=%v ", Q)
		//发现终点就退出
		if cur == end {
			break
		}

		for _, dir := range dirs { //开始进行 上 左 下 右 的行走路线
			next := cur.add(dir) //拿到下个节点的值
			//撞墙检测
			val, ok := next.at(maze)
			if !ok || val == 1 {
				continue
			}
			//检测是否走过了
			val, ok = next.at(steps)
			if !ok || val != 0 {
				continue
			}
			//检测是否回到原点了不能探索
			if next == start {
				continue
			}
			//开始探索 ,将值存入 steps
			curStep, _ := cur.at(steps)
			steps[next.i][next.j] = curStep + 1

			Q = append(Q, next)
		}
	}
	return steps
}

func (p point) at(grid [][]int) (int, bool) {
	//往上越界,往下越界
	if p.i < 0 || p.i >= len(grid) {
		return 0, false
	}

	if p.j < 0 || p.j >= len(grid[p.i]) {
		return 0, false
	}

	return grid[p.i][p.j], true

}

// 添加节点
func (p point) add(r point) point {
	return point{p.i + r.i, p.j + r.j}
}

func main() {
	maze := readMaze("maze.in") //读取地图数据
	steps := walk(maze, point{0, 0}, point{len(maze) - 1, len(maze[0]) - 1})
	for _, rows := range steps {
		for _, val := range rows {
			fmt.Printf("%3d ", val)
		}
		fmt.Println()
	}

}
