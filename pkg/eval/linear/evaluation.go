package eval

import (
	"errors"

	"github.com/DARKEMPIRESL/GODARKCHESS/internal/domain"

	. "github.com/DARKEMPIRESL/GODARKCHESS/pkg/common"
)

const totalPhase = 24

var errAddComplexFeature = errors.New("errAddComplexFeature")

type EvaluationService struct {
	tuning             bool
	score              Score
	kingpawnTable      []kingPawnEntry
	features           []int
	weights            []Score
	pieceCount         [COLOUR_NB][PIECE_NB]int
	phase              int
	passedPawns        uint64
	mobilityArea       [COLOUR_NB]uint64
	attacked           [COLOUR_NB]uint64
	attackedBy2        [COLOUR_NB]uint64
	attackedBy         [COLOUR_NB][PIECE_NB]uint64
	pawnAttacksBy2     [COLOUR_NB]uint64
	kingAttacksCount   [COLOUR_NB]int
	kingSq             [COLOUR_NB]int
	kingAttackersCount [COLOUR_NB]int
	kingAreas          [COLOUR_NB]uint64
}

type kingPawnEntry struct {
	wpawns, bpawns uint64
	wking, bking   int
	score          Score
	passed         uint64
}

func NewEvaluationService() *EvaluationService {
	var e = &EvaluationService{
		kingpawnTable: make([]kingPawnEntry, 1<<16),
		features:      make([]int, totalFeatureSize),
		weights:       make([]Score, totalFeatureSize),
	}
	e.initWeights()
	return e
}

func (e *EvaluationService) EnableTuning() {
	e.tuning = true
}

func (e *EvaluationService) initWeights() {
	if 2*len(e.features) != len(w) {
		return
	}
	for i := range e.weights {
		e.weights[i] = S(w[2*i], w[2*i+1])
	}
}

func (e *EvaluationService) Evaluate(p *Position) int {
	e.init(p)

	var pawnKingKey = murmurMix(p.Pawns&p.White,
		murmurMix(p.Pawns&p.Black,
			murmurMix(p.Kings&p.White,
				p.Kings&p.Black)))
	var pke = &e.kingpawnTable[pawnKingKey%uint64(len(e.kingpawnTable))]
	if e.tuning ||
		!(pke.wpawns == p.Pawns&p.White &&
			pke.bpawns == p.Pawns&p.Black &&
			pke.wking == e.kingSq[SideWhite] &&
			pke.bking == e.kingSq[SideBlack]) {
		pke.wpawns = p.Pawns & p.White
		pke.bpawns = p.Pawns & p.Black
		pke.wking = e.kingSq[SideWhite]
		pke.bking = e.kingSq[SideBlack]

		e.evalKingAndPawns(p)
		pke.score = e.score
		pke.passed = e.passedPawns
	} else {
		e.score = pke.score
		e.passedPawns = pke.passed
	}

	e.evalFirstPass(p)
	e.evalSecondPass(p)

	e.addFeature(fPawnValue, e.pieceCount[SideWhite][Pawn]-e.pieceCount[SideBlack][Pawn])
	e.addFeature(fKnightValue, e.pieceCount[SideWhite][Knight]-e.pieceCount[SideBlack][Knight])
	e.addFeature(fBishopValue, e.pieceCount[SideWhite][Bishop]-e.pieceCount[SideBlack][Bishop])
	e.addFeature(fRookValue, e.pieceCount[SideWhite][Rook]-e.pieceCount[SideBlack][Rook])
	e.addFeature(fQueenValue, e.pieceCount[SideWhite][Queen]-e.pieceCount[SideBlack][Queen])

	if p.WhiteMove {
		e.addFeature(fTempo, 1)
	} else {
		e.addFeature(fTempo, -1)
	}

	var phase = e.pieceCount[SideWhite][Knight] + e.pieceCount[SideBlack][Knight] +
		e.pieceCount[SideWhite][Bishop] + e.pieceCount[SideBlack][Bishop] +
		2*(e.pieceCount[SideWhite][Rook]+e.pieceCount[SideBlack][Rook]) +
		4*(e.pieceCount[SideWhite][Queen]+e.pieceCount[SideBlack][Queen])
	if phase > totalPhase {
		phase = totalPhase
	}
	e.phase = phase

	var result = (e.score.Mg()*phase + e.score.Eg()*(totalPhase-phase)) / (totalPhase * 100)
	var strongSide int
	if result > 0 {
		strongSide = SideWhite
	} else {
		strongSide = SideBlack
	}
	result = result * e.computeFactor(strongSide, p) / scaleNormal

	if !p.WhiteMove {
		result = -result
	}

	return result
}

