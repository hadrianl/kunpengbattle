package main

import (
	"math/rand"

	kpb "./kunpengBattle"
)

type hadrianlStrategy struct {
	TeamID         int
	TeamName       string
	Allies         map[int]*kpb.KunPengPlayer
	Enemies        map[int]*kpb.KunPengPlayer
	Teams          map[int]*kpb.KunPengTeam
	Map            kpb.KunPengMap
	CurrentRoundID int
	MatrixMap      [][][]int
}

func (s *hadrianlStrategy) Registrate(registration kpb.KunPengRegistration) error {
	s.TeamID = registration.TeamID
	s.TeamName = registration.TeamName
	return nil
}

func (s *hadrianlStrategy) LegStart(legStart kpb.KunPengLegStart) error {
	s.Teams = make(map[int]*kpb.KunPengTeam)
	s.Allies = make(map[int]*kpb.KunPengPlayer)
	s.Enemies = make(map[int]*kpb.KunPengPlayer)
	for _, t := range legStart.Teams {
		s.Teams[t.ID] = &t

		if s.TeamID == t.ID {
			for _, playerID := range t.Players {
				s.Allies[playerID] = &kpb.KunPengPlayer{ID: playerID, Team: t.ID}
			}
		} else {
			for _, playerID := range t.Players {
				s.Enemies[playerID] = &kpb.KunPengPlayer{ID: playerID, Team: t.ID}
			}
		}

	}

	s.Map = legStart.Map
	return nil
}

func (s *hadrianlStrategy) LegEnd(legEnd kpb.KunPengLegEnd) error {
	for _, t := range legEnd.Teams {
		team := s.Teams[t.ID]
		team.Point = t.Point
	}
	return nil
}

func (s *hadrianlStrategy) React(round kpb.KunPengRound) (kpb.KunPengAction, error) {
	s.CurrentRoundID = round.ID

	action := new(kpb.KunPengAction)
	action.ID = s.CurrentRoundID
	action.Actions = make([]kpb.KunPengMove, 0, len(s.Allies))
	move := [4]string{"up", "down", "right", "left"}
	// players := round.Players
	// powers := round.Power

	for _, player := range s.Allies {

		ac := kpb.KunPengMove{Team: s.TeamID, PlayerID: player.ID, Move: []string{move[rand.Intn(4)]}}
		action.Actions = append(action.Actions, ac)
	}
	return *action, nil

}

// func (s *hadrianlStrategy) initMatrixMap(kpMap kpb.KunPengMap) [][][6]int {
// 	width := kpMap.Width
// 	height := kpMap.Height
// 	// matrixMap := [width][height][6]int  // xoffset, yoffset
// 	var matrixMap [][][6]int
// 	for x := 0; x < width; x++ {
// 		for y := 0; y < height; y++ {
// 			matrixMap[x][y] = [2]int{0, 0}
// 			// switch {
// 			// case x == 0:
// 			// 	matrixMap[x][y][0] = 0
// 			// case x == width:
// 			// 	matrixMap[x][y][2] = 0
// 			// }

// 			// switch {
// 			// case y == 0:
// 			// 	matrixMap[x][y][3] = 0
// 			// case y == height:
// 			// 	matrixMap[x][y][1] = 0
// 			// }

// 		}
// 	}

// 	for _, m := range kpMap.Meteor {
// 		matrixMap[m.X][m.Y] = [6]int{-1, -1, -1, -1, -1, 0}
// 	}

// 	for _, t := range kpMap.Tunnel {
// 		switch t.Direction {
// 		case "up":
// 			matrixMap[t.X][t.Y] = [6]int{0, 1}
// 		case "right":
// 			matrixMap[t.X][t.Y] = [6]int{1, 0}
// 		case "down":
// 			matrixMap[t.X][t.Y] = [6]int{0, -1}
// 		case "left":
// 			matrixMap[t.X][t.Y] = [6]int{-1, 0}
// 		}

// 	}

// 	return matrixMap

// }

// func calcWeight(power kpb.KunPengPower, player kpb.KunPengPlayer) float32 {
// 	for x:=player.X;
// }
