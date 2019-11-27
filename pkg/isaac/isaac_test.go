package isaac

import "testing"

func TestUint8(t *testing.T) {
	expectedValues := []uint8{91, 55, 20, 250, 148, 244, 57, 43, 89, 240, 137, 200, 245, 183, 207, 215, 99, 217, 16, 158, 52, 172, 180, 231, 252, 76, 9, 169, 101, 175, 202, 133}
	rng := New(make([]uint64, 256))
	for i, v := range expectedValues {
		next := rng.Uint8()
		if next != v {
			t.Errorf("*ISAAC.Uint8() failed, expected[%v]='%v', got '%v'", i, v, next)
		}
	}
}

func TestUint32(t *testing.T) {
	expectedValues := []uint32{2499033387, 4122464215, 883733735, 1706019461, 4061323680, 227296070, 2918044857, 2842052479, 2540425808, 252178577, 769461287, 2586658236, 1202489466, 2898528068, 3398674424, 623191030, 1645169217, 3068652366, 1390644596, 3122921300, 1220769609, 2431414785, 2790010922, 4223011841, 1323799079, 1960788473, 2287947294, 3078519590, 4015995844, 1490747927, 3861692908, 3257894546}
	rng := New(make([]uint64, 256))
	for i, v := range expectedValues {
		next := rng.Uint32()
		if next != v {
			t.Errorf("*ISAAC.Uint32() failed, expected[%v]='%v', got '%v'", i, v, next)
		}
	}
}

func TestUint16(t *testing.T) {
	expectedValues := []uint16{23351, 5370, 38132, 14635, 23024, 35272, 62903, 53207, 25561, 4254, 13484, 46311, 64588, 2473, 26031, 51845, 4330, 42730, 61970, 57760, 47918, 65317, 3468, 17222, 62775, 49126, 44525, 54457, 35922, 54116, 43366, 18303}
	rng := New(make([]uint64, 256))
	for i, v := range expectedValues {
		next := rng.Uint16()
		if next != v {
			t.Errorf("*ISAAC.Uint16() failed, expected[%v]='%v', got '%v'", i, v, next)
		}
	}
}

func TestUint64(t *testing.T) {
	expectedValues := []uint64{6572745247643941163, 6480831359995072471, 7194800151375688935, 18179916418947992197, 1218970177171022240, 13487998468627383110, 17669802660492858553, 10111376542692427647, 10807113810723983952, 674802807710675089, 6014335206068915239, 12612219875867051452, 13096941448193673338, 1624593600476350276, 6710029832301029368, 2967518739381429238, 13041473590980987457, 8006771451810540366, 3066781666356989300, 14325611978660511572, 12487995719355168585, 7613387012756500993, 5572840478758154282, 8554484584701889537, 291146565780214311, 16113903992955683321, 16547285933118934558, 2857822541202750246, 7280023713093729220, 3518176724902412823, 15660299115607933420, 17457197768181253778}
	rng := New(make([]uint64, 256))
	for i, v := range expectedValues {
		next := rng.Uint64()
		if next != v {
			t.Errorf("*ISAAC.Uint16() failed, expected[%v]='%v', got '%v'", i, v, next)
		}
	}
	for i := 0; i < 225; i++ {
		rng.Uint64()
	}
}

func TestInt(t *testing.T) {
	expectedValues := []int{6572745247643941163, 6480831359995072471, 7194800151375688935, 8956544382093216389, 1218970177171022240, 4264626431772607302, 8446430623638082745, 888004505837651839, 1583741773869208144, 674802807710675089, 6014335206068915239, 3388847839012275644, 3873569411338897530, 1624593600476350276, 6710029832301029368, 2967518739381429238, 3818101554126211649, 8006771451810540366, 3066781666356989300, 5102239941805735764, 3264623682500392777, 7613387012756500993, 5572840478758154282, 8554484584701889537, 291146565780214311, 6890531956100907513, 7323913896264158750, 2857822541202750246, 7280023713093729220, 3518176724902412823, 6436927078753157612, 8233825731326477970}
	rng := New(make([]uint64, 256))
	for i, v := range expectedValues {
		next := rng.Int()
		if next != v {
			t.Errorf("*ISAAC.Int() failed, expected[%v]='%v', got '%v'", i, v, next)
		}
	}
}