func (e *EvaluationService) init(p *Position) {
	e.score = S(0, 0)

	for pt := Pawn; pt <= King; pt++ {
		e.pieceCount[SideWhite][pt] = 0
		e.pieceCount[SideBlack][pt] = 0

		e.attackedBy[SideWhite][pt] = 0
		e.attackedBy[SideBlack][pt] = 0
	}

	e.pieceCount[SideWhite][Pawn] = PopCount(p.Pawns & p.White)
	e.pieceCount[SideBlack][Pawn] = PopCount(p.Pawns & p.Black)

	e.passedPawns = 0

	e.attacked[SideWhite] = 0
	e.attacked[SideBlack] = 0
	e.attackedBy2[SideWhite] = 0
	e.attackedBy2[SideBlack] = 0
	e.kingAttackersCount[SideWhite] = 0
	e.kingAttackersCount[SideBlack] = 0
	e.kingAttacksCount[SideWhite] = 0
	e.kingAttacksCount[SideBlack] = 0

	e.kingSq[SideWhite] = FirstOne(p.Kings & p.White)
	e.kingSq[SideBlack] = FirstOne(p.Kings & p.Black)

	e.kingAreas[SideWhite] = kingAreaMasks[SideWhite][e.kingSq[SideWhite]]
	e.kingAreas[SideBlack] = kingAreaMasks[SideBlack][e.kingSq[SideBlack]]

	e.attackedBy[SideWhite][Pawn] = AllWhitePawnAttacks(p.Pawns & p.White)
	e.attackedBy[SideBlack][Pawn] = AllBlackPawnAttacks(p.Pawns & p.Black)

	e.pawnAttacksBy2[SideWhite] = UpLeft(p.Pawns&p.White) & UpRight(p.Pawns&p.White)
	e.pawnAttacksBy2[SideBlack] = DownLeft(p.Pawns&p.Black) & DownRight(p.Pawns&p.Black)

	e.attacked[SideWhite] |= e.attackedBy[SideWhite][Pawn]
	e.attacked[SideBlack] |= e.attackedBy[SideBlack][Pawn]

	e.mobilityArea[SideWhite] = ^(p.Pawns&p.White | e.attackedBy[SideBlack][Pawn])
	e.mobilityArea[SideBlack] = ^(p.Pawns&p.Black | e.attackedBy[SideWhite][Pawn])
}

