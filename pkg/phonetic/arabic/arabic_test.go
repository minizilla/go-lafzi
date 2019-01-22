package arabic_test

import (
	"testing"

	ar "github.com/billyzaelani/go-lafzi/pkg/arabic"
	"github.com/billyzaelani/go-lafzi/pkg/phonetic/arabic"
)

func TestNormalizedUthmani(t *testing.T) {
	tables := []struct {
		s        []byte
		expected string
	}{
		{[]byte(string([]rune{ar.MaddahA, ar.AlefA, ar.SHLigatureSad, ar.SHLigatureQaf,
			ar.SHMeemInit, ar.SHLamAlef, ar.AlefWasla, ar.SHJeem, ar.SHThreeDots,
			ar.SHSeen, ar.RubElHizb, ar.SHURectZero, ar.SWaw, ar.SHMeemIsolated,
			ar.SLSeen, ar.Sajdah, ar.ECHStop, ar.HamzaA, ar.RHFCStop, ar.SLMeem, ar.Tatweel})),
			string([]rune{ar.Alef, ar.Hamza})},
		{[]byte("اقْتَرَبَ"), "إِقْتَرَبَ"},
		{[]byte("اقْرَ"), "إِقْرَ"},
		// Adz-Dzariyat(51) verse: 47
		{[]byte("بِأَيْي۟دٍۢ"), "بِأَيْدٍ"},
	}

	for _, table := range tables {
		actual := string(arabic.NormalizedUthmani(table.s)[:])
		if actual != table.expected {
			t.Errorf("expected: %s, actual: %s", table.expected, actual)
		}
	}
}

func TestRemoveSpace(t *testing.T) {
	tables := []struct {
		s        []byte
		expected string
	}{
		// Al-Fatihah(1) verse: 1
		{[]byte("بِسْمِ ٱللَّهِ ٱلرَّحْمَـٰنِ ٱلرَّحِيمِ"), "بِسْمِاللَّهِالرَّحْمَنِالرَّحِيمِ"},
		// Yasin(36) verse: 2
		{[]byte("وَٱلْقُرْءَانِ ٱلْحَكِيمِ"), "وَالْقُرْءَانِالْحَكِيمِ"},
		// Al-Baqarah(2) verse: 249
		{[]byte("وَٱللَّهُ مَعَ ٱلصَّـٰبِرِينَ"), "وَاللَّهُمَعَالصَّبِرِينَ"},
	}

	for _, table := range tables {
		table.s = arabic.NormalizedUthmani(table.s)
		actual := string(arabic.RemoveSpace(table.s)[:])
		if actual != table.expected {
			t.Errorf("expected: %s, actual: %s", table.expected, actual)
		}
	}
}

func TestRemoveShadda(t *testing.T) {
	tables := []struct {
		s        []byte
		expected string
	}{
		// Al-Fatihah(1) verse: 1
		{[]byte("بِسْمِ ٱللَّهِ ٱلرَّحْمَـٰنِ ٱلرَّحِيمِ"), "بِسْمِاللَهِالرَحْمَنِالرَحِيمِ"},
		// Yasin(36) verse: 2
		{[]byte("وَٱلْقُرْءَانِ ٱلْحَكِيمِ"), "وَالْقُرْءَانِالْحَكِيمِ"},
		// Al-Baqarah(2) verse: 249
		{[]byte("وَٱللَّهُ مَعَ ٱلصَّـٰبِرِينَ"), "وَاللَهُمَعَالصَبِرِينَ"},
	}

	for _, table := range tables {
		table.s = arabic.NormalizedUthmani(table.s)
		table.s = arabic.RemoveSpace(table.s)
		actual := string(arabic.RemoveShadda(table.s)[:])
		if actual != table.expected {
			t.Errorf("expected: %s, actual: %s", table.expected, actual)
		}
	}
}

