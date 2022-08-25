package pgn

import (
	"testing"
)

func TestPgn(t *testing.T) {
	var game, err = ParseGame(testPgn)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(game)
}

const testPgn = `[Event "CCRL 40/15"]
[Site "CCRL"]
[Date "2021.10.06"]
[Round "792.1.311"]
[White "Demolito 2021-07-09 64-bit"]
[Black "Counter 4.0 64-bit"]
[Result "1/2-1/2"]
[ECO "B36"]
[Opening "Sicilian"]
[Variation "accelerated fianchetto, Maroczy bind"]
[PlyCount "195"]
[WhiteElo "3101"]
[BlackElo "3091"]

1. e4 {+0.00/1 0s} c5 {+0.00/1 0s} 2. Nf3 {+0.00/1 0s} Nc6 {+0.00/1 0s} 3. d4
{+0.00/1 0s} cxd4 {+0.00/1 0s} 4. Nxd4 {+0.00/1 0s} g6 {+0.00/1 0s} 5. c4
{+0.00/1 0s} Nf6 {+0.00/1 0s} 6. Nc3 {+0.00/1 0s} d6 {+0.00/1 0s} 7. Be2
{+0.00/1 0s} Nxd4 {+0.00/1 0s} 8. Qxd4 {+0.00/1 0s} Bg7 {+0.00/1 0s} 9. Be3
{+0.36/25 56s} O-O {(O-O) -0.35/22 26s} 10. Qd2 {(Qd2) +0.38/22 19s} Be6 {(Ng4)
-0.50/21 35s} 11. O-O {(O-O) +0.52/21 23s} Qa5 {(Ng4) -0.49/22 36s} 12. Rac1
{(Rad1) +0.50/22 30s} a6 {(Rac8) -0.46/22 32s} 13. Rfd1 {(b3) +0.26/20 64s}
Rfc8 {(Rfc8) -0.16/22 24s} 14. h3 {(b3) +0.23/21 61s} b5 {(b5) -0.30/20 22s}
15. Bf3 {(cxb5) +0.04/21 22s} bxc4 {(Rab8) +0.00/24 23s} 16. e5 {(e5) +0.16/23
16s} Qxe5 {(Qxe5) +0.00/25 22s} 17. Bxa8 {(Bxa8) +0.18/24 19s} Rxa8 {(Rxa8)
+0.06/27 29s} 18. Bd4 {(Re1) +0.03/23 21s} Qh5 {(Qa5) +0.00/25 38s} 19. Ne2
{(Ne2) +0.02/22 35s} Qb5 {(Qb5) +0.00/25 32s} 20. Qe3 {(Bxf6) -0.03/21 26s} Rb8
{(Re8) +0.20/24 39s} 21. Nc3 {(a4) -0.10/24 43s} Qd7 {(Qa5) +0.10/21 37s} 22.
Rd2 {(Rd2) +0.10/23 17s} Qc6 {(a5) +0.00/24 27s} 23. Re1 {(Re1) +0.08/23 33s}
Rb4 {(a5) +0.00/26 28s} 24. g4 {(Rde2) +0.21/23 54s} h5 {(h5) +0.00/22 47s} 25.
f3 {(f3) +0.18/20 16s} Qb7 {(Qb7) +0.00/21 20s} 26. Ree2 {(Ree2) +0.16/21 16s}
Kf8 {(Kf8) -0.15/22 23s} 27. Kf2 {(Rg2) +0.16/19 24s} a5 {(Kg8) +0.00/22 27s}
28. Qg5 {(Bxf6) +0.07/19 29s} Qa8 {(Qa8) +0.09/21 38s} 29. Qf4 {(Kg3) +0.19/20
35s} hxg4 {(hxg4) +0.07/21 35s} 30. hxg4 {(hxg4) +0.11/21 18s} Rb8 {(Kg8)
+0.00/23 45s} 31. Kg1 {(Bxf6) +0.25/22 36s} Kg8 {(a4) +0.00/21 24s} 32. Rh2
{(Rh2) +0.06/21 18s} Qc6 {(a4) +0.00/23 29s} 33. Rd1 {(Qe3) +0.31/20 28s} a4
{(Nd5) +0.00/22 24s} 34. a3 {(Ne4) +0.09/22 58s} Rb3 {(Rb3) +0.00/22 24s} 35.
Rdd2 {(Qe3) +0.10/21 20s} Qb7 {(Rb8) +0.25/21 18s} 36. Rhe2 {(Rhe2) -0.02/23
33s} Kf8 {+0.20/22 30s} 37. Rf2 {(Rf2) -0.02/23 31s} Ke8 {(Kg8) +0.00/22 24s}
38. Qg5 {(Qg5) +0.34/21 27s} Qa8 {(Qc6) +0.00/23 21s} 39. Qh4 {(Qh4) +0.22/21
11s} Kd7 {(Kd7) +0.00/25 39s} 40. Qh1 {(Qh2) +0.21/19 6s} Qa5 {(Qc6) +0.00/22
17s} 41. Rfe2 {(Rfe2) +0.15/19 16s} Kd8 {(Rb8) +0.00/22 69s} 42. Qh4 {(Qh2)
+0.26/22 41s} Kd7 {(Kd7) -0.30/21 28s} 43. Rh2 {(f4) +0.36/20 49s} Rb8 {(Nxg4)
+0.00/22 18s} 44. Qf2 {(Rd1) +0.46/19 24s} Ke8 {(Qa8) +0.00/23 17s} 45. Kf1
{(Rd1) +0.35/21 50s} Rb3 {(Kd7) +0.00/23 18s} 46. Qh4 {(Kg1) +0.39/21 35s} g5
{(Kd8) +0.00/21 20s} 47. Qe1 {(Qe1) +0.56/20 13s} Qa8 {(Bc8) +0.00/24 28s} 48.
Qe2 {(Rhf2) +0.32/21 13s} Kf8 {(Qb7) +0.00/22 19s} 49. Rh3 {(Rh3) +0.00/21 40s}
Qb7 {(Qb7) +0.00/24 27s} 50. Kf2 {(Qh2) +0.06/21 29s} Kg8 {(Bd7) +0.00/22 17s}
51. Ke1 {(Kf1) +0.42/21 13s} Kf8 {(Kf8) +0.00/23 18s} 52. Kd1 {(Qh2) +0.23/23
47s} Bd7 {(Nd5) +0.00/20 21s} 53. Kc1 {(Kc1) +0.26/22 26s} Bc6 {(e5) +0.00/21
21s} 54. f4 {(Qxc4) +0.89/20 15s} gxf4 {(gxf4) -1.24/21 22s} 55. g5 {(g5)
+0.68/21 12s} Ng8 {(Ng8) -0.56/22 20s} 56. Bxg7+ {(Bxg7) +0.65/21 12s} Kxg7
{(Kxg7) -0.57/23 19s} 57. Qxc4 {(Qxc4) +0.63/21 16s} Be4 {(Be4) -0.38/23 20s}
58. Qd4+ {(Rdh2) +0.95/20 21s} f6 {(f6) -0.21/20 19s} 59. Rhh2 {(Rhh2) +1.19/22
41s} Qc6 {(Bg6) -0.52/22 36s} 60. Kd1 {(Kd1) +0.22/19 13s} Bf3+ {(Bf3) -0.45/22
30s} 61. Ke1 {(Ke1) +0.49/20 13s} e5 {(e5) -0.68/22 22s} 62. Qd3 {(Qd3)
+0.57/21 34s} e4 {(e4) -0.81/22 19s} 63. Qd4 {(Qd4) +0.42/22 14s} Qb6 {(Qb6)
-0.83/23 22s} 64. gxf6+ {(gxf6) +0.42/22 15s} Nxf6 {(Nxf6) -0.83/25 21s} 65.
Qxb6 {(Qxb6) +0.47/24 14s} Rxb6 {(Rxb6) -0.81/25 28s} 66. Rh4 {(Rh4) +0.54/22
19s} e3 {(Nh5) -0.84/25 46s} 67. Rdh2 {(Rdh2) +1.39/25 13s} Ng4 {(Ng4) -0.84/25
21s} 68. Rxg4+ {(Rxg4) +1.45/26 14s} Bxg4 {(Bxg4) -0.84/27 37s} 69. Rg2 {(Rg2)
+1.52/26 13s} Kh6 {(Kh6) -0.85/25 15s} 70. Rxg4 {(Rxg4) +1.52/27 14s} Rxb2
{(Rxb2) -0.86/26 17s} 71. Rxf4 {(Rxf4) +1.60/27 14s} Rb3 {(Rb3) -0.88/28 32s}
72. Nxa4 {(Nxa4) +1.60/29 15s} Kg5 {(Rxa3) -0.89/28 29s} 73. Rb4 {(Rf8)
+1.60/28 16s} Rxa3 {(Rxa3) -0.88/25 15s} 74. Nb6 {(Ke2) +1.60/25 18s} Kf5
{(Ra2) -0.87/24 19s} 75. Ke2 {(Ke2) +1.60/27 21s} Rc3 {(Rc3) -0.88/24 15s} 76.
Rh4 {(Rh4) +1.60/29 21s} Ke6 {(Ke6) -0.88/25 18s} 77. Nc4 {(Nc4) +1.60/32 21s}
Kd5 {(Kd5) -0.89/26 21s} 78. Nxe3+ {(Nxe3) +1.60/33 24s} Kc6 {(Kc6) -0.89/27
20s} 79. Nf5 {(Nf5) +1.60/32 29s} Kd5 {(Rc5) -0.89/26 20s} 80. Kd2 {(Rh5)
+1.60/30 41s} Rf3 {(Rc7) -0.89/27 24s} 81. Ne3+ {(Rh5) +1.61/30 14s} Kc6 {(Kc6)
-0.89/27 18s} 82. Rh6 {(Rh5) +1.60/28 14s} Rf4 {(Rf4) -0.90/25 20s} 83. Kd3
{(Kd3) +1.61/31 15s} Kc5 {(Kc5) -0.91/28 34s} 84. Rh5+ {(Rh5) +1.61/29 14s} Kc6
{(Kc6) -0.93/28 23s} 85. Nf5 {(Nf5) +1.61/30 15s} Rf1 {(Rf1) -0.91/28 23s} 86.
Nd4+ {(Nd4) +1.66/28 15s} Kb6 {(Kb6) -0.91/29 25s} 87. Nb3 {(Kc4) +1.66/28 15s}
Re1 {(Kc6) -0.91/28 27s} 88. Rh6 {(Nd4) +1.69/29 15s} Kc6 {(Kc6) -0.91/29 30s}
89. Kc4 {(Kc4) +1.66/31 16s} Re5 {(Re5) -0.91/29 30s} 90. Nd4+ {(Nd4) +1.80/30
16s} Kd7 {(Kd7) -0.97/29 38s} 91. Kd3 {(Rh8) +1.76/27 16s} Re1 {(Re1) -1.02/27
22s} 92. Nf5 {(Nf5) +1.84/31 16s} Rd1+ {(Rd1) -1.05/25 18s} 93. Kc3 {(Kc3)
+1.84/33 16s} Kc7 {(Rd5) -1.08/25 18s} 94. Nd4 {(Rf6) +1.84/31 16s} Rc1+ {(Rc1)
-1.08/25 20s} 95. Kd2 {(Kb3) +1.89/31 17s} Rb1 {(Rb1) -1.09/26 18s} 96. Kc2
{(Kc2) +1.95/28 17s} Rg1 {(Rb6) -1.10/27 30s} 97. Nb5+ {(Nb5) +2.01/27 18s} Kc6
{(Kc6) -1.10/29 20s} 98. Nxd6 {(Nxd6) +2.02/27 18s} 1/2-1/2

`