func (e *EvaluationService) evalKingAndPawns(p *Position) {
	var x uint64
	var sq int

	for side := SideWhite; side <= SideBlack; side++ {
		var sign int
		var forward int
		if side == SideWhite {
			sign = 1
			forward = 8
		} else {
			sign = -1
			forward = -8
		}
		var US = side
		var THEM = side ^ 1
		var friendlyPawns = p.Colours(US) & p.Pawns
		var enemyPawns = p.Colours(THEM) & p.Pawns

		for x = friendlyPawns; x != 0; x &= x - 1 {
			sq = FirstOne(x)
			e.addComplexFeature(fPawnPST, relativeSq32(side, sq), sign)

			if PawnAttacksNew(THEM, sq)&friendlyPawns != 0 {
				e.addComplexFeature(fPawnProtected, relativeSq32(side, sq), sign)
			}
			if adjacentFilesMask[File(sq)]&ranks[Rank(sq)]&friendlyPawns != 0 {
				e.addComplexFeature(fPawnDuo, relativeSq32(side, sq), sign)
			}

			if adjacentFilesMask[File(sq)]&friendlyPawns == 0 {
				e.addFeature(fPawnIsolated, sign)
			}
			if FileMask[File(sq)]&^SquareMask[sq]&friendlyPawns != 0 {
				e.addFeature(fPawnDoubled, sign)
			}

			var stoppers = enemyPawns & passedPawnMasks[side][sq]
			// passed pawn
			if stoppers == 0 && upperRankMasks[US][Rank(sq)]&FileMask[File(sq)]&p.Pawns == 0 {
				e.passedPawns |= SquareMask[sq]
				var r = Max(0, relativeRankOf(side, sq)-Rank3)
				e.addComplexFeature(fPassedPawn, r, sign)
				var keySq = sq + forward
				e.addComplexFeature(fPassedEnemyKing, 8*r+distanceBetween[keySq][e.kingSq[THEM]], sign)
				e.addComplexFeature(fPassedOwnKing, 8*r+distanceBetween[keySq][e.kingSq[US]], sign)
			}
		}

		{
			// KING
			sq = e.kingSq[US]
			e.addComplexFeature(fKingPST, relativeSq32(side, sq), sign)

			for x = kingShieldMasks[US][sq] & friendlyPawns; x != 0; x &= x - 1 {
				var sq = FirstOne(x)
				e.addPst12(fKingShield, side, sq, sign)
			}

			/*for file := Max(FileA, File(sq)-1); file <= Min(FileH, File(sq)+1); file++ {
				var ours = friendly & p.Pawns & FileMask[file] & forwardRanksMasks[US][Rank(sq)]
				var ourDist int
				if ours == 0 {
					ourDist = 7
				} else {
					ourDist = Rank(sq) - Rank(Backmost(US, ours))
					if ourDist < 0 {
						ourDist = -ourDist
					}
				}
				e.addComplexFeature(fKingShield, 8*file+ourDist, sign)
			}*/
		}
	}
}