func TestJoinConsonant(t *testing.T) {
	tables := []struct {
		s        []byte
		expected string
	}{
		// Al-Kahfi(19) verse: 16 (idgham mutamatsilain: meem + sukun + meem -> meem)
		{[]byte("لَكُم مِّنْ أَمْرِكُم مِّرْفَقًا"), "لَكُمِنْأَمْرِكُمِرْفَقًا"},
		// Al-A'raf(7) verse: 160
		{[]byte("ٱضْرِب بِّعَصَاكَ ٱلْحَجَر"), "اضْرِبِعَصَاكَالْحَجَر"},
	}

	for _, table := range tables {
		table.s = arabic.NormalizedUthmani(table.s)
		table.s = arabic.RemoveSpace(table.s)
		table.s = arabic.RemoveShadda(table.s)
		actual := string(arabic.JoinConsonant(table.s)[:])
		if actual != table.expected {
			t.Errorf("expected: %s, actual: %s", table.expected, actual)
		}
	}
}

func TestFixBoundary(t *testing.T) {
	tables := []struct {
		s        []byte
		expected string
	}{
		// Al-Fatihah(1) verse: 1 (harakat ending -> sukun)
		{[]byte("بِسْمِ ٱللَّهِ ٱلرَّحْمَـٰنِ ٱلرَّحِيمِ"), "بِسْمِالَهِالرَحْمَنِالرَحِيمْ"},
		// Al-Kahfi(19) verse: 16 (alef ending then fathatan -> fatha)
		{[]byte("لَكُم مِّنْ أَمْرِكُم مِّرْفَقًا"), "لَكُمِنْأَمْرِكُمِرْفَقَ"},
		// Al-A'raf(7) verse: 160 (consonant ending -> same, alef start -> alef hamza above + fatha + alef)
		{[]byte("ٱضْرِب بِّعَصَاكَ ٱلْحَجَر"), "أَاضْرِبِعَصَاكَالْحَجَر"},
		// Ar-Rum(30) verse: 21 (harakat ending then teh marbuta -> heh + sukun)
		{[]byte("مَّوَدَّةً وَرَحْمَةً"), "مَوَدَةًوَرَحْمَهْ"},
	}

	for _, table := range tables {
		table.s = arabic.NormalizedUthmani(table.s)
		table.s = arabic.RemoveSpace(table.s)
		table.s = arabic.RemoveShadda(table.s)
		table.s = arabic.JoinConsonant(table.s)
		actual := string(arabic.FixBoundary(table.s)[:])
		if actual != table.expected {
			t.Errorf("expected: %s, actual: %s", table.expected, actual)
		}
	}
}

func TestTanwinSub(t *testing.T) {
	tables := []struct {
		s        []byte
		expected string
	}{
		// Ar-Rum(30) verse: 21 (fathatan -> fatha + noon + sukun)
		{[]byte("مَّوَدَّةً وَرَحْمَةً"), "مَوَدَةَنْوَرَحْمَهْ"},
		// Al-Baqarah(2) verse: 25 (kasratan -> kasra + noon + sukun)
		{[]byte("جَنَّـٰتٍ تَجْرِ"), "جَنَتِنْتَجْرْ"},
		// Al-Baqarah(2) verse: 143 (dammatan -> damma + noon + sukun)
		{[]byte("لَرَءُوفٌ رَّحِيمٌ"), "لَرَءُوفُنْرَحِيمْ"},
	}

	for _, table := range tables {
		table.s = arabic.NormalizedUthmani(table.s)
		table.s = arabic.RemoveSpace(table.s)
		table.s = arabic.RemoveShadda(table.s)
		table.s = arabic.JoinConsonant(table.s)
		table.s = arabic.FixBoundary(table.s)
		actual := string(arabic.TanwinSub(table.s)[:])
		if actual != table.expected {
			t.Errorf("expected: %s, actual: %s", table.expected, actual)
		}
	}
}