func TestInt31n(t *testing.T) {
	expectedValues := []int32{5818, 9598, 2057, 3972, 9456, 529, 6794, 6617, 5914, 587, 1791, 6022, 2799, 6748, 7913, 1450, 3830, 7144, 3237, 7271, 2842, 5661, 6496, 9832, 3082, 4565, 5327, 7167, 9350, 3470, 8991, 7585}
	rng := New(make([]uint64, 256))
	for i, v := range expectedValues {
		next := rng.Int31n(1e4)
		if next != v {
			t.Errorf("*ISAAC.Int31n(10000) failed, expected[%v]='%v', got '%v'", i, v, next)
		}
	}
}

func TestInt63n(t *testing.T) {
	expectedValues := []int64{941163, 72471, 688935, -559419, 22240, -168506, -693063, -123969, -567664, 675089, 915239, -500164, -878278, 350276, 29368, 429238, -564159, 540366, 989300, -40044, -383031, 500993, 154282, 889537, 214311, -868295, -617058, 750246, 729220, 412823, -618196, -297838}
	rng := New(make([]uint64, 256))
	for i, v := range expectedValues {
		next := rng.Int63n(1e6)
		if next != v {
			t.Errorf("*ISAAC.Int63n(1000000) failed, expected[%v]='%v', got '%v'", i, v, next)
		}
	}
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("*ISAAC.Int63n(-1) did not panic when expected!")
		}
	}()
	rng.Int63n(-1)
}

func TestInt63(t *testing.T) {
	expectedValues := []int64{6572745247643941163, 6480831359995072471, 7194800151375688935, -266827654761559419, 1218970177171022240, -4958745605082168506, -776941413216693063, -8335367531017123969, -7639630262985567664, 674802807710675089, 6014335206068915239, -5834524197842500164, -5349802625515878278, 1624593600476350276, 6710029832301029368, 2967518739381429238, -5405270482728564159, 8006771451810540366, 3066781666356989300, -4121132095049040044, -5958748354354383031, 7613387012756500993, 5572840478758154282, 8554484584701889537, 291146565780214311, -2332840080753868295, -1899458140590617058, 2857822541202750246, 7280023713093729220, 3518176724902412823, -2786444958101618196, -989546305528297838}
	rng := New(make([]uint64, 256))
	for i, v := range expectedValues {
		next := rng.Int63()
		if next != v {
			t.Errorf("*ISAAC.Int63() failed, expected[%v]='%v', got '%v'", i, v, next)
		}
	}
}

func TestIntn(t *testing.T) {
	expectedValues := []int{58185, 95983, 20576, 39721, 94560, 5292, 67941, 66171, 59148, 5871, 17915, 60225, 27997, 67486, 79131, 14509, 38304, 71447, 32378, 72711, 28423, 56610, 64960, 98324, 30822, 45653, 53270, 71677, 93504, 34709, 89912, 75853}
	rng := New(make([]uint64, 256))
	for i, v := range expectedValues {
		next := rng.Intn(1e5)
		if next != v {
			t.Errorf("*ISAAC.Intn(100000) failed, expected[%v]='%v', got '%v'", i, v, next)
		}
	}
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("*ISAAC.Int63n(-1) did not panic when expected!")
		}
	}()
	rng.Intn(-1)
}

func TestNextChar(t *testing.T) {
	expectedValues := []byte{87, 123, 51, 69, 121, 37, 96, 94, 88, 37, 49, 89, 58, 96, 107, 45, 68, 99, 62, 101, 59, 85, 93, 125, 61, 75, 82, 100, 120, 64, 117, 104}
	rng := New(make([]uint64, 256))
	for i, v := range expectedValues {
		next := rng.NextChar()
		if next != v {
			t.Errorf("*ISAAC.NextChar() failed, expected[%v]='%v', got '%v'", i, v, next)
		}
	}
}