func (e *EvaluationService) evalFirstPass(p *Position) {
	var x, attacks uint64
	var sq int

	var occ = p.AllPieces()

	for side := SideWhite; side <= SideBlack; side++ {
		var sign int
		if side == SideWhite {
			sign = 1
		} else {
			sign = -1
		}
		var US = side
		var THEM = side ^ 1
		var friendly = p.Colours(US)
		var enemy = p.Colours(THEM)

		for x = p.Knights & friendly; x != 0; x &= x - 1 {
			sq = FirstOne(x)
			e.pieceCount[US][Knight]++
			e.addComplexFeature(fKnightPST, relativeSq32(side, sq), sign)

			attacks = KnightAttacks[sq]
			e.addComplexFeature(fKnightMobility, PopCount(attacks&e.mobilityArea[US]), sign)

			e.attackedBy2[US] |= e.attacked[US] & attacks
			e.attacked[US] |= attacks
			e.attackedBy[US][Knight] |= attacks

			attacks &= e.kingAreas[THEM] &^ e.pawnAttacksBy2[THEM]
			if attacks != 0 {
				e.kingAttackersCount[THEM]++
				e.kingAttacksCount[THEM] += PopCount(attacks)
			}

			if outpostSquares[side]&SquareMask[sq] != 0 &&
				outpostSquareMasks[US][sq]&enemy&p.Pawns == 0 {
				e.addFeature(fKnightOutpost, sign)
			}
		}

		for x = p.Bishops & friendly; x != 0; x &= x - 1 {
			sq = FirstOne(x)
			e.pieceCount[US][Bishop]++
			e.addComplexFeature(fBishopPST, relativeSq32(side, sq), sign)

			attacks = BishopAttacks(sq, occ)
			e.addComplexFeature(fBishopMobility, PopCount(attacks&e.mobilityArea[US]), sign)

			e.attackedBy2[US] |= e.attacked[US] & attacks
			e.attacked[US] |= attacks
			e.attackedBy[US][Bishop] |= attacks

			attacks &= e.kingAreas[THEM] &^ e.pawnAttacksBy2[THEM]
			if attacks != 0 {
				e.kingAttackersCount[THEM]++
				e.kingAttacksCount[THEM] += PopCount(attacks)
			}

			if side == SideWhite {
				e.addFeature(fBishopRammedPawns,
					PopCount(sameColorSquares(sq)&p.Pawns&p.White&Down(p.Pawns&p.Black)))
			} else {
				e.addFeature(fBishopRammedPawns,
					-PopCount(sameColorSquares(sq)&p.Pawns&p.Black&Up(p.Pawns&p.White)))
			}
		}

		for x = p.Rooks & friendly; x != 0; x &= x - 1 {
			sq = FirstOne(x)
			e.pieceCount[US][Rook]++
			e.addComplexFeature(fRookPST, relativeSq32(side, sq), sign)

			attacks = RookAttacks(sq, occ&^(friendly&p.Rooks))
			e.addComplexFeature(fRookMobility, PopCount(attacks&e.mobilityArea[US]), sign)

			e.attackedBy2[US] |= e.attacked[US] & attacks
			e.attacked[US] |= attacks
			e.attackedBy[US][Rook] |= attacks

			attacks &= e.kingAreas[THEM] &^ e.pawnAttacksBy2[THEM]
			if attacks != 0 {
				e.kingAttackersCount[THEM]++
				e.kingAttacksCount[THEM] += PopCount(attacks)
			}

			attacks = FileMask[File(sq)]
			if (attacks & friendly & p.Pawns) == 0 {
				if (attacks & p.Pawns) == 0 {
					e.addFeature(fRookOpen, sign)
				} else {
					e.addFeature(fRookSemiopen, sign)
				}
			}
		}

		for x = p.Queens & friendly; x != 0; x &= x - 1 {
			sq = FirstOne(x)
			e.pieceCount[US][Queen]++
			e.addComplexFeature(fQueenPST, relativeSq32(side, sq), sign)

			attacks = QueenAttacks(sq, occ)
			e.addComplexFeature(fQueenMobility, PopCount(attacks&e.mobilityArea[US]), sign)

			e.attackedBy2[US] |= e.attacked[US] & attacks
			e.attacked[US] |= attacks
			e.attackedBy[US][Queen] |= attacks

			attacks &= e.kingAreas[THEM] &^ e.pawnAttacksBy2[THEM]
			if attacks != 0 {
				e.kingAttackersCount[THEM]++
				e.kingAttacksCount[THEM] += PopCount(attacks)
			}

			e.addFeature(fKingQueenTropism, sign*distanceBetween[sq][e.kingSq[THEM]])
		}

		{
			// KING
			sq = e.kingSq[US]

			attacks = KingAttacks[sq]
			e.attackedBy2[US] |= e.attacked[US] & attacks
			e.attacked[US] |= attacks
			e.attackedBy[US][King] |= attacks
		}

		if e.pieceCount[US][Bishop] >= 2 {
			e.addFeature(fBishopPair, sign)
		}
	}

	e.addFeature(fMinorBehindPawn,
		PopCount((p.Knights|p.Bishops)&p.White&Down(p.Pawns))-
			PopCount((p.Knights|p.Bishops)&p.Black&Up(p.Pawns)))

	e.addFeature(fMinorProtected,
		PopCount((p.Knights|p.Bishops)&p.White&e.attackedBy[SideWhite][Pawn])-
			PopCount((p.Knights|p.Bishops)&p.Black&e.attackedBy[SideBlack][Pawn]))
}

var kingAttackWeight = [...]int{2, 4, 8, 12, 13, 14, 15, 16}