func TestRemoveMadda(t *testing.T) {
	tables := []struct {
		s        []byte
		expected string
	}{
		// Al-Baqarah(2) verse: 143 (fatha + alef + non-harakat -> fatha + non-harakat)
		{[]byte("عَلَى ٱلنَّاسِ"), "عَلَىالنَسْ"},
		// Ali-Imran(3) verse: 44 (kasra + yeh + non-harakat -> kasra + non-harakat)
		{[]byte("غَلِيظَ"), "غَلِظْ"},
		// Al-Baqarah(2) verse: 143 (damma + waw + non-harakat -> damma + non-harakat)
		{[]byte("لَرَءُوفٌ رَّحِيمٌ"), "لَرَءُفُنْرَحِمْ"},
		// Al-Fatihah(1) verse: 5
		{[]byte("إِيَّاكَ نَعْبُدُ وَإِيَّاكَ نَسْتَعِينُ"), "إِيَكَنَعْبُدُوَإِيَكَنَسْتَعِنْ"},
	}

	for _, table := range tables {
		table.s = arabic.NormalizedUthmani(table.s)
		table.s = arabic.RemoveSpace(table.s)
		table.s = arabic.RemoveShadda(table.s)
		table.s = arabic.JoinConsonant(table.s)
		table.s = arabic.FixBoundary(table.s)
		table.s = arabic.TanwinSub(table.s)
		actual := string(arabic.RemoveMadda(table.s)[:])
		if actual != table.expected {
			t.Errorf("expected: %s, actual: %s", table.expected, actual)
		}
	}
}

func TestRemoveUnreadConsonant(t *testing.T) {
	tables := []struct {
		s        []byte
		expected string
	}{
		// Al-Baqarah(2) verse: 143 (double remove: alef & lam)
		{[]byte("عَلَى ٱلنَّاسِ"), "عَلَنَسْ"},
		// Ali-Imran(3) verse: 44 (single remove: yeh)
		{[]byte("غَلِيظَ"), "غَلِظْ"},
		// Al-Baqarah(2) verse: 102
		{[]byte("تَتْلُوا۟ ٱلشَّيَـٰطِينُ"), "تَتْلُشَيَطِنْ"},
	}

	for _, table := range tables {
		table.s = arabic.NormalizedUthmani(table.s)
		table.s = arabic.RemoveSpace(table.s)
		table.s = arabic.RemoveShadda(table.s)
		table.s = arabic.JoinConsonant(table.s)
		table.s = arabic.FixBoundary(table.s)
		table.s = arabic.TanwinSub(table.s)
		table.s = arabic.RemoveMadda(table.s)
		actual := string(arabic.RemoveUnreadConsonant(table.s)[:])
		if actual != table.expected {
			t.Errorf("expected: %s, actual: %s", table.expected, actual)
		}
	}
}

func TestIqlabSub(t *testing.T) {
	tables := []struct {
		s        []byte
		expected string
	}{
		// Al-Baqarah(2) verse: 253 (noon + sukun + beh -> meem + sukun + beh)
		{[]byte("مِنۢ بَعْدِهِم"), "مِمْبَعْدِهِم"},
	}

	for _, table := range tables {
		table.s = arabic.NormalizedUthmani(table.s)
		table.s = arabic.RemoveSpace(table.s)
		table.s = arabic.RemoveShadda(table.s)
		table.s = arabic.JoinConsonant(table.s)
		table.s = arabic.FixBoundary(table.s)
		table.s = arabic.TanwinSub(table.s)
		table.s = arabic.RemoveMadda(table.s)
		table.s = arabic.RemoveUnreadConsonant(table.s)
		actual := string(arabic.IqlabSub(table.s)[:])
		if actual != table.expected {
			t.Errorf("expected: %s, actual: %s", table.expected, actual)
		}
	}
}

