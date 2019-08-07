package main

import (
	// "log"
	"math"
	"math/rand"
	"fmt"

	kpb "./kunpengBattle"
)

type simpleStrategy struct {
	TeamID         int
	TeamName       string
	Allies         map[int]*kpb.KunPengPlayer
	Enemies        map[int]*kpb.KunPengPlayer
	Teams          map[int]*kpb.KunPengTeam
	Map            kpb.KunPengMap
	CurrentRoundID int
	MatrixMap      [][][]int
}

func (s *simpleStrategy) Registrate(registration kpb.KunPengRegistration) error {
	s.TeamID = registration.TeamID
	s.TeamName = registration.TeamName
	return nil
}

func (s *simpleStrategy) LegStart(legStart kpb.KunPengLegStart) error {
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

func (s *simpleStrategy) LegEnd(legEnd kpb.KunPengLegEnd) error {
	for _, t := range legEnd.Teams {
		team := s.Teams[t.ID]
		team.Point = t.Point
	}
	return nil
}

var movementOffset = map[string][2]int{"up": [2]int{0, -1}, "right": [2]int{1, 0}, "down": [2]int{0, 1}, "left": [2]int{-1, 0}, "": [2]int{0, 0}}

func (s *simpleStrategy) React(round kpb.KunPengRound) (kpb.KunPengAction, error) {
	fmt.Printf("Round: mode: %v  force: %v teamID: %v", round.Mode, s.Teams[s.TeamID].Force, s.TeamID)
	s.CurrentRoundID = round.ID

	action := new(kpb.KunPengAction)
	action.ID = s.CurrentRoundID
	action.Actions = make([]kpb.KunPengMove, 0, len(s.Allies))
	// move := [5]string{"up", "down", "right", "left", ""}
	players := round.Players
	enemyInView := []*kpb.KunPengPlayer{}
	powers := round.Power

	for _, p := range players {
		if ally, ok := s.Allies[p.ID]; ok {
			ally.Score = p.Score
			ally.Sleep = p.Sleep
			ally.X = p.X
			ally.Y = p.Y
		}

		if Enemy, ok := s.Enemies[p.ID]; ok {
			Enemy.Score = p.Score
			Enemy.Sleep = p.Sleep
			Enemy.X = p.X
			Enemy.Y = p.Y
			enemyInView = append(enemyInView, Enemy)
		}
	}

	for _, player := range s.Allies {
		// offset := float64(0)
		// log.Println(player)
		movementWeight := map[string]float64{"up": 0, "right": 0, "down": 0, "left": 0, "": 0}
		// m := make([]string, 1)
		if len(powers) > 0 {
			for _, power := range powers {
				for k, w := range movementWeight {
					xoffset := float64(power.X - (player.X + movementOffset[k][0]))
					yoffset := float64(power.Y - (player.Y + movementOffset[k][1]))
					coffset := (math.Abs(xoffset) + math.Abs(yoffset)) + 1
					if coffset < 2*float64(s.Map.Vision)+1.0 {
						cweightPoint := float64(power.Point) / coffset
						movementWeight[k] = w + cweightPoint
					}

				}

				// carea := math.Abs(xoffset + 1) * math.Abs(yoffset +1)
			}
		}

		if len(enemyInView) > 0 {
			// log.Println("enemyInView:", enemyInView)
			// log.Printf("Mode: %v Force: %v", round.Mode, s.Teams[s.TeamID].Force)
			for _, enemy := range enemyInView {
				for k, w := range movementWeight {
					xoffset := float64(enemy.X - (player.X + movementOffset[k][0]))
					yoffset := float64(enemy.Y - (player.Y + movementOffset[k][1]))
					coffset := (math.Abs(xoffset) + math.Abs(yoffset)) + 1
					var cweightPoint float64
					if round.Mode != s.Teams[s.TeamID].Force {
						cweightPoint = float64(enemy.Score+10) / coffset
					} else {
						cweightPoint = float64(-15-player.Score) / coffset
					}
					movementWeight[k] = w + cweightPoint
				}
			}
		}

		for _, ally := range s.Allies {
			for k, w := range movementWeight {
				if ally.ID != player.ID {
					xoffset := float64(ally.X - (player.X + movementOffset[k][0]))
					yoffset := float64(ally.Y - (player.Y + movementOffset[k][1]))
					coffset := (math.Abs(xoffset) + math.Abs(yoffset)) + 1
					if coffset <= 2 {
						cweightPoint := float64(-1) / coffset
						movementWeight[k] = w + cweightPoint
					}
				}
			}
		}

		for _, meteor := range s.Map.Meteor {
			for k, w := range movementWeight {
				xoffset := float64(meteor.X - (player.X + movementOffset[k][0]))
				yoffset := float64(meteor.Y - (player.Y + movementOffset[k][1]))
				coffset := (math.Abs(xoffset) + math.Abs(yoffset)) + 1
				if coffset == 1 {
					cweightPoint := float64(-10) / coffset
					movementWeight[k] = w + cweightPoint
				}

			}
		}

		for x := range []int{-1, s.Map.Width} {
			for y := range []int{-1, s.Map.Height} {
				for k, w := range movementWeight {
					xoffset := float64(x - (player.X + movementOffset[k][0]))
					yoffset := float64(y - (player.Y + movementOffset[k][1]))
					coffset := (math.Abs(xoffset) + math.Abs(yoffset)) + 1
					if coffset == 1 {
						cweightPoint := float64(-10) / coffset
						movementWeight[k] = w + cweightPoint
					}

				}
			}
		}

		skipMove := []string{}
		for _, m := range s.Map.Tunnel {
			switch {
			case m.Y == player.Y && m.X-player.X == 1 && m.Direction == "left":
				skipMove = append(skipMove, "right")
				fallthrough
			case m.Y == player.Y && m.X-player.X == -1 && m.Direction == "right":
				skipMove = append(skipMove, "left")
				fallthrough
			case m.X == player.X && m.Y-player.Y == -1 && m.Direction == "down":
				skipMove = append(skipMove, "up")
				fallthrough
			case m.X == player.X && m.Y-player.Y == 11 && m.Direction == "up":
				skipMove = append(skipMove, "down")
				fallthrough
			case (player.X == 0 || player.X == s.Map.Width-1) && (player.Y == 0 || player.Y == s.Map.Height-1):
				skipMove = append(skipMove, "")

				// default:
				// 	skipMove = append(skipMove, "")
			}
		}

		ac := kpb.KunPengMove{Team: s.TeamID, PlayerID: player.ID, Move: choiceMovement(movementWeight, skipMove...)}
		action.Actions = append(action.Actions, ac)
	}

	return *action, nil

}

func bool2Int(b bool) uint {
	if b {
		return 1
	}

	return 0
}

func choiceMovement(mw map[string]float64, skipMove ...string) []string {
	var move []string
	var weight = -math.MaxFloat64
	// mwLoog:
	for k, w := range mw {
		// for _, sm := range skipMove {
		// 	if sm == k {
		// 		continue mwLoog
		// 	}
		// }

		if w >= weight {
			weight = w
			move = []string{k}
		}
	}

	if weight == 0 && move[0] == "" {
		move[0] = []string{"up", "down", "right", "left"}[rand.Intn(4)]
	}

	return move

}
