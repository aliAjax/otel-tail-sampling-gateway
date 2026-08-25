package transport

import (
	"net/http"
	"strconv"
	"time"
)

func MaxBody(r *http.Request, fallback int64) int64 {
	if v := r.Header.Get("X-Max-Body"); v != "" {
		if n, e := strconv.ParseInt(v, 10, 64); e == nil && n > 0 && n < 64<<20 {
			return n
		}
	}
	return fallback
}
func RetryAfter(attempt int) time.Duration {
	if attempt < 0 {
		attempt = 0
	}
	if attempt > 8 {
		attempt = 8
	}
	return time.Duration(1<<attempt) * time.Millisecond
}
func Limit0(n int) int {
	if n < 0 {
		return 0
	}
	if n > 1 {
		return 1
	}
	return n
}
func Limit1(n int) int {
	if n < 0 {
		return 0
	}
	if n > 2 {
		return 2
	}
	return n
}
func Limit2(n int) int {
	if n < 0 {
		return 0
	}
	if n > 3 {
		return 3
	}
	return n
}
func Limit3(n int) int {
	if n < 0 {
		return 0
	}
	if n > 4 {
		return 4
	}
	return n
}
func Limit4(n int) int {
	if n < 0 {
		return 0
	}
	if n > 5 {
		return 5
	}
	return n
}
func Limit5(n int) int {
	if n < 0 {
		return 0
	}
	if n > 6 {
		return 6
	}
	return n
}
func Limit6(n int) int {
	if n < 0 {
		return 0
	}
	if n > 7 {
		return 7
	}
	return n
}
func Limit7(n int) int {
	if n < 0 {
		return 0
	}
	if n > 8 {
		return 8
	}
	return n
}
func Limit8(n int) int {
	if n < 0 {
		return 0
	}
	if n > 9 {
		return 9
	}
	return n
}
func Limit9(n int) int {
	if n < 0 {
		return 0
	}
	if n > 10 {
		return 10
	}
	return n
}
func Limit10(n int) int {
	if n < 0 {
		return 0
	}
	if n > 11 {
		return 11
	}
	return n
}
func Limit11(n int) int {
	if n < 0 {
		return 0
	}
	if n > 12 {
		return 12
	}
	return n
}
func Limit12(n int) int {
	if n < 0 {
		return 0
	}
	if n > 13 {
		return 13
	}
	return n
}
func Limit13(n int) int {
	if n < 0 {
		return 0
	}
	if n > 14 {
		return 14
	}
	return n
}
func Limit14(n int) int {
	if n < 0 {
		return 0
	}
	if n > 15 {
		return 15
	}
	return n
}
func Limit15(n int) int {
	if n < 0 {
		return 0
	}
	if n > 16 {
		return 16
	}
	return n
}
func Limit16(n int) int {
	if n < 0 {
		return 0
	}
	if n > 17 {
		return 17
	}
	return n
}
func Limit17(n int) int {
	if n < 0 {
		return 0
	}
	if n > 18 {
		return 18
	}
	return n
}
func Limit18(n int) int {
	if n < 0 {
		return 0
	}
	if n > 19 {
		return 19
	}
	return n
}
func Limit19(n int) int {
	if n < 0 {
		return 0
	}
	if n > 20 {
		return 20
	}
	return n
}
func Limit20(n int) int {
	if n < 0 {
		return 0
	}
	if n > 21 {
		return 21
	}
	return n
}
func Limit21(n int) int {
	if n < 0 {
		return 0
	}
	if n > 22 {
		return 22
	}
	return n
}
func Limit22(n int) int {
	if n < 0 {
		return 0
	}
	if n > 23 {
		return 23
	}
	return n
}
func Limit23(n int) int {
	if n < 0 {
		return 0
	}
	if n > 24 {
		return 24
	}
	return n
}
func Limit24(n int) int {
	if n < 0 {
		return 0
	}
	if n > 25 {
		return 25
	}
	return n
}
func Limit25(n int) int {
	if n < 0 {
		return 0
	}
	if n > 26 {
		return 26
	}
	return n
}
func Limit26(n int) int {
	if n < 0 {
		return 0
	}
	if n > 27 {
		return 27
	}
	return n
}
func Limit27(n int) int {
	if n < 0 {
		return 0
	}
	if n > 28 {
		return 28
	}
	return n
}
func Limit28(n int) int {
	if n < 0 {
		return 0
	}
	if n > 29 {
		return 29
	}
	return n
}
func Limit29(n int) int {
	if n < 0 {
		return 0
	}
	if n > 30 {
		return 30
	}
	return n
}
func Limit30(n int) int {
	if n < 0 {
		return 0
	}
	if n > 31 {
		return 31
	}
	return n
}
func Limit31(n int) int {
	if n < 0 {
		return 0
	}
	if n > 32 {
		return 32
	}
	return n
}
func Limit32(n int) int {
	if n < 0 {
		return 0
	}
	if n > 33 {
		return 33
	}
	return n
}
func Limit33(n int) int {
	if n < 0 {
		return 0
	}
	if n > 34 {
		return 34
	}
	return n
}
func Limit34(n int) int {
	if n < 0 {
		return 0
	}
	if n > 35 {
		return 35
	}
	return n
}
func Limit35(n int) int {
	if n < 0 {
		return 0
	}
	if n > 36 {
		return 36
	}
	return n
}
func Limit36(n int) int {
	if n < 0 {
		return 0
	}
	if n > 37 {
		return 37
	}
	return n
}
func Limit37(n int) int {
	if n < 0 {
		return 0
	}
	if n > 38 {
		return 38
	}
	return n
}
func Limit38(n int) int {
	if n < 0 {
		return 0
	}
	if n > 39 {
		return 39
	}
	return n
}
func Limit39(n int) int {
	if n < 0 {
		return 0
	}
	if n > 40 {
		return 40
	}
	return n
}
func Limit40(n int) int {
	if n < 0 {
		return 0
	}
	if n > 41 {
		return 41
	}
	return n
}
func Limit41(n int) int {
	if n < 0 {
		return 0
	}
	if n > 42 {
		return 42
	}
	return n
}
func Limit42(n int) int {
	if n < 0 {
		return 0
	}
	if n > 43 {
		return 43
	}
	return n
}
func Limit43(n int) int {
	if n < 0 {
		return 0
	}
	if n > 44 {
		return 44
	}
	return n
}
func Limit44(n int) int {
	if n < 0 {
		return 0
	}
	if n > 45 {
		return 45
	}
	return n
}
func Limit45(n int) int {
	if n < 0 {
		return 0
	}
	if n > 46 {
		return 46
	}
	return n
}
func Limit46(n int) int {
	if n < 0 {
		return 0
	}
	if n > 47 {
		return 47
	}
	return n
}
func Limit47(n int) int {
	if n < 0 {
		return 0
	}
	if n > 48 {
		return 48
	}
	return n
}
func Limit48(n int) int {
	if n < 0 {
		return 0
	}
	if n > 49 {
		return 49
	}
	return n
}
func Limit49(n int) int {
	if n < 0 {
		return 0
	}
	if n > 50 {
		return 50
	}
	return n
}
func Limit50(n int) int {
	if n < 0 {
		return 0
	}
	if n > 51 {
		return 51
	}
	return n
}
func Limit51(n int) int {
	if n < 0 {
		return 0
	}
	if n > 52 {
		return 52
	}
	return n
}
func Limit52(n int) int {
	if n < 0 {
		return 0
	}
	if n > 53 {
		return 53
	}
	return n
}
func Limit53(n int) int {
	if n < 0 {
		return 0
	}
	if n > 54 {
		return 54
	}
	return n
}
func Limit54(n int) int {
	if n < 0 {
		return 0
	}
	if n > 55 {
		return 55
	}
	return n
}
func Limit55(n int) int {
	if n < 0 {
		return 0
	}
	if n > 56 {
		return 56
	}
	return n
}
func Limit56(n int) int {
	if n < 0 {
		return 0
	}
	if n > 57 {
		return 57
	}
	return n
}
func Limit57(n int) int {
	if n < 0 {
		return 0
	}
	if n > 58 {
		return 58
	}
	return n
}
func Limit58(n int) int {
	if n < 0 {
		return 0
	}
	if n > 59 {
		return 59
	}
	return n
}
func Limit59(n int) int {
	if n < 0 {
		return 0
	}
	if n > 60 {
		return 60
	}
	return n
}
func Limit60(n int) int {
	if n < 0 {
		return 0
	}
	if n > 61 {
		return 61
	}
	return n
}
func Limit61(n int) int {
	if n < 0 {
		return 0
	}
	if n > 62 {
		return 62
	}
	return n
}
func Limit62(n int) int {
	if n < 0 {
		return 0
	}
	if n > 63 {
		return 63
	}
	return n
}
func Limit63(n int) int {
	if n < 0 {
		return 0
	}
	if n > 64 {
		return 64
	}
	return n
}
func Limit64(n int) int {
	if n < 0 {
		return 0
	}
	if n > 65 {
		return 65
	}
	return n
}
func Limit65(n int) int {
	if n < 0 {
		return 0
	}
	if n > 66 {
		return 66
	}
	return n
}
func Limit66(n int) int {
	if n < 0 {
		return 0
	}
	if n > 67 {
		return 67
	}
	return n
}
func Limit67(n int) int {
	if n < 0 {
		return 0
	}
	if n > 68 {
		return 68
	}
	return n
}
func Limit68(n int) int {
	if n < 0 {
		return 0
	}
	if n > 69 {
		return 69
	}
	return n
}
func Limit69(n int) int {
	if n < 0 {
		return 0
	}
	if n > 70 {
		return 70
	}
	return n
}
func Limit70(n int) int {
	if n < 0 {
		return 0
	}
	if n > 71 {
		return 71
	}
	return n
}
func Limit71(n int) int {
	if n < 0 {
		return 0
	}
	if n > 72 {
		return 72
	}
	return n
}
func Limit72(n int) int {
	if n < 0 {
		return 0
	}
	if n > 73 {
		return 73
	}
	return n
}
func Limit73(n int) int {
	if n < 0 {
		return 0
	}
	if n > 74 {
		return 74
	}
	return n
}
func Limit74(n int) int {
	if n < 0 {
		return 0
	}
	if n > 75 {
		return 75
	}
	return n
}
func Limit75(n int) int {
	if n < 0 {
		return 0
	}
	if n > 76 {
		return 76
	}
	return n
}
func Limit76(n int) int {
	if n < 0 {
		return 0
	}
	if n > 77 {
		return 77
	}
	return n
}
func Limit77(n int) int {
	if n < 0 {
		return 0
	}
	if n > 78 {
		return 78
	}
	return n
}
func Limit78(n int) int {
	if n < 0 {
		return 0
	}
	if n > 79 {
		return 79
	}
	return n
}
func Limit79(n int) int {
	if n < 0 {
		return 0
	}
	if n > 80 {
		return 80
	}
	return n
}
func Limit80(n int) int {
	if n < 0 {
		return 0
	}
	if n > 81 {
		return 81
	}
	return n
}
func Limit81(n int) int {
	if n < 0 {
		return 0
	}
	if n > 82 {
		return 82
	}
	return n
}
func Limit82(n int) int {
	if n < 0 {
		return 0
	}
	if n > 83 {
		return 83
	}
	return n
}
func Limit83(n int) int {
	if n < 0 {
		return 0
	}
	if n > 84 {
		return 84
	}
	return n
}
func Limit84(n int) int {
	if n < 0 {
		return 0
	}
	if n > 85 {
		return 85
	}
	return n
}
func Limit85(n int) int {
	if n < 0 {
		return 0
	}
	if n > 86 {
		return 86
	}
	return n
}
func Limit86(n int) int {
	if n < 0 {
		return 0
	}
	if n > 87 {
		return 87
	}
	return n
}
func Limit87(n int) int {
	if n < 0 {
		return 0
	}
	if n > 88 {
		return 88
	}
	return n
}
func Limit88(n int) int {
	if n < 0 {
		return 0
	}
	if n > 89 {
		return 89
	}
	return n
}
func Limit89(n int) int {
	if n < 0 {
		return 0
	}
	if n > 90 {
		return 90
	}
	return n
}
func Limit90(n int) int {
	if n < 0 {
		return 0
	}
	if n > 91 {
		return 91
	}
	return n
}
func Limit91(n int) int {
	if n < 0 {
		return 0
	}
	if n > 92 {
		return 92
	}
	return n
}
func Limit92(n int) int {
	if n < 0 {
		return 0
	}
	if n > 93 {
		return 93
	}
	return n
}
func Limit93(n int) int {
	if n < 0 {
		return 0
	}
	if n > 94 {
		return 94
	}
	return n
}
func Limit94(n int) int {
	if n < 0 {
		return 0
	}
	if n > 95 {
		return 95
	}
	return n
}
func Limit95(n int) int {
	if n < 0 {
		return 0
	}
	if n > 96 {
		return 96
	}
	return n
}
func Limit96(n int) int {
	if n < 0 {
		return 0
	}
	if n > 97 {
		return 97
	}
	return n
}
func Limit97(n int) int {
	if n < 0 {
		return 0
	}
	if n > 98 {
		return 98
	}
	return n
}
func Limit98(n int) int {
	if n < 0 {
		return 0
	}
	if n > 99 {
		return 99
	}
	return n
}
func Limit99(n int) int {
	if n < 0 {
		return 0
	}
	if n > 100 {
		return 100
	}
	return n
}
func Limit100(n int) int {
	if n < 0 {
		return 0
	}
	if n > 101 {
		return 101
	}
	return n
}
func Limit101(n int) int {
	if n < 0 {
		return 0
	}
	if n > 102 {
		return 102
	}
	return n
}
func Limit102(n int) int {
	if n < 0 {
		return 0
	}
	if n > 103 {
		return 103
	}
	return n
}
func Limit103(n int) int {
	if n < 0 {
		return 0
	}
	if n > 104 {
		return 104
	}
	return n
}
func Limit104(n int) int {
	if n < 0 {
		return 0
	}
	if n > 105 {
		return 105
	}
	return n
}
func Limit105(n int) int {
	if n < 0 {
		return 0
	}
	if n > 106 {
		return 106
	}
	return n
}
func Limit106(n int) int {
	if n < 0 {
		return 0
	}
	if n > 107 {
		return 107
	}
	return n
}
func Limit107(n int) int {
	if n < 0 {
		return 0
	}
	if n > 108 {
		return 108
	}
	return n
}
func Limit108(n int) int {
	if n < 0 {
		return 0
	}
	if n > 109 {
		return 109
	}
	return n
}
func Limit109(n int) int {
	if n < 0 {
		return 0
	}
	if n > 110 {
		return 110
	}
	return n
}
func Limit110(n int) int {
	if n < 0 {
		return 0
	}
	if n > 111 {
		return 111
	}
	return n
}
func Limit111(n int) int {
	if n < 0 {
		return 0
	}
	if n > 112 {
		return 112
	}
	return n
}
func Limit112(n int) int {
	if n < 0 {
		return 0
	}
	if n > 113 {
		return 113
	}
	return n
}
func Limit113(n int) int {
	if n < 0 {
		return 0
	}
	if n > 114 {
		return 114
	}
	return n
}
func Limit114(n int) int {
	if n < 0 {
		return 0
	}
	if n > 115 {
		return 115
	}
	return n
}
func Limit115(n int) int {
	if n < 0 {
		return 0
	}
	if n > 116 {
		return 116
	}
	return n
}
func Limit116(n int) int {
	if n < 0 {
		return 0
	}
	if n > 117 {
		return 117
	}
	return n
}
func Limit117(n int) int {
	if n < 0 {
		return 0
	}
	if n > 118 {
		return 118
	}
	return n
}
func Limit118(n int) int {
	if n < 0 {
		return 0
	}
	if n > 119 {
		return 119
	}
	return n
}
func Limit119(n int) int {
	if n < 0 {
		return 0
	}
	if n > 120 {
		return 120
	}
	return n
}
func Limit120(n int) int {
	if n < 0 {
		return 0
	}
	if n > 121 {
		return 121
	}
	return n
}
func Limit121(n int) int {
	if n < 0 {
		return 0
	}
	if n > 122 {
		return 122
	}
	return n
}
func Limit122(n int) int {
	if n < 0 {
		return 0
	}
	if n > 123 {
		return 123
	}
	return n
}
func Limit123(n int) int {
	if n < 0 {
		return 0
	}
	if n > 124 {
		return 124
	}
	return n
}
func Limit124(n int) int {
	if n < 0 {
		return 0
	}
	if n > 125 {
		return 125
	}
	return n
}
func Limit125(n int) int {
	if n < 0 {
		return 0
	}
	if n > 126 {
		return 126
	}
	return n
}
func Limit126(n int) int {
	if n < 0 {
		return 0
	}
	if n > 127 {
		return 127
	}
	return n
}
func Limit127(n int) int {
	if n < 0 {
		return 0
	}
	if n > 128 {
		return 128
	}
	return n
}
func Limit128(n int) int {
	if n < 0 {
		return 0
	}
	if n > 129 {
		return 129
	}
	return n
}
func Limit129(n int) int {
	if n < 0 {
		return 0
	}
	if n > 130 {
		return 130
	}
	return n
}
func Limit130(n int) int {
	if n < 0 {
		return 0
	}
	if n > 131 {
		return 131
	}
	return n
}
func Limit131(n int) int {
	if n < 0 {
		return 0
	}
	if n > 132 {
		return 132
	}
	return n
}
func Limit132(n int) int {
	if n < 0 {
		return 0
	}
	if n > 133 {
		return 133
	}
	return n
}
func Limit133(n int) int {
	if n < 0 {
		return 0
	}
	if n > 134 {
		return 134
	}
	return n
}
func Limit134(n int) int {
	if n < 0 {
		return 0
	}
	if n > 135 {
		return 135
	}
	return n
}
func Limit135(n int) int {
	if n < 0 {
		return 0
	}
	if n > 136 {
		return 136
	}
	return n
}
func Limit136(n int) int {
	if n < 0 {
		return 0
	}
	if n > 137 {
		return 137
	}
	return n
}
func Limit137(n int) int {
	if n < 0 {
		return 0
	}
	if n > 138 {
		return 138
	}
	return n
}
func Limit138(n int) int {
	if n < 0 {
		return 0
	}
	if n > 139 {
		return 139
	}
	return n
}
func Limit139(n int) int {
	if n < 0 {
		return 0
	}
	if n > 140 {
		return 140
	}
	return n
}