func TestIdghamSub(t *testing.T) {
	tables := []struct {
		s        []byte
		expected string
	}{
		// Ali-Imran(3) verse: 145 (exception: dunya -> same)
		{[]byte("دُّنْيَا نُۭ"), "دُنْيَنْ"},
		// As-Saff(61) verse: 4 (exception: bunyan -> same)
		{[]byte("بُنْيَـٰنٌۭ مَّر"), "بُنْيَنُمَر"},
		// Ar-Ra'd(13) verse: 4 (exception: sinwan -> same)
		{[]byte("صِنْوَانٌۭ"), "صِنْوَنْ"},
		// Al-An'am(6) verse: 99 (exception: qinwan -> same)
		{[]byte("قِنْوَانٌۭ"), "قِنْوَنْ"},
		// Al-Qalam(68) verse: 1 (exception: nunwalqalam -> same)
		{[]byte("نٌ ۚ وَٱلْقَلَمِ وَمَا يَسْطُرُونَ"), "نُنْوَلْقَلَمِوَمَيَسْطُرُنْ"},
		// Al-Baqarah(2) verse: 2 (idgham bighunnah: noon + sukun + meem -> meem)
		// 						  (idgham bilaghunnah: noon + sukun + reh -> reh)
		{[]byte("هُدًۭى مِّن رَّبِّهِمْ ۖ"), "هُدَمِرَبِهِمْ"},
		// Al-Baqarah(2) verse: 7 (idgham bighunnah: noon + sukun + waw -> waw)
		{[]byte("غِشَـٰوَةٌۭ ۖ وَلَهُم"), "غِشَوَةُوَلَهُم"},
		// Al-Baqarah(2) verse: 8, (idgham bighunnah: noon + sukun + yeh -> yeh)
		{[]byte("مَن يَقُولُ"), "مَيَقُلْ"},
		// Al-Baqarah(2) verse: 12, (idgham bilaghunnah: noon + sukun + lam -> lam)
		{[]byte("وَلَـٰكِن لَّا يَشْعُرُونَ"), "وَلَكِلَيَشْعُرُنْ"},
		// Al-Baqarah(2) verse: 48 (idgham bighunnah: noon + sukun + noon -> noon)
		// note: noon + sukun + noon may be filtered out in joinConsonant
		{[]byte("عَنْ نَّفْسٍۢ"), "عَنَفْسْ"},
	}

	for _, table := range tables {
		table.s = arabic.NormalizedUthmani(table.s)
		table.s = arabic.RemoveSpace(table.s)
		table.s = arabic.RemoveShadda(table.s)
		table.s = arabic.JoinConsonant(table.s)
		table.s = arabic.FixBoundary(table.s)
		table.s = arabic.TanwinSub(table.s)
		table.s = arabic.RemoveMadda(table.s)
		table.s = arabic.RemoveUnreadConsonant(table.s)
		table.s = arabic.IqlabSub(table.s)
		actual := string(arabic.IdghamSub(table.s)[:])
		if actual != table.expected {
			t.Errorf("expected: %s, actual: %s", table.expected, actual)
		}
	}
}

func TestRemoveHarakat(t *testing.T) {
	tables := []struct {
		s        []byte
		expected string
	}{
		// Ar-Rum(30) verse: 21 (fathatan -> fatha + noon + sukun)
		{[]byte("مَّوَدَّةً وَرَحْمَةً"), "مودةورحمه"},
		// Al-Baqarah(2) verse: 25 (kasratan -> kasra + noon + sukun)
		{[]byte("جَنَّـٰتٍ تَجْرِ"), "جنتنتجر"},
		// Al-Baqarah(2) verse: 143 (dammatan -> damma + noon + sukun)
		{[]byte("لَرَءُوفٌ رَّحِيمٌ"), "لرءفرحم"},
	}

	for _, table := range tables {
		table.s = arabic.NormalizedUthmani(table.s)
		table.s = arabic.RemoveSpace(table.s)
		table.s = arabic.RemoveShadda(table.s)
		table.s = arabic.JoinConsonant(table.s)
		table.s = arabic.FixBoundary(table.s)
		table.s = arabic.TanwinSub(table.s)
		table.s = arabic.RemoveMadda(table.s)
		table.s = arabic.RemoveUnreadConsonant(table.s)
		table.s = arabic.IqlabSub(table.s)
		table.s = arabic.IdghamSub(table.s)
		actual := string(arabic.RemoveHarakat(table.s)[:])
		if actual != table.expected {
			t.Errorf("expected: %s, actual: %s", table.expected, actual)
		}
	}
}