func (e *EvaluationService) evalSecondPass(p *Position) {
	var occ = p.AllPieces()

	for side := SideWhite; side <= SideBlack; side++ {
		var sign int
		if side == SideWhite {
			sign = 1
		} else {
			sign = -1
		}
		var US = side
		var THEM = side ^ 1
		var friendly = p.Colours(US)
		var enemy = p.Colours(THEM)

		//if e.kingAttackersCount[THEM] > 1-e.pieceCount[US][Queen]
		{
			// king safety

			var val = sign * kingAttackWeight[Min(len(kingAttackWeight)-1, e.kingAttackersCount[THEM])]
			//var val = sign

			weak := e.attacked[US] & ^e.attackedBy2[THEM] & (^e.attacked[THEM] | e.attackedBy[THEM][Queen] | e.attackedBy[THEM][King])
			safe := ^friendly & (^e.attacked[THEM] | (weak & e.attackedBy2[US]))

			knightThreats := KnightAttacks[e.kingSq[THEM]]
			bishopThreats := BishopAttacks(e.kingSq[THEM], occ)
			rookThreats := RookAttacks(e.kingSq[THEM], occ)
			queenThreats := bishopThreats | rookThreats

			e.addFeature(fSafetyKnightCheck, val*PopCount(knightThreats&safe&e.attackedBy[US][Knight]))
			e.addFeature(fSafetyBishopCheck, val*PopCount(bishopThreats&safe&e.attackedBy[US][Bishop]))
			e.addFeature(fSafetyRookCheck, val*PopCount(rookThreats&safe&e.attackedBy[US][Rook]))
			e.addFeature(fSafetyQueenCheck, val*PopCount(queenThreats&safe&e.attackedBy[US][Queen]))
			e.addFeature(fSafetyWeakSquares, val*PopCount(e.kingAreas[THEM]&weak))

			/*var safety = 0
			if e.pieceCount[US][Queen] == 0 {
				safety += -2
			}
			safety += 4*PopCount(knightThreats&safe&e.attackedBy[US][Knight]) +
				1*PopCount(bishopThreats&safe&e.attackedBy[US][Bishop]) +
				4*PopCount(rookThreats&safe&e.attackedBy[US][Rook]) +
				2*PopCount(queenThreats&safe&e.attackedBy[US][Queen]) +
				1*PopCount(e.kingAreas[THEM]&weak)
			if safety > 0 {
				e.addFeature(fSafetyMain, sign*safety*safety)
			}*/
		}

		{
			// threats

			var knights = friendly & p.Knights
			var bishops = friendly & p.Bishops
			var rooks = friendly & p.Rooks
			var queens = friendly & p.Queens

			var attacksByPawns = e.attackedBy[THEM][Pawn]
			var attacksByMinors = e.attackedBy[THEM][Knight] | e.attackedBy[THEM][Bishop]
			var attacksByMajors = e.attackedBy[THEM][Rook] | e.attackedBy[THEM][Queen]

			var poorlyDefended = (e.attacked[THEM] & ^e.attacked[US]) |
				(e.attackedBy2[THEM] & ^e.attackedBy2[US] & ^e.attackedBy[US][Pawn])

			var weakMinors = (knights | bishops) & poorlyDefended

			e.addFeature(fThreatWeakPawn, sign*PopCount(friendly&p.Pawns & ^attacksByPawns & poorlyDefended))
			e.addFeature(fThreatMinorAttackedByPawn, sign*PopCount((knights|bishops)&attacksByPawns))
			e.addFeature(fThreatMinorAttackedByMinor, sign*PopCount((knights|bishops)&attacksByMinors))
			e.addFeature(fThreatMinorAttackedByMajor, sign*PopCount(weakMinors&attacksByMajors))
			e.addFeature(fThreatRookAttackedByLesser, sign*PopCount(rooks&(attacksByPawns|attacksByMinors)))
			e.addFeature(fThreatMinorAttackedByKing, sign*PopCount(weakMinors&e.attackedBy[THEM][King]))
			e.addFeature(fThreatRookAttackedByKing, sign*PopCount(rooks&poorlyDefended&e.attackedBy[THEM][King]))
			e.addFeature(fThreatQueenAttackedByOne, sign*PopCount(queens&e.attacked[THEM]))
		}

		for temp := e.passedPawns & friendly; temp != 0; temp &= temp - 1 {
			var sq = FirstOne(temp)
			var keySq = sq + 8*sign

			var r = Max(0, relativeRankOf(side, sq)-Rank3)

			if enemy&SquareMask[keySq] == 0 {
				e.addComplexFeature(fPassedCanMove, r, sign)
			}

			if (SquareMask[keySq] & e.attacked[THEM]) == 0 {
				e.addComplexFeature(fPassedSafeMove, r, sign)
			}
		}
	}
}

const (
	scaleNormal = 128
)

const (
	QueenSideBB = FileAMask | FileBMask | FileCMask | FileDMask
	KingSideBB  = FileEMask | FileFMask | FileGMask | FileHMask
)

