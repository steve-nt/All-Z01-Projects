package tetris

//sortAreas rearranges the tetrominoes so that complex shapes (with larger areas) come first
func sortAreas(tetrominoes [][][]string) [][][]string {

	for i := 0; i < len(tetrominoes); i++ {
		for j := 0; j < len(tetrominoes)-1-i; j++ {
			if len(tetrominoes[j])*len(tetrominoes[j][0]) < len(tetrominoes[j+1])*len(tetrominoes[j+1][0]) {
				tetrominoes[j], tetrominoes[j+1] = tetrominoes[j+1], tetrominoes[j]
			}
		}
	}
	return tetrominoes
}