func TestRead(t *testing.T) {
	expectedValues := []byte{91, 55, 20, 250, 148, 244, 57, 43, 89, 240, 137, 200, 245, 183, 207, 215, 99, 217, 16, 158, 52, 172, 180, 231, 252, 76, 9, 169, 101, 175, 202, 133, 16, 234, 166, 234, 242, 18, 225, 160, 187, 46, 255, 37, 13, 140, 67, 70, 245, 55, 191, 230, 173, 237, 212, 185, 140, 82, 211, 100, 169, 102, 71, 127}
	rng := New(make([]uint64, 256))
	buf := make([]byte, 64)
	if n, err := rng.Read(buf); n != 64 || err != nil {
		t.Errorf("*ISAAC.Read([64]byte) failed, length=%v, expected 64; err=%v, expected nil", n, err)
	}
	for i, v := range expectedValues {
		if buf[i] != v {
			t.Errorf("*ISAAC.Read([64]byte) failed, expected '%v', got '%v'", expectedValues, buf)
		}
	}
	var emptyBuf []byte
	if n, err := rng.Read(emptyBuf); n != 0 || err == nil {
		t.Errorf("*ISAAC.Read([0]byte) failed, length=%v, expected 0, err='%v'", n, err)
	}
}

func TestString(t *testing.T) {
	expected := "W{3Ey%`^x%1Y:`k-Dc>e;U]}=KRdx@uhr}*<T4ia7)2&m'g($npst]jCgJBa"
	rng := New(make([]uint64, 256))
	next := rng.String(60)
	if expected != next {
		t.Errorf("*ISAAC.String() failed, expected '%v', got '%v'", expected, next)
	}
}

func TestUint8n(t *testing.T) {
	expectedValues := []byte{91, 55, 20, 57, 43, 89, 16, 52, 76, 9, 16, 18, 46, 37, 13, 67, 70, 55, 82, 71, 80, 9, 93, 15, 7, 83, 54, 1, 45, 12, 39, 7}
	rng := New(make([]uint64, 256))
	for i, v := range expectedValues {
		next := rng.Uint8n(95)
		if next != v {
			t.Errorf("*ISAAC.Uint8n(95) failed, expected[%v]='%v', got '%v'", i, v, next)
		}
	}
}

func TestNew(t *testing.T) {
	expectedValues := []uint64{6572745247643941163, 6480831359995072471, 7194800151375688935, 18179916418947992197, 1218970177171022240, 13487998468627383110, 17669802660492858553, 10111376542692427647}
	rng := New(make([]uint64, 256))
	for i := 0; i < 8; i++ {
		next := rng.Uint64()
		if expectedValues[i] != next {
			t.Errorf("New([256]uint64{0,0,...}) failed, expected[%v]='%v', got '%v'", i, expectedValues[i], next)
		}
	}

	expectedValues = []uint64{14911831006260245106, 11813233674846717172, 17170476074744812420, 8670859143016805681, 1752840964919745038, 4876486772559178780, 14329438922335658793, 4960405231656654252}
	rng = New([]uint64{0})
	for i := 0; i < 8; i++ {
		next := rng.Uint64()
		if expectedValues[i] != next {
			t.Errorf("New([]uint64{0}) failed, expected[%v]='%v', got '%v'", i, expectedValues[i], next)
		}
	}

	expectedValues = []uint64{11365293634481549611, 2167809020880128771, 383993775173253251, 6310894772139438801, 8366073314422365142, 14250004781945377944, 14394499189877871481, 4641423784921619051}
	rng = New([]uint64{13371337, 73317331, 0xDEADBEEF, 0xBADDAD5AD50})
	for i := 0; i < 8; i++ {
		next := rng.Uint64()
		if expectedValues[i] != next {
			t.Errorf("New([]uint64{13371337, 73317331, 0xDEADBEEF, 0xBADDAD5AD50}) failed, expected[%v]='%v', got '%v'", i, expectedValues[i], next)
		}
	}
}
