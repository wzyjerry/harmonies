package finder

import (
	"github.com/wzyjerry/harmonies/pkg/cube"
	"github.com/wzyjerry/harmonies/pkg/pattern"
)

// 因为潜在的解空间很大，采用启发式搜索
// 1. 考虑地图 A 面 5 x 5 共 23 个格子；B 面 4 x 7 共 25 个格子，最大填充 25 次。
//    特别地，这里只考虑 A 面。
// 2. 简单起见不考虑生成的形状，但打分时对长条形状有所惩罚。
// 3. 考虑到每种卡牌都是由主要地形和次要地形两种不同类型组成，因此惩罚最多种类的地形，平方加权。
// 4. 鼓励更少的占地，特别地，当发现第一个解时搜索空间不超过当前最小占地 + 3。
// 5. 区域总是相邻的。
// 6. 以第一张卡牌的默认形状作为起点。

// 正向计分：
// 1. 完成的任务数量，每个任务记 50 分

// 负向计分：
// 1. 占用格子数量，每个格子扣 3 分
// 2. 考虑 6，我们总是将第一张卡牌的动物位置标记为中心。
//    扣除 距离^2 分数
// 3. 对于使用最多的 Perfab，惩罚 count^2
// 4. 使用 token 的数量，每个token 扣 5 分

type Scene struct {
	Score      int
	pattern    *pattern.Pattern
	usage      int
	usageToken int
	perfabs    map[string]int
	distance   int
	animals    int
	animalsAt  [][]cube.Hex
}

type PriorityQueue []Scene

func (pq PriorityQueue) Len() int { return len(pq) }

// Less 大根堆
func (pq PriorityQueue) Less(i, j int) bool {
	return pq[i].Score > pq[j].Score
}

func (pq PriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
}

func (pq *PriorityQueue) Push(x any) {
	*pq = append(*pq, x.(Scene))
}

func (pq *PriorityQueue) Pop() any {
	old := *pq
	n := len(old)
	x := old[n-1]
	*pq = old[0 : n-1]
	return x
}

func (s *Scene) CalcScore() {
	score := s.animals * 50
	score -= s.usage * 3
	score -= s.usageToken * 5
	score -= s.distance * s.distance
	largest := 0
	for _, count := range s.perfabs {
		if count > largest {
			largest = count
		}
	}
	score -= largest * largest
	s.Score = score
}