func (e *EvaluationService) computeFactor(own int, p *Position) int {
	var them = own ^ 1
	var ownPawns = e.pieceCount[own][Pawn]
	if ownPawns <= 1 {
		var ownForce = computeForce(e, own)
		var theirForce = computeForce(e, own^1)
		if ownPawns == 0 {
			if ownForce <= 4 {
				return scaleNormal * 1 / 16
			}
			if ownForce-theirForce <= 4 {
				return scaleNormal * 1 / 4
			}
		} else if ownPawns == 1 {
			var theirMinor = e.pieceCount[them][Knight]+e.pieceCount[them][Bishop] != 0
			if ownForce <= 4 && theirMinor {
				return scaleNormal * 1 / 8
			}
			if ownForce == theirForce && theirMinor {
				return scaleNormal * 1 / 2
			}
		}
	}
	var strong = p.Colours(own)

	var strongPawnCount = e.pieceCount[own][Pawn]
	var x = 8 - strongPawnCount
	var pawnScale = 128 - x*x

	if strong&p.Pawns&QueenSideBB == 0 ||
		strong&p.Pawns&KingSideBB == 0 {
		pawnScale -= 20
	}

	//var pawnScale = scaleNormal

	if e.pieceCount[SideWhite][Bishop] == 1 &&
		e.pieceCount[SideBlack][Bishop] == 1 &&
		onlyOne(p.Bishops&darkSquares) {
		if p.Knights|p.Rooks|p.Queens == 0 {
			pawnScale = Min(pawnScale, scaleNormal*1/2)
		}
	}

	return pawnScale
}

func computeForce(e *EvaluationService, side int) int {
	return 4*(e.pieceCount[side][Knight]+e.pieceCount[side][Bishop]) +
		6*e.pieceCount[side][Rook] +
		12*e.pieceCount[side][Queen]
}

func (e *EvaluationService) StartingWeights() []float64 {
	var material = []float64{100, 100, 325, 325, 325, 325, 500, 500, 1000, 1000}
	var result = make([]float64, 2*totalFeatureSize)
	copy(result, material)
	return result
}

func (e *EvaluationService) ComputeFeatures(pos *Position) domain.TuneEntry {
	for i := range e.features {
		e.features[i] = 0
	}
	e.Evaluate(pos)
	var size int
	for _, v := range e.features {
		if v != 0 {
			size++
		}
	}
	var features = make([]domain.FeatureInfo, 0, size)
	for i, v := range e.features {
		if v != 0 {
			features = append(features, domain.FeatureInfo{Index: int16(i), Value: int16(v)})
		}
	}
	var result = domain.TuneEntry{
		Features:         features,
		MgPhase:          float32(e.phase) / totalPhase,
		WhiteStrongScale: float32(e.computeFactor(SideWhite, pos)) / scaleNormal,
		BlackStrongScale: float32(e.computeFactor(SideBlack, pos)) / scaleNormal,
	}
	result.EgPhase = 1 - result.MgPhase
	return result
}

func (e *EvaluationService) addPst12(feature, side, sq, value int) {
	e.addComplexFeature(feature, file4(sq), value)
	e.addComplexFeature(feature, 4+relativeRankOf(side, sq), value)
}

func (e *EvaluationService) addMobility(feature, side, sq, mobility, value int) {
	value *= sqrtInt[mobility]
	e.addComplexFeature(feature, file4(sq), value)
	e.addComplexFeature(feature, 4+relativeRankOf(side, sq), value)
}

func (e *EvaluationService) addFeature(feature, value int) {
	e.addComplexFeature(feature, 0, value)
}

func (e *EvaluationService) addComplexFeature(feature, featureIndex, value int) {
	var info = &infos[feature]
	var index = info.StartIndex + featureIndex
	var w = e.weights[index]
	//e.score.mg += value * w.mg
	//e.score.eg += value * w.eg
	e.score += Score(value) * w
	if e.tuning {
		if featureIndex >= info.Size {
			panic(errAddComplexFeature)
		}
		e.features[index] += value
	}
}