func TestEncode(t *testing.T) {
	tables := []struct {
		s        []byte
		expected string
	}{
		// Al-Fatihah(1) verse: 1-7 non-vowel
		{[]byte("بِسْمِ ٱللَّهِ ٱلرَّحْمَـٰنِ ٱلرَّحِيمِ"), "BSMLHRHMNRHM"},
		{[]byte("ٱلْحَمْدُ لِلَّهِ رَبِّ ٱلْعَـٰلَمِين"), "XLHMDLLHRBLXLMN"},
		{[]byte("ٱلرَّحْمَـٰنِ ٱلرَّحِيم"), "XRHMNRHM"},
		{[]byte("مَـٰلِكِ يَوْمِ ٱلدِّين"), "MLKYWMDN"},
		{[]byte("إِيَّاكَ نَعْبُدُ وَإِيَّاكَ نَسْتَعِينُ"), "XYKNXBDWXYKNSTXN"},
		{[]byte("ٱهْدِنَا ٱلصِّرَٰطَ ٱلْمُسْتَقِيمَ"), "XHDNSRTLMSTKM"},
		{[]byte("صِرَٰطَ ٱلَّذِينَ أَنْعَمْتَ عَلَيْهِمْ غَيْرِ ٱلْمَغْضُوبِ عَلَيْهِمْ وَلَا ٱلضَّآلِّينَ"), "SRTLZNXNXMTXLYHMGYRLMGDBXLYHMWLDLN"},
	}

	for _, table := range tables {
		table.s = arabic.NormalizedUthmani(table.s)
		table.s = arabic.RemoveSpace(table.s)
		table.s = arabic.RemoveShadda(table.s)
		table.s = arabic.JoinConsonant(table.s)
		table.s = arabic.FixBoundary(table.s)
		table.s = arabic.TanwinSub(table.s)
		table.s = arabic.RemoveMadda(table.s)
		table.s = arabic.RemoveUnreadConsonant(table.s)
		table.s = arabic.IqlabSub(table.s)
		table.s = arabic.IdghamSub(table.s)
		table.s = arabic.RemoveHarakat(table.s)
		actual := string(arabic.Encode(table.s)[:])
		if actual != table.expected {
			t.Errorf("expected: %s, actual: %s", table.expected, actual)
		}
	}

	tables = []struct {
		s        []byte
		expected string
	}{
		// Al-Fatihah(1) verse: 1-7 with vowel
		{[]byte("بِسْمِ ٱللَّهِ ٱلرَّحْمَـٰنِ ٱلرَّحِيمِ"), "BSMLHRHMNRHM"},
		{[]byte("ٱلْحَمْدُ لِلَّهِ رَبِّ ٱلْعَـٰلَمِين"), "XLHMDLLHRBLXLMN"},
		{[]byte("ٱلرَّحْمَـٰنِ ٱلرَّحِيم"), "XRHMNRHM"},
		{[]byte("مَـٰلِكِ يَوْمِ ٱلدِّين"), "MLKYWMDN"},
		{[]byte("إِيَّاكَ نَعْبُدُ وَإِيَّاكَ نَسْتَعِينُ"), "XYKNXBDWXYKNSTXN"},
		{[]byte("ٱهْدِنَا ٱلصِّرَٰطَ ٱلْمُسْتَقِيمَ"), "XHDNSRTLMSTKM"},
		{[]byte("صِرَٰطَ ٱلَّذِينَ أَنْعَمْتَ عَلَيْهِمْ غَيْرِ ٱلْمَغْضُوبِ عَلَيْهِمْ وَلَا ٱلضَّآلِّينَ"), "SRTLZNXNXMTXLYHMGYRLMGDBXLYHMWLDLN"},
		// Al-Baqarah(2) verse: 20 (alef maksura bug)
		{[]byte("إِنَّ ٱللَّهَ عَلَىٰ كُلِّ شَىْءٍۢ قَدِيرٌۭ"), "XNLHXLKLSYXNKDR"},
		// Adz-Dzariyat(51) verse: 47
		{[]byte("وَٱلسَّمَآءَ بَنَيْنَـٰهَا بِأَيْي۟دٍۢ وَإِنَّا لَمُوسِعُونَ"), "WSMXBNYNHBXYDWXNLMSXN"},
	}

	for _, table := range tables {
		table.s = arabic.NormalizedUthmani(table.s)
		table.s = arabic.RemoveSpace(table.s)
		table.s = arabic.RemoveShadda(table.s)
		table.s = arabic.JoinConsonant(table.s)
		table.s = arabic.FixBoundary(table.s)
		table.s = arabic.TanwinSub(table.s)
		table.s = arabic.RemoveMadda(table.s)
		table.s = arabic.RemoveUnreadConsonant(table.s)
		table.s = arabic.IqlabSub(table.s)
		table.s = arabic.IdghamSub(table.s)
		table.s = arabic.RemoveHarakat(table.s)
		actual := string(arabic.Encode(table.s)[:])
		if actual != table.expected {
			t.Errorf("expected: %s, actual: %s", table.expected, actual)
		}
	}
}
