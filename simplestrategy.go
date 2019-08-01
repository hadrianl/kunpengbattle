package main

import (
	"log"
	"math"
	"math/rand"

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

func (s *simpleStrategy) React(round kpb.KunPengRound) (kpb.KunPengAction, error) {
	log.Printf("Round: mode: %v  force: %v teamID: %v", round.Mode, s.Teams[s.TeamID].Force, s.TeamID)
	s.CurrentRoundID = round.ID

	action := new(kpb.KunPengAction)
	action.ID = s.CurrentRoundID
	action.Actions = make([]kpb.KunPengMove, 0, len(s.Allies))
	move := [5]string{"up", "down", "right", "left", ""}
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
		m := make([]string, 1)
		powerWeightPoint := float64(0)
		pm := ""
		if len(powers) > 0 {
			for _, power := range powers {
				xoffset := float64(power.X - player.X)
				yoffset := float64(power.Y - player.Y)
				// carea := math.Abs(xoffset + 1) * math.Abs(yoffset +1)
				coffset := (math.Abs(xoffset) + math.Abs(yoffset))
				cweightPoint := float64(power.Point) / coffset
				if cweightPoint > powerWeightPoint {
					// area = carea
					// offset = coffset
					powerWeightPoint = cweightPoint
					if math.Abs(xoffset) >= math.Abs(yoffset) {
						if xoffset >= 0 {
							pm = "right"
						} else {
							pm = "left"
						}
					} else {
						if yoffset >= 0 {
							pm = "down"
						} else {
							pm = "up"
						}
					}
				}
			}
		}

		enemyWeightPoint := float64(0)
		em := ""
		if len(enemyInView) > 0 {
			log.Println("enemyInView:", enemyInView)
			// log.Printf("Mode: %v Force: %v", round.Mode, s.Teams[s.TeamID].Force)
			for _, enemy := range enemyInView {
				xoffset := float64(enemy.X - player.X)
				yoffset := float64(enemy.Y - player.Y)

				coffset := (math.Abs(xoffset) + math.Abs(yoffset))

				if round.Mode != s.Teams[s.TeamID].Force {
					cweightPoint := float64(enemy.Score+10) / coffset
					if cweightPoint > enemyWeightPoint {
						enemyWeightPoint = cweightPoint
						if math.Abs(xoffset) >= math.Abs(yoffset) {
							if xoffset >= 0 {
								em = "right"
							} else {
								em = "left"
							}
						} else {
							if yoffset >= 0 {
								em = "down"
							} else {
								em = "up"
							}
						}
					}
				} else {
					cweightPoint := float64(-10-player.Score) / coffset
					if cweightPoint < enemyWeightPoint {
						enemyWeightPoint = cweightPoint
						if math.Abs(xoffset) >= math.Abs(yoffset) {
							if xoffset >= 0 {
								em = "left"
							} else {
								em = "right"
							}
						} else {
							if yoffset >= 0 {
								em = "up"
							} else {
								em = "down"
							}
						}
					}
				}
			}
		}

		if powerWeightPoint != 0 || enemyWeightPoint != 0 {
			log.Println(powerWeightPoint, enemyWeightPoint)
			if math.Abs(enemyWeightPoint) >= math.Abs(powerWeightPoint) {
				m[0] = em
			} else {
				m[0] = pm
			}
		} else {
			m[0] = move[rand.Intn(5)]
		}

		ac := kpb.KunPengMove{Team: s.TeamID, PlayerID: player.ID, Move: m}
		action.Actions = append(action.Actions, ac)

	}

	return *action, nil

}
